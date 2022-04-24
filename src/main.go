package main

import (
	"commander/src/osutil"
	"fmt"
	"os"
)

func main() {
	src := "./src"
	dst := "./dst"
	var i int = 0
	co := osutil.CopyOptions{
		Recursive:     true,
		Depth:         &(i),
		FileMode:      nil,
		DirectoryMode: nil,
		FilterFunc:    nil,
	}

	c := osutil.NewCopy(src, dst, &co)
	_ = c

	fmt.Println(os.Getenv("Niyazi"))
}
