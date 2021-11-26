package main

import (
	// "fmt"
	"sort"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"regexp"
	"strings"
	"math"
	// "reflect"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

//Euclid's algorithm, which is extremely efficient.
func gcd2(n1, n2 int) int {
	for {
	  t := n2;
	  n1, n2 = t, n1 % n2
	  if n2 == 0 {
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
	result := []gaussFactor{} //map[string][2]int{}
	gStr = regexp.MustCompile(" ").ReplaceAllString(gStr, "")
	if len(gStr) == 0 {
		return result, "Expression is missing."
	}
	// Create array of strings, each representing a gaussian integer
	gsStr := strings.Split(gStr, ",")
	for _, gStr := range gsStr {
		gaussianInt, message := gaussianParse(gStr)
		if len(message) != 0 {
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
	return result, ""
}

func gcd2Complex(gaussa, gaussb map[string][2]int) map[string][2]int {
	gauss := map[string][2]int{}
	if len(gaussa) > len(gaussb) {
		gaussa, gaussb = gaussb, gaussa // better to iterate over a shorter map
	}
	for prime, paira := range gaussa {
		mod2, exponenta := paira[0], paira[1]
		pairb, found := gaussb[prime]
		if found {
			exponentb := pairb[1]
			gauss[prime] = [2]int{mod2, int(math.Min(float64(exponenta), float64(exponentb)))}
		}
	}
	return gauss
}

func gcdParse(nStr string) (int, string) {
	var result int
	var ns []int
	nStr = regexp.MustCompile(" ").ReplaceAllString(nStr, "")
	if len(nStr) == 0 {
		return result, "Expression is missing."
	}
	nsStr := strings.Split(nStr, ",")
	for _, nStr := range nsStr {
		_, err := strconv.ParseFloat(nStr, 64)
		if err != nil {
			return result, "There is a problem with this number: " + nStr
		}
		var n int
		n, err = strconv.Atoi(nStr)
		if err != nil {
			return result, "The following number should be an integer, not a decimal: " + nStr
		} else {
			if n > 0 {
				ns = append(ns, n)
			} else {
				return result, "The following number is not positive: " + nStr
			}
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

func factorizeParse(numberStr string) (int, string) {
	var number int
	if numberStr[0:1] == "-" {
		numberStr = numberStr[1:]
	}
	if numberStr == "1" {
		return number, "This number is neither prime nor composite."
	}
	badNumber := "There is something wrong with your number."
	if len(numberStr) > 18 {
		tooLarge := "Your number is too large."
		if len(numberStr) > 20 {
			return number, tooLarge
		}
		if len(numberStr) == 19 {
			numTrunc, err := strconv.Atoi(numberStr[0:6])
			if err != nil {
				return number, badNumber
			}
			if numTrunc > 922336 {
				return number, tooLarge
			}
		}
	}
	_, err := strconv.ParseFloat(numberStr, 64)
	if err != nil {
		return number, badNumber
	}
	number, err = strconv.Atoi(numberStr)
	if err != nil {
		return number, "Note that the number may not be a decimal."
	}
	return number, ""
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
	_, err := strconv.ParseFloat(str, 64)
	if err == nil {
		integer, err := strconv.Atoi(str)
		if err == nil {
			return integer, ""
		} else {
			return 0, "Your " + part + " part (" + str + ") does not seem to be an integer"
		}
	} else {
		return 0, "Your " + part + " part (" + str + ") does not seem to be a number."
	}
}

func gaussianParse(zStr string) ([2]int, string) {
	z := [2]int{1., 0.}
	noNumber := "You need to input a Gaussian integer."
	neither := "This number is neither prime nor composite."
	zStr = regexp.MustCompile(" ").ReplaceAllString(zStr, "")
	zStr = regexp.MustCompile("j").ReplaceAllString(zStr, "i")
	if zStr[0:1] == "+" {
		zStr = zStr[1:]
	}
	if len(zStr) == 0 {
		return z, noNumber
	}
	if zStr[len(zStr) - 1:] == "i" {
		// Number has an imaginary part
		zStr = zStr[0: len(zStr) - 1]
		if len(zStr) == 0 {
			return z, neither
		} else if zStr == "-" {
			return z, neither
		}
		zSlice := strings.Split(zStr, "+")
		if len(zSlice) == 2 {
			// Number's real part is nonzero and imaginary part is positive.
			int, message := partParse(zSlice[0], "real")
			if len(message) > 0 {
				return z, message
			} else {
				z[0] = int
			}
			int, message = partParse(zSlice[1], "imaginary")
			if len(message) > 0 {
				return z, message
			} else {
				z[1] = int
			}
		} else {
			zSlice = strings.Split(zStr, "-")
			if zStr[0:1] != "-" && len(zSlice) == 2 {
				// Numbers real part is nonzero and imaginary part is negative.
				int, message := partParse(zSlice[0], "real")
				if len(message) > 0 {
					return z, message
				} else {
					z[0] = int
				}
				int, message = partParse(zSlice[1], "imaginary")
				if len(message) > 0 {
					return z, message
				} else {
					z[1] = -int
				}
			} else {
				// Number's real part is zero.
				z[0] = 0
				int, message := partParse(zStr, "imaginary")
				if len(message) > 0 {
					return z, message
				} else {
					if math.Abs(float64(int)) == 1. {
						return z, neither
					}
					z[1] = int
				}
			}
		}
	} else {
		// Number is purely real.
		z[1] = 0
		int, message := partParse(zStr, "real")
		if len(message) > 0 {
			return z, message
		} else {
			if math.Abs(float64(int)) == 1. {
				return z, neither
			}
			z[0] = int
		}
	}
	return z, ""
}

func gaussian(z [2]int) (bool, int, map[string][2]int) {
	gaussianFactors := map[string][2]int{}
	_, factors := factorize(modulus(z))
	for _, pair := range factors {
		prime := pair[0]
		exponent := pair[1]
		// Here are the factors of 1 + i
		if prime == 2 {
			gaussianFactors["1+i"] = [2]int{2, exponent}
			// gaussianFactors = append(gaussianFactors, gaussFactor{"1+i", 2, exponent})
			for count := 0; count < exponent; count++ {
				_, z = modulo(z, [2]int{1, 1})
			}
		} else {
			// Here are the (irreducible) real prime factors, which occur in pairs.
			if prime % 4 == 3 {
				gaussianFactors[strconv.Itoa(prime)] = [2]int{prime, exponent / 2}
				// gaussianFactors = append(gaussianFactors, gaussFactor{strconv.Itoa(prime), prime, exponent / 2})
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
							if even == 1 {
								im = ""
							}
							gaussianFactors[strconv.Itoa(odd) + "+" + im + "i"] = [2]int{prime, count}
							// gaussianFactors = append(gaussianFactors, gaussFactor{strconv.Itoa(odd) + "+" + im + "i", prime, count})
						}
						break
					}
				}
				// For the remaining factors, the real component must be the even one.
				count2 := exponent - count
				if count2 > 0 {
					im := strconv.Itoa(odd)
					if odd == 1 {
						im = ""
					}
					gaussianFactors[strconv.Itoa(even) + "+" + im + "i"] = [2]int{prime, count2}
					// gaussianFactors = append(gaussianFactors, gaussFactor{strconv.Itoa(even) + "+" + im + "i", prime, count2})
				}
				for count = 0; count < count2; count++ {
					_, z = modulo(z, [2]int{2 * n, 2 * m + 1})
				}
			}
		}
	}
	// sort.Slice(gaussianFactors, func(i, j int) bool {
		// return gaussianFactors[i].mod2 < gaussianFactors[j].mod2
	// })
	// Determine exponent of i, based upon what is left after dividing by all Gaussian primes.
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
		// for _, gaussianFactor := range gaussianFactors {
			if pair[1] > 1 {
				isPrime = false
				break
			}
		}
	}
	return isPrime, n, gaussianFactors
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	// I opted not to use this version of router, for technical reasons.
	// router := gin.New()
	router := gin.Default()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})
	// router.GET("/:input", func(c *gin.Context) {
		// inputStr := c.Param("input")
		// if len(strings.Split(inputStr, ",")) > 1  {
			// result, message := gcdParse(inputStr)
			// _, results := factorize(result)
			// factors := [][2]string{}
			// var isPrime bool
			// if result > 1 {
				// isPrime = false
				// for _, pair := range results {
					// prime, exponent := pair[0], pair[1]
					// factors = append(factors, [2]string{strconv.Itoa(prime), strconv.Itoa(exponent)})
				// }
			// } else {
				// isPrime = true
				// factors = append(factors, [2]string{"1", "1"})
			// }
			// c.HTML(http.StatusOK, "result.tmpl.html", gin.H{
				// "input": inputStr,
				// "factors": factors,
				// "message": message,
				// "isPrime": isPrime,
				// "type": "GCD",
				// "title": "Real GCD",
			// })
		// } else {
			// number, message := factorizeParse(inputStr)
			// isPrime, results := factorize(number)
			// Convert from map to array of 2-pls so that 0-th element can be handled separately in results.html.
			// factors := [][2]string{}
			// for _, pair := range results {
				// prime, exponent := pair[0], pair[1]
				// factors = append(factors, [2]string{strconv.Itoa(prime), strconv.Itoa(exponent)})
			// }
			// c.HTML(http.StatusOK, "result.tmpl.html", gin.H{
				// "input": inputStr,
				// "isPrime": isPrime,
				// "factors": factors,
				// "message": message,
				// "type": "integer",
				// "title": "Real factorization",
			// })
		// }
	// })
	// router.GET("/json/:input", func(c *gin.Context) {
		// inputStr := c.Param("input")
		// var resultStr string
		// if len(strings.Split(inputStr, ",")) > 1 {
			// result, message := gcdParse(inputStr)
			// resultStr = "{\"input\": " + inputStr
			// if len(message) > 0 {
				// resultStr += ", \"message\": " + message
			// } else {
				// resultStr += ", \"gcd\": " + strconv.Itoa(result)
			// }
		// } else {
			// number, message := factorizeParse(inputStr)
			// isPrime, result := factorize(number)
			// resultStr = "{\"input\": " + inputStr
			// if len(message) > 0 {
				// resultStr += ", \"message\": " + message
			// } else {
				// resultStr += ", \"isPrime\": " + strconv.FormatBool(isPrime)
				// if !isPrime {
					// factorStr, _ := json.Marshal(result)
					// resultStr += ", \"factors\": " + string(factorStr)
				// }
			// }
		// }
		// c.String(http.StatusOK, resultStr + "}")
	// })
	router.GET("/complex/:input", func(c *gin.Context) {
		inputStr := c.Param("input")
		// if len(strings.Split(inputStr, ",")) > 1 {
			// results, message := gcdComplexParse(inputStr)
			// if len(message) == 0 {
				// factors := [][2]string{}
				// var isPrime bool
				// if len(results) > 0 {
					// isPrime = false
					// for _, pair := range results {
						// prime, exponent := pair[0], pair[1]
						// factors = append(factors, [2]string{"(" + prime + ")", exponent})
					// }
				// } else {
					// isPrime = true
					// factors = append(factors, [2]string{"1", "1"})
				// }
				// c.HTML(http.StatusOK, "result.tmpl.html", gin.H{
					// "input": inputStr,
					// "factors": factors,
					// "message": message,
					// "isPrime": isPrime,
					// "type": "GCD",
					// "title": "Complex GCD",
				// })
			// }
		// } else {
			z, message := gaussianParse(inputStr)
			if len(message) == 0 {
				factors := [][2]string{}
				var isPrime bool
				var number int
				var resultsUnsorted map[string][2]int
				isPrime, number, resultsUnsorted = gaussian(z)
				results := []gaussFactor{} //:= [][2]string{}
				for prime, pair := range resultsUnsorted {
					results = append(results, gaussFactor{prime, pair[0], pair[1]})
				}
				sort.Slice(results, func(i, j int) bool {
					return results[i].mod2 < results[j].mod2
				})
				PREFACTOR := [4]string{"", "i", "-", "-i"}
				// Transform from results (map) to factors (array of 2-ples) to enable me to treat 0-th element differently in results.html.
				firstFactor := true
				for _, singleFactor := range results {
					prime, exponent := singleFactor.prime, strconv.Itoa(singleFactor.exponent)
					factor := ""
					if firstFactor {
						coef := PREFACTOR[number]
						if number % 2 == 0 || strings.Contains(prime, "i") {
							// No multiplication symbol is required, so I just modify the first factor.
							factor += coef
						} else {
							// Multiplication symbol is required, so I prepend one factor.
							factors = append(factors, [2]string{coef, "1"})
						}
					}
					firstFactor = false
					if !strings.Contains(prime, "i") {
						factor += prime
					} else {
						factor += "(" + prime + ")"
					}
					factors = append(factors, [2]string{factor, exponent})
				}
				c.HTML(http.StatusOK, "result.tmpl.html", gin.H{
					"input": inputStr,
					"factors": factors,
					"message": message,
					"isPrime": isPrime,
					"type": "Gaussian",
					"title": "Complex factorization",
				})
			}
		// }
	})
	router.GET("/complex/json/:input", func(c *gin.Context) {
		inputStr := c.Param("input")
		var resultStr string
		// if len(strings.Split(inputStr, ",")) > 1 {
			// result, message := gcdComplexParse(inputStr)
			// resultStr = "{\"input\": " + inputStr
			// if len(message) > 0 {
				// resultStr += ", \"message\": " + message
			// } else {
				// gcdResult, _ := json.Marshal(result)
				// resultStr += ", \"gcd\": " + string(gcdResult)
			// }
		// } else {
			z, message := gaussianParse(inputStr)
			resultStr = "{\"input\": " + inputStr
			if len(message) > 0 {
				resultStr += ", \"message\": " + message
			} else {
				isPrime, n, resultsUnsorted := gaussian(z)

				results := []gaussFactor{}
				for prime, pair := range resultsUnsorted {
					results = append(results, gaussFactor{prime, pair[0], pair[1]})
				}
				sort.Slice(results, func(i, j int) bool {
					return results[i].mod2 < results[j].mod2
				})

				twoFields := [][2]string{}
				for _, result := range results {
					twoFields = append(twoFields, [2]string{result.prime, strconv.Itoa(result.exponent)})
				}
				resultStr += ", \"exponent\": " + strconv.Itoa(n)
				resultStr += ", \"isPrime\": " + strconv.FormatBool(isPrime)
				factorStr, _ := json.Marshal(twoFields)
				resultStr += ", \"factors\": " + string(factorStr)
			}
		// }
		c.String(http.StatusOK, resultStr + "}")
	})
//
	router.Run(":" + port)
	// Use the space below when testing app as CLI./
	// input := "2147483646"
	// fmt.Println(factorizeParse(input))
	// number, _ := factorizeParse(input)
	// fmt.Println(number)
	// fmt.Println(factorize(number))
	// input := "1+3i"
	// z, message := gaussianParse(input)
	// fmt.Println(z, message)
	// fmt.Println(gaussian(z))
}
