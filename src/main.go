package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	src := "./src/a.yml"
	abs, err := filepath.Abs(src)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	dir := filepath.Dir(abs)
	println(dir)
}
