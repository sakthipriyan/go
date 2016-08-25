package main

/*
import (
	"fmt"
	"github.com/sakthipriyan/go/gofile"
)

func main() {
    gofile.Start("/tmp/gofile")
	fmt.Println(gofile.NewDataIn(123, []byte("key"), []byte("value")))
}

*/

import (
    "os"
    "log"
)

func main() {
    // Open a new file for writing only
    file, err := os.OpenFile(
        "test.txt",
        os.O_RDWR|os.O_CREATE,
        0666,
    )
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    // Write bytes to file
    byteSlice := []byte("Bytes hello!\n")
    bytesWritten, err := file.Write(byteSlice)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Wrote %d bytes.\n", bytesWritten)
	newPosition, err := file.Seek(0,0)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Just moved to 5:", newPosition)

	byteSlice = make([]byte, 16)
    bytesRead, err := file.Read(byteSlice)
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Number of bytes read: %d\n", bytesRead)
    log.Printf("Data read: %s\n", byteSlice)


}
