package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	timeout := flag.Int("timeout", 10, "timeout")

	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		log.Fatal("Missing host or port")
	}

	host := args[0]
	port := args[1]
	address := net.JoinHostPort(host, port)

	conn, err := net.DialTimeout("tcp", address, time.Duration(*timeout)*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to", address)

	defer func() {
		err = conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	done := make(chan struct{})
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGQUIT)

	go func() {
		reader := bufio.NewReader(conn)
		buf := make([]byte, 1024)

		for {
			n, err := reader.Read(buf)
			if err != nil {
				if err == io.EOF {
					fmt.Println("Server closed connection")
				} else {
					fmt.Fprintf(os.Stderr, "Error reading from server: %v\n", err)
				}
				close(done)
				return
			}

			if n > 0 {
				_, err = os.Stdout.Write(buf[:n])
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error writing to stdout: %v\n", err)
				}
			}
		}
	}()

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		writer := bufio.NewWriter(conn)

		for scanner.Scan() {
			text := scanner.Text() + "\n"
			_, err := writer.WriteString(text)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error writing: %v\n", err)
				close(done)
				return
			}
			err = writer.Flush()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error flushing: %v\n", err)
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
		}

		close(done)
	}()

	select {
	case <-done:
	case <-sigChan:
		fmt.Println("Received SIGQUIT, exiting...")
	}
}
