# Evo-Sim

Evo-Sim is a simple particle life simulation with evolutionary algorithms implemented in Go using the Pixel game library.

## Overview

Evo-Sim simulates the life of particles in a 2D environment. Each particle has its own set of properties such as position, velocity, radius, and genetic makeup. The simulation includes features like particle movement, pheromone trails, and evolutionary algorithms for trait selection.

## Features

- Particle life simulation
- Random generation of particle properties
- Particle movement with collision detection
- Pheromone trails left by particles
- Evolutionary algorithms for trait selection over time

## Requirements

To run Evo-Sim, you need to have Go installed on your system. You can download and install Go from the [official Go website](https://golang.org/).

## Installation

Clone the Evo-Sim repository:

```bash
git clone https://github.com/sourceaura/evo-sim.git

evo-sim/
├── README.md
├── main.go
└── frontend/
    └── main.go

Navigate to the project directory:

cd evo-sim 

Run the simulation:

go run main.go
