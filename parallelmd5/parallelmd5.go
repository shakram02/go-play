package main

import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
)

func walk(rootDir string) <-chan string {
	out := make(chan string)

	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.Mode().IsRegular() {
			return nil
		}

		out <- path
		return nil
	}

	go func() {
		defer close(out)
		filepath.Walk(rootDir, walkFunc)
	}()

	return out
}

type File struct {
	path    string
	content []byte
}

type FileMd5 struct {
	path   string
	digest string
}

func readFile(pathChan <-chan string) <-chan File {
	out := make(chan File)

	go func() {
		defer close(out)

		for path := range pathChan {
			content, err := os.ReadFile(path)
			if err != nil {
				continue
			}

			out <- File{
				path:    path,
				content: content,
			}
		}
	}()

	return out
}

func md5File(fileChan <-chan File) <-chan FileMd5 {
	out := make(chan FileMd5)

	go func() {
		defer close(out)

		for f := range fileChan {
			digest := md5.Sum(f.content)
			out <- FileMd5{
				path:   f.path,
				digest: fmt.Sprintf("%x", digest),
			}
		}
	}()

	return out
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage: parrallelmd5.go PATH\n")
		os.Exit(1)
	}

	rootDir := os.Args[1]
	out := md5File(readFile(walk(rootDir)))

	for o := range out {
		fmt.Printf("%s\t%s\n", o.digest, o.path)
	}

}
