package main

import (
	"sort"
	"regexp"
	"strings"
	"math/big"
)

type gaussFactor struct {
	prime string
	mod2 *big.Int
	exponent int
}

func gcdComplexParse(gStr string) ([]gaussFactor, string) {
	gs := []map[string]modExp{}
	result := []gaussFactor{}
	var message string
	if gStr = regexp.MustCompile(" ").ReplaceAllString(gStr, ""); len(gStr) == 0 {
		message = "Expression is missing."
	} else {
		// Create array of strings, each representing a gaussian integer
		gsStr := strings.Split(gStr, ",")
		for _, gStr := range gsStr {
			if gaussianInt, message := gaussianParse(gStr); len(message) != 0 {
				return result, message
			} else {
				_, _, gaussianFactors := gaussian(gaussianInt)
				gs = append(gs, gaussianFactors)
			}
		}
		mapResult := gcdComplex(gs)
		// Convert result from map to slice of structs, to enable sorting by squared modulus
		for prime, pair := range mapResult {
			mod2, exponent := pair[0], pair[1]
			result = append(result, gaussFactor{prime, mod2, exponent})
		}
		sort.Slice(result, func(i, j int) bool {
			return result[i].mod2 < result[j].mod2
		})
	}
	return result, message
}
