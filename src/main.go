package main

import (
	"commander/src/osutil"
	"fmt"
)

func main() {
	src := "./src"
	dst := "./dst"
	var i int = 2
	co := osutil.CopyOptions{
		Recursive:     true,
		Depth:         &(i),
		FileMode:      nil,
		DirectoryMode: nil,
		FilterFunc:    nil,
	}

	c := osutil.NewCopy(src, dst, &co)
	_ = c

	fmt.Println(c.Copy())

}
