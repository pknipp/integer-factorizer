package main

import (
	"sort"
	"math/big"
)

func factorize(number *big.Int) (bool, [][2]*big.Int) {
	j := 1
	jBig := big.NewInt(int64(j))
	// Use a map for storing factors, to facilitate the algo.
	// key = prime, value = exponent
	factors := map[int]*big.Int{}
	// One only needs to search up until the square root of number.
	j2Big := big.NewInt(0).Mul(jBig, jBig)
	for j2Big.Cmp(number) == -1 {
		// After 3, only odd numbers can be prime (so step-size = 2, then)
		if jBig.Cmp(big.NewInt(3)) == -1 {
			jBig.Add(jBig, big.NewInt(1))
		} else {
			jBig.Add(jBig, big.NewInt(2))
		}
		for {
			// Continue dividing out (and tabulating) j until j is no longer a factor of number.
			if big.NewInt(0).Mod(number, jBig).Cmp(big.NewInt(0)) == 0 {
				_, facFound := factors[j]
				if facFound {
					factors[j].Add(factors[j], big.NewInt(1))
				} else {
					factors[j] = big.NewInt(1)
				}
				number.Div(number, jBig)
			} else {
				// Go to next possible factor.
				break
			}
		}
	}
	// The last factor is needed if the largest factor occurs by itself.
	if number.Cmp(big.NewInt(1)) != 0 {
		factors[int(number.Int64())] = big.NewInt(1)
	}
	// Below is a necessary - but not sufficient - condition of primacy.
	isPrime := len(factors) == 1
	// The condition below is required to make it "sufficient".
	if isPrime {
		for _, exponent := range factors {
			if exponent.Cmp(big.NewInt(1)) == 1 {
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
