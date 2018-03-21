package main

import (
	"fmt"

	"github.com/rokf/go-dump/gc/marksweep"
)

func main() {
	vm := &marksweep.VM{}
	vm.Init()
	vm.PushINT(1)
	vm.PushINT(2)
	vm.PushINT(3)

	vm.GC()
	if vm.Allocated != 3 {
		fmt.Println("Should've let them live")
	}
}
