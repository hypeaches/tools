package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Counter struct {
	code                 int32
	space                int32
	annotation           int32
	multiAnnotationStart bool
}

var counter Counter = Counter{0, 0, 0, false}

func main() {
	flag.Parse()
	start := time.Now()
	err := filepath.Walk(flag.Arg(0), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			panic(err)
		}
		if info.IsDir() {
			return nil
		}
		if FilteFile(path) {
			CountCode(path)
		}
		return nil
	})
	end := time.Now()
	elapsed := end.Sub(start)
	if err != nil {
		panic(err)
	}
	fmt.Println("      code: ", counter.code)
	fmt.Println("annotation: ", counter.annotation)
	fmt.Println("     space: ", counter.space)
	fmt.Println("     total: ", counter.code+counter.annotation+counter.space)
	fmt.Println("      time: ", elapsed)
}

func FilteFile(path string) bool {
	return strings.HasSuffix(path, ".h") ||
		strings.HasSuffix(path, ".hpp") ||
		strings.HasSuffix(path, ".cpp") ||
		strings.HasSuffix(path, ".cc") ||
		strings.HasSuffix(path, ".c") ||
		strings.HasSuffix(path, ".go") ||
		strings.HasSuffix(path, ".java")
}

func CountCode(path string) {
	ReadLine(path, func(linebytes []byte) {
		linebytes = bytes.TrimSpace(linebytes)
		line := string(linebytes[:])
		if counter.multiAnnotationStart {
			counter.annotation += 1
			if strings.HasSuffix(line, "*/") {
				counter.multiAnnotationStart = false
			}
			return
		}
		if strings.HasPrefix(line, "/*") {
			counter.annotation += 1
			counter.multiAnnotationStart = true
			return
		}
		if strings.HasPrefix(line, "//") {
			counter.annotation += 1
			return
		}
		if line == "" {
			counter.space += 1
			return
		}
		counter.code += 1
	})
}

func ReadLine(filename string, hookfn func([]byte)) {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	rd := bufio.NewReader(file)
	for {
		line, err := rd.ReadBytes('\n')
		if err == io.EOF {
			return
		}
		if err != nil {
			panic(err)
		}

		hookfn(line)
	}
}
