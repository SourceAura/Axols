package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	windowWidth  = 1100
	windowHeight = 800

	numParticles = 500
	minRadius    = 1
	maxRadius    = 5
	centerSize   = 2
	outlineWidth = 0.5
)

// Particle represents a single particle in the simulation
type Particle struct {
	pos    pixel.Vec // Position
	vel    pixel.Vec // Velocity
	radius float64   // Radius
	gene   color.RGBA
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Particle Life Simulation",
		Bounds: pixel.R(0, 0, windowWidth, windowHeight),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	imd := imdraw.New(nil)

	// Create particles
	particles := make([]*Particle, numParticles)
	for i := range particles {
		particles[i] = NewParticle()
	}

	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		win.Clear(color.Black)
		imd.Clear()

		for _, p := range particles {
			p.Update(dt)
			p.Draw(imd)
		}

		imd.Draw(win)
		win.Update()

		// Print debugging/logging information
		logInfo(particles)
	}
}

// NewParticle creates a new particle with random properties
func NewParticle() *Particle {
	return &Particle{
		pos:    pixel.V(rand.Float64()*windowWidth, rand.Float64()*windowHeight),
		vel:    pixel.V(rand.Float64()*100-50, rand.Float64()*100-50),
		radius: rand.Float64()*(maxRadius-minRadius) + minRadius,
		gene:   randomColor(),
	}
}

// Update updates the position of the particle
func (p *Particle) Update(dt float64) {
	p.pos = p.pos.Add(p.vel.Scaled(dt))
	if p.pos.X < 0 || p.pos.X > windowWidth {
		p.vel.X = -p.vel.X
	}
	if p.pos.Y < 0 || p.pos.Y > windowHeight {
		p.vel.Y = -p.vel.Y
	}
}

// Draw draws the particle
func (p *Particle) Draw(imd *imdraw.IMDraw) {
	imd.Color = p.gene

	// Draw the center
	imd.Push(p.pos)
	imd.Circle(centerSize, 0) // smaller center
	imd.Polygon(0)            // Start a new polygon

	// Draw the outline
	imd.Color = colornames.White
	imd.Push(p.pos)
	imd.Circle(p.radius, outlineWidth)
}

// randomColor generates a random color
func randomColor() color.RGBA {
	return color.RGBA{
		R: uint8(rand.Intn(256)),
		G: uint8(rand.Intn(256)),
		B: uint8(rand.Intn(256)),
		A: 150, // Less transparency
	}
}

// logInfo prints debugging/logging information
func logInfo(particles []*Particle) {
	// Clear the console
	print("\033[H\033[2J")

	// Print the number of particles
	fmt.Println("Number of particles:", len(particles))
	// Add more debugging/logging information as needed
}

func main() {
	rand.Seed(time.Now().UnixNano())
	pixelgl.Run(run)
}
