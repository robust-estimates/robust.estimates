package main

import (
	"go/format"
	"go/parser"
	"go/token"
	"os"
)

func main() {
	for _, path := range os.Args[1:] {
		if path == "--" {
			continue
		}

		f, e := parser.ParseFile(token.NewFileSet(), path, nil, 0)
		if e != nil {
			panic(e)
		}

		out, e := os.Create(path)
		if e != nil {
			panic(e)
		}

		if e = format.Node(out, token.NewFileSet(), f); e != nil {
			panic(e)
		}

		if e = out.Close(); e != nil {
			panic(e)
		}
	}
}
