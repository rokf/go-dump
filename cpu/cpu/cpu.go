package cpu

// Registers
const (
	R0 = iota
	R1
	R2
	R3
	R4
	R5
	R6
	R7
)

// instructions
const (
	CLF = iota
	MOV
	STI
	LDI
	LII
	PSH
	POP
	INC
	DEC
	// math
	ADD
	SUB
	MUL
	DIV
	// jump
	JLZ
	JGZ
	JEZ
	JNZ
	JMP
	// shift
	SHL
	SHR
	// binary ops
	BAND
	BOR
	BNOT
	BXOR
	HLT // halt
)

// The CPU structure holds the state
type CPU struct {
	// the loaded program
	Memory []int64
	// program counter
	pc int64
	// stack pointer
	sp int64
	// Register array
	Reg [8]int64
	// current instruction
	inst int64
	// destination
	dest int64
	// source
	src int64
	// flags
	zero bool
	ltz  bool
	gtz  bool
}

func (c *CPU) clearFlags() {
	c.zero = false
	c.ltz = false
	c.gtz = false
}

func (c *CPU) setFlags(d int64) {
	c.zero = d == 0
	c.ltz = d < 0
	c.gtz = d > 0
}

// NewCPU creates a new CPU instance
func NewCPU(mem []int64) (c CPU) {
	c = CPU{
		Memory: mem,
		pc:     -1,
		sp:     int64(len(mem) - 1),
	}
	return
}

func (c *CPU) fetch() {
	c.pc++
	c.inst = c.Memory[c.pc]
	if c.pc+1 < int64(len(c.Memory)) {
		c.dest = c.Memory[c.pc+1]
		c.src = c.Memory[c.pc+2]
	} else {
		c.dest = 0
		c.src = 0
	}
}

func (c *CPU) execute() {
	switch c.inst {
	case CLF:
		c.clearFlags()
	case MOV:
		c.Reg[c.dest] = c.Reg[c.src]
		c.pc += 2
	case STI:
		c.Memory[c.dest] = c.Reg[c.src]
		c.pc += 2
	case LDI:
		c.Reg[c.dest] = c.Memory[c.src]
		c.pc += 2
	case LII:
		c.Reg[c.dest] = c.src
		c.pc += 2
	case PSH:
		c.sp--
		c.pc++
		c.Memory[c.sp] = c.Reg[c.Memory[c.pc]]
	case POP:
		c.pc++
		c.Reg[c.Memory[c.pc]] = c.Memory[c.sp]
		c.sp++
	case INC:
		c.Reg[c.dest]++
		c.pc++
	case DEC:
		c.Reg[c.dest]--
		c.pc++
	case ADD:
		c.Reg[c.dest] += c.Reg[c.src]
		c.pc += 2
	case SUB:
		c.Reg[c.dest] -= c.Reg[c.src]
		c.pc += 2
	case MUL:
		c.Reg[c.dest] *= c.Reg[c.src]
		c.pc += 2
	case DIV:
		c.Reg[c.dest] /= c.Reg[c.src]
		c.pc += 2
	case JLZ:
		if c.ltz {
			c.pc = c.Memory[c.pc+1]
		} else {
			c.pc++
		}
	case JGZ:
		if c.gtz {
			c.pc = c.Memory[c.pc+1]
		} else {
			c.pc++
		}
	case JEZ:
		if c.zero {
			c.pc = c.Memory[c.pc+1]
		} else {
			c.pc++
		}
	case JNZ:
		if !c.zero {
			c.pc = c.Memory[c.pc+1]
		} else {
			c.pc++
		}
	case JMP:
		c.pc = c.Memory[c.pc+1]
	case SHL:
		c.Reg[c.dest] <<= uint64(c.Reg[c.src])
		c.pc += 2
	case SHR:
		c.Reg[c.dest] >>= uint64(c.Reg[c.src])
		c.pc += 2
	case BAND:
		c.Reg[c.dest] &= c.Reg[c.src]
		c.pc += 2
	case BOR:
		c.Reg[c.dest] |= c.Reg[c.src]
		c.pc += 2
	case BNOT:
		c.Reg[c.dest] = ^c.Reg[c.dest]
		c.pc++
	case BXOR:
		c.Reg[c.dest] ^= c.Reg[c.src]
		c.pc += 2
	}

}

// Run runs the loaded program
func (c *CPU) Run() {
	for c.inst != HLT {
		c.fetch()
		c.execute()
	}
}
