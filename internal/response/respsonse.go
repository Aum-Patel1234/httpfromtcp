package response

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"httpfromtcp/internal/request"
	"io"
)

type Response struct {
}
type HandlerError struct {
	StatusCode StatusCode
	Message    string
}

// Handler has a type of func which returns *HandlerError
type Handler func(w *Writer, req *request.Request)

type StatusCode int

const (
	StatusOk                  = 200
	StatusBadRequest          = 400
	StatusInternalServerError = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	var err error
	switch statusCode {
	case StatusOk:
		_, err = w.Write([]byte("HTTP/1.1 200 OK\r\n"))
	case StatusBadRequest:
		_, err = w.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
	case StatusInternalServerError:
		_, err = w.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
	default:
		_, err = w.Write([]byte(""))
	}

	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", fmt.Sprint(contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return *h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	var err error = nil
	for k, v := range headers.All() {
		_, err = fmt.Fprintf(w, "%s: %s\r\n", k, v)
	}
	fmt.Fprintf(w, "\r\n")
	return err
}

// NOTE: we are writing our own writer to make the users of our api do what they want
type Writer struct {
	writer io.Writer
}

func NewWriter(writer io.Writer) *Writer {
	return &Writer{
		writer: writer,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	var err error
	switch statusCode {
	case StatusOk:
		_, err = w.writer.Write([]byte("HTTP/1.1 200 OK\r\n"))
	case StatusBadRequest:
		_, err = w.writer.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
	case StatusInternalServerError:
		_, err = w.writer.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
	default:
		_, err = w.writer.Write([]byte(""))
	}

	return err
}
func (w *Writer) WriteHeaders(headers headers.Headers) error {
	var err error = nil
	for k, v := range headers.All() {
		_, err = fmt.Fprintf(w.writer, "%s: %s\r\n", k, v)
	}
	fmt.Fprintf(w.writer, "\r\n")
	return err
}
func (w *Writer) WriteBody(p []byte) (int, error) {
	n, err := w.writer.Write(p)

	return n, err
}
