package filesystem

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

const basePath = "fileSystem/"

func GetLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string, 1)

	go func() {
		defer f.Close()
		defer close(out)

		line := ""
		for {
			data := make([]byte, 8)
			_, err := f.Read(data)
			if err != nil {
				break
			}

			if idx := bytes.IndexByte(data, '\n'); idx != -1 {
				line += string(data[:idx])
				out <- line
				data = data[idx+1:]
				line = ""
			}

			line += string(data)
		}
		if len(line) != 0 {
			out <- line
		}
	}()

	return out
}

func FileIO() {
	file, err := os.Open(basePath + "messages.txt")
	if err != nil {
		log.Fatal("failed to read the file")
	}

	for {
		data := make([]byte, 8)
		n, err := file.Read(data)
		if err != nil {
			fmt.Printf(err.Error()) // EOF
			break
		}
		fmt.Printf("Read : %s %d\n", data, n)
	}
}

func ReadLine() {
	file, err := os.Open(basePath + "messages.txt")
	if err != nil {
		log.Fatal("failed to read the file", err.Error())
	}
	defer file.Close()

	line := ""
	for {
		data := make([]byte, 8)
		_, err := file.Read(data)
		if err != nil {
			break
		}

		if idx := bytes.IndexByte(data, '\n'); idx != -1 {
			line += string(data[:idx])
			fmt.Printf("read: %s\n", line)
			data = data[idx+1:]
			line = ""
		}

		line += string(data)
	}
	if len(line) != 0 {
		fmt.Printf("read: %s\n", line)
	}
}
