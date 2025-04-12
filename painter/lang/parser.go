package lang

import (
	"bufio"
	"fmt"
	"io"
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
		return nil, fmt.Errorf(commandName + " is not a valid command")
	}

	operation, err := getOperationFunc(args)
	if err != nil {
		return nil, err
	}

	return operation, nil
}
