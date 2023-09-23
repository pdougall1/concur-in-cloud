package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"cloud.google.com/go/pubsub"
)

type MessageData struct {
	Target         string  `json:"target"`
	PopulationSize int     `json:"population_size"`
	MaxGenerations int     `json:"max_generations"`
	MutationRate   float64 `json:"mutation_rate"`
}

func main() {
	fmt.Println("Hello, World!")

	projectID := "development-366613"
	topicID := "concur-in-cloud"

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		panic(err)
	}

	t := client.Topic(topicID)
	defer t.Stop()

	m := MessageData{
		Target:         "SomeString",
		PopulationSize: 10000,
		MaxGenerations: 50000,
		MutationRate:   0.25,
	}

	data, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	keepPublishing := true
	for keepPublishing {
		jitter := rand.Intn(100)
		time.Sleep(time.Duration(jitter) * time.Millisecond)

		result := t.Publish(ctx, &pubsub.Message{
			Data: data,
		})

		go func(r *pubsub.PublishResult) {
			id, err := r.Get(ctx) // Block until the result is confirmed

			if err != nil {
				fmt.Printf("Failed to publish: %v\n", err)
				return
			}

			fmt.Printf("Published message; msg ID: %v\n", id)
		}(result)
	}

	fmt.Println("Done")
}
