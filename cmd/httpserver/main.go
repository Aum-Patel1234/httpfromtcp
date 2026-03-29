package main

import (
	"crypto/sha256"
	"fmt"
	"httpfromtcp/internal/headers"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const port = 42069
const aumBody = "<html><head><title>400 Bad Request</title></head><body><h1>Aum is the GOAT</h1></body></html>"

const aumLangBody = "<html><head><title>500 Internal Server Error</title></head><body><h1>aumlang is the GOAT programming language</h1></body></html>"
const defaultBody = `<html>
<head><title>200 OK</title></head>
<body>
<h1>Success!</h1>
</body>
</html>`

func toStr(bytes []byte) string {
	str := ""
	for _, b := range bytes {
		str += fmt.Sprintf("%02x", b)
	}
	return str
}

func main() {
	server, err := server.Serve(port, func(w *response.Writer, req *request.Request) {
		h := response.GetDefaultHeaders(0)
		h.Replace("Content-Type", "text/html")

		var body []byte
		var status response.StatusCode

		path := req.RequestLine.RequestTarget
		switch {
		case path == "/aum":
			status = response.StatusBadRequest
			body = []byte(aumBody)

		case path == "/aumlang":
			status = response.StatusInternalServerError
			body = []byte(aumLangBody)

		// IMPORTANT: video RequestTarget
		case path == "/video":
			f, _ := os.ReadFile("assets/vim.mp4")

			h.Replace("Content-Type", "video/mp4")
			h.Replace("Content-Length", fmt.Sprintf("%d", len(f)))

			w.WriteStatusLine(response.StatusOk)
			w.WriteHeaders(h)
			w.WriteBody(f)
			return

		case strings.HasPrefix(path, "/httpbin/stream"):
			res, err := http.Get("https://httpbin.org/" + path[len("/httpbin/"):])
			if err != nil {
				body = []byte(defaultBody)
				status = response.StatusInternalServerError
			} else {
				w.WriteStatusLine(response.StatusOk)

				h.Delete("Content-Length")
				h.Set("transfer-encoding", "chunked")
				h.Replace("Content-Type", "text/plain")
				w.WriteHeaders(h)

				for {
					data := make([]byte, 32)
					n, err := res.Body.Read(data)
					if err != nil {
						break
					}

					w.WriteBody(fmt.Appendf(nil, "%x\r\n", n))
					w.WriteBody(data[:n])
					w.WriteBody([]byte("\r\n"))
				}
				w.WriteBody([]byte("0\r\n\r\n"))
				return
			}
		// NOTE: trailers included here
		case strings.HasPrefix(path, "/httpbin/html"):
			res, err := http.Get("https://httpbin.org/" + path[len("/httpbin/"):])
			if err != nil {
				body = []byte(defaultBody)
				status = response.StatusInternalServerError
			} else {
				w.WriteStatusLine(response.StatusOk)

				h.Delete("Content-Length")
				h.Set("transfer-encoding", "chunked")
				h.Set("Trailer", "X-Content-SHA256")
				h.Set("Trailer", "X-Content-Length")
				h.Replace("Content-Type", "text/plain")
				w.WriteHeaders(h)

				fullBody := []byte{}
				for {
					data := make([]byte, 32)
					n, err := res.Body.Read(data)
					if err != nil {
						break
					}

					fullBody = append(fullBody, data[:n]...)
					w.WriteBody(fmt.Appendf(nil, "%x\r\n", n))
					w.WriteBody(data[:n])
					w.WriteBody([]byte("\r\n"))
				}
				w.WriteBody([]byte("0\r\n"))
				trailers := headers.NewHeaders()
				sha := sha256.Sum256(fullBody)
				trailers.Set("X-Content-SHA256", toStr(sha[:]))
				// trailers.Set("X-Content-SHA256", hex.EncodeToString(sha[:]))
				trailers.Set("X-Content-Length", fmt.Sprintf("%d", len(fullBody)))
				w.WriteHeaders(*trailers)
				return
			}

		default:
			status = response.StatusOk
			body = []byte(defaultBody)
		}

		h.Replace("Content-Length", fmt.Sprintf("%d", len(body)))

		w.WriteStatusLine(status)
		w.WriteHeaders(h)
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
