package queue

import (
	"log"
	"os"
	"path/filepath"
)

type Queue struct {
	baseDir    string
	readOffset uint64
	indexFile  *os.File
	dataFile   *os.File
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

func (q *Queue) Read() {
	log.Println("Reading the queue file")
	byteSlice := make([]byte, 16)
	bytesRead, err := q.dataFile.Read(byteSlice)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Number of bytes read: %d\n", bytesRead)
	log.Printf("Data read: %s\n", byteSlice)
	q.readOffset++
}

func (q *Queue) Write() {
	log.Println("Writing the queue file")
	// Write bytes to file
	byteSlice := []byte("Bytes hello!\n")
	log.Println(q)
	bytesWritten, err := q.dataFile.Write(byteSlice)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Wrote %d bytes.\n", bytesWritten)
	newPosition, err := q.dataFile.Seek(0, 0)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Just moved to 5:", newPosition)
}

func (q *Queue) Close() {
	q.indexFile.Close()
	q.dataFile.Close()
	offset := filepath.Join(q.baseDir, "offset")
	log.Println(offset)
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
