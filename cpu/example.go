package main

import (
	"fmt"

	"github.com/rokf/go-dump/cpu/cpu"
)

func main() {

	program := []int64{
		cpu.LII, cpu.R0, 1,
		cpu.LII, cpu.R1, -1,
		cpu.ADD, cpu.R0, cpu.R1,
		cpu.HLT,
	}

	c := cpu.NewCPU(program)
	c.Run()

	fmt.Println(c.Reg[cpu.R0])
}
