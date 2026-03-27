package networking

import (
	"fmt"
	"httpfromtcp/internal/request"
	"log"
	"net"
)

func Listen() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("error", "error", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("error", "error", err)
		}

		// lines := filesystem.GetLinesChannel(conn)
		lines, err := request.RequestFromReader(conn)
		if err != nil {
			log.Println("Error while reading from tcp connection: ", err)
			conn.Close()
			return
		}
		// for line := range lines {
		// 	fmt.Printf("read: %s\n", line)
		// }
		fmt.Println("Request Line:")
		fmt.Printf(" - Method: %s\n - Target: %s\n - Version: %s\n", lines.RequestLine.Method, lines.RequestLine.RequestTarget, lines.RequestLine.HttpVersion)
		fmt.Println("Headers:")
		for k, v := range lines.Headers.All() {
			fmt.Printf("- %s : %s\n", k, v)
		}
		fmt.Println("Body:")
		fmt.Println(string(lines.Body))
	}
}
