package queue

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"fmt"
)

type QueueRead struct {
    Resp chan []byte
}

type QueueWrite struct {
    Data  []byte
    Resp chan bool
}

type Queue struct {
	Read chan QueueRead
	Write chan QueueWrite
	Close chan chan bool
	baseDir    string
	offsetId   uint64
	nextOffset uint64
	indexFile  *os.File
	dataFile   *os.File
}

func NewQueue(dir string) (Queue, error) {
	q := Queue{}
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0700)
		if err != nil {
			return q, err
		}
		q.nextOffset = 0
	} else {
		q.nextOffset = ReadOffset(filepath.Join(dir, "offset"))
	}
	q.baseDir = dir
	q.indexFile = openFile(dir, "index")
	q.dataFile = openFile(dir, "data")
	q.Read = make(chan QueueRead)
	q.Write = make(chan QueueWrite)
	q.Close = make(chan chan bool)
	go q.run()
	log.Println("Opening the queue file")
	return q,nil
}

func (q *Queue) read() ([]byte, error) {

	log.Println("Reading the queue file")
	buf := make([]byte, 24)
	_, err = q.dataFile.ReadAt(buf,int64(q.nextOffset))
	if err != nil {
		return nil, err
	}
	data := readHeader(buf)
	log.Println("Data read:", *data)

	buf = make([]byte, data.size)
	_, err = q.dataFile.Read(buf)
	if err != nil {
		return nil, err
	}
	log.Println("Data read:", buf)

	q.nextOffset += uint64(24) + uint64(data.size)
	return buf, nil
}

func (q *Queue) write(data []byte) {
	log.Println("Writing to queue", data)
	data = binaryData(q.offsetId, data)
	q.offsetId++
	log.Println("Writing to file", data)
	_, err := q.dataFile.Seek(0, io.SeekEnd)
	if err != nil {
		log.Fatal(err)
	}
	_, err = q.dataFile.Write(data)
	if err != nil {
		log.Fatal(err)
	}
}

func (q *Queue) close() {
	q.indexFile.Close()
	q.dataFile.Close()
	close(q.Read)
	close(q.Write)
	close(q.Close)
	offset := filepath.Join(q.baseDir, "offset")
	WriteOffset(offset, q.nextOffset)
}

func (q *Queue) run(){
	for {
		select {
		case r := <- q.Read:
			data, err := q.read()
			if err != nil {
				fmt.Println(err)
				r.Resp <- nil
			} else {
				r.Resp <- data
			}
		case w := <- q.Write:
			log.Println("Writing")
			q.write(w.Data)
			w.Resp <- true
		case c:= <- q.Close:
			q.close()
			c <- true
			return
		}
	}
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
