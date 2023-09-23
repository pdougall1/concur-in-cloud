package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"cloud.google.com/go/pubsub"
)

const times = 1

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
		MutationRate:   0.5,
	}

	data, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	results := make([]*pubsub.PublishResult, times)

	for i := 0; i < times; i++ {
		results[i] = t.Publish(ctx, &pubsub.Message{
			Data: data,
		})
	}

	wg := sync.WaitGroup{}
	for _, r := range results {
		wg.Add(1)

		go func(r *pubsub.PublishResult) {
			id, err := r.Get(ctx) // Block until the result is confirmed

			if err != nil {
				fmt.Printf("Failed to publish: %v\n", err)
				return
			}

			fmt.Printf("Published message; msg ID: %v\n", id)

			wg.Done()
		}(r)
	}

	wg.Wait()

	fmt.Println("Done")
}
