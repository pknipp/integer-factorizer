package main

import (
	"sort"
	"strconv"
	"regexp"
	"strings"
	"math"
)

var TOOLARGE string = "Your number is too large."

//Euclid's algorithm, which is extremely efficient.
func gcd2(n1, n2 int) int {
	for {
	  t := n2;
	  if n1, n2 = t, n1 % n2; n2 == 0 {
		return n1
	  }
	}
}

type gaussFactor struct {
	prime string
	mod2 int
	exponent int
}

func gcdComplexParse(gStr string) ([]gaussFactor, string) {
	gs := []map[string][2]int{}
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

func gcd2Complex(gaussa, gaussb map[string][2]int) map[string][2]int {
	gauss := map[string][2]int{}
	if len(gaussa) > len(gaussb) {
		gaussa, gaussb = gaussb, gaussa // better to iterate over a shorter map
	}
	for prime, paira := range gaussa {
		mod2, exponenta := paira[0], paira[1]
		if pairb, found := gaussb[prime]; found {
			exponentb := pairb[1]
			gauss[prime] = [2]int{mod2, int(math.Min(float64(exponenta), float64(exponentb)))}
		}
	}
	return gauss
}

func checkIntStr(nStr string) (int, string) {
	nStr = regexp.MustCompile(" ").ReplaceAllString(nStr, "")
	var number int
	missingNumber := "Number is missing."
	if nStr == "" {
		return number, missingNumber
	} else {
		if nStr[0:1] == "-" {
			nStr = nStr[1:]
		}
	}
	if nStr == "" {
		return number, missingNumber
	}
	badNumber := "There is something wrong with your number."
	if len(nStr) > 18 {
		if len(nStr) > 19 {
			return number, TOOLARGE
		}
		if len(nStr) == 19 {
			if numTrunc, err := strconv.Atoi(nStr[0:6]); err != nil {
				return number, badNumber
			} else if numTrunc > 922336 {
				return number, TOOLARGE
			}
		}
	}
	if _, err := strconv.ParseFloat(nStr, 64); err != nil {
		return number, badNumber
	}
	if number, err := strconv.Atoi(nStr); err != nil {
		return number, "Note that the number may not be a decimal."
	} else {
		return number, ""
	}
}

func gcdParse(nStr string) (int, string) {
	var message string
	var ns []int
	var n int
	nsStr := strings.Split(nStr, ",")
	for _, nStr := range nsStr {
		if n, message = checkIntStr(nStr); len(message) == 0 || n == 1 {
			ns = append(ns, n)
		} else {
			return 0, message
		}
	}
	return gcd(ns), ""
}

func gcd(ns []int) int {
	if len(ns) == 1 {
		return ns[0]
	} else {
		return gcd(append(ns[2:], gcd2(ns[0], ns[1])))
	}
}

func gcdComplex(gausss []map[string][2]int) map[string][2]int {
	if len(gausss) == 1 {
		return gausss[0]
	} else {
		return gcdComplex(append(gausss[2:], gcd2Complex(gausss[0], gausss[1])))
	}
}

func factorizeParse(nStr string) (int, string) {
	number, message := checkIntStr(nStr)
	if number == 1 {
		message = "This number is neither prime nor composite."
	}
	return number, message
}

func factorize(number int) (bool, [][2]int) {
	j := 1
	factors := map[int]int{}
	// One only needs to search up until the square root of number.
	for j * j < number {
		// After 3, only odd numbers can be prime.
		if j < 3 {
			j++
		} else {
			j += 2
		}
		for {
			// Continue dividing out (and counting) j until j is no longer a factor of number.
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
	// Below is a necessary - but not sufficient - condition.
	isPrime := len(factors) == 1
	// The condition below is required to make it "sufficient".
	if isPrime {
		for _, exponent := range factors {
			if exponent > 1 {
				isPrime = false
				break
			}
		}
	}
	factorsSorted := [][2]int{}
	for prime, exponent := range factors {
		factorsSorted = append(factorsSorted, [2]int{prime, exponent})
	}
	sort.Slice(factorsSorted, func(i, j int) bool {
		return factorsSorted[i][0] < factorsSorted[j][0]
	})
	return isPrime, factorsSorted
}

func modulo(z0, z1 [2]int) (bool, [2]int) {
	var returnIsFactor bool
	returnQuotient := [2]int{0, 0}
	den := z1[0] * z1[0] + z1[1] * z1[1]
	numR:= z0[0] * z1[0] + z0[1] * z1[1]
	numI:= z0[1] * z1[0] - z0[0] * z1[1]
	if numR % den == 0 && numI % den == 0 {
		returnQuotient = [2]int{numR / den, numI / den}
		returnIsFactor = true
	}
	return returnIsFactor, returnQuotient
}

func modulus(z [2]int) int {
	return z[0] * z[0] + z[1] * z[1]
}

func partParse(str, part string) (int, string) {
	if str == "" && part == "imaginary" {
		return 1, ""
	}
	if _, err := strconv.ParseFloat(str, 64); err == nil {
		if integer, err := strconv.Atoi(str); err == nil {
			return integer, ""
		} else {
			return 0, "Your " + part + " part (" + str + ") does not seem to be an integer"
		}
	} else {
		return 0, "Your " + part + " part (" + str + ") does not seem to be a number."
	}
}
