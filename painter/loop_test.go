package painter

import (
	"image"
	"image/color"
	"image/draw"
	"reflect"
	"testing"

	"golang.org/x/exp/shiny/screen"
)

func TestOrder(t *testing.T) {
	var (
		loop     Loop
		receiver testReceiver
	)
	loop.Receiver = &receiver
	var operations []int
	loop.Start(mockScreen{})
	loop.Post(OperationFunc(func(screen.Texture) {
		operations = append(operations, 1)
		loop.Post(OperationFunc(func(screen.Texture) {
			operations = append(operations, 3)
		}))
	}))
	loop.Post(OperationFunc(func(screen.Texture) {
		operations = append(operations, 2)
	}))
	loop.StopAndWait()
	if !reflect.DeepEqual([]int{1, 2, 3}, operations) {
		t.Error("Bad order:", operations)
	}
}

func TestMakeWhiteFillOp(t *testing.T) {
	var (
		loop     Loop
		receiver testReceiver
	)
	loop.Receiver = &receiver
	loop.Start(mockScreen{})
	loop.Post(MakeWhiteFillOp(&State{}))
	loop.Post(UpdateOp)
	loop.StopAndWait()
	if receiver.lastTexture == nil {
		t.Fatal("Texture was not updated")
	}
	texture, ok := receiver.lastTexture.(*mockTexture)
	if !ok {
		t.Fatal("Unexpected texture", receiver.lastTexture)
	}
	if len(texture.Colors) != 1 {
		t.Fatal("Invalid length", len(texture.Colors))
	}
	if texture.Colors[0] != color.White {
		t.Fatal("Invalid color, should be white")
	}
}

func TestMakeGreenFillOp(t *testing.T) {
	var (
		loop     Loop
		receiver testReceiver
	)
	loop.Receiver = &receiver
	loop.Start(mockScreen{})
	loop.Post(MakeGreenFillOp(&State{}))
	loop.Post(UpdateOp)
	loop.StopAndWait()
	if receiver.lastTexture == nil {
		t.Fatal("Texture was not updated")
	}
	texture, ok := receiver.lastTexture.(*mockTexture)
	if !ok {
		t.Fatal("Unexpected texture", receiver.lastTexture)
	}
	if len(texture.Colors) != 1 {
		t.Fatal("Invalid length", len(texture.Colors))
	}
	texturesColor := texture.Colors[0]
	_, g, _, a := texturesColor.RGBA()
	if g != 0xffff || a != 0xffff {
		t.Fatal("Invalid color, should be green", g, a)
	}
}

func TestMakeBgRectOp(t *testing.T) {
	var (
		loop     Loop
		receiver testReceiver
	)
	loop.Receiver = &receiver
	loop.Start(mockScreen{})
	loop.Post(MakeGreenFillOp(&State{}))
	loop.Post(MakeBgRectOp(&State{}, 0.1, 0.1, 0.3, 0.3))
	loop.Post(UpdateOp)
	loop.StopAndWait()
	if receiver.lastTexture == nil {
		t.Fatal("Texture was not updated")
	}
	texture, ok := receiver.lastTexture.(*mockTexture)
	if !ok {
		t.Fatal("Unexpected texture", receiver.lastTexture)
	}
	if len(texture.Colors) != 3 {
		t.Fatal("Invalid length", len(texture.Colors))
	}
	if !reflect.DeepEqual(texture.Colors[1:], []color.Color{nil, color.Black}) {
		t.Fatal("Invalid colors")
	}
	if texture.Rectangle.Eq(image.Rect(80, 80, 80, 80)) {
		t.Fatal("Invalid texture bounds")
	}
}

func TestMakeFigureOp(t *testing.T) {
	var (
		loop     Loop
		receiver testReceiver
	)
	loop.Receiver = &receiver
	loop.Start(mockScreen{})
	loop.Post(MakeFigureOp(&State{}, 0.1, 0.2))
	loop.Post(UpdateOp)
	loop.StopAndWait()
	if receiver.lastTexture == nil {
		t.Fatal("Texture was not updated")
	}
	texture, ok := receiver.lastTexture.(*mockTexture)
	if !ok {
		t.Fatal("Unexpected texture", receiver.lastTexture)
	}
	if len(texture.Colors) != 3 {
		t.Fatal("Invalid length", len(texture.Colors))
	}
	shapeColor := color.RGBA{R: 255, G: 230, B: 69, A: 255}
	if !reflect.DeepEqual(texture.Colors, []color.Color{nil, shapeColor, shapeColor}) {
		t.Fatal("Invalid colors", texture.Colors)
	}
}

func TestMakeMoveOp(t *testing.T) {
	var (
		loop     Loop
		receiver testReceiver
		state    State
	)
	loop.Receiver = &receiver
	loop.Start(mockScreen{})
	loop.Post(MakeFigureOp(&state, 0.1, 0.2))
	loop.Post(MakeBgRectOp(&state, 0.1, 0.2, 0.3, 0.3))
	loop.Post(MakeMoveOp(&state, 0.1, 0.2))
	loop.Post(UpdateOp)
	loop.StopAndWait()
	if receiver.lastTexture == nil {
		t.Fatal("Texture was not updated")
	}
	texture, ok := receiver.lastTexture.(*mockTexture)
	if !ok {
		t.Fatal("Unexpected texture", receiver.lastTexture)
	}
	if len(texture.Colors) != 11 {
		t.Fatal("Invalid length", len(texture.Colors))
	}
}

func TestMakeResetOp(t *testing.T) {
	var (
		loop     Loop
		receiver testReceiver
		state    State
	)
	loop.Receiver = &receiver
	loop.Start(mockScreen{})
	loop.Post(MakeFigureOp(&state, 0.1, 0.2))
	loop.Post(MakeBgRectOp(&state, 0.1, 0.2, 0.3, 0.3))
	loop.Post(MakeMoveOp(&state, 0.1, 0.2))
	loop.Post(MakeResetOp(&state))
	loop.Post(UpdateOp)
	loop.StopAndWait()
	if receiver.lastTexture == nil {
		t.Fatal("Texture was not updated")
	}
	texture, ok := receiver.lastTexture.(*mockTexture)
	if !ok {
		t.Fatal("Unexpected texture", receiver.lastTexture)
	}
	if len(texture.Colors) != 12 {
		t.Fatal("Invalid length", len(texture.Colors))
	}
	if texture.Colors[11] != color.Black {
		t.Fatal("Invalid last color, should be black")
	}
	if state.BgRect != nil || state.Figures != nil {
		t.Fatal("Sate is not nulled")
	}
}

type testReceiver struct {
	lastTexture screen.Texture
}

func (tr *testReceiver) Update(t screen.Texture) {
	tr.lastTexture = t
	tr.lastTexture = t
}

type mockScreen struct{}

func (m mockScreen) NewBuffer(size image.Point) (screen.Buffer, error) {
	panic("implement me")
}

func (m mockScreen) NewTexture(size image.Point) (screen.Texture, error) {
	return new(mockTexture), nil
}

func (m mockScreen) NewWindow(opts *screen.NewWindowOptions) (screen.Window, error) {
	panic("implement me")
}

type mockTexture struct {
	Colors    []color.Color
	Rectangle image.Rectangle
}

func (m *mockTexture) Release() {}

func (m *mockTexture) Size() image.Point { return size }

func (m *mockTexture) Bounds() image.Rectangle {
	return image.Rectangle{Max: m.Size()}
}

func (m *mockTexture) Upload(dp image.Point, src screen.Buffer, sr image.Rectangle) {
	panic("implement me")
}

func (m *mockTexture) Fill(dr image.Rectangle, src color.Color, op draw.Op) {
	m.Colors = append(m.Colors, src)
	m.Rectangle = dr
}
