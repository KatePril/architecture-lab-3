package lang

import (
	"io"

	"github.com/KatePril/architecture-lab-3/painter"
)

type Parser interface {
	Parse(in io.Reader) ([]painter.Operation, error)
}

type OperationParser struct {
}

func (p *OperationParser) Parse(in io.Reader) ([]painter.Operation, error) {
	var res []painter.Operation
	state := &painter.State{}

	res = append(res, painter.MakeWhiteFillOp(state))
	res = append(res, painter.MakeBgRectOp(state, 55, 55, 75, 75))
	res = append(res, painter.MakeFigureOp(state, 201, 165))
	res = append(res, painter.MakeGreenFillOp(state))
	res = append(res, painter.MakeFigureOp(state, 300, 300))
	res = append(res, painter.MakeBgRectOp(state, 60, 100, 300, 300))
	res = append(res, painter.MakeMoveOp(state, 100, 100))
	res = append(res, painter.MakeResetOp(state))
	res = append(res, painter.UpdateOp)

	return res, nil
}
