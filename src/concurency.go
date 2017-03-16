package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"sync"
)

func main() {
	// go echo(os.Stdin, os.Stdout)
	// time.Sleep(30 * time.Second)
	// fmt.Println("Timed out.")
	// os.Exit(0)

	// fmt.Println("Outside a goroutine.")
	// go func() {
	// 	fmt.Println("Inside a goroutine")
	// }()
	// fmt.Println("Outside again.")
	// runtime.Gosched()
	var wg sync.WaitGroup
	var i = -1
	var file string
	for i, file = range os.Args[1:] {
		wg.Add(1)
		go func(filename string) {
			compress(filename)
			wg.Done()
		}(file)
	}
	wg.Wait()
	fmt.Println("Compressed %d files\n", i+1)
}

func echo(in io.Reader, out io.Writer) {
	io.Copy(out, in)
}

func compress(filename string) error {
	in, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(filename + ".gz")
	if err != nil {
		return err
	}
	defer out.Close()
	gzout := gzip.NewWriter(out)
	_, err = io.Copy(gzout, in)
	gzout.Close()
	return err
}
