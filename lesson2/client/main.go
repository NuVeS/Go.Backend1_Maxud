package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	go func() {
		io.Copy(os.Stdout, conn)
	}()
	fmt.Fprintln(conn, "I did it")
	io.Copy(conn, os.Stdin)
	fmt.Fprintf(conn, "exit")
}
