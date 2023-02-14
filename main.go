package main

import (
	"fmt"
	"runtime"

	"github.com/moqsien/gvc/pkgs/vctrl"
)

func main() {
	fmt.Println("hello world")
	h := vctrl.NewGoVersion()
	h.Run()
	fmt.Println(runtime.GOOS, runtime.GOARCH)
}
