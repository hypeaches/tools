package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type CommandLine struct {
	input string
	recursive bool
	quote string
	lineEnd string
	diff string //todo: 实现diff
	suffix []string
}

var cmd CommandLine = CommandLine{".", false, "", ",", "", []string{}}

func main() {
	initCommandLine()

	if cmd.recursive {
		walkDirectoryRecursively()
	} else {
		walkDirectory()
	}
}

func initCommandLine() {
	flag.StringVar(&cmd.input, "input", ".", "搜索目录")
	flag.BoolVar(&cmd.recursive, "recursive", true, "是否递归搜索目录下的所有文件")
	flag.StringVar(&cmd.quote, "quote", "'", "结果用哪种引号包裹，默认单引号")
	flag.StringVar(&cmd.lineEnd, "line-end", ",", "行尾结束符，默认逗号")

	var flagSuffix string
	flag.StringVar(&flagSuffix, "suffix", ".cpp .cc .c", "搜索文件后缀名")
	flag.Parse()

	cmd.suffix = strings.Fields(flagSuffix)

	sep := string(os.PathSeparator)
	if !strings.HasSuffix(cmd.input, sep) {
		cmd.input = cmd.input + sep
	}
}

func walkDirectory() {
	dir, err := ioutil.ReadDir(cmd.input)
	if err != nil {
		panic(err)
	}

	for _, fi := range dir {
		if fi.IsDir() {
			continue
		}
		outputFileName(fi.Name())
	}
}

func walkDirectoryRecursively() {
	filepath.Walk(cmd.input, func(filename string, fi os.FileInfo, err error) error {
		if fi.IsDir() {
			return nil
		}
		outputFileName(filename)
		return nil
	})
}

func outputFileName(filename string) {
	for _, suffix := range cmd.suffix {
		if strings.HasSuffix(filename, suffix) {
			filename = cmd.quote + strings.Replace(filename, cmd.input, "", -1) + cmd.quote + cmd.lineEnd
			fmt.Println(filename)
		}
	}
}
