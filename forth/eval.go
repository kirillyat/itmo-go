//go:build !solution

package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Evaluator struct {
	stack          []int
	userCmd        map[string][]string
	userCmdBuilder *CmdBuilder
	builtInCmds    []string
}

type CmdBuilder struct {
	name string
	cmds []string
}

func NewEvaluator() *Evaluator {
	builtInCmds := []string{"+", "-", "*", "/", "dup", "drop", "swap", "over", ":", ";"}
	return &Evaluator{
		userCmd:     make(map[string][]string),
		builtInCmds: builtInCmds,
	}
}

func newCmdBuilder() *CmdBuilder {
	return &CmdBuilder{}
}

func (e *Evaluator) Process(input string) ([]int, error) {
	tokens := strings.Fields(input)

	for _, token := range tokens {
		if err := e.execute(token); err != nil {
			return nil, err
		}
	}
	return e.stack, nil
}

func (e *Evaluator) execute(token string) error {
	token = strings.ToLower(token)
	if e.userCmdBuilder != nil && token != ";" {
		if e.userCmdBuilder.name == "" {
			e.userCmdBuilder.name = token
		} else {
			if definition, found := e.userCmd[token]; found {
				e.userCmdBuilder.cmds = append(e.userCmdBuilder.cmds, definition...)
			} else {
				e.userCmdBuilder.cmds = append(e.userCmdBuilder.cmds, token)
			}
		}
		return nil
	}

	if definition, found := e.userCmd[token]; found {
		for _, tok := range definition {
			if err := e.executeBuildIn(tok); err != nil {
				return err
			}
		}
		return nil
	}
	return e.executeBuildIn(token)
}

func (e *Evaluator) executeBuildIn(token string) error {

	switch token {
	case "+", "-", "*", "/":
		if len(e.stack) < 2 {
			return fmt.Errorf("not enough arguments for '%s'", token)
		}
		a, b := e.stack[len(e.stack)-2], e.stack[len(e.stack)-1]
		e.stack = e.stack[:len(e.stack)-2]
		switch token {
		case "+":
			e.stack = append(e.stack, a+b)
		case "-":
			e.stack = append(e.stack, a-b)
		case "*":
			e.stack = append(e.stack, a*b)
		case "/":
			if b == 0 {
				return fmt.Errorf("division by zero")
			}
			e.stack = append(e.stack, a/b)
		}
	case "dup", "drop", "swap", "over":
		if err := e.performStandardOperation(token); err != nil {
			return err
		}
	case ":":
		if e.userCmdBuilder == nil {
			e.userCmdBuilder = newCmdBuilder()
		} else {
			return fmt.Errorf("command definition already started")
		}
	case ";":
		if e.userCmdBuilder != nil && e.userCmdBuilder.name != "" {
			// Prevent overwriting numbers
			if _, err := strconv.Atoi(e.userCmdBuilder.name); err == nil {
				return fmt.Errorf("cannot redefine number")
			}
			e.userCmd[e.userCmdBuilder.name] = e.userCmdBuilder.cmds
			e.userCmdBuilder = nil
		} else {
			return fmt.Errorf("invalid end of command definition")
		}

	default:
		value, err := strconv.Atoi(token)
		if err != nil {
			return fmt.Errorf("unknown word: %s", token)
		}
		e.stack = append(e.stack, value)
	}
	return nil
}

func (e *Evaluator) performStandardOperation(operation string) error {
	switch operation {
	case "dup":
		if len(e.stack) < 1 {
			return fmt.Errorf("not enough arguments for 'dup'")
		}
		last := e.stack[len(e.stack)-1]
		e.stack = append(e.stack, last)
	case "drop":
		if len(e.stack) < 1 {
			return fmt.Errorf("not enough arguments for 'drop'")
		}
		e.stack = e.stack[:len(e.stack)-1]
	case "swap":
		if len(e.stack) < 2 {
			return fmt.Errorf("not enough arguments for 'swap'")
		}
		e.stack[len(e.stack)-1], e.stack[len(e.stack)-2] = e.stack[len(e.stack)-2], e.stack[len(e.stack)-1]
	case "over":
		if len(e.stack) < 2 {
			return fmt.Errorf("not enough arguments for 'over'")
		}
		secondLast := e.stack[len(e.stack)-2]
		e.stack = append(e.stack, secondLast)
	default:
		return fmt.Errorf("unknown operation: %s", operation)
	}
	return nil
}
