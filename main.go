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
	foodSpawnRate      = 0.01 // Probability of food spawning per frame
	consumeRadius      = 10   // Radius within which a particle can consume food
	nutritionPerFood   = 1    // Amount of nutrition gained per unit of food
	evolutionSpeedup   = 0.1  // Speedup factor for evolution due to consuming food
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
	pos           pixel.Vec // Position
	vel           pixel.Vec // Velocity
	radius        float64   // Radius
	gene          color.RGBA
	nucleusPos    pixel.Vec // Position of the nucleus
	biome         Biome     // Biome the particle belongs to
	consumedFood  int       // Counter for consumed food
	timeSinceLast float64   // Time since last food consumption
}

// Food represents a source of nutrition for particles
type Food struct {
	pos       pixel.Vec // Position
	radius    float64   // Radius
	color     color.RGBA
	biome     Biome   // Biome the food belongs to
	nutrition float64 // Nutrition value of the food
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

	// Create food sources
	foods := make([]Food, 0)

	generation := 0
	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		win.Clear(colornames.White)
		imd.Clear()

		// Update and draw particles
		for _, p := range particles {
			p.Update(dt, foods)
			p.Draw(imd)
		}

		// Update and draw food sources
		updateFoodSources(&foods, dt)
		drawFoodSources(imd, foods)

		imd.Draw(win)
		win.Update()

		// Call debugDisplay function with relevant data
		debugDisplay(particles, foods, generation)

		// Apply genetic algorithm
		particles = evolvePopulation(particles, foods)

		generation++
	}
}

// NewParticle creates a new particle with random properties
func NewParticle() *Particle {
	p := &Particle{
		pos:           pixel.V(rand.Float64()*windowWidth, rand.Float64()*windowHeight),
		radius:        rand.Float64()*(maxRadius-minRadius) + minRadius,
		gene:          randomColor(),
		nucleusPos:    pixel.V(0, 0),
		consumedFood:  0,
		timeSinceLast: 0,
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
func (p *Particle) Update(dt float64, foods []Food) {
	p.pos = p.pos.Add(p.vel.Scaled(dt))
	if p.pos.X < 0 || p.pos.X > windowWidth {
		p.vel.X = -p.vel.X
	}
	if p.pos.Y < 0 || p.pos.Y > windowHeight {
		p.vel.Y = -p.vel.Y
	}

	// Update nucleus position
	p.nucleusPos = p.pos.Add(pixel.V(-p.radius/2, p.radius/2))

	// Consume food if within range
	for _, food := range foods {
		if p.pos.To(food.pos).Len() <= consumeRadius {
			p.consumeFood(food)
			break
		}
	}

	// Update time since last food consumption
	p.timeSinceLast += dt
}

// Draw draws the particle
func (p *Particle) Draw(imd *imdraw.IMDraw) {
	imd.Color = p.gene

	// Draw the main body
	imd.Push(p.pos)
	imd.Circle(p.radius, 0)

	// Draw the nucleus
	imd.Color = colornames.Yellow
	imd.Push(p.nucleusPos)
	imd.Circle(nucleusRadius, 0)
}

// Update the food sources, generating new ones over time
func updateFoodSources(foods *[]Food, dt float64) {
	// Generate new food sources with a certain probability
	if rand.Float64() < foodSpawnRate {
		*foods = append(*foods, NewFood())
	}

	// Update existing food sources (not implemented here)
}

// Draw the food sources
func drawFoodSources(imd *imdraw.IMDraw, foods []Food) {
	for _, f := range foods {
		imd.Color = f.color
		imd.Push(f.pos)
		imd.Circle(f.radius, 0)
	}
}

// NewFood creates a new food source with random properties
func NewFood() Food {
	pos := pixel.V(rand.Float64()*windowWidth, rand.Float64()*windowHeight)
	radius := rand.Float64()*5 + 3
	color := randomColor()
	nutrition := rand.Float64() * 10

	// Assign a random biome
	biome := Biome(rand.Intn(4))

	return Food{pos, radius, color, biome, nutrition}
}

// debugDisplay displays debugging information
func debugDisplay(particles []*Particle, foods []Food, generation int) {
	// Clear the console
	fmt.Print("\033[H\033[2J")

	// Print debugging information
	fmt.Println("Generation:", generation)
	fmt.Println("Number of particles:", len(particles))
	fmt.Println("Number of food sources:", len(foods))
	// Add more relevant debugging information here
}

// evolvePopulation applies genetic algorithms to evolve the population
func evolvePopulation(particles []*Particle, foods []Food) []*Particle {
	for _, p := range particles {
		// Speed up evolution based on food consumption
		if p.consumedFood > 0 {
			for i := 0; i < p.consumedFood; i++ {
				p.mutate()
			}
		}
	}

	return particles
}

// consumeFood consumes a food source and increases nutrition
func (p *Particle) consumeFood(food Food) {
	p.consumedFood++
	p.timeSinceLast = 0
}

// mutate applies mutation to the particle's traits
func (p *Particle) mutate() {
	// Mutation can be implemented here based on specific requirements
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

func main() {
	rand.Seed(time.Now().UnixNano())
	pixelgl.Run(run)
}
