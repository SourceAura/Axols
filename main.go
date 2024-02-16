package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

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
	Population Population
	Generation int
}

// Init initializes the application
func (m Model) Init() tea.Cmd {
	m.Population.Init(10, 10)
	return func() tea.Msg {
		return tickMsg{}
	}
}

// Update updates the application state
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit
	case tickMsg:
		m.Population.Evolve(0.01) // Set mutation rate to 0.01
		m.Generation++
		return m, func() tea.Msg {
			return tickMsg{}
		}
	}
	return m, nil
}

// View renders the application
func (m Model) View() string {
	return fmt.Sprintf("Generation: %d\nPopulation: %+v", m.Generation, m.Population.Creatures)
}

// Main entry point
func main() {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Start the program
	if err := tea.NewProgram(Model{}).Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting program: %v", err)
		os.Exit(1)
	}
}

// Custom tick message
type tickMsg struct{}

func (tickMsg) String() string {
	return "tick"
}

func (tickMsg) Erase() bool {
	return true
}
