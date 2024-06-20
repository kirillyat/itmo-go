//go:build !solution

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	lines := make(map[string]int)
	for _, file := range os.Args[1:] {
		f, _ := os.Open(file)
		input := bufio.NewScanner(f)

		for input.Scan() {
			lines[input.Text()]++
		}

		f.Close()
	}
	for k, v := range lines {
		if v > 1 {
			fmt.Printf("%d\t%s\n", v, k)
		}
	}
}
