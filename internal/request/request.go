package request

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"strings"
)

const HTTP_LINE_SEPERATOR = "\r\n"

type RequestState int

const (
	StateInitialized RequestState = iota
	StateHeaders
	StateDone
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	Headers     *headers.Headers
	state       RequestState
}

func parseRequestLine(b string) (*RequestLine, int, error) {
	idx := strings.Index(b, HTTP_LINE_SEPERATOR)
	if idx == -1 {
		return nil, 0, nil
	}

	startLine := b[:idx]
	// restLine := b[idx+len(HTTP_LINE_SEPERATOR):]

	parts := strings.Split(startLine, " ")
	if len(parts) != 3 {
		return nil, idx + len(HTTP_LINE_SEPERATOR), fmt.Errorf("Malformed request line.")
	}

	if !strings.HasPrefix(parts[2], "HTTP/") {
		return nil, idx + len(HTTP_LINE_SEPERATOR), fmt.Errorf("invalid HTTP version")
	}
	version := strings.TrimPrefix(parts[2], "HTTP/")

	return &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   version,
	}, idx + len(HTTP_LINE_SEPERATOR), nil
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0

	for {
		currentData := data[read:]

		switch r.state {
		case StateInitialized:
			rl, n, err := parseRequestLine(string(currentData))
			if err != nil {
				return 0, err
			}

			// not enought data yet
			if n == 0 {
				return read, err
			}

			r.RequestLine = *rl
			read += n
			r.state = StateHeaders
		case StateHeaders:
			n, done, err := r.Headers.Parse(currentData)
			if err != nil {
				return 0, err
			}
			if n == 0 {
				return read, nil
			}

			read += n
			if done {
				r.state = StateDone
			}

		case StateDone:
			return read, nil
		default:
			panic("Something bad happened")
		}
	}

}

func (r *Request) StateDone() bool {
	return r.state == StateDone
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := Request{RequestLine{}, headers.NewHeaders(), StateInitialized}
	// NOTE: can overrun buffer with len more than 1kb
	buf := make([]byte, 1024)
	bufLen := 0

	for !request.StateDone() {
		n, err := reader.Read(buf[bufLen:])
		if err != nil {
			return nil, err
		}

		bufLen += n
		readN, err := request.parse(buf[:bufLen])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[readN:bufLen])
		bufLen -= readN
	}

	return &request, nil
}
