# Axols - An Evolution Simulation

Axols is an evolution simulation implemented in Go using the Pixel game library. It simulates the life and evolution of particle-like creatures called "axols" in a 2D environment.

## Overview

Axols simulates the life of small creatures in a 2D environment. Each axol has its own set of properties such as position, velocity, size, speed, sense radius, and genetic makeup. The simulation includes features like creature movement, food sources, and evolutionary algorithms for trait selection.

## Features

- Two distinct species of axols with different initial properties
- Random generation of axol properties and food sources
- Axol movement with boundary collision detection and food-seeking behavior
- Food sources that axols can consume for nutrition
- Evolutionary algorithms for trait selection over generations
- Mutation of axol properties (size, speed, sense radius, and color)
- Visualization of axols with translucent bodies, nuclei, and wiggling tails
- Multiple biomes (currently not fully implemented)
- Debug display for simulation statistics (not shown in the provided code)

## Simulation Details

- The simulation runs in generations, with each generation lasting 10 seconds
- Axols seek and consume food within their sense radius
- At the end of each generation, the population evolves based on food consumption
- The top 50% of axols survive and reproduce, creating offspring with mixed traits
- Mutation can occur during reproduction, slightly altering axol properties

## Requirements

To run Axols, you need to have Go installed on your system. You can download and install Go from the [official Go website](https://golang.org/).

You'll also need to install the Pixel game library and its dependencies. You can do this by running:

## License

This project is licensed under the Creative Commons Attribution 4.0 International License. To view a copy of this license, visit http://creativecommons.org/licenses/by/4.0/ or send a letter to Creative Commons, PO Box 1866, Mountain View, CA 94042, USA.