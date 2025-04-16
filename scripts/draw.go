//go:build ignoretest
// +build ignoretest

package main

import (
	"bytes"
	"fmt"
	"math"
	"net/http"
	"time"
)

func sendPostRequest(data string) error {
	url := "http://localhost:17000/"
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func main() {
	// Initialize canvas
	commands := []string{
		"reset",
		"green",
		"bgrect 100 100 300 300",
	}

	for _, cmd := range commands {
		if err := sendPostRequest(cmd); err != nil {
			fmt.Printf("Error sending command '%s': %v\n", cmd, err)
			return
		}
	}

	// Create figure at center
	if err := sendPostRequest("figure 200 200"); err != nil {
		fmt.Printf("Error sending command 'figure 200 200': %v\n", err)
		return
	}

	if err := sendPostRequest("bgrect 0.25 0.25 0.75 0.75"); err != nil {
		fmt.Printf("Error sending command 'figure 200 200': %v\n", err)
		return
	}

	if err := sendPostRequest("update"); err != nil {
		fmt.Printf("Error sending command 'update': %v\n", err)
		return
	}

	// Animation parameters
	count := 0
	radius := 100.0 // Smaller radius to stay within visible area
	centerX, centerY := 200.0, 200.0
	timeStep := 50 * time.Millisecond // 20 frames per second
	angularSpeed := 2.0               // degrees per frame
	angle := 0.0

	for count < 360 { // Complete a full circle (360 degrees)
		count++
		angle += angularSpeed
		if angle >= 360 {
			angle -= 360
		}

		// Convert angle to radians
		radians := angle * math.Pi / 180

		// Calculate new position
		x := centerX + radius*math.Cos(radians)
		y := centerY + radius*math.Sin(radians)

		moveCommand := fmt.Sprintf("move %.0f %.0f", x, y)
		if err := sendPostRequest(moveCommand); err != nil {
			fmt.Printf("Error sending command '%s': %v\n", moveCommand, err)
			return
		}

		if err := sendPostRequest("update"); err != nil {
			fmt.Printf("Error sending command 'update': %v\n", err)
			return
		}

		time.Sleep(timeStep)
	}
}
