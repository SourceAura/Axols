package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Define tickMsg type for handling tick messages
type tickMsg struct{}

// Creature represents an individual in the population
type Creature struct {
	Genes []int
}

// Population represents a group of creatures
type Population struct {
	Creatures []Creature
}

// Init initializes the initial population
func (p *Population) Init(size, geneLength int) {
	p.Creatures = make([]Creature, size)
	for i := range p.Creatures {
		p.Creatures[i].Genes = make([]int, geneLength)
		for j := range p.Creatures[i].Genes {
			p.Creatures[i].Genes[j] = rand.Intn(2)
		}
	}
}

// Evolve performs the evolution process
func (p *Population) Evolve(mutationRate float64) {
	// For simplicity, let's assume random mating
	for i := 0; i < len(p.Creatures); i += 2 {
		parent1 := p.Creatures[i]
		parent2 := p.Creatures[i+1]
		child1Genes := crossover(parent1.Genes, parent2.Genes)
		child2Genes := crossover(parent2.Genes, parent1.Genes)
		p.Creatures[i].Genes = mutate(child1Genes, mutationRate)
		p.Creatures[i+1].Genes = mutate(child2Genes, mutationRate)
	}
}

// crossover performs crossover between two parents' genes
func crossover(parent1, parent2 []int) []int {
	midpoint := rand.Intn(len(parent1))
	childGenes := make([]int, len(parent1))
	copy(childGenes[:midpoint], parent1[:midpoint])
	copy(childGenes[midpoint:], parent2[midpoint:])
	return childGenes
}

// mutate performs mutation on the genes
func mutate(genes []int, mutationRate float64) []int {
	for i := range genes {
		if rand.Float64() < mutationRate {
			genes[i] = 1 - genes[i] // Flip the bit
		}
	}
	return genes
}

// Model defines the application's data and logic
type Model struct {
	Population     Population
	Generation     int
	MutationRate   float64 // Mutation rate of the population
	PopulationSize int     // Size of the population
}

// Add interactive controls to adjust simulation parameters
func (m Model) Init() tea.Cmd {
	// Initialize the population with the new parameters
	m.Population.Init(m.PopulationSize, 10) // Adjust gene length as needed
	return func() tea.Msg {
		return tickMsg{}
	}
}

// Update method to handle key presses to adjust simulation parameters
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "+":
			m.MutationRate += 0.01
		case "-":
			m.MutationRate -= 0.01
			if m.MutationRate < 0 {
				m.MutationRate = 0
			}
		}
		// Ensure mutation rate is within a valid range
		if m.MutationRate < 0 {
			m.MutationRate = 0
		}
		// Evolution step
		m.Population.Evolve(m.MutationRate)
		m.Generation++
		return m, func() tea.Msg {
			return tickMsg{}
		}
	case tickMsg:
		// Trigger evolution periodically
		m.Population.Evolve(m.MutationRate)
		m.Generation++
		return m, func() tea.Msg {
			return tickMsg{}
		}
	}
	return m, nil
}

// View method to display ASCII art representation of creatures and information panel
func (m Model) View() string {
	var sb strings.Builder

	// ASCII art representation of creatures
	for _, creature := range m.Population.Creatures {
		for _, gene := range creature.Genes {
			if gene == 0 {
				sb.WriteString("o")
			} else {
				sb.WriteString("X")
			}
		}
		sb.WriteString("\n")
	}

	// Information panel
	sb.WriteString(fmt.Sprintf("\nGeneration: %d\n", m.Generation))
	sb.WriteString(fmt.Sprintf("Mutation Rate: %.2f\n", m.MutationRate))
	sb.WriteString(fmt.Sprintf("Population Size: %d\n", m.PopulationSize))

	// Display additional statistics
	bestFitness := calculateBestFitness(m.Population)
	averageFitness := calculateAverageFitness(m.Population)
	sb.WriteString(fmt.Sprintf("Best Fitness: %d\n", bestFitness))
	sb.WriteString(fmt.Sprintf("Average Fitness: %.2f\n", averageFitness))
	// Add more statistics as needed

	return sb.String()
}

// calculateBestFitness calculates the best fitness score in the population
func calculateBestFitness(pop Population) int {
	bestFitness := 0
	for _, creature := range pop.Creatures {
		fitness := calculateFitness(creature)
		if fitness > bestFitness {
			bestFitness = fitness
		}
	}
	return bestFitness
}

// calculateAverageFitness calculates the average fitness score of the population
func calculateAverageFitness(pop Population) float64 {
	totalFitness := 0
	for _, creature := range pop.Creatures {
		totalFitness += calculateFitness(creature)
	}
	return float64(totalFitness) / float64(len(pop.Creatures))
}

// calculateFitness calculates the fitness score of a creature
func calculateFitness(creature Creature) int {
	// Here you can define your fitness function
	// For example, count the number of 'X's in the genes
	fitness := 0
	for _, gene := range creature.Genes {
		if gene == 1 {
			fitness++
		}
	}
	return fitness
}

// Main entry point
func main() {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Initialize the Model with default parameters
	initialModel := Model{
		PopulationSize: 10, // Default population size
		MutationRate:   0.01,
	}

	// Start the program
	if err := tea.NewProgram(initialModel).Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting program: %v", err)
		os.Exit(1)
	}
}
