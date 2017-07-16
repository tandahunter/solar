package main

import (
	"math"

	"github.com/tandahunter/solarutil"
)

//G Gravitational constant
const G = 6.67384e-11

//SecondsPerStep is the number of seconds in the sim per step
const SecondsPerStep = 8

//StepsPerFrame determines how many calculations should be calculated per cycle
const StepsPerFrame = 10000

//MetersPerUnit is a gigameter
const MetersPerUnit = 1000000000

func getAcceleration(distance, starMass float64) float64 {
	return G * starMass / (math.Pow(distance, float64(2)))
}

func main() {
	solarBodies := *solarutil.NewPlanetArray()
	solarBodies = append(solarBodies, *solarutil.NewStar("Sun", 1.988435e30, 0, 0, 0))
	solarBodies = append(solarBodies, *solarutil.NewPlanet("Earth", 5.9721986e24, 150, 0, 0, solarutil.NewVector(0, 0, 2.963e-5)))
	star := solarBodies.GetStars().FirstOrDefault()

	planets := solarBodies.GetPlanets()
	planets.PrintNames()

	for x := 0; x < 2000; x++ {
		for _, i := range *planets {
			performOrbitalManoeuvre(star, &i)
		}
	}

	earth := solarBodies.GetPlanetsByName("Earth").FirstOrDefault()
	earth.PrintVector()
}

func performOrbitalManoeuvre(star, planet *solarutil.Planet) {
	vel := solarutil.NewVector(0, 0, 0)
	speed := float64(0)

	for i := 0; i < StepsPerFrame; i++ {
		distance := star.DistanceTo(planet)
		speed = getAcceleration(distance*MetersPerUnit, star.Mass) * SecondsPerStep

		vel.SubVectors(star.Vector, planet.Vector)
		vel.SetLength(speed / MetersPerUnit)
		planet.Velocity.Add(vel)

		planet.Vector.X += planet.Velocity.X * SecondsPerStep
		planet.Vector.Y += planet.Velocity.Y * SecondsPerStep
		planet.Vector.Z += planet.Velocity.Z * SecondsPerStep
	}
}
