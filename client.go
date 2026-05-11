package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/quic-go/quic-go"
)

func main() {

	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-chat"},
	}

	conn, err := quic.DialAddr(
		context.Background(),
		"127.0.0.1:4242",
		tlsConf,
		&quic.Config{
			KeepAlivePeriod: 20 * time.Second,
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	stream, err := conn.OpenStreamSync(
		context.Background(),
	)

	if err != nil {
		log.Fatal(err)
	}

	go receiveMessages(stream)

	scanner := bufio.NewScanner(os.Stdin)

	for {

		scanner.Scan()

		text := scanner.Text()

		_, err := stream.Write([]byte(text + "\n"))

		if err != nil {
			fmt.Println("Send Error:", err)
			return
		}
	}
}

func receiveMessages(stream *quic.Stream) {

	reader := bufio.NewReader(stream)

	for {

		msg, err := reader.ReadString('\n')

		if err != nil {
			fmt.Println("Disconnected from server")
			return
		}

		fmt.Print(msg)
	}
}
