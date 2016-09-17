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

func main() {

	serverIP := flag.String("listen", "127.0.0.1:8081", "host:port")
	flag.Parse()

	doneChan := make(chan bool)
	signalChan := make(chan os.Signal, 2)
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
		conns := &Connections{m: make(map[int]net.Conn)}
		for i := 0; ; i++ {
			conn, err := ln.Accept()
			if err != nil {
				fmt.Println(err)
				continue
			}
			go handleConnection(conns, conn, i)
		}
	}()

	<-doneChan

}

func handleConnection(conns *Connections, conn net.Conn, id int) {
	scanner := bufio.NewScanner(conn)
	if !scanner.Scan() {
		return
	}
	name := scanner.Text()
	conns.add(id, name, conn)
	for {
		if !scanner.Scan() {
			conns.remove(id, name)
			return
		}
		text := scanner.Text()
		conns.publish(id, name, text)
	}
}

type Connections struct {
	sync.RWMutex
	m map[int]net.Conn
}

func (conns *Connections) add(id int, name string, conn net.Conn) {
	conns.Lock()
	conns.m[id] = conn
	conns.Unlock()
	conns.publish(id, name, "\b Joined")
	fmt.Println(conns.m)
}

func (conns *Connections) remove(id int, name string) {
	conns.Lock()
	conns.m[id].Close()
	delete(conns.m, id)
	conns.Unlock()
	conns.publish(id, name, "\b Left")
}

func (conns *Connections) publish(id int, name, text string) {
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
