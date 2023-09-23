package main

import (
	"math/rand"
	"sort"
	"sync"
	"time"
)

func NewEvolverConcurrent(target string, populationSize int, maxGenerations int, mutationRate float64) Evolver {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano() + int64(rand.Intn(1000000))))

	return EvolverConcurrent{
		target:         target,
		populationSize: populationSize,
		maxGenerations: maxGenerations,
		mutationRate:   mutationRate,
		seededRand:     seededRand,
	}
}

type EvolverConcurrent struct {
	target         string
	populationSize int
	maxGenerations int
	mutationRate   float64
	seededRand     *rand.Rand
}

func (e EvolverConcurrent) Evolve() (int, string) {

}

func randomChars(length int, seededRand *rand.Rand) chan<- byte {
	c := make(chan<- byte, length)

	for i := 0; i < length; i++ {
		go func(index int) {
			c <- charset[seededRand.Intn(len(charset))]
		}(i)
	}
	return c
}

func generatePopulation(size int, length int, seededRand *rand.Rand) []string {
	population := make([]string, size)
	wg := sync.WaitGroup{}

	for i := range population {
		wg.Add(1)
		index := i
		go func(index int) {
			population[index] = randomString(length, seededRand)
		}(index)
	}

	wg.Wait()
	return population
}

func fitness(s1, s2 string) int {
	count := 0
	for i := range s1 {
		if s1[i] != s2[i] {
			count++
		}
	}
	return count
}

// Select two parents with a preference for strings with better fitness
func selectParents(population []string, target string) (string, string) {
	sort.SliceStable(population, func(i, j int) bool {
		return fitness(population[i], target) < fitness(population[j], target)
	})
	return population[0], population[1]
}

func crossover(parent1, parent2 string) string {
	half := len(parent1) / 2
	return parent1[:half] + parent2[half:]
}

func mutate(s string, mutationRate float64, seededRand *rand.Rand) string {
	if seededRand.Float64() < mutationRate {
		index := seededRand.Intn(len(s))
		char := charset[seededRand.Intn(len(charset))]
		return s[:index] + string(char) + s[index+1:]
	}
	return s
}

func replaceLeastFit(population []string, newIndividuals []string, target string) {
	sort.SliceStable(population, func(i, j int) bool {
		return fitness(population[i], target) > fitness(population[j], target)
	})
	copy(population[:len(newIndividuals)], newIndividuals)
}

func reachedTarget(population []string, target string) (bool, string) {
	for _, s := range population {
		if s == target {
			return true, s
		}
	}
	return false, ""
}

func EvolveStrConcurrent(target string, populationSize int, maxGenerations int, mutationRate float64) (int, string) {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano() + int64(rand.Intn(1000000))))
	population := generatePopulation(populationSize, len(target), seededRand)
	var parent1, parent2 string

	generation := 0
	reached, reachedStr := reachedTarget(population, target)
	for !reached && generation < maxGenerations {
		parent1, parent2 = selectParents(population, target)
		child := crossover(parent1, parent2)
		child = mutate(child, mutationRate, seededRand)
		replaceLeastFit(population, []string{child}, target)
		generation++
		reached, reachedStr = reachedTarget(population, target)
	}

	return generation, reachedStr
}
