package main

import (
//	"flag"
	"fmt"
	"github.com/sakthipriyan/go/queue"
)

func main() {
	/*
	listen := flag.String("listen", "127.0.0.1:64001", "host:port")
	dir := flag.String("dir", "/tmp/go-queue/", "go queue directory")
	flag.Parse()
	queue.Serve(*listen, *dir)
	*/

	q,_ := queue.NewQueue("/tmp/queue/test1/url")
	wr := make(chan bool)
	rr := make(chan []byte)
	cr := make(chan bool)

	done := make(chan bool,2)
	go func(){
		for i := 1; i <= 10000; i++ {
			q.Write <- queue.QueueWrite{[]byte("Bytes hello!"), wr}
			fmt.Println(<- wr)
		}
		done <- true
	}()

	go func(){
		for i := 1; i <= 100; i++ {
			q.Read <- queue.QueueRead{rr}
			fmt.Println(string(<- rr))
		}
		done <- true
	}()
	<- done
	<- done
	q.Close <- cr
	<- cr
	fmt.Println("Queue Closed")
}
