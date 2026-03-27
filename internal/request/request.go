package request

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"strconv"
	"strings"
)

const HTTP_LINE_SEPERATOR = "\r\n"

type RequestState int

const (
	StateInitialized RequestState = iota
	StateHeaders
	StateBody
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
	Body        []byte
	state       RequestState
}

func getInt(headers headers.Headers, name string, defautlValue int) int {
	val, exists := headers.Get(name)
	if !exists {
		return defautlValue
	}

	value, err := strconv.Atoi(val)
	if err != nil {
		return defautlValue
	}

	return value
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
		if len(currentData) == 0 {
			break
		}

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
				if r.hasBody() {
					r.state = StateBody
				} else {
					r.state = StateDone
				}
			}

		case StateBody:
			contentLength := getInt(*r.Headers, "content-length", 0)
			if contentLength == 0 {
				panic("chunk encoding not impelemented")
			}

			remaining := min(contentLength-len(r.Body), len(currentData))
			r.Body = append(r.Body, currentData[:remaining]...)
			read += remaining

			if len(r.Body) == contentLength {
				r.state = StateDone
			}

		case StateDone:
			return read, nil
		default:
			panic("Something bad happened")
		}
	}

	return read, nil
}

func (r *Request) hasBody() bool {
	// TODO: when doing chunk encoding update this
	length := getInt(*r.Headers, "content-length", 0)
	return length > 0
}
func (r *Request) StateDone() bool {
	return r.state == StateDone
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := Request{RequestLine{}, headers.NewHeaders(), []byte(""), StateInitialized}
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
