package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
)

type WarmUp struct {
	Exercises []Exercise
	TotalTime int
	Name      string
}
type Exercises []Exercise

type Exercise struct {
	Name     string `json:"name"`
	MaxRepts int    `json:"maxRepts"`
}

const ExercisesFile string = "exercises.json"
const MaxTime int = 180
const MinReps int = 8

var Animals = [...]string{"Bear", "Dog", "Ant", "Dolphin", "Cat", "Bird"}
var Extra = [...]string{"Fast", "Heavy", "Radioactive", "Fun", "Busy", "Super"}

func main() {
	exercises := loadExercises(ExercisesFile)
	warmUp := getWarmUp(exercises)
	displayWarmUp(warmUp)
}

func getWarmUp(e Exercises) WarmUp {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(e), func(i, j int) { e[i], e[j] = e[j], e[i] })
	totalTime := 0
	done := false
	w := WarmUp{}
	w.Name = generateTrainingName()
	for i := 0; i < len(e); i++ {
		ex := e[i]
		exRepts := rand.Intn(ex.MaxRepts-MinReps) + MinReps

		exTime := exRepts * 2
		if totalTime+exTime > MaxTime {
			exRepts = (MaxTime - totalTime) / 2
			if exRepts < MinReps {
				exRepts = 0
			}
			exTime = exRepts * 2
			done = true
		}

		if exRepts > 0 {
			totalTime += exTime
			exercise := Exercise{Name: ex.Name, MaxRepts: exRepts}
			w.Exercises = append(w.Exercises, exercise)
		}

		if done == true {
			break
		}
	}
	w.TotalTime = totalTime

	return w
}

func displayWarmUp(w WarmUp) {
	fmt.Println("The training is called:", w.Name)
	fmt.Println("It takes", w.TotalTime, "seconds or", secondsToMinutes(w.TotalTime))
	fmt.Println("The exercises are:")
	for i := 0; i < len(w.Exercises); i++ {
		displayExercise(w.Exercises[i])
	}
}

func displayExercise(e Exercise) {
	fmt.Println(e.MaxRepts, e.Name)
}

func loadExercises(file string) Exercises {
	jsonFile, err := os.Open(file)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Exercises array
	var exercises Exercises

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'exercises' which we defined above
	json.Unmarshal(byteValue, &exercises)

	return exercises
}

func generateTrainingName() string {

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	posA := r1.Intn(len(Animals))
	posE := r1.Intn(len(Extra))

	name := Extra[posE] + " " + Animals[posA]

	return name
}

func secondsToMinutes(inSeconds int) string {
	strTime := time.Duration(inSeconds) * time.Second
	return shortDur(strTime)
}

func shortDur(d time.Duration) string {
	s := d.String()
	if strings.HasSuffix(s, "m0s") {
		s = s[:len(s)-2]
	}
	if strings.HasSuffix(s, "h0m") {
		s = s[:len(s)-2]
	}
	return s
}
