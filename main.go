package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

type PubSubMessage struct {
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

type Evolver interface {
	Evolve() (int, string)
}

func main() {
	logger.Info("Starting server...")
	http.HandleFunc("/", HandleTheMessage)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		logger.Info("Defaulting to port", "port", port)
	}

	logger.Info("Listening on port", "port", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logger.Error(err.Error())
		return
	}
}

func HandleTheMessage(w http.ResponseWriter, r *http.Request) {
	logger.Info("Handling message concur in cloud")

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error("Error reading body", "error", err.Error())
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	var pubsubMessage PubSubMessage

	if err := json.Unmarshal(bodyBytes, &pubsubMessage); err != nil {
		logger.Error("json.Unmarshal DATA", "error", err.Error(), "data", string(pubsubMessage.Message.Data))
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	logger.Info("Received raw data AFTER", "data", pubsubMessage.Message.Data)

	var md MessageData
	if err := json.Unmarshal(pubsubMessage.Message.Data, &md); err != nil {
		logger.Error("Error unmarshaling MessageData", "error", err.Error(), "data", string(pubsubMessage.Message.Data))
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	logger.Info("Message data", "message_data", md)

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

	evolver := NewEvolverConcurrent(target, populationSize, maxGenerations, mutationRate)
	// evolver := NewEvolverSequential(target, populationSize, maxGenerations, mutationRate)
	evolver.Evolve()

	logger.Info("Success", "generation_count", genCount, "closest", closest)

	w.WriteHeader(http.StatusOK)
}
