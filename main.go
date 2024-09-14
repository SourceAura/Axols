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
)

const (
	windowWidth        = 900
	windowHeight       = 700
	numAxols           = 100 // Reduced number for clarity
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

// Axol represents a single particle in the simulation
type Axol struct {
	pos           pixel.Vec
	vel           pixel.Vec
	genome        Genome
	species       int
	tailAngle     float64
	consumedFood  int     // Add this field
	timeSinceLast float64 // Add this field
}

type Genome struct {
	size        float64
	speed       float64
	senseRadius float64
	color       color.RGBA
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
		Title:  "Axol Simulation",
		Bounds: pixel.R(0, 0, windowWidth, windowHeight),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	imd := imdraw.New(nil)

	// Create axols of two species
	axols := make([]*Axol, numAxols)
	for i := range axols {
		species := i % 2 // Alternates between 0 and 1
		axols[i] = NewAxol(species)
	}

	// Create food sources
	foods := make([]Food, 0)

	deepPurple := color.RGBA{R: 20, G: 0, B: 30, A: 255}

	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		win.Clear(deepPurple)
		imd.Clear()

		// Update and draw axols
		for _, a := range axols {
			a.Update(dt, foods)
			a.Draw(imd)
		}

		// Update and draw food sources
		updateFoodSources(&foods, dt)
		drawFoodSources(imd, foods)

		imd.Draw(win)
		win.Update()
	}
}

// NewAxol creates a new axol with random properties
func NewAxol(species int) *Axol {
	var genome Genome
	if species == 0 {
		genome = Genome{
			size:        5,
			speed:       50,
			senseRadius: 30,
			color:       color.RGBA{R: 100, G: 200, B: 255, A: 150},
		}
	} else {
		genome = Genome{
			size:        7,
			speed:       40,
			senseRadius: 40,
			color:       color.RGBA{R: 255, G: 100, B: 200, A: 150},
		}
	}

	return &Axol{
		pos:           pixel.V(rand.Float64()*windowWidth, rand.Float64()*windowHeight),
		vel:           pixel.V(rand.Float64()*2-1, rand.Float64()*2-1).Unit().Scaled(genome.speed),
		genome:        genome,
		species:       species,
		tailAngle:     0,
		consumedFood:  0, // Initialize this field
		timeSinceLast: 0, // Initialize this field
	}
}

// Update updates the position of the axol
func (a *Axol) Update(dt float64, foods []Food) {
	a.pos = a.pos.Add(a.vel.Scaled(dt))
	if a.pos.X < 0 || a.pos.X > windowWidth {
		a.vel.X = -a.vel.X
	}
	if a.pos.Y < 0 || a.pos.Y > windowHeight {
		a.vel.Y = -a.vel.Y
	}
	a.tailAngle += 6 * dt // Reduced from 10 to 6 to slow down the animation
	a.timeSinceLast += dt

	// Consume food if within range
	for _, food := range foods {
		if a.pos.To(food.pos).Len() <= consumeRadius {
			a.consumeFood(food)
			break
		}
	}
}

// Draw draws the axol
func (a *Axol) Draw(imd *imdraw.IMDraw) {
	// Draw translucent body
	imd.Color = a.genome.color
	imd.Push(a.pos)
	imd.Circle(a.genome.size, 0)

	// Draw nucleus
	nucleusColor := color.RGBA{R: 255, G: 255, B: 255, A: 200}
	imd.Color = nucleusColor
	imd.Push(a.pos)
	imd.Circle(a.genome.size/3, 0)

	// Draw wiggling tail
	tailLength := a.genome.size * 3 // Reduced from 4 to 3
	tailSegments := 20
	waveFrequency := 2.0
	maxAmplitude := a.genome.size * 0.25 // Reduced from 0.3 to 0.25

	imd.Color = a.genome.color
	for i := 0; i <= tailSegments; i++ {
		t := float64(i) / float64(tailSegments)
		segmentPos := a.pos.Add(a.vel.Unit().Scaled(-t * tailLength))

		// Calculate wiggle offset (reversed t in the sine function)
		wiggleOffset := math.Sin(a.tailAngle+(1-t)*waveFrequency*math.Pi) * maxAmplitude * t

		// Apply offset perpendicular to tail direction
		segmentPos = segmentPos.Add(a.vel.Normal().Scaled(wiggleOffset))

		imd.Push(segmentPos)
	}
	imd.Line(a.genome.size * 0.2) // Adjust line thickness based on Axol size
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
func debugDisplay(particles []*Axol, foods []Food, generation int) {
	// Clear the console
	fmt.Print("\033[H\033[2J")

	// Print debugging information
	fmt.Println("Generation:", generation)
	fmt.Println("Number of particles:", len(particles))
	fmt.Println("Number of food sources:", len(foods))
	// Add more relevant debugging information here
}

// evolvePopulation applies genetic algorithms to evolve the population
func evolvePopulation(axols []*Axol, foods []Food) []*Axol {
	for _, a := range axols {
		// Speed up evolution based on food consumption
		if a.consumedFood > 0 {
			for i := 0; i < a.consumedFood; i++ {
				a.mutate()
			}
		}
	}

	return axols
}

// consumeFood consumes a food source and increases nutrition
func (a *Axol) consumeFood(food Food) {
	a.consumedFood++
	a.timeSinceLast = 0
}

// mutate applies mutation to the particle's traits
func (p *Axol) mutate() {
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
