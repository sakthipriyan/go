package main

import (
	"fmt"
	"github.com/sakthipriyan/go/queue"
)

/*
func main() {
    gofile.Start("/tmp/gofile")
	fmt.Println(gofile.NewDataIn(123, []byte("key"), []byte("value")))
}

*/

func main() {
	dir := "/tmp/goqueue"
	q := queue.Queue{}
	q.Open(dir)
	q.Write([]byte("Bytes hello!"))
	fmt.Println(string(q.Read()))
	q.Close()
}
