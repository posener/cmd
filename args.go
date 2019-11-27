package cmd

import (
	"fmt"
	"strconv"
)

// ArgsStr are string positional arguments. If it is created with cap > 0, it will be used to define
// the number of required arguments.
//
// Usage
//
// To get a list of arbitrary number of arguments:
//
// 	root := cmd.Root()
//
// 	var cmd.ArgsStr args
// 	root.ArgsVar(&args, "[arg...]", "list of arguments")
//
// To get a list of specific number of arguments:
//
// 	root := cmd.Root()
//
// 	args := make(cmd.ArgsStr, 3)
// 	root.ArgsVar(&args, "[arg1] [arg2] [arg3]", "list of 3 arguments")
type ArgsStr []string

// Set implements the ArgsValue interface.
func (a *ArgsStr) Set(args []string) error {
	if cap(*a) > 0 && len(args) != cap(*a) {
		return fmt.Errorf("required %d positional args, got %v", cap(*a), args)
	}
	*a = args
	return nil
}

// ArgsInt are int positional arguments. If it is created with cap > 0, it will be used to define
// the number of required arguments.
//
// Usage
//
// To get a list of arbitrary number of integers:
//
// 	root := cmd.Root()
//
// 	var cmd.ArgsInt args
// 	root.ArgsVar(&args, "[int...]", "list of integer args")
//
// To get a list of specific number of integers:
//
// 	root := cmd.Root()
//
// 	args := make(cmd.ArgsInt, 3)
// 	root.ArgsVar(&args, "[int1] [int2] [int3]", "list of 3 integers")
type ArgsInt []int

// Set implements the ArgsValue interface.
func (a *ArgsInt) Set(args []string) error {
	if cap(*a) > 0 && len(args) != cap(*a) {
		return fmt.Errorf("required %d positional args, got %v", cap(*a), args)
	}
	*a = (*a)[:0] // Reset length to 0.
	for i, arg := range args {
		v, err := strconv.Atoi(arg)
		if err != nil {
			return fmt.Errorf("invalid int positional argument at position %d with value %v", i, arg)
		}
		*a = append(*a, v)
	}
	return nil
}
