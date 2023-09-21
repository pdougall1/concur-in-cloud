package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

type TheMessage struct {
	Message struct {
		Attributes struct {
			Target         string  `json:"target"`
			PopulationSize int     `json:"populationSize"`
			MaxGenerations int     `json:"maxGenerations"`
			MutationRate   float64 `json:"mutationRate"`
		} `json:"attributes"`

		ID string `json:"id"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

func main() {
	http.HandleFunc("/", HandleTheMessage)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

	// target := "SomeString"
	// populationSize := 10000
	// maxGenerations := 50000
	// mutationRate := 0.5

	// genCount, closest := EvolveStrSeq(target, populationSize, maxGenerations, mutationRate)
	// maxFitness := fitness(closest, target)
	// fmt.Printf("Success: %d : %d : %s", maxFitness, genCount, closest)
}

func HandleTheMessage(w http.ResponseWriter, r *http.Request) {
	var m TheMessage
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &m); err != nil {
		log.Printf("json.Unmarshal: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	target := string(m.Message.Attributes.Target)
	if target == "" {
		target = "SomeStringToEvolveTo"
	}

	populationSize := int(m.Message.Attributes.PopulationSize)
	if populationSize == 0 {
		populationSize = 1000
	}

	maxGenerations := int(m.Message.Attributes.MaxGenerations)
	if maxGenerations == 0 {
		maxGenerations = 5000
	}

	mutationRate := float64(m.Message.Attributes.MutationRate)
	if mutationRate == 0 {
		mutationRate = 0.01
	}

	genCount, closest := EvolveStrSeq(target, populationSize, maxGenerations, mutationRate)

	log.Printf("Success: %d : %s", genCount, closest)
}
