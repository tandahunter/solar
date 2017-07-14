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

func getAcceleration(distance, starMass float64) float64 {
	return G * starMass / (math.Pow(distance, float64(2)))
}

func main() {
	solarBodies := *solarutil.NewPlanetArray()
	solarBodies = append(solarBodies, *solarutil.NewStar("Sun", 50000, 0, 0, 0))
	solarBodies = append(solarBodies, *solarutil.NewPlanet("Mercury", 500, 600, 0, 0))
	solarBodies = append(solarBodies, *solarutil.NewPlanet("Venus", 600, 800, 0, 0))
	solarBodies = append(solarBodies, *solarutil.NewPlanet("Earth", 1000, 1000, 0, 0))
	solarBodies = append(solarBodies, *solarutil.NewPlanet("Mars", 800, 2000, 0, 0))

	star := solarBodies.GetStars().FirstOrDefault()

	planets := solarBodies.GetPlanets()
	planets.PrintNames()

	for _, i := range *planets {
		performOrbitalManoeuvre(star, &i)
	}
}

func performOrbitalManoeuvre(star, planet *solarutil.Planet) {
	d := star.DistanceTo(planet)
	a := getAcceleration(d, star.Mass)

	fmt.Printf("%f:%f\n", a, d)
}
