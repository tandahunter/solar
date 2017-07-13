package main

import (
	"fmt"
	"math"
)

//G Gravitational constant
var G = 6.67384e-11

//SEC_PER_STEP is the number of seconds in the sim per step
var SEC_PER_STEP = 8

type planetArray []planet

type planet struct {
	Name   string
	Mass   float64
	Vector vector
	IsStar bool
}

type vector struct {
	X float64
	Y float64
	Z float64
}

func newVector(x, y, z float64) vector {
	vector := vector{}
	vector.X = x
	vector.Y = y
	vector.Z = z

	return vector
}

func newPlanet(name string, mass float64, x, y, z float64) planet {
	planet := planet{}
	planet.Name = name
	planet.Mass = mass
	planet.Vector = newVector(x, y, z)
	planet.IsStar = false

	return planet
}

func newStar(name string, mass float64, x, y, z float64) planet {
	planet := newPlanet(name, mass, x, y, z)
	planet.IsStar = true
	return planet
}

func (p *planet) printName() {
	fmt.Print(p.Name)
	fmt.Println()
}

func (p *planet) distanceTo(p2 *planet) float64 {
	return 5
}

func (p planetArray) printNames() {
	for _, planet := range p {
		planet.printName()
	}
}

func (p planetArray) getPlanets() planetArray {
	return p.getFilteredPlanets(false)
}

func (p planetArray) getStars() planetArray {
	return p.getFilteredPlanets(true)
}

func (p planetArray) firstOrDefault() planet {
	if len(p) > 0 {
		return p[0]
	}

	return planet{}
}

func (p planetArray) getFilteredPlanets(isStar bool) planetArray {
	toReturn := planetArray{}

	for _, planet := range p {
		if planet.IsStar == isStar {
			toReturn = append(toReturn, planet)
		}
	}

	return toReturn
}

func getPlanets() planetArray {

	planets := planetArray{}
	planets = append(planets, newPlanet("Earth", 1000, 1000, 0, 0))
	planets = append(planets, newStar("Sun", 50000, 0, 0, 0))

	return planets
}

func getAcceleration(distance, starMass float64) float64 {
	return G * starMass / (math.Pow(distance, float64(2)))
}

func main() {
	systemBodies := getPlanets()
	planets := systemBodies.getPlanets()
	star := systemBodies.getStars().firstOrDefault()

	for _, i := range planets {
		performOrbitalManoeuvre(&star, &i)
	}
}

func performOrbitalManoeuvre(star, planet *planet) {
	d := star.distanceTo(planet)
	a := getAcceleration(d, star.Mass)

	fmt.Printf("%f", a)
}
