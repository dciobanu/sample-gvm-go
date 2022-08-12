package main

import (
	"reflect"
	"testing"
)

func Test_Errors(t *testing.T) {

	vm := NewGenesysVM([]uint16{}, 2000)
	if err := vm.Execute(); err == nil {
		t.Errorf("Expecting failure after 1000 instructions")
	}

	if vm.ic != 1000 || vm.ip != 1000 {
		t.Errorf("Unexpected IC/IP")
	}

	vm = NewGenesysVM([]uint16{123}, 2000)
	if err := vm.Execute(); err == nil {
		t.Errorf("HALT with parameters must return an error")
	}

	vm = NewGenesysVM([]uint16{9999}, 2000)
	if err := vm.Execute(); err == nil {
		t.Errorf("Instructions over 1,000 are invalid and must return an error")
	}

}

func Test_Halting(t *testing.T) {
	vm := NewGenesysVM([]uint16{100}, 1000)
	vm.Execute()

	if vm.ic != 1 || vm.ip != 0 {
		t.Errorf("Halting failed: %v", vm)
	}
}

func Test_SetRegister(t *testing.T) {
	code := make([]uint16, 0, 1000)

	for i := 0; i < 10; i++ {
		code = append(code, uint16(200+10*i+i))
	}

	code = append(code, 100)

	vm := NewGenesysVM(code, 1000)

	var reg2 [10]uint16
	for i := 0; i < 10; i++ {
		vm.Step()
		reg2[i] = uint16(i)
		if !reflect.DeepEqual(reg2, vm.regs) {
			t.Fatalf("Invalid register value at step %d", i)
		}

		if vm.ic != uint64(i+1) || vm.ip != uint16(i+1) {
			t.Fatalf("Invalid IP %d or IC %d at step %d", vm.ip, vm.ic, i)
		}
	}

	vm.Step()
	if vm.running {
		t.Fatalf("VM not halted")
	}
}

func Test_AddNtoR(t *testing.T) {
	vm := NewGenesysVM([]uint16{356, 343, 359, 329, 100}, 1000)
	var xregs [10]uint16

	checkState := func() {
		if !reflect.DeepEqual(xregs, vm.regs) {
			t.Fatalf("Unexpected registers state: %v, expected %v", vm.regs, xregs)
		}
	}

	vm.Step()
	xregs[5] = 6
	checkState()

	vm.Step()
	xregs[4] = 3
	checkState()

	vm.Step()
	xregs[5] = 15
	checkState()

	// Test value overflow
	vm.regs[2] = 997
	vm.Step()
	xregs[2] = 6
	checkState()

	vm.Step()
	if vm.running {
		t.Fatalf("VM not halted")
	}
}

func Test_MulRbyN(t *testing.T) {
	vm := NewGenesysVM([]uint16{402, 417, 429, 439, 445, 451, 460, 0, 0, 499, 100}, 1000)

	var xregs [10]uint16 = [10]uint16{0, 1, 100, 999, 200, 13, 14, 0, 0, 998}
	copy(vm.regs[:], xregs[:])

	expected := []uint16{0, 7, 900, 991, 0, 13, 0, 0, 0, 982}

	for i := range expected {
		xregs[i] = expected[i]
		vm.Step()

		if !reflect.DeepEqual(xregs, vm.regs) {
			t.Fatalf("Unexpected registers state: %v, expected %v at step %d", vm.regs, xregs, i)
		}
	}

	vm.Step()
	if vm.running {
		t.Fatalf("VM not halted")
	}
}

func Test_MovRtoR(t *testing.T) {
	vm := NewGenesysVM([]uint16{519, 510, 594, 100}, 1000)

	var xregs [10]uint16 = [10]uint16{0, 100, 200, 300, 400, 500, 600, 700, 800, 900}
	copy(vm.regs[:], xregs[:])

	checkState := func() {
		if !reflect.DeepEqual(xregs, vm.regs) {
			t.Fatalf("Unexpected registers state: %v, expected %v", vm.regs, xregs)
		}
	}

	vm.Step()
	xregs[1] = 900
	checkState()

	vm.Step()
	xregs[1] = 0
	checkState()

	vm.Step()
	xregs[9] = 400
	checkState()

	vm.Step()
	if vm.running {
		t.Fatalf("VM not halted")
	}
}

func Test_AddRtoR(t *testing.T) {
	vm := NewGenesysVM([]uint16{601, 646, 699, 100}, 1000)

	var xregs [10]uint16 = [10]uint16{0, 100, 200, 300, 400, 500, 600, 700, 800, 900}
	copy(vm.regs[:], xregs[:])

	checkState := func() {
		if !reflect.DeepEqual(xregs, vm.regs) {
			t.Fatalf("Unexpected registers state: %v, expected %v", vm.regs, xregs)
		}
	}

	vm.Step()
	xregs[0] = 100
	checkState()

	vm.Step()
	xregs[4] = 0
	checkState()

	vm.Step()
	xregs[9] = 800
	checkState()

	vm.Step()
	if vm.running {
		t.Fatalf("VM not halted")
	}
}

func Test_MulRtoR(t *testing.T) {
	vm := NewGenesysVM([]uint16{701, 732, 734, 777, 798, 100}, 1000)

	var xregs [10]uint16 = [10]uint16{3, 100, 2, 300, 0, 500, 666, 777, 999, 999}
	copy(vm.regs[:], xregs[:])

	checkState := func() {
		if !reflect.DeepEqual(xregs, vm.regs) {
			t.Fatalf("Unexpected registers state: %v, expected %v (IP: %d)", vm.regs, xregs, vm.ip)
		}
	}

	vm.Step()
	xregs[0] = 300
	checkState()

	vm.Step()
	xregs[3] = 600
	checkState()

	vm.Step()
	xregs[3] = 0
	checkState()

	vm.Step()
	xregs[7] = 729
	checkState()

	vm.Step()
	xregs[9] = 1
	checkState()

	vm.Step()
	if vm.running {
		t.Fatalf("VM not halted")
	}
}

func Test_MoveMemToR(t *testing.T) {
	vm := NewGenesysVM([]uint16{850, 861, 872, 883, 893, 100, 666, 777, 888, 999}, 1000)

	var xregs [10]uint16 = [10]uint16{6, 7, 8, 9, 0, 0, 0, 0, 0, 0}
	copy(vm.regs[:], xregs[:])

	checkState := func() {
		if !reflect.DeepEqual(xregs, vm.regs) {
			t.Fatalf("Unexpected registers state: %v, expected %v (IP: %d)", vm.regs, xregs, vm.ip)
		}
	}

	vm.Step()
	xregs[5] = 666
	checkState()

	vm.Step()
	xregs[6] = 777
	checkState()

	vm.Step()
	xregs[7] = 888
	checkState()

	vm.Step()
	xregs[8] = 999
	checkState()

	vm.Step()
	xregs[9] = 999
	checkState()

	vm.Step()
	if vm.running {
		t.Fatalf("VM not halted")
	}
}

func Test_MoveRToMem(t *testing.T) {
	vm := NewGenesysVM([]uint16{901, 912, 923, 934, 945, 956, 967, 978, 989, 999, 100}, 1000)

	var xregs [10]uint16 = [10]uint16{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}
	copy(vm.regs[:], xregs[:])

	var xmem [1000]uint16
	copy(xmem[:], vm.memory[:])

	checkState := func() {
		if !reflect.DeepEqual(xregs, vm.regs) {
			t.Fatalf("Unexpected registers state: %v, expected %v (IP: %d)", vm.regs, xregs, vm.ip)
		}

		if !reflect.DeepEqual(xmem, vm.memory) {
			t.Fatalf("Unexpected memory state: %v, expected %v (IP: %d)", vm.memory, xmem, vm.ip)
		}
	}

	for i := 1; i <= 9; i++ {
		vm.Step()
		xmem[10*(i+1)] = uint16(i * 10)
		checkState()
	}

	vm.Step()
	xmem[100] = 100
	checkState()

	vm.Step()
	if vm.running {
		t.Fatalf("VM not halted")
	}
}

func Test_Jmp(t *testing.T) {
	vm := NewGenesysVM([]uint16{11, 0, 9, 100}, 1000)

	var xregs [10]uint16 = [10]uint16{0, 2, 1, 1, 1, 1, 1, 1, 1, 0}
	copy(vm.regs[:], xregs[:])

	expect := func(ip uint16, ic uint64) {
		if vm.ip != ip || vm.ic != ic {
			t.Fatalf("Unpexpected IP/IC state %d/%d expecting %d/%d", vm.ip, vm.ic, ip, ic)
		}
	}

	expect(0, 0)

	vm.Step()
	expect(2, 1)

	vm.Step()
	expect(3, 2)

	vm.Step()
	if vm.running {
		t.Fatalf("VM not halted")
	}

	expect(3, 3)
}

func Test_Sample1(t *testing.T) {
	var sample_code1 []uint16 = []uint16{299, 492, 495, 399, 492, 495, 399, 283, 279, 689, 78, 100, 0, 0, 0}

	vm := NewGenesysVM(sample_code1, 1000)
	vm.Execute()
	count, running := vm.GetStats()

	if running {
		t.Errorf("VM still running, expected to be halted")
	}

	if count != 16 {
		t.Errorf("Instructions count for sample1 expected to be 16, actual %d", count)
	}
}
