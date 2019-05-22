package channel

import "fmt"

// Pipline 的组合用法
func piplineExample() {
	multiply := func(values []int, multiplier int) []int {
		multipiedValues := make([]int, len(values))
		for i, v := range values {
			multipiedValues[i] = v * multiplier
		}
		return multipiedValues
	}

	add := func(values []int, additive int) []int {
		addedValues := make([]int, len(values))
		for i, v := range values {
			addedValues[i] = v + additive
		}
		return addedValues
	}

	ints := []int{1, 2, 3, 4}
	for _, v := range multiply(add(multiply(ints, 2), 1), 2) {
		fmt.Println(v)
	}

}
