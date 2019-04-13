package main

/*
**使用方式：
**计算文件摘要
**    -path=文件完整路径
**计算目录dir下文件摘要，不递归
**    -path=dir/*
**    -path=dir\*
**计算目录dir下所有文件摘要，递归
**    -path=dir
**    -path=dir/
**    -path=dir\
**支持的算法
**    md5 sha1
*/

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"flag"
	"fmt"
	"hash"
	"os"
	"path/filepath"
	"strings"
)

func getMd5Hash() hash.Hash {
	return md5.New()
}

func getSha1Hash() hash.Hash {
	return sha1.New()
}

var hashGenMap map[string]func()hash.Hash = make(map[string]func()hash.Hash, 10)

func initHashGenMap() {
	hashGenMap["md5"] = getMd5Hash
	hashGenMap["sha1"] = getSha1Hash
}

var flagPath string
var flagAlgo string
var flagDiff string
func initFlags() {
	flag.StringVar(&flagPath, "path", "", "计算摘要的文件或目录。用法：\n" +
		"\t计算单个文件的摘要：\n" +
		"\t\t-path=文件\n" +
		"\t计算目录下所有文件的摘要，不递归：\n" +
		"\t\t-path=dir/*\n" +
		"\t\t-path=dir\\*\n" +
		"\t计算目录下所有文件的摘要，递归：\n" +
		"\t\t-path=dir\n" +
		"\t\t-path=dir/\n" +
		"\t\t-path=dir\\")
	flag.StringVar(&flagAlgo, "algo", "sha1", "摘要算法，支持md5和sha1")
	flag.StringVar(&flagDiff, "diff", "", "文件摘要与给定的字符串进行比较")
	flag.Parse()
}

func main() {
	initFlags()
	initHashGenMap()
	hash, ok := hashGenMap[flagAlgo]
	if !ok {
		fmt.Printf("error:unsuported hash algorithm:%s\n", flagAlgo)
		return
	}

	walkFn, root := GetWalkFunc(flagPath)
	walkFn(root, func(path string){
		hashVal := HashHandler(path, hash())
		hashStr := fmt.Sprintf("%x", string(hashVal))
		fmt.Printf("file:%s\n", path)
		fmt.Printf("hash:%s\n", hashStr)
		CompareHash(hashStr, flagDiff)
	})
}

type WalkFunc func(string, func(string))error
func GetWalkFunc(path string) (WalkFunc, string) {
	var walkFn WalkFunc = nil
	if strings.HasSuffix(path, "/*") {
		walkFn = WalkDirTree
		path = strings.TrimSuffix(path, "/*")
	} else if strings.HasSuffix(path, "\\*") {
		walkFn = WalkDirTree
		path = strings.TrimSuffix(path, "\\*")
	} else {
		walkFn = WalkDirTreeRecursive
	}
	return walkFn, path
}

func WalkDirTree(root string, walkFn func(string)) error {
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {return err}
		if info.IsDir() {
			if root == path {
				return nil
			} else {
				return filepath.SkipDir
			}
		}
		walkFn(path)
		return nil
	})
	return err
}

func WalkDirTreeRecursive(root string, walkFn func(string)) error {
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {return err}
		if info.IsDir() {return nil}
		walkFn(path)
		return nil
	})
	return err
}

func ReadFile(filename string, readFn func(val []byte, sz int)) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	buf := make([]byte, 3)
	bi := bufio.NewReader(file)
	for sz, _ := bi.Read(buf); sz > 0; {
		readFn(buf, sz)
		sz, _ = bi.Read(buf)
	}
	return nil
}

func HashHandler(filename string, hash hash.Hash) string {
	ReadFile(filename, func(val []byte, sz int){
		hash.Write(val)
	})
	val := hash.Sum(nil)
	return string(val)
}

func CompareHash(digHashStr string, diffStr string) {
	if diffStr == "" {
		return
	}
	if digHashStr == diffStr {
		fmt.Println("comp:pass")
	} else {
		fmt.Println("comp:diff")
	}
}