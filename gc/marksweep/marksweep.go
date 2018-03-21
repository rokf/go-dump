package marksweep

import "fmt"

const SMAX = 256

const (
	INT = iota
	PAIR
)

type Object struct {
	Type   uint
	Marked bool
	Next   *Object
	Value  int
	Head   *Object
	Tail   *Object
}

type VM struct {
	Stack     [SMAX]*Object
	StackSize int
	First     *Object
	Allocated int
	Max       int
}

func (vm *VM) Init() {
	vm.StackSize = 0
	vm.First = nil
	vm.Allocated = 0
	vm.Max = 8
}

func (vm *VM) Push(obj *Object) {
	if vm.StackSize >= SMAX {
		panic("Overflow!")
	}
	vm.Stack[vm.StackSize] = obj
	vm.StackSize++
}

func (vm *VM) Pop() *Object {
	if vm.StackSize <= 0 {
		panic("Underflow!")
	}
	vm.StackSize--
	obj := vm.Stack[vm.StackSize]
	return obj
}

func (obj *Object) Mark() {
	if obj.Marked {
		return
	}

	obj.Marked = true

	if obj.Type == PAIR {
		obj.Head.Mark()
		obj.Tail.Mark()
	}
}

func (vm *VM) MarkAll() {
	for i := 0; i < vm.StackSize; i++ {
		vm.Stack[i].Mark()
	}
}

func (vm *VM) Sweep() {
	objPtr := vm.First
	for {
		if objPtr == nil {
			return
		}
		if !objPtr.Marked {
			vm.Allocated--
		} else {
			objPtr.Marked = false
		}
		objPtr = objPtr.Next
	}
}

func (vm *VM) GC() {
	allocated := vm.Allocated
	vm.MarkAll()
	vm.Sweep()
	vm.Max = vm.Allocated * 2
	fmt.Println("collected:", allocated-vm.Allocated, "remaining:", vm.Allocated)
}

func (vm *VM) NewObject(objType uint) *Object {
	if vm.Allocated == vm.Max {
		vm.GC()
	}

	obj := &Object{
		Type:   objType,
		Next:   vm.First,
		Marked: false,
	}

	vm.First = obj

	vm.Allocated++

	return obj
}

func (vm *VM) PushINT(value int) {
	obj := vm.NewObject(INT)
	obj.Value = value

	vm.Push(obj)
}

func (vm *VM) PushPAIR() {
	obj := vm.NewObject(PAIR)
	h := vm.Pop()
	t := vm.Pop()
	obj.Head = h
	obj.Tail = t

	vm.Push(obj)
}

// temporary
func (vm *VM) PrintStack() {
	for i := 0; i < vm.StackSize; i++ {
		fmt.Println(vm.Stack[i], vm.Stack[i].Next)
	}
}

func (vm *VM) Free() {
	vm.StackSize = 0
	vm.GC()
}
