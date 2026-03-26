package request

import (
	"fmt"
	"io"
	"strings"
)

const HTTP_LINE_SEPERATOR = "\r\n"

type RequestState int

const (
	initialized RequestState = iota
	done
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
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
		switch r.state {
		case initialized:
			rl, n, err := parseRequestLine(string(data[read:]))
			if err != nil {
				return 0, err
			}

			// not enought data yet
			if n == 0 {
				return 0, err
			}

			r.RequestLine = *rl
			read += n
			r.state = done
		case done:
			return read, nil
		}
	}

}

func (r *Request) done() bool {
	return r.state == done
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := Request{RequestLine{}, initialized}
	// NOTE: can overrun buffer with len more than 1kb
	buf := make([]byte, 1024)
	bufLen := 0

	for !request.done() {
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
