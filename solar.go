package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os/exec"
	"time"

	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/tandahunter/solarutil"
)

const (
	//G is the gravitational constant
	G = 6.67384e-11
	//SecondsPerStep is the number of seconds in the sim per step
	//more seconds per step = less accuracy.  Originally 8
	SecondsPerStep = 200 //8
	//StepsPerFrame determines how many calculations should be calculated per cycle
	//more steps per frame = less performance.  Originally 10000
	StepsPerFrame = 250 // 10000
	//MetersPerUnit is a gigameter
	MetersPerUnit = 1000000000
	//WebServerPort is the port on which to open the web router
	WebServerPort = 8080
	//WebSocketPort is the port on which to open the web socket
	WebSocketPort = 8081
)

var planets *solarutil.PlanetArray
var sun *solarutil.Star
var manoeuvreCount = 0

func main() {
	//Create some planets, and a sun
	initSun()
	initPlanets()

	//Initiate a ticker to start generating some data / orbiting the planets
	initFrameTicker()

	//Start listening for requests
	initRouter()
}

func initRouter() {
	go func() {
		router := mux.NewRouter().StrictSlash(true)
		router.HandleFunc("/Sun/", getSun)
		router.HandleFunc("/Client/", getClient)
		router.HandleFunc("/Planets/", getPlanets)
		router.HandleFunc("/Planets/{id}/", getPlanet)
		router.HandleFunc("/Planets/{id}/Vector/", getPlanetVector)

		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", WebServerPort), router))
	}()

	exec.Command("cmd", "/c", "start", "http://localhost:8080/Client/").Start()
	fmt.Printf("Solar is listening on port:%d", WebServerPort)

	http.HandleFunc("/Planets/", streamPlanets)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", WebSocketPort), nil))
}

func initSun() {
	sun = solarutil.NewStar("Sun", 1.988435e30, 0, 0, 0)
}

func initPlanets() {
	planets = solarutil.NewPlanetArray()
	planets.AddNewPlanet(1, "Mercury", 3.30104e23, 50.32, 4.74e-5)
	planets.AddNewPlanet(2, "Venus", 4.86732e24, 108.8, 3.5e-5)
	planets.AddNewPlanet(3, "Earth", 5.9721986e24, 150, 2.963e-5)
	planets.AddNewPlanet(4, "Mars", 6.41693e23, 227.94, 0.0000228175)
	planets.AddNewPlanet(5, "Jupiter", 1.89813e27, 778.33, 0.0000129824)
}

func initFrameTicker() {
	ticker := time.NewTicker(time.Millisecond * 20)

	go func() {
		for range ticker.C {
			for _, planet := range *planets {
				go func(s *solarutil.Star, p *solarutil.Planet) {
					performOrbitalManoeuvre(s, p)
				}(sun, planet)
			}
			manoeuvreCount++
		}
	}()
}

func validateGet(w http.ResponseWriter, r *http.Request) bool {
	return validateRestMethod(w, r, http.MethodGet)
}

func validatePost(w http.ResponseWriter, r *http.Request) bool {
	return validateRestMethod(w, r, http.MethodPost)
}

func validateRestMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method != method {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return false
	}
	return true
}

func setCORSHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func getClient(w http.ResponseWriter, r *http.Request) {
	if validateGet(w, r) {
		setCORSHeader(w)

		f, err := ioutil.ReadFile("resources/client.html")
		if err != nil {
			http.Error(w, "Failed dependency", http.StatusFailedDependency)
		} else {
			w.Write(f)
		}
	}
}

func getSun(w http.ResponseWriter, r *http.Request) {
	if validateGet(w, r) {
		setCORSHeader(w)
		json.NewEncoder(w).Encode(sun)
	}
}

func getPlanets(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		setCORSHeader(w)
		json.NewEncoder(w).Encode(planets)

	} else if r.Method == http.MethodPost {
		setCORSHeader(w)

		var planet solarutil.Planet

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&planet)

		if err == nil {
			*planets = append(*planets, &planet)
			planet.ID = len(*planets)
			json.NewEncoder(w).Encode(planet.ID)

			log.Println("Added planet")
		} else {
			http.Error(w, err.Error(), http.StatusExpectationFailed)
			log.Println(err.Error())
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getPlanetVector(w http.ResponseWriter, r *http.Request) {
	if validateGet(w, r) {
		setCORSHeader(w)

		sid := mux.Vars(r)["id"]
		id, err := strconv.Atoi(sid)

		var planet *solarutil.Planet
		if err == nil {
			planet = planets.GetPlanetByID(id)
		}

		if planet != nil {
			json.NewEncoder(w).Encode(planet.Vector)
		} else {
			http.Error(w, "Planet not found", http.StatusNotFound)
		}
	}
}

func getPlanet(w http.ResponseWriter, r *http.Request) {
	if validateGet(w, r) {
		setCORSHeader(w)

		sid := mux.Vars(r)["id"]
		id, err := strconv.Atoi(sid)

		var planet *solarutil.Planet
		if err == nil {
			planet = planets.GetPlanetByID(id)
		}

		if planet != nil {
			json.NewEncoder(w).Encode(planet)
		} else {
			http.Error(w, "Planet not found", http.StatusNotFound)
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func streamPlanets(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	c, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("Could not open websocket connection")
		return
	}
	defer c.Close()

	for {
		if err = c.WriteJSON(planets); err != nil {
			log.Println(err)
			break
		}

		time.Sleep(1000)
	}
}

func performOrbitalManoeuvre(star *solarutil.Star, planet *solarutil.Planet) {
	vel := solarutil.NewVector(0, 0, 0)
	speed := float64(0)

	for i := 0; i < StepsPerFrame; i++ {
		distance := planet.DistanceTo(sun)
		speed = getAcceleration(distance*MetersPerUnit, star.Mass) * SecondsPerStep

		vel.SubVectors(star.Vector, planet.Vector)
		vel.SetLength(speed / MetersPerUnit)
		planet.Velocity.Add(vel)

		planet.Vector.X += planet.Velocity.X * SecondsPerStep
		planet.Vector.Y += planet.Velocity.Y * SecondsPerStep
		planet.Vector.Z += planet.Velocity.Z * SecondsPerStep
	}
}

func getAcceleration(distance, starMass float64) float64 {
	return G * starMass / (math.Pow(distance, float64(2)))
}
