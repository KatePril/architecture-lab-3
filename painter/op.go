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

func repaint(state *State, t screen.Texture) {
	t.Fill(t.Bounds(), state.BackgroundColor, screen.Src)
	if state.BgRect != nil {
		t.Fill(image.Rect(state.BgRect.X0, state.BgRect.Y0, state.BgRect.X1, state.BgRect.Y1), color.Black, screen.Src)
	}

	for i := range state.Figures {
		ui.DrawShape(t.Fill, state.Figures[i].X, state.Figures[i].Y, 0.5)
	}
}

// MakeWhiteFillOp зафарбовує текстуру у білий колір.
func MakeWhiteFillOp(state *State) Operation {
	return OperationFunc(func(t screen.Texture) {
		state.BackgroundColor = color.White
		repaint(state, t)
	})
}

// MakeGreenFillOp зафарбовує текстуру у зелений колір.
func MakeGreenFillOp(state *State) Operation {
	return OperationFunc(func(t screen.Texture) {
		state.BackgroundColor = color.RGBA{G: 0xff, A: 0xff}
		repaint(state, t)
	})
}

// MakeBgRectOp малює на фоні прямокутник чорного кольору у вказаних координатах.
func MakeBgRectOp(state *State, x0, y0, x1, y1 float32) Operation {
	return OperationFunc(func(t screen.Texture) {

		state.BgRect = &FigureRect{
			X0: int(float32(t.Size().X) * x0),
			Y0: int(float32(t.Size().Y) * y0),
			X1: int(float32(t.Size().X) * x1),
			Y1: int(float32(t.Size().Y) * y1)}
		repaint(state, t)
	})
}

// MakeFigureOp малює нову фігуру Т з центром у вказаних координатах поверх сформованого фону
func MakeFigureOp(state *State, x, y float32) Operation {
	return OperationFunc(func(t screen.Texture) {
		state.Figures = append(state.Figures, FigureT{
			X: int(float32(t.Size().X) * x),
			Y: int(float32(t.Size().Y) * y)})
		repaint(state, t)
	})
}

// MakeMoveOp переміщує усі фігури попередньо намальовані за допомогою команди figure у вказані координати
func MakeMoveOp(state *State, x, y float32) Operation {
	return OperationFunc(func(t screen.Texture) {
		xCoord := int(float32(t.Size().X) * x)
		yCoord := int(float32(t.Size().Y) * y)

		newFigures := make([]FigureT, len(state.Figures))
		for i := range state.Figures {
			newFigures[i] = FigureT{xCoord, yCoord}
		}
		state.Figures = newFigures

		repaint(state, t)
	})
}

// MakeResetOp очищує весь поточний стан текстури
func MakeResetOp(state *State) Operation {
	return OperationFunc(func(t screen.Texture) {
		state.BackgroundColor = color.Black
		state.BgRect = nil
		state.Figures = nil
		t.Fill(t.Bounds(), state.BackgroundColor, screen.Src)
	})
}
