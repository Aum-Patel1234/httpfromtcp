package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers struct {
	headers map[string]string
}

const SEPERATOR = "\r\n"

func isToken(s []byte) bool {
	if len(s) == 0 {
		return false
	}

	for i := range s {
		c := s[i]

		// ALPHA (A-Z, a-z)
		if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') {
			continue
		}

		// DIGIT (0-9)
		if c >= '0' && c <= '9' {
			continue
		}

		// tchar special characters
		switch c {
		case '!', '#', '$', '%', '&', '\'', '*',
			'+', '-', '.', '^', '_', '`', '|', '~':
			continue
		}

		// anything else → invalid
		return false
	}

	return true
}

func parserHeader(fieldLine []byte) (string, string, error) { // key, val, err
	parts := bytes.SplitN(fieldLine, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("Malformed Header")
	}

	key := parts[0]
	value := bytes.TrimSpace(parts[1])
	if bytes.HasSuffix(key, []byte(" ")) {
		return "", "", fmt.Errorf("Malformed Field Line")
	}

	return string(key), string(value), nil
}

func (h *Headers) Get(key string) (string, bool) {
	str, ok := h.headers[strings.ToLower(key)]
	return str, ok
}

func (h *Headers) Replace(key, value string) {
	name := strings.ToLower(key)
	h.headers[name] = value
}

func (h *Headers) Delete(key string) {
	name := strings.ToLower(key)
	delete(h.headers, name)
}

func (h *Headers) Set(key, value string) {
	name := strings.ToLower(key)
	val, ok := h.headers[name]
	// NOTE: if header aldready present add it and seperate via comma(csv)
	if ok {
		value = val + "," + value
	}
	h.headers[name] = value
}
func (h *Headers) All() map[string]string {
	return h.headers
}

func (h Headers) Parse(data []byte) (int, bool, error) {
	read := 0
	done := false
	for {
		idx := bytes.Index(data[read:], []byte(SEPERATOR))
		if idx == -1 {
			break
		}
		// EMPTY HEADER
		if idx == 0 {
			done = true
			read += len(SEPERATOR)
			break
		}

		key, value, err := parserHeader(data[read : read+idx])
		if err != nil {
			return 0, done, err
		}

		if !isToken([]byte(key)) {
			return 0, false, fmt.Errorf("Malformed header name")
		}

		read += idx + len(SEPERATOR)
		h.Set(key, value)
	}

	return read, done, nil
}

func NewHeaders() *Headers {
	return &Headers{
		headers: map[string]string{},
	}
}
