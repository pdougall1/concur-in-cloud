package main

// import (
// 	"math/rand"
// 	"sort"
// 	"time"
// )

// const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

// func randomString(length int) string {
// 	b := make([]byte, length)
// 	for i := range b {
// 		b[i] = charset[seededRand.Intn(len(charset))]
// 	}
// 	return string(b)
// }

// func generatePopulation(size int, length int) []string {
// 	population := make([]string, size)
// 	for i := range population {
// 		population[i] = randomString(length)
// 	}
// 	return population
// }

// func fitness(s1, s2 string) int {
// 	count := 0
// 	for i := range s1 {
// 		if s1[i] != s2[i] {
// 			count++
// 		}
// 	}
// 	return count
// }

// // Select two parents with a preference for strings with better fitness
// func selectParents(population []string, target string) (string, string) {
// 	sort.SliceStable(population, func(i, j int) bool {
// 		return fitness(population[i], target) < fitness(population[j], target)
// 	})
// 	return population[0], population[1]
// }

// func crossover(parent1, parent2 string) string {
// 	half := len(parent1) / 2
// 	return parent1[:half] + parent2[half:]
// }

// func mutate(s string, mutationRate float64) string {
// 	if seededRand.Float64() < mutationRate {
// 		index := seededRand.Intn(len(s))
// 		char := charset[seededRand.Intn(len(charset))]
// 		return s[:index] + string(char) + s[index+1:]
// 	}
// 	return s
// }

// func replaceLeastFit(population []string, newIndividuals []string, target string) {
// 	sort.SliceStable(population, func(i, j int) bool {
// 		return fitness(population[i], target) > fitness(population[j], target)
// 	})
// 	copy(population[:len(newIndividuals)], newIndividuals)
// }

// func reachedTarget(population []string, target string) (bool, string) {
// 	for _, s := range population {
// 		if s == target {
// 			return true, s
// 		}
// 	}
// 	return false, ""
// }

// func EvolveStrSeq(target string, populationSize int, maxGenerations int, mutationRate float64) (int, string) {
// 	population := generatePopulation(populationSize, len(target))
// 	var parent1, parent2 string

// 	generation := 0
// 	reached, reachedStr := reachedTarget(population, target)
// 	for !reached && generation < maxGenerations {
// 		parent1, parent2 = selectParents(population, target)
// 		child := crossover(parent1, parent2)
// 		child = mutate(child, mutationRate)
// 		replaceLeastFit(population, []string{child}, target)
// 		generation++
// 		reached, reachedStr = reachedTarget(population, target)
// 	}

// 	return generation, reachedStr
// }
