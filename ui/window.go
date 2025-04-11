package ui

import (
	"image"
	"image/color"
	"log"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/imageutil"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/draw"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

type Visualizer struct {
	Title         string
	Debug         bool
	OnScreenReady func(s screen.Screen)

	w    screen.Window
	tx   chan screen.Texture
	done chan struct{}

	sz  size.Event
	pos image.Rectangle
}

func (pw *Visualizer) Main() {
	pw.tx = make(chan screen.Texture)
	pw.done = make(chan struct{})
	pw.pos.Max.X = 200
	pw.pos.Max.Y = 200
	driver.Main(pw.run)
}

func (pw *Visualizer) Update(t screen.Texture) {
	pw.tx <- t
}

func (pw *Visualizer) run(s screen.Screen) {
	w, err := s.NewWindow(&screen.NewWindowOptions{
		Title:  pw.Title,
		Width:  800,
		Height: 800,
	})
	if err != nil {
		log.Fatal("Failed to initialize the app window:", err)
	}
	defer func() {
		w.Release()
		close(pw.done)
	}()

	if pw.OnScreenReady != nil {
		pw.OnScreenReady(s)
	}

	pw.w = w

	events := make(chan any)
	go func() {
		for {
			e := w.NextEvent()
			if pw.Debug {
				log.Printf("new event: %v", e)
			}
			if detectTerminate(e) {
				close(events)
				break
			}
			events <- e
		}
	}()

	var t screen.Texture

	for {
		select {
		case e, ok := <-events:
			if !ok {
				return
			}
			pw.handleEvent(e, t)

		case t = <-pw.tx:
			w.Send(paint.Event{})
		}
	}
}

func detectTerminate(e any) bool {
	switch e := e.(type) {
	case lifecycle.Event:
		if e.To == lifecycle.StageDead {
			return true // Window destroy initiated.
		}
	case key.Event:
		if e.Code == key.CodeEscape {
			return true // Esc pressed.
		}
	}
	return false
}

func (pw *Visualizer) handleEvent(e any, t screen.Texture) {
	switch e := e.(type) {

	case size.Event: // Оновлення даних про розмір вікна.
		pw.sz = e

	case error:
		log.Printf("ERROR: %s", e)

	case mouse.Event:
		if t == nil {
			if e.Button == mouse.ButtonRight {
				log.Printf("Clicked: %v, %v", e.X, e.Y)
				x := int(e.X)
				y := int(e.Y)
				pw.drawDefaultUI(&x, &y)
			}
		}

	case paint.Event:
		// Малювання контенту вікна.
		if t == nil {
			pw.drawDefaultUI(nil, nil)
		} else {
			// Використання текстури отриманої через виклик Update.
			pw.w.Scale(pw.sz.Bounds(), t, t.Bounds(), draw.Src, nil)
		}
		pw.w.Publish()
	}
}

func (pw *Visualizer) drawDefaultUI(x, y *int) {
	pw.w.Fill(pw.sz.Bounds(), color.RGBA{R: 174, G: 255, B: 168}, draw.Src)

	if x == nil {
		defaultX := pw.sz.WidthPx / 2
		x = &defaultX
	}
	if y == nil {
		defaultY := pw.sz.HeightPx / 2
		y = &defaultY
	}

	DrawShape(pw.w.Fill, *x, *y, 1)

	// Малювання білої рамки.
	for _, br := range imageutil.Border(pw.sz.Bounds(), 10) {
		pw.w.Fill(br, color.White, draw.Src)
	}
}

func DrawShape(Fill func(dr image.Rectangle, src color.Color, op draw.Op), x, y int, scale float64) {
	s := 100
	shapeColor := color.RGBA{R: 255, G: 230, B: 69, A: 255}
	scaledS := int(float64(s) * scale)

	rotate := func(px, py int) (int, int) {
		px, py = px-x, py-y
		return x + py, y - px
	}

	x00, y00 := rotate(x-scaledS, y-scaledS)
	x10, y10 := rotate(x+scaledS, y)
	rect0 := image.Rect(x00, y00, x10, y10).Canon()

	x01, y01 := rotate(x-(scaledS/2), y)
	x11, y11 := rotate(x+(scaledS/2), y+scaledS)
	rect1 := image.Rect(x01, y01, x11, y11).Canon()

	Fill(rect0, shapeColor, draw.Src)
	Fill(rect1, shapeColor, draw.Src)
}
