package main

import (
	"fmt"
	"github.com/disintegration/gift"
	"runtime"
)

func main() {
	_ = gift.New()
	fmt.Println(runtime.NumCPU())
}
