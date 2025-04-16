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
		loop  Loop
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
	if !reflect.DeepEqual([]int{ 1, 2, 3 }, operations) {
		t.Error("Bad order:", operations)
	}
}

func TestMakeWhiteFillOp(t *testing.T) {
	var (
		loop  Loop
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
		loop  Loop
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
	if g != 0xffff || a != 0xffff  {
		t.Fatal("Invalid color, should be green", g, a)
	}
}

type testReceiver struct {
	lastTexture screen.Texture
}

func (tr *testReceiver) Update(t screen.Texture) {
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
	Colors []color.Color
}

func (m *mockTexture) Release() {}

func (m *mockTexture) Size() image.Point { return size }

func (m *mockTexture) Bounds() image.Rectangle {
	return image.Rectangle{Max: m.Size()}
}

func (m *mockTexture) Upload(dp image.Point, src screen.Buffer, sr image.Rectangle) {}
func (m *mockTexture) Fill(dr image.Rectangle, src color.Color, op draw.Op) {
	m.Colors = append(m.Colors, src)
}
