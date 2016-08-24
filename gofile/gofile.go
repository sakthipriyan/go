package gofile

import (
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

type Index struct {
	id     uint64
	offset uint64
	size   uint32
}

type DataIn struct {
	id        uint64
	ts        uint64
	crc       uint32
	keySize   uint32
	valueSize uint32
	key       []byte
	value     []byte
}

type DataOut struct {
	id    uint64
	ts    uint64
	key   []byte
	value []byte
}

func NewDataIn(id uint64, key, value []byte) (*DataIn, *Index) {
	ts := uint64(time.Now().UnixNano() / 1000000)
	crc := crc32.ChecksumIEEE(value)
	keySize := uint32(len(key))
	valueSize := uint32(len(value))
	dataIn := DataIn{id, ts, crc, keySize, valueSize, key, value}
	index := Index{id, 0, 28 + keySize + valueSize}
	return &dataIn, &index
}

func NewDataOut(data []byte) {

}

func write() {

}

func read() {

}

func Start(dir string) {
	prepareFiles(dir)
	offset := GetOffset(offsetFileName(dir))
    fmt.Println(offset)

}

func prepareFiles(dir string) {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		createDir(dir)
	}
}

func createDir(dir string) {
	err := os.MkdirAll(dir, 0700)
	if err != nil {
		log.Fatal("Unable to use directory ", dir, err)
	}
	log.Println("Created directory", dir)

	for _, filename := range []string{"index", "data"} {
		file, err := os.Create(filename)
		if err != nil {
			log.Fatal(err)
		}
		file.Close()
	}
	SetOffset(offsetFileName(dir), 0)
}

func SetOffset(filename string, offset uint64) {
	str := strconv.FormatUint(offset, 10)
	err := ioutil.WriteFile(filename, []byte(str), 0600)
	if err != nil {
		log.Fatal(err)
	}
    log.Println("Set Offset", offset)
}

func GetOffset(filename string) uint64 {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	offset, err := strconv.ParseUint(string(data), 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Get Offset", offset)
	return offset
}

func offsetFileName(dir string) string {
	return dir + string(os.PathSeparator) + "offset"
}
