package main

import (
	"image/color"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	windowWidth        = 900
	windowHeight       = 700
	numParticles       = 369
	minRadius          = 1
	maxRadius          = 5
	centerSize         = 2
	outlineWidth       = 0.5
	pheromoneSpread    = 2
	pheromoneIntensity = 1.0  // Initial intensity of the pheromone trail
	pheromoneDecay     = 0.05 // Rate at which pheromone intensity decreases over time
	pheromoneAlpha     = 200  // Alpha value for pheromone trail color
)

// Particle represents a single particle in the simulation
type Particle struct {
	pos    pixel.Vec // Position
	vel    pixel.Vec // Velocity
	radius float64   // Radius
	gene   color.RGBA
}

// Pheromone represents a pheromone trail left by particles
type Pheromone struct {
	pos       pixel.Vec // Position
	intensity float64   // Intensity
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

	// Create pheromone trail
	pheromoneTrail := make([]Pheromone, 0)

	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		win.Clear(color.Black)
		imd.Clear()

		// Update and draw particles
		for _, p := range particles {
			p.Update(dt)
			p.Draw(imd)
		}

		// Update and draw pheromone trail
		updatePheromoneTrail(&pheromoneTrail, dt)
		drawPheromoneTrail(imd, pheromoneTrail)

		imd.Draw(win)
		win.Update()

		// Run evolutionary algorithm
		EvolutionaryAlgorithm(particles)
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

// Update the pheromone trail, reducing intensity and removing trails with very low intensity
func updatePheromoneTrail(trail *[]Pheromone, dt float64) {
	// Update intensity of existing pheromone trails
	for i := range *trail {
		(*trail)[i].intensity -= pheromoneDecay * dt
	}
	// Remove trails with very low intensity
	j := 0
	for _, p := range *trail {
		if p.intensity > 0 {
			(*trail)[j] = p
			j++
		}
	}
	*trail = (*trail)[:j]
}

// Draw the pheromone trail with a fading effect
func drawPheromoneTrail(imd *imdraw.IMDraw, trail []Pheromone) {
	for _, p := range trail {
		// Calculate color based on intensity
		alpha := uint8(p.intensity * pheromoneAlpha)
		c := color.RGBA{R: 255, G: 255, B: 255, A: alpha}
		imd.Color = c

		// Draw pheromone trail as circles with decreasing intensity
		imd.Push(p.pos)
		imd.Circle(pheromoneSpread, 0)
	}
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

// EvolutionaryAlgorithm simulates trait selection over time
func EvolutionaryAlgorithm(particles []*Particle) {
	selected := selection(particles)
	crossover(selected)
}

// selection selects particles based on some criteria
func selection(particles []*Particle) []*Particle {
	// Add your selection logic here
	return particles
}

// crossover performs crossover operation on selected particles
func crossover(selected []*Particle) {
	// Add your crossover logic here
}

func main() {
	rand.Seed(time.Now().UnixNano())
	pixelgl.Run(run)
}
