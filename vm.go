package main

import (
	"fmt"
)

type GenesysVM struct {
	memory  [1000]uint16
	regs    [10]uint16
	ip      uint16
	ic      uint64
	ilimit  uint64
	running bool
}

func NewGenesysVM(code []uint16, limit uint64) *GenesysVM {
	var vm GenesysVM

	if len(code) > len(vm.memory) {
		panic("Executable does not fit in memory")
	}

	copy(vm.memory[:len(code)], code)
	vm.ilimit = limit
	vm.running = true

	return &vm
}

func (vm *GenesysVM) GetStats() (ic uint64, running bool) {
	return vm.ic, vm.running
}

func (vm *GenesysVM) Panic(s string) {
	var err error

	if vm.ip < 1000 {
		err = fmt.Errorf("%s: registers: %v, ip = %d, ic = %d, instruction = %d", s, vm.regs, vm.ip, vm.ic, vm.memory[vm.ip])
	} else {
		err = fmt.Errorf("%s: registers: %v, ip = %d, ic = %d, INVALID IP", s, vm.regs, vm.ip, vm.ic)
	}

	vm.running = false
	panic(err)
}

func (vm *GenesysVM) Step() {
	var instruction, d1, d2, d3, addr uint16

	if vm.ic > vm.ilimit {
		vm.Panic("Execution timeout - max limit reached")
	}

	if vm.ip > 999 {
		vm.Panic("IP out of bounds")
	}

	vm.ic++
	instruction = vm.memory[vm.ip]

	if instruction > 999 {
		vm.Panic("Invalid instruction")
	}

	d1 = instruction / 100
	d2 = (instruction / 10) % 10
	d3 = instruction % 10

	switch d1 {
	case 1:
		// HALT
		if instruction == 100 {
			vm.running = false
			return
		} else {
			vm.Panic("Invalid instruction - HALT does not take parameters")
		}
	case 2:
		// MOV register, value
		vm.regs[d2] = d3
		vm.ip++
	case 3:
		// ADD register, value
		vm.regs[d2] = (vm.regs[d2] + d3) % 1000
		vm.ip++
	case 4:
		// MUL register, value
		vm.regs[d2] = (vm.regs[d2] * d3) % 1000
		vm.ip++
	case 5:
		// MOV register1, register2
		vm.regs[d2] = vm.regs[d3]
		vm.ip++
	case 6:
		// ADD register1, register2
		vm.regs[d2] = (vm.regs[d2] + vm.regs[d3]) % 1000
		vm.ip++
	case 7:
		// MUL register1, register2
		vm.regs[d2] = uint16((int(vm.regs[d2]) * int(vm.regs[d3])) % 1000)
		vm.ip++
	case 8:
		// MOV register1, memory[register2]
		addr = vm.regs[d3]
		if addr > 999 {
			vm.Panic("Invalid memory reference")
		} else {
			vm.regs[d2] = vm.memory[addr]
		}
		vm.ip++
	case 9:
		// MOV memory[register2], register1
		addr = vm.regs[d3]
		if addr > 999 {
			vm.Panic("Invalid memory reference")
		} else {
			vm.memory[addr] = vm.regs[d2]
		}
		vm.ip++
	case 0:
		// CMP register2,0
		// JNZ register1

		if vm.regs[d3] == 0 {
			vm.ip++
		} else {
			addr = vm.regs[d2]
			vm.ip = addr
		}

	default:
		vm.Panic("Invalid instruction")
	}
}

func (vm *GenesysVM) Execute() (res error) {
	defer func() {
		if e := recover(); e != nil {
			res = e.(error)
		}
	}()

	for vm.running {
		vm.Step()
	}

	return nil
}
