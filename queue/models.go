package queue

import (
	"hash/crc32"
	"time"
    "bytes"
    "encoding/binary"
    "fmt"
)

type Index struct {
	id     uint64
	offset uint64
}

type Data struct {
	id   uint64
	ts   uint64
	crc  uint32
	size uint32
	data []byte
}

func readHeader(b []byte) *Data{
    var id, ts uint64
    var crc, size uint32
    buf := bytes.NewReader(b)

    err := binary.Read(buf, binary.LittleEndian, &id)
    if err != nil {
        fmt.Println("binary.Read failed:", err)
    }

    err = binary.Read(buf, binary.LittleEndian, &ts)
    if err != nil {
        fmt.Println("binary.Read failed:", err)
    }

    err = binary.Read(buf, binary.LittleEndian, &crc)
    if err != nil {
        fmt.Println("binary.Read failed:", err)
    }

    err = binary.Read(buf, binary.LittleEndian, &size)
    if err != nil {
        fmt.Println("binary.Read failed:", err)
    }

    return &Data {id,ts,crc,size,[]byte{}}
}

func newData(id uint64, data []byte) *Data {
	ts := uint64(time.Now().UnixNano() / 1000000)
	crc := crc32.ChecksumIEEE(data)
	size := uint32(len(data))
	return &Data{id, ts, crc, size, data}
}

func binaryData(id uint64, src []byte) []byte {
        data := newData(id, src)
        fmt.Println("Data:", *data)
        buf := new(bytes.Buffer)
        col := []interface{} {
            data.id,data.ts,data.crc,data.size,data.data,
        }

        for _,c := range col {
            err := binary.Write(buf, binary.LittleEndian, c)
            if err != nil {
                fmt.Println("binary.Write failed:", err)
            }
        }
        return buf.Bytes()
}
