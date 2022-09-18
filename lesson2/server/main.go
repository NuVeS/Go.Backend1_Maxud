package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
)

var (
	status   = make(chan string)
	messages = make(chan string)
)
var connections = make(map[net.Conn]bool)

func main() {
	listener, err := net.Listen("tcp", ":8080")

	if err != nil {
		fmt.Println("error listerning")
	} else {
		fmt.Println("listerning")
	}

	// signal.Notify(stop, os.Interrupt)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("wrong connection %e", err)
			continue
		} else {
			fmt.Println("Connected")
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	connections[conn] = true

	ctx, cancelFunc := context.WithCancel(context.Background())
	context := context.WithValue(ctx, "conn", conn)
	go broadcasting(context)

	who := conn.RemoteAddr().String()
	status <- who + " joined"

	fmt.Printf("New connection %s", who)

	input := bufio.NewScanner(conn)
	for {
		input.Scan()
		text := input.Text()
		if text == "exit" {
			delete(connections, conn)
			status <- who + " left"
			cancelFunc()
			conn.Close()
		} else {
			messages <- who + ": " + input.Text()
		}
	}
}

func broadcasting(context context.Context) {
	conn := context.Value("conn").(net.Conn)
	for {
		select {
		case <-context.Done():
			return
		case msg := <-status:
			fmt.Println(msg)
			fmt.Fprintln(conn, msg)
		case msg := <-messages:
			fmt.Println(msg)
			fmt.Fprintln(conn, msg)
		}
	}
}
