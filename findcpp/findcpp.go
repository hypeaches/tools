package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type CommandLine struct {
	src string
	recursive bool
	quote string
	lineEnd string
	diff string //todo: 实现diff
	suffix []string
}

var cmd CommandLine = CommandLine{".", false, "", ",", "", []string{}}
var diffFile map[string]bool = make(map[string]bool)

func main() {
	initCommandLine()
	readDiffFile()

	if cmd.recursive {
		walkDirectoryRecursively()
	} else {
		walkDirectory()
	}
}

func initCommandLine() {
	flag.StringVar(&cmd.src, "src", ".", "搜索目录")
	flag.BoolVar(&cmd.recursive, "recursive", true, "是否递归搜索目录下的所有文件")
	flag.StringVar(&cmd.quote, "quote", "'", "结果用哪种引号包裹，默认单引号")
	flag.StringVar(&cmd.lineEnd, "line-end", ",", "行尾结束符，默认逗号")
	flag.StringVar(&cmd.diff, "diff", "", "与指定文件进行比较，找出不同的行并输出")

	var flagSuffix string
	flag.StringVar(&flagSuffix, "suffix", ".cpp .cc .c", "搜索文件后缀名")
	flag.Parse()

	cmd.suffix = strings.Fields(flagSuffix)

	sep := string(os.PathSeparator)
	if !strings.HasSuffix(cmd.src, sep) {
		cmd.src = cmd.src + sep
	}
}

func readDiffFile() {
	if cmd.diff == "" {
		return;
	}
	f, err := os.Open(cmd.diff)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n')

		line = strings.TrimRight(line, "\n")
		diffFile[line] = true

		if err == io.EOF {
			return
		}
		if err != nil {
			panic(err)
		}
	}
}

func walkDirectory() {
	dir, err := ioutil.ReadDir(cmd.src)
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
	filepath.Walk(cmd.src, func(filename string, fi os.FileInfo, err error) error {
		if err != nil {
			panic(err)
		}
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
			filename = cmd.quote + strings.Replace(filename, cmd.src, "", -1) + cmd.quote + cmd.lineEnd
			_, ok := diffFile[filename]
			if !ok {
				fmt.Println(filename)
			}
		}
	}
}
