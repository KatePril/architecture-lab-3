package lang

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/KatePril/architecture-lab-3/painter"
)

type Parser interface {
	Parse(in io.Reader) ([]painter.Operation, error)
}

type GetOperationFunc func(args []string) (painter.Operation, error)

type OperationParser struct {
	dictionary map[string]GetOperationFunc
}

func CreateOperationParser(state *painter.State) *OperationParser {
	dictionary := setOperationDictionary(state)
	return &OperationParser{dictionary: dictionary}
}

func (p *OperationParser) Parse(in io.Reader) ([]painter.Operation, error) {
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)

	var res []painter.Operation

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		commandLine := scanner.Text()
		operation, err := p.parseLine(commandLine)
		if err != nil {
			return nil, err
		}
		if operation == nil {
			continue
		}

		res = append(res, operation)
	}

	return res, nil
}

func (p *OperationParser) parseLine(commandLine string) (painter.Operation, error) {
	commandLineSlice := strings.Fields(commandLine)
	length := len(commandLineSlice)
	if length == 0 {
		return nil, nil
	}

	commandName := strings.ToLower(commandLineSlice[0])
	var args []string
	if length > 1 {
		args = commandLineSlice[1:]
	} else {
		args = []string{}
	}

	getOperationFunc, ok := p.dictionary[commandName]
	if !ok {
		return nil, fmt.Errorf("%s", commandName+" is not a valid command")
	}

	operation, err := getOperationFunc(args)
	if err != nil {
		return nil, err
	}

	return operation, nil
}

func setOperationDictionary(state *painter.State) map[string]GetOperationFunc {
	dictionary := make(map[string]GetOperationFunc)

	dictionary["update"] = func(args []string) (painter.Operation, error) {
		return painter.UpdateOp, nil
	}

	dictionary["white"] = func(args []string) (painter.Operation, error) {
		return painter.MakeWhiteFillOp(state), nil
	}

	dictionary["green"] = func(args []string) (painter.Operation, error) {
		return painter.MakeGreenFillOp(state), nil
	}

	dictionary["bgrect"] = func(args []string) (painter.Operation, error) {
		const rightNumOfArgs int = 4
		if len(args) != rightNumOfArgs {
			return nil, fmt.Errorf("invalid arguments")
		}

		var parsedArgs [rightNumOfArgs]int

		for i := 0; i < rightNumOfArgs; i++ {
			parsed, err := strconv.Atoi(args[i])
			if err != nil {
				return nil, fmt.Errorf("invalid arguments")
			}

			parsedArgs[i] = parsed
		}

		return painter.MakeBgRectOp(state, parsedArgs[0], parsedArgs[1], parsedArgs[2], parsedArgs[3]), nil
	}

	dictionary["figure"] = func(args []string) (painter.Operation, error) {
		const rightNumOfArgs int = 2
		if len(args) != rightNumOfArgs {
			return nil, fmt.Errorf("invalid arguments")
		}

		var parsedArgs [rightNumOfArgs]int

		for i := 0; i < rightNumOfArgs; i++ {
			parsed, err := strconv.Atoi(args[i])
			if err != nil {
				return nil, fmt.Errorf("invalid arguments")
			}

			parsedArgs[i] = parsed
		}

		return painter.MakeFigureOp(state, parsedArgs[0], parsedArgs[1]), nil
	}

	dictionary["move"] = func(args []string) (painter.Operation, error) {
		const rightNumOfArgs int = 2
		if len(args) != rightNumOfArgs {
			return nil, fmt.Errorf("invalid arguments")
		}

		var parsedArgs [rightNumOfArgs]int

		for i := 0; i < rightNumOfArgs; i++ {
			parsed, err := strconv.Atoi(args[i])
			if err != nil {
				return nil, fmt.Errorf("invalid arguments")
			}

			parsedArgs[i] = parsed
		}

		return painter.MakeMoveOp(state, parsedArgs[0], parsedArgs[1]), nil
	}

	dictionary["reset"] = func(args []string) (painter.Operation, error) {
		return painter.MakeResetOp(state), nil
	}

	return dictionary
}
