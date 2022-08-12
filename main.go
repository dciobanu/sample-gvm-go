package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func loadInputs(fname string) (res [][]uint16) {
	f, err := os.Open(fname)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	res = append(res, nil)

	for scanner.Scan() {
		if scanner.Text() == "" {
			res = append(res, nil)
		} else {
			if value, err := strconv.Atoi(scanner.Text()); err != nil {
				panic("Can't parse int")
			} else {
				res[len(res)-1] = append(res[len(res)-1], uint16(value))
			}
		}
	}

	return
}

func main() {
	input := loadInputs("doc/samples.txt")
	fmt.Println(input)

	for _, mem := range input {
		vm := NewGenesysVM(mem, 10000)

		err := vm.Execute()
		if err != nil {
			fmt.Println(err)
		} else {
			count, _ := vm.GetStats()
			fmt.Println(count)
		}
	}
}
