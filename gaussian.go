package main

import (
	"strconv"
	"math"
)

var ZERO [2]int = [2]int{0, 0}
var MAXINT float64 = 0.999999 * math.Pow(2., 63.)

func gaussian(z [2]int) (bool, int, map[string][2]int) {
	gaussianFactors := map[string][2]int{}
	_, factors := factorize(modulus(z))
	for _, pair := range factors {
		prime := pair[0]
		exponent := pair[1]
		// Here are the factors of 1 + i
		if prime == 2 {
			gaussianFactors["1+i"] = [2]int{2, exponent}
			for count := 0; count < exponent; count++ {
				_, z = modulo(z, [2]int{1, 1})
			}
		} else {
			// Here are the (irreducible) real prime factors, which occur in pairs.
			if prime % 4 == 3 {
				gaussianFactors[strconv.Itoa(prime)] = [2]int{prime, exponent / 2}
				for count := 0; count < exponent / 2; count++ {
					for i, _ := range z {
						z[i] /= prime
					}
				}
			} else {
				// Here are Gaussian integers for which one component is odd and the other is even.
				// Find ints m, n such that (2m+1)^2 + (2n)^2 = mod4
				mod4 := (prime - 1) / 4
				// Now this becomes m*(m+1) + n^2 = mod4
				m := 0
				var n int
				for {
					n64 := math.Sqrt(float64(mod4 - m * (m + 1)))
					nm := int(math.Floor(n64))
					np := int(math.Ceil(n64))
					if m * (m + 1) + nm * nm == mod4 {
						n = nm
						break
					} else if m * (m + 1) + np * np == mod4 {
						n = np
						break
					}
					m++
				}
				odd := 2 * m + 1
				even := 2 * n
				// First, let's consider possibility that the real component is the odd one.
				count := 0
				for {
					isFactor, quotient := modulo(z, [2]int{odd, even})
					if isFactor {
						z = quotient
						count++
					} else {
						if count > 0 {
							im := strconv.Itoa(even)
							gaussianFactors[strconv.Itoa(odd) + "+" + im + "i"] = [2]int{prime, count}
						}
						break
					}
				}
				// For the remaining factors, the real component must be the even one.
				count2 := exponent - count
				if count2 > 0 {
					im := strconv.Itoa(odd)
					if im == "1" {
						im = ""
					}
					gaussianFactors[strconv.Itoa(even) + "+" + im + "i"] = [2]int{prime, count2}
				}
				for count = 0; count < count2; count++ {
					_, z = modulo(z, [2]int{2 * n, 2 * m + 1})
				}
			}
		}
	}
	// The following logic is a bit obtuse, but it determines exponent of i, based upon what is left after dividing by all Gaussian primes.
	var n int
	if math.Abs(float64(z[0])) == 1 {
		n = 1 - z[0]
	} else {
		n = 2 - z[1]
	}
	// Below is a necessary - but not sufficient - condition.
	isPrime := len(gaussianFactors) == 1
	// The next condition is required to make it "sufficient".
	if isPrime {
		for _, pair := range gaussianFactors {
			if pair[1] > 1 {
				isPrime = false
				break
			}
		}
	}
	return isPrime, n, gaussianFactors
}
