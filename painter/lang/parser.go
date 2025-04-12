package lang

import (
	"bufio"
	"io"

	"github.com/KatePril/architecture-lab-3/painter"
)

type Parser interface {
	Parse(in io.Reader) ([]painter.Operation, error)
}

type OperationParser struct {
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
	return nil, nil
}
