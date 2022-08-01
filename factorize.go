package main

import (
	"sort"
	"math/big"
)

type factor struct {
	prime *big.Int
	exponent int
}

func factorize(number *big.Int) (bool, []factor) {
	j := big.NewInt(int64(1))
	// Use a map for storing factors, to facilitate the algo.
	// key = prime, value = exponent
	factors := map[*big.Int]int{}
	// One only needs to search up until the square root of number.
	j2 := big.NewInt(0).Mul(j, j)
	for j2.Cmp(number) == -1 {
		// After 3, only odd numbers can be prime (so step-size = 2, then)
		if j.Cmp(big.NewInt(3)) == -1 {
			j.Add(j, big.NewInt(1))
		} else {
			j.Add(j, big.NewInt(2))
		}
		for {
			// Continue dividing out (and tabulating) j until j is no longer a factor of number.
			if big.NewInt(0).Mod(number, j).Cmp(big.NewInt(0)) == 0 {
				_, facFound := factors[j]
				if facFound {
					factors[j]++
				} else {
					factors[j] = 1
				}
				number.Div(number, j)
			} else {
				// Go to next possible factor.
				break
			}
		}
	}
	// The last factor is needed if the largest factor occurs by itself.
	if number.Cmp(big.NewInt(1)) != 0 {
		factors[number] = 1
	}
	// Below is a necessary - but not sufficient - condition of primacy.
	isPrime := len(factors) == 1
	// The condition below is required to make it "sufficient".
	if isPrime {
		for _, exponent := range factors {
			if exponent == 1 {
				isPrime = false
			}
		}
	}
	// Change data structure from map to slice (of 2-component arrays), to facilitate sorting.

	factorsSorted := []factor{}
	for prime, exponent := range factors {
		factorsSorted = append(factorsSorted, factor{prime, exponent})
	}
	sort.Slice(factorsSorted, func(i, j int) bool {
		return factorsSorted[i].prime.Cmp(factorsSorted[j].prime) == -1
	})
	return isPrime, factorsSorted
}
