package main

import (
	"image/color"
	"net/http"

	"github.com/KatePril/architecture-lab-3/painter"
	"github.com/KatePril/architecture-lab-3/painter/lang"
	"github.com/KatePril/architecture-lab-3/ui"
)

func main() {
	var (
		pv ui.Visualizer // Візуалізатор створює вікно та малює у ньому.

		// Потрібні для частини 2.
		opLoop painter.Loop // Цикл обробки команд.
		parser lang.Parser  // Парсер команд.
	)

	pv.Debug = true
	pv.Title = "Simple painter"

	pv.OnScreenReady = opLoop.Start
	opLoop.Receiver = &pv
	state := &painter.State{
		BackgroundColor: color.White,
		BgRect:          &painter.FigureRect{X0: 0, Y0: 0, X1: 100, Y1: 100},
		Figures:         []painter.FigureT{},
	}
	parser = lang.CreateOperationParser(state)

	go func() {
		http.Handle("/", lang.HttpHandler(&opLoop, parser))
		_ = http.ListenAndServe("localhost:17000", nil)
	}()

	pv.Main()
	opLoop.StopAndWait()
}
