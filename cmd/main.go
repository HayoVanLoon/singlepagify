package main

import (
	"fmt"
	"github.com/HayoVanLoon/singlepagify"
	"golang.org/x/net/html"
	"os"
	"strings"
)

func usage() {
	fmt.Printf(`Usage:
  %s <input html file> <output hmtl file>
`, os.Args[0])
}

func main() {
	if len(os.Args) != 3 {
		usage()
		os.Exit(3)
	}
	dir, file := splitInput(os.Args[1])
	out, err := singlepagify.Process(dir, file)
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
	output, err := os.Create(os.Args[2])
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
	if err = html.Render(output, out); err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
	fmt.Println("Results written to " + os.Args[2])
}

func splitInput(s string) (string, string) {
	idx := strings.LastIndex(s, string(os.PathSeparator))
	if idx < 0 {
		return ".", s
	}
	return s[:idx], s[idx+1:]
}
