package main

import (
	"fmt"
	"github.com/sakthipriyan/go/gofile"
)

func main() {
    gofile.Start("/tmp/gofile")
	fmt.Println(gofile.NewDataIn(123, []byte("key"), []byte("value")))
}
