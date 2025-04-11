package painter

import (
	"github.com/KatePril/architecture-lab-3/ui"
	"golang.org/x/exp/shiny/screen"
	"image"
	"image/color"
)

// Operation змінює вхідну текстуру.
type Operation interface {
	// Do виконує зміну операції, повертаючи true, якщо текстура вважається готовою для відображення.
	Do(t screen.Texture) (ready bool)
}

type FigureT struct {
	X int
	Y int
}

type FigureRect struct {
	X0 int
	Y0 int
	X1 int
	Y1 int
}

type State struct {
	BackgroundColor color.Color
	BgRect          *FigureRect
	Figures         []FigureT
}

// OperationList групує список операції в одну.
type OperationList []Operation

func (ol OperationList) Do(t screen.Texture) (ready bool) {
	for _, o := range ol {
		ready = o.Do(t) || ready
	}
	return
}

// UpdateOp операція, яка не змінює текстуру, але сигналізує, що текстуру потрібно розглядати як готову.
var UpdateOp = updateOp{}

type updateOp struct{}

func (op updateOp) Do(t screen.Texture) bool { return true }

// OperationFunc використовується для перетворення функції оновлення текстури в Operation.
type OperationFunc func(t screen.Texture)

func (f OperationFunc) Do(t screen.Texture) bool {
	f(t)
	return false
}

// TODO: GreenFill and WhiteFill are unnecessary
// WhiteFill зафарбовує тестуру у білий колір. Може бути викоистана як Operation через OperationFunc(WhiteFill).
func WhiteFill(t screen.Texture) {
	t.Fill(t.Bounds(), color.White, screen.Src)
}

// GreenFill зафарбовує тестуру у зелений колір. Може бути викоистана як Operation через OperationFunc(GreenFill).
func GreenFill(t screen.Texture) {
	t.Fill(t.Bounds(), color.RGBA{G: 0xff, A: 0xff}, screen.Src)
}

func repaint(state *State, t screen.Texture) {
	t.Fill(t.Bounds(), state.BackgroundColor, screen.Src)
	if state.BgRect != nil {
		t.Fill(image.Rect(state.BgRect.X0, state.BgRect.Y0, state.BgRect.X1, state.BgRect.Y1), color.Black, screen.Src)
	}

	for i := range state.Figures {
		ui.DrawShape(t.Fill, state.Figures[i].X, state.Figures[i].Y, 0.5)
	}
}

func MakeWhiteFillOp(state *State) Operation {
	return OperationFunc(func(t screen.Texture) {
		state.BackgroundColor = color.White
		repaint(state, t)
	})
}

func MakeGreenFillOp(state *State) Operation {
	return OperationFunc(func(t screen.Texture) {
		state.BackgroundColor = color.RGBA{G: 0xff, A: 0xff}
		repaint(state, t)
	})
}

func MakeBgRectOp(state *State, x0, y0, x1, y1 int) Operation {
	return OperationFunc(func(t screen.Texture) {
		state.BgRect = &FigureRect{x0, y0, x1, y1}
		repaint(state, t)
	})
}

func MakeFigureOp(state *State, x, y int) Operation {
	return OperationFunc(func(t screen.Texture) {
		state.Figures = append(state.Figures, FigureT{x, y})
		repaint(state, t)
	})
}

func MakeMoveOp(state *State, x, y int) Operation {
	return OperationFunc(func(t screen.Texture) {
		newFigures := make([]FigureT, len(state.Figures))
		for i := range state.Figures {
			newFigures[i] = FigureT{x, y}
		}
		state.Figures = newFigures

		sizeX := (state.BgRect.X1 - state.BgRect.X0) / 2
		sizeY := (state.BgRect.Y1 - state.BgRect.Y0) / 2

		state.BgRect = &FigureRect{
			X0: x - sizeX,
			Y0: y - sizeY,
			X1: x + sizeX,
			Y1: y + sizeY,
		}

		repaint(state, t)
	})
}

func MakeResetOp(state *State) Operation {
	return OperationFunc(func(t screen.Texture) {
		state.BackgroundColor = color.Black
		state.BgRect = nil
		state.Figures = nil
		t.Fill(t.Bounds(), state.BackgroundColor, screen.Src)
	})
}
