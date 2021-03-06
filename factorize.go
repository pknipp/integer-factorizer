package main

import (
	"sort"
)

func factorize(number int) (bool, [][2]int) {
	j := 1
	// Use a map for storing factors, to facilitate the algo.
	// key = prime, value = exponent
	factors := map[int]int{}
	// One only needs to search up until the square root of number.
	for j * j < number {
		// After 3, only odd numbers can be prime (so step-size = 2, then)
		if j < 3 {
			j++
		} else {
			j += 2
		}
		for {
			// Continue dividing out (and tabulating) j until j is no longer a factor of number.
			if number % j == 0 {
				_, facFound := factors[j]
				if facFound {
					factors[j]++
				} else {
					factors[j] = 1
				}
				number /= j
			} else {
				// Go to next possible factor.
				break
			}
		}
	}
	// The last factor is needed if the largest factor occurs by itself.
	if number != 1 {
		factors[number] = 1
	}
	// Below is a necessary - but not sufficient - condition of primacy.
	isPrime := len(factors) == 1
	// The condition below is required to make it "sufficient".
	if isPrime {
		for _, exponent := range factors {
			if exponent > 1 {
				isPrime = false
			}
		}
	}
	// Change data structure from map to slice (of 2-component arrays), to facilitate sorting.
	factorsSorted := [][2]int{}
	for prime, exponent := range factors {
		factorsSorted = append(factorsSorted, [2]int{prime, exponent})
	}
	sort.Slice(factorsSorted, func(i, j int) bool {
		return factorsSorted[i][0] < factorsSorted[j][0]
	})
	return isPrime, factorsSorted
}
