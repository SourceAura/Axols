package main

import (
	"math/rand"
	"time"

	"axols/sim"

	"github.com/faiface/pixel/pixelgl"
)

// startDebugDisplay prints simulation debug info periodically in a separate goroutine.
// This is a minimal implementation that prints a message every few seconds.
func startDebugDisplay() {
	go func() {
		for {
			time.Sleep(5 * time.Second)
			println("[DEBUG] Simulation running... (implement stats display as needed)")
		}
	}()
}

func main() {
	rand.Seed(time.Now().UnixNano())
	startDebugDisplay()
	pixelgl.Run(sim.Run)
}
