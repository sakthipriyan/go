package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var conns = struct {
	sync.RWMutex
	m map[int]net.Conn
}{m: make(map[int]net.Conn)}

func main() {

	serverIP := flag.String("server", "127.0.0.1:8081", "ListenIP:Port")
	flag.Parse()

	signalChan := make(chan os.Signal, 2)
	doneChan := make(chan bool)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	fmt.Printf("Launching server at %s\n", *serverIP)
	ln, err := net.Listen("tcp", *serverIP)

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		<-signalChan
		fmt.Println("Stopping server...")
		ln.Close()
		doneChan <- true
	}()

	go func() {
		for i := 0; ; i++ {
			conn, err := ln.Accept()
			if err != nil {
				fmt.Println(err)
				continue
			}
			go handleConnection(conn, i)
		}
	}()

	<-doneChan

}

func handleConnection(conn net.Conn, id int) {

	scanner := bufio.NewScanner(conn)
	if !scanner.Scan() {
		return
	}

	conns.Lock()
	conns.m[id] = conn
	conns.Unlock()

	name := scanner.Text()
	publish(id, name, "\b Joined")

	for {
		if !scanner.Scan() {
			conn.Close()
			conns.Lock()
			delete(conns.m, id)
			conns.Unlock()
			publish(id, name, "\b Left")
			return
		}
		text := scanner.Text()
		publish(id, name, text)
	}
}

func publish(id int, name, text string) {
	ts := time.Now().Format(time.Kitchen)
	fmt.Printf("[%s/%s]:%s \n", name, ts, text)
	conns.RLock()
	for k, v := range conns.m {
		if k != id {
			fmt.Fprintf(v, "[%s/%s]:%s \n", name, ts, text)
		}
	}
	conns.RUnlock()
}
