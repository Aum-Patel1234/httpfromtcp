package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func main() {
	server, err := server.Serve(port, func(w *response.Writer, req *request.Request) {
		headers := response.GetDefaultHeaders(0)
		headers.Replace("Content-Type", "text/html")

		var body []byte
		var status response.StatusCode

		switch req.RequestLine.RequestTarget {
		case "/aum":
			status = response.StatusBadRequest
			body = []byte(`<html>
				<head><title>400 Bad Request</title></head>
				<body>
					<h1>Aum is the GOAT</h1>
				</body>
			</html>`)

		case "/aumlang":
			status = response.StatusInternalServerError
			body = []byte(`<html>
				<head><title>500 Internal Server Error</title></head>
				<body>
					<h1>aumlang is the GAOT programming language</h1>
				</body>
			</html>`)

		default:
			status = response.StatusOk
			body = []byte(`<html>
				<head><title>200 OK</title></head>
				<body>
					<h1>Success!</h1>
				</body>
			</html>`)
		}

		headers.Replace("Content-Length", fmt.Sprintf("%d", len(body)))

		w.WriteStatusLine(status)
		w.WriteHeaders(headers)
		w.WriteBody(body)
	})
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
