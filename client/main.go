package main

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/iterator"
)

func list(projectID string) ([]*pubsub.Topic, error) {
	// projectID := "my-project-id"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("pubsub.NewClient: %w", err)
	}
	defer client.Close()

	var topics []*pubsub.Topic

	it := client.Topics(ctx)
	for {
		topic, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Next: %w", err)
		}
		topics = append(topics, topic)
	}

	return topics, nil
}

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

	list, err := list(projectID)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Topics: %v\n", list[0].ID())

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

	result := t.Publish(ctx, &pubsub.Message{
		Data: []byte(data),
	})

	id, err := result.Get(ctx)
	if err != nil {
		fmt.Printf("Failed to publish: %v\n", err)
		return
	}
	fmt.Printf("Published message; msg ID: %v\n", id)
}
