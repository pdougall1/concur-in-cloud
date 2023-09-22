package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

type TheMessage struct {
	Message struct {
		Data       []byte `json:"data"`
		Attributes struct {
			AccountID string `json:"account_id"`
		} `json:"attributes"`

		ID string `json:"id"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

type MessageData struct {
	Target         string  `json:"target"`
	PopulationSize int     `json:"population_size"`
	MaxGenerations int     `json:"max_generations"`
	MutationRate   float64 `json:"mutation_rate"`
}

func main() {
	logger.Info("Starting server...")
	http.HandleFunc("/", HandleTheMessage)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		logger.Info("Defaulting to port %s", port)
	}

	logger.Info("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logger.Error(err.Error())
		return
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
		logger.Error("ioutil.ReadAll: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &m); err != nil {
		logger.Error("json.Unmarshal body: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var md MessageData
	if err := json.Unmarshal(m.Message.Data, &md); err != nil {
		logger.Error("json.Unmarshal DATA: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	target := string(md.Target)
	if target == "" {
		target = "SomeStringToEvolveTo"
	}

	populationSize := int(md.PopulationSize)
	if populationSize == 0 {
		populationSize = 1000
	}

	maxGenerations := int(md.MaxGenerations)
	if maxGenerations == 0 {
		maxGenerations = 5000
	}

	mutationRate := float64(md.MutationRate)
	if mutationRate == 0 {
		mutationRate = 0.01
	}

	logger.Info("Evolving", "target", target, "populationSize", populationSize, "maxGenerations", maxGenerations, "mutationRate", mutationRate)
	genCount, closest := EvolveStrSeq(target, populationSize, maxGenerations, mutationRate)

	logger.Info("Success: %d : %s", genCount, closest)
}
