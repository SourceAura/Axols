package main

import (
	"fmt"
	"image/color"
	"math"
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
	outlineWidth       = 0.5
	pheromoneSpread    = 2
	pheromoneIntensity = 1.0  // Initial intensity of the pheromone trail
	pheromoneDecay     = 0.05 // Rate at which pheromone intensity decreases over time
	pheromoneAlpha     = 200  // Alpha value for pheromone trail color
	initialSpeedFactor = 0.5  // Factor to slow down initial movement speed
	nucleusRadius      = 0.3  // Radius of the nucleus
	mutationRate       = 0.1  // Rate of mutation
)

// Biome represents a separate environment within the simulation
type Biome int

const (
	Overworld Biome = iota
	BubbleBiome1
	BubbleBiome2
	BubbleBiome3
)

// Particle represents a single particle in the simulation
type Particle struct {
	pos        pixel.Vec // Position
	vel        pixel.Vec // Velocity
	radius     float64   // Radius
	gene       color.RGBA
	nucleusPos pixel.Vec // Position of the nucleus
	biome      Biome     // Biome the particle belongs to
}

// Pheromone represents a pheromone trail left by particles
type Pheromone struct {
	pos       pixel.Vec // Position
	intensity float64   // Intensity
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Evo-Sim",
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

	generation := 0
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

		// Call debugDisplay function with relevant data
		debugDisplay(particles, pheromoneTrail, generation)

		// Apply genetic algorithm
		particles = evolvePopulation(particles)

		generation++
	}
}

// NewParticle creates a new particle with random properties
func NewParticle() *Particle {
	p := &Particle{
		pos:        pixel.V(rand.Float64()*windowWidth, rand.Float64()*windowHeight),
		radius:     rand.Float64()*(maxRadius-minRadius) + minRadius,
		gene:       randomColor(),
		nucleusPos: pixel.V(0, 0),
	}

	// Set initial velocity (slowed by half)
	p.vel = pixel.V(rand.Float64()*100-50, rand.Float64()*100-50).Scaled(initialSpeedFactor)

	// Set nucleus position relative to the particle
	p.nucleusPos = p.pos.Add(pixel.V(-p.radius/2, p.radius/2))

	// Assign a random biome
	p.biome = Biome(rand.Intn(4))

	return p
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

	// Update nucleus position
	p.nucleusPos = p.pos.Add(pixel.V(-p.radius/2, p.radius/2))
}

// Draw draws the particle
func (p *Particle) Draw(imd *imdraw.IMDraw) {
	imd.Color = p.gene

	// Draw the body
	imd.Push(p.pos)
	imd.Circle(p.radius, outlineWidth)

	// Draw the tail
	tailLength := p.radius * 2
	tailAngle := p.vel.Angle() - math.Pi
	tailStart := pixel.V(
		p.pos.X-tailLength*math.Cos(tailAngle),
		p.pos.Y-tailLength*math.Sin(tailAngle),
	)
	imd.Push(tailStart)
	imd.Push(p.pos)
	imd.Line(outlineWidth)

	// Draw the head
	headSize := p.radius * 1.5
	headPos := pixel.V(
		p.pos.X+headSize*math.Cos(tailAngle),
		p.pos.Y+headSize*math.Sin(tailAngle),
	)
	imd.Push(headPos)
	imd.Circle(headSize, outlineWidth)

	// Draw the eyes
	eyeSize := p.radius * 0.2
	leftEyePos := pixel.V(
		headPos.X+headSize/2*math.Cos(tailAngle-math.Pi/4),
		headPos.Y+headSize/2*math.Sin(tailAngle-math.Pi/4),
	)
	rightEyePos := pixel.V(
		headPos.X+headSize/2*math.Cos(tailAngle+math.Pi/4),
		headPos.Y+headSize/2*math.Sin(tailAngle+math.Pi/4),
	)
	imd.Push(leftEyePos)
	imd.Circle(eyeSize, outlineWidth)
	imd.Push(rightEyePos)
	imd.Circle(eyeSize, outlineWidth)

	// Draw the nucleus
	imd.Color = colornames.White
	imd.Push(p.nucleusPos)
	imd.Circle(nucleusRadius, 0)
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

// debugDisplay displays debugging information
func debugDisplay(particles []*Particle, pheromoneTrail []Pheromone, generation int) {
	// Clear the console
	fmt.Print("\033[H\033[2J")

	// Print debugging information
	fmt.Println("Generation:", generation)
	fmt.Println("Number of particles:", len(particles))
	fmt.Println("Number of pheromone trails:", len(pheromoneTrail))
	// Add more relevant debugging information here
}

// evolvePopulation applies genetic algorithms to evolve the population
func evolvePopulation(particles []*Particle) []*Particle {
	// Perform crossover and mutation
	// (Not implemented here; can be added based on specific requirements)

	return particles
}

func main() {
	rand.Seed(time.Now().UnixNano())
	pixelgl.Run(run)
}
