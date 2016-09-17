package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var doneChan = make(chan bool)

func main() {

	serverIP := flag.String("dial", "127.0.0.1:8081", "host:port")
	flag.Parse()

	signalChan := make(chan os.Signal, 2)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	fmt.Println("To exit type 'Ctrl + C'")

	conn, err := net.DialTimeout("tcp", *serverIP, 3*time.Second)

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		<-signalChan
		conn.Close()
		fmt.Println("CLIENT DISCONNECTED")
		doneChan <- true
	}()
	go clientInput(conn)
	<-doneChan
}

func clientInput(conn net.Conn) {
	fmt.Print("Enter Name:")
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	name := scanner.Text()
	fmt.Fprintf(conn, name+"\n")
	go serverOutput(name, conn)
	for scanner.Scan() {
		text := scanner.Text()
		if len(text) != 0 {
			fmt.Fprintf(conn, text+"\n")
		}
	}
}

func serverOutput(name string, conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for {
		if !scanner.Scan() {
			fmt.Println("SERVER DISCONNECTED")
			doneChan <- true
			return
		}
		fmt.Println(scanner.Text())
	}
}
