package queue

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

type Queue struct {
	baseDir     string
	readOffset  uint64
	writeOffset uint64
	indexFile   *os.File
	dataFile    *os.File
}

func (q *Queue) Open(dir string) error {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0700)
		if err != nil {
			log.Fatal(err)
		}
		q.readOffset = 0
	} else {
		q.readOffset = ReadOffset(filepath.Join(dir, "offset"))
	}
	q.baseDir = dir
	q.indexFile = openFile(dir, "index")
	q.dataFile = openFile(dir, "data")

	log.Println("Opening the queue file")
	return nil
}

func (q *Queue) Read() []byte {

	log.Println("Reading the queue file")
	_, err := q.dataFile.Seek(io.SeekStart, 0)
	if err != nil {
		log.Fatal(err)
	}
	buf := make([]byte, 24)
	_, err = q.dataFile.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Data read:", buf)
	data := readHeader(buf)
	log.Println("Data read:", *data)

	buf = make([]byte, data.size)
	_, err = q.dataFile.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Data read:", buf)

	q.readOffset++
	return buf
}

func (q *Queue) Write(data []byte) {
	log.Println("Writing to queue", data)
	data = binaryData(q.writeOffset, data)
	log.Println("Writing to file", data)
	_, err := q.dataFile.Seek(0, io.SeekStart)
	if err != nil {
		log.Fatal(err)
	}
	_, err = q.dataFile.Write(data)
	if err != nil {
		log.Fatal(err)
	}
}

func (q *Queue) Close() {
	q.indexFile.Close()
	q.dataFile.Close()
	offset := filepath.Join(q.baseDir, "offset")
	WriteOffset(offset, q.readOffset)
}

func openFile(baseDir string, filename string) *os.File {
	file, err := os.OpenFile(
		filepath.Join(baseDir, filename),
		os.O_RDWR|os.O_CREATE,
		0666,
	)
	if err != nil {
		log.Fatal(err)
	}
	return file
}
