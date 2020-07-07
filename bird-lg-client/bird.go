package main

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"net/http"
)

func birdCopy(dst io.Writer, src io.Reader) {
	buf := make([]byte, 65536) // Too small for `show route`
	src.Read(buf)
	index := bytes.IndexByte(buf, 0)
	if index > 0 {
		buf = buf[:index]
	}
	str := bytes.NewReader(buf)
	br := bufio.NewReader(str)
	for {
		line, err := br.ReadBytes('\n')
		if err != nil && len(line) == 0 {
			break
		}
		if dst == nil {
			continue
		}
		if len(line) < 4 {
			dst.Write(line[1:])
		} else {
			status_code := line[:4]
			if isDigit(string(status_code)) {
				dst.Write(line[5:])
			} else {
				dst.Write(line[1:])
			}
		}
	}
}

func birdHandler(httpW http.ResponseWriter, httpR *http.Request) {
	query := string(httpR.URL.Query().Get("q"))
	if query == "" {
		invalidHandler(httpW, httpR)
	} else {
		bird, err := net.Dial("unix", birdSocket)
		if err != nil {
			panic(err)
		}
		defer bird.Close()
		birdCopy(nil, bird)
		bird.Write([]byte("restrict\n"))
		birdCopy(nil, bird)
		bird.Write([]byte(query + "\n"))
		birdCopy(httpW, bird)
	}
}
