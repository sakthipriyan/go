package queue

import (
	"io/ioutil"
	"log"
	"strconv"
)

func WriteOffset(filepath string, offset uint64) {
	str := strconv.FormatUint(offset, 10)
	err := ioutil.WriteFile(filepath, []byte(str), 0600)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Wrote Offset", offset, "in", filepath)
}

func ReadOffset(filepath string) uint64 {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}

	offset, err := strconv.ParseUint(string(data), 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Read Offset", offset, "from", filepath)
	return offset
}
