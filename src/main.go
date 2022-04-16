package main

import (
	"commander/src/osutil"
	"fmt"
)

func main() {
	src := "./src/main.go"
	dst := "./dst/main.go.bak"

	err := osutil.CopyAll(src, dst, false, true)
	if err != nil {
		fmt.Println(err)
	}

}
