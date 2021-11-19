package main

import (
	// "fmt"
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

func gcd2(n1, n2 int) int {
	for {
	  t := n2;
	  n1, n2 = t, n1 % n2
	  if n2 == 0 {
		return n1
	  }
	}
}

func gcdParse(nStr string) (int, string) {
	var result int
	var ns []int
	if len(nStr) == 0 {
		return result, "Array is missing."
	}
	if nStr[len(nStr) - 1 :] != "]" {
		return result, "Array should end with a close bracket."
	} else {
		nStr = nStr[0: len(nStr) - 1]
		nStr = regexp.MustCompile(" ").ReplaceAllString(nStr, "")
		if len(nStr) == 0 {
			return result, "Your brackets don't contain any numbers."
		} else {
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

func factorize(number int) (bool, map[int]int) {
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
	return isPrime, factors
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

func parsePart(str, part string) (int, string) {
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
	z := [2]int{0., 0.}
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
			int, message := parsePart(zSlice[0], "real")
			if len(message) > 0 {
				return z, message
			} else {
				z[0] = int
			}
			int, message = parsePart(zSlice[1], "imaginary")
			if len(message) > 0 {
				return z, message
			} else {
				z[1] = int
			}
		} else {
			zSlice = strings.Split(zStr, "-")
			if zStr[0:1] != "-" && len(zSlice) == 2 {
				// Numbers real part is nonzero and imaginary part is negative.
				int, message := parsePart(zSlice[0], "real")
				if len(message) > 0 {
					return z, message
				} else {
					z[0] = int
				}
				int, message = parsePart(zSlice[1], "imaginary")
				if len(message) > 0 {
					return z, message
				} else {
					z[1] = -int
				}
			} else {
				// Number's real part is zero.
				z[0] = 0
				int, message := parsePart(zStr, "imaginary")
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
		int, message := parsePart(zStr, "real")
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

func gaussian(z [2]int) (bool, int, map[string]int) {
	gaussianFactors := map[string]int{}
	_, factors := factorize(modulus(z))
	for prime, exponent := range factors {
		// Here are the factors of 1 + i
		if prime == 2 {
			gaussianFactors["1+i"] = exponent
			for count := 0; count < exponent; count++ {
				_, z = modulo(z, [2]int{1, 1})
			}
		} else {
			// Here are the (irreducible) real prime factors, which occur in pairs.
			if prime % 4 == 3 {
				gaussianFactors[strconv.Itoa(prime)] = exponent / 2
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
							gaussianFactors[strconv.Itoa(odd) + "+" + im + "i"] = count
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
					gaussianFactors[strconv.Itoa(even) + "+" + im + "i"] = count2
				}
				for count = 0; count < count2; count++ {
					_, z = modulo(z, [2]int{2 * n, 2 * m + 1})
				}
			}
		}
	}
	// Determine exponent of i.
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
		for _, exponent := range gaussianFactors {
			if exponent > 1 {
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
	router.GET("/:input", func(c *gin.Context) {
		inputStr := c.Param("input")
		if inputStr[0:1] == "[" {
			result, message := gcdParse(inputStr[1:])
			c.HTML(http.StatusOK, "result.tmpl.html", gin.H{
				"numbers": inputStr,
				"result": result,
				"message": message,
				"type": "GCD",
				"title": "GCD",
			})
		} else {
			number, message := factorizeParse(inputStr)
			isPrime, results := factorize(number)
			// Convert from map to array of 2-pls so that 0-th element can be handled separately in results.html.
			factors := [][2]string{}
			for prime, exponent := range results {
				factors = append(factors, [2]string{strconv.Itoa(prime), strconv.Itoa(exponent)})
			}
			c.HTML(http.StatusOK, "result.tmpl.html", gin.H{
					"number": inputStr,
					"isPrime": isPrime,
					"factors": factors,
					"message": message,
					"type": "integer",
					"title": "prime factorization",
			})
		}
	})
	router.GET("/json/:input", func(c *gin.Context) {
		inputStr := c.Param("input")
		var resultStr string
		if inputStr[0:1] == "[" {
			result, message := gcdParse(inputStr[1:])
			resultStr = "{\"numbers\": " + inputStr
			if len(message) > 0 {
				resultStr += ", \"message\": " + message
			} else {
				resultStr += ", \"gcd\": " + strconv.Itoa(result)
			}
		} else {
			number, message := factorizeParse(inputStr)
			isPrime, result := factorize(number)
			resultStr = "{\"number\": " + inputStr
			if len(message) > 0 {
				resultStr += ", \"message\": " + message
			} else {
				resultStr += ", \"isPrime\": " + strconv.FormatBool(isPrime)
				if !isPrime {
					factorStr, _ := json.Marshal(result)
					resultStr += ", \"factors\": " + string(factorStr)
				}
			}
		}
		c.String(http.StatusOK, resultStr + "}")
	})

	router.GET("/complex/:input", func(c *gin.Context) {
		inputStr := c.Param("input")
		z, message := gaussianParse(inputStr)
		factors := [][2]string{}
		var isPrime bool
		var number int
		results := map[string]int{}
		if len(message) == 0 {
			isPrime, number, results = gaussian(z)
			PREFACTOR := [4]string{"", "i", "-", "-i"}
			// Transform from results (map) to factors (array of 2-ples) to enable me to treat 0-th element differently in results.html.
			firstFactor := true
			for prime, exponent := range results {
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
				factors = append(factors, [2]string{factor, strconv.Itoa(exponent)})
			}
		}
		c.HTML(http.StatusOK, "result.tmpl.html", gin.H{
				"number": inputStr,
				"factors": factors,
				"message": message,
				"isPrime": isPrime,
				"type": "Gaussian",
				"title": "Gaussian-prime factorization",
		})
	})

	router.GET("/complex/json/:input", func(c *gin.Context) {
		inputStr := c.Param("input")
		z, message := gaussianParse(inputStr)
		resultStr := "{\"number\": " + inputStr
		if len(message) > 0 {
			resultStr += ", \"message\": " + message
		} else {
			isPrime, n, result := gaussian(z)
			resultStr += ", \"exponent\": " + strconv.Itoa(n)
			resultStr += ", \"isPrime\": " + strconv.FormatBool(isPrime)
			factorStr, _ := json.Marshal(result)
			resultStr += ", \"factors\": " + string(factorStr)
		}
		c.String(http.StatusOK, resultStr + "}")
	})

	router.Run(":" + port)
	// Use the space below when testing app as CLI./
	// input := "1-5i"
	// fmt.Println(input)
	// z, message := gaussianParse(input)
	// fmt.Println(z, message)
	// if len(message) == 0 {
		// fmt.Println(gaussian(z))
	// }
}
