package main

import (
	"fmt"
	"math"

	"github.com/tandahunter/solarutil"
)

//G Gravitational constant
const G = 6.67384e-11

//SEC_PER_STEP is the number of seconds in the sim per step
const SEC_PER_STEP = 8

func getPlanets() *solarutil.PlanetArray {
	planets := *solarutil.NewPlanetArray()
	planets = append(planets, *solarutil.NewPlanet("Earth", 1000, 1000, 0, 0))
	planets = append(planets, *solarutil.NewStar("Sun", 50000, 0, 0, 0))

	return &planets
}

func getAcceleration(distance, starMass float64) float64 {
	return G * starMass / (math.Pow(distance, float64(2)))
}

func main() {
	systemBodies := getPlanets()
	planets := systemBodies.GetPlanets()
	star := systemBodies.GetStars().FirstOrDefault()

	for _, i := range *planets {
		performOrbitalManoeuvre(star, &i)
	}
}

func performOrbitalManoeuvre(star, planet *solarutil.Planet) {
	d := star.DistanceTo(planet)
	a := getAcceleration(d, star.Mass)

	fmt.Printf("%f:%f", a, d)
}
