package queue

type Message struct {
	action string // WRITE, READ
	ts uint64 //Timestamp
	crc uint32 //crc
	size uint32 //size
	data []byte //payload
}

type 
