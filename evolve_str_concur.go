package main

import (
	"math/rand"
	"sort"
	"strings"
	"sync"
	"time"
)

func NewEvolverConcurrent(md MessageData) Evolver {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano() + int64(rand.Intn(1000000))))

	evolver := &EvolverConcurrent{
		target:         md.Target,
		populationSize: md.PopulationSize,
		maxGenerations: md.MaxGenerations,
		mutationRate:   md.MutationRate,
		seededRand:     seededRand,
	}

	evolver.generatePopulation()

	return evolver
}

type EvolverConcurrent struct {
	target         string
	populationSize int
	maxGenerations int
	mutationRate   float64
	seededRand     *rand.Rand
	population     []string
	generation     int
}

func (e *EvolverConcurrent) GetGenerations() int {
	return e.generation
}

func (e *EvolverConcurrent) randomString(length int) string {
	mu := sync.Mutex{}
	var b strings.Builder

	for i := 0; i < length; i++ {
		mu.Lock()
		b.WriteByte(charset[e.seededRand.Intn(len(charset))])
		mu.Unlock()
	}

	return b.String()
}

func (e *EvolverConcurrent) generatePopulation() {
	e.population = make([]string, e.populationSize)

	length := len(e.target)
	wg := sync.WaitGroup{}

	for i := range e.population {
		i := i
		wg.Add(1)
		go func(i int) {
			e.population[i] = e.randomString(length)
			wg.Done()
		}(i)
	}

	wg.Wait()

	sort.SliceStable(e.population, func(i, j int) bool {
		return fitness(e.population[i], e.target) < fitness(e.population[j], e.target)
	})
}

func (e *EvolverConcurrent) reachedTarget() bool {
	for _, s := range e.population {
		if s == e.target {
			return true
		}
	}

	return false
}

func (e *EvolverConcurrent) Evolve() {
	var parent1, parent2 string

	// It was at this point that I realized this can't really be done concurrently :sob:
	for !e.reachedTarget() && e.generation < e.maxGenerations {
		parent1, parent2 = e.selectParents()
		child := crossover(parent1, parent2)
		child = mutate(child, e.mutationRate, e.seededRand)
		e.replaceLeastFit([]string{child})
		e.generation++
	}
}

// Select two parents with a preference for strings with better fitness
func (e *EvolverConcurrent) selectParents() (string, string) {
	return e.population[0], e.population[1]
}

func (e *EvolverConcurrent) replaceLeastFit(newIndividuals []string) {
	sort.SliceStable(e.population, func(i, j int) bool {
		return fitness(e.population[i], e.target) > fitness(e.population[j], e.target)
	})
	copy(e.population[:len(newIndividuals)], newIndividuals)
}
