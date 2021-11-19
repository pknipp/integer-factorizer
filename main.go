package main

import (
	"fmt"
	// "encoding/json"
	// "log"
	// "net/http"
	// "os"
	"strconv"
	"regexp"
	"strings"
	"math"
	// "reflect"
	// "github.com/gin-gonic/gin"
	// _ "github.com/heroku/x/hmetrics/onload"
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
			// Continue dividing out (and counting) j until it's no longer a factor of number.
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
			if exponent == 1 {
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

func gaussianFactorize(zStr string) (bool, int, [][3]int, string) {
	gaussianFactors := [][3]int{}
	var z [2]int
	noNumber := "You need to input a Gaussian integer."
	neither := "This number is neither prime nor composite."
	zStr = regexp.MustCompile(" ").ReplaceAllString(zStr, "")
	zStr = regexp.MustCompile("j").ReplaceAllString(zStr, "i")
	if zStr[0:1] == "+" {
		zStr = zStr[1:]
	}
	if len(zStr) == 0 {
		return false, 0, gaussianFactors, noNumber
	}
	if zStr[len(zStr) - 1:] == "i" {
		// Number has an imaginary part
		zStr = zStr[0: len(zStr) - 1]
		if len(zStr) == 0 {
			return false, 0, gaussianFactors, neither
		} else if zStr == "-" {
			return false, 0, gaussianFactors, neither
		}
		zSlice := strings.Split(zStr, "+")
		if len(zSlice) == 2 {
			// Number's real part is nonzero and imaginary part is positive.
			int, message := parsePart(zSlice[0], "real")
			if len(message) > 0 {
				return false, 0, gaussianFactors, message
			} else {
				z[0] = int
			}
			int, message = parsePart(zSlice[1], "imaginary")
			if len(message) > 0 {
				return false, 0, gaussianFactors, message
			} else {
				z[1] = int
			}
		} else {
			zSlice = strings.Split(zStr, "-")
			if zStr[0:1] != "-" && len(zSlice) == 2 {
				// Number's real part is nonzero and imaginary part is negative.
				int, message := parsePart(zSlice[0], "real")
				if len(message) > 0 {
					return false, 0, gaussianFactors, message
				} else {
					z[0] = int
				}
				int, message = parsePart(zSlice[1], "imaginary")
				if len(message) > 0 {
					return false, 0, gaussianFactors, message
				} else {
					z[1] = -int
				}
			} else {
				// Number's real part is zero.
				z[0] = 0
				int, message := parsePart(zStr, "imaginary")
				if len(message) > 0 {
					return false, 0, gaussianFactors, message
				} else {
					if math.Abs(float64(int)) == 1. {
						return false, 0, gaussianFactors, neither
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
			return false, 0, gaussianFactors, message
		} else {
			if math.Abs(float64(int)) == 1. {
				return false, 0, gaussianFactors, neither
			}
			z[0] = int
		}
	}
	_, factors := factorize(modulus(z))
	for prime, exponent := range factors {
		// Here are the factors of 1 + i
		if prime == 2 {
			gaussianFactors = append(gaussianFactors, [3]int{1, 1, exponent})
			for count := 0; count < exponent; count++ {
				_, z = modulo(z, [2]int{1, 1})
			}
		} else {
			mod4 := prime
			// Here are the (irreducible) real prime factors, which occur in pairs.
			if mod4 % 4 == 3 {
				gaussianFactors = append(gaussianFactors, [3]int{prime, 0, exponent / 2})
				for count := 0; count < exponent / 2; count++ {
					for i, _ := range z {
						z[i] /= prime
					}
				}
			} else {
				// Here are Gaussian integers for which one component is odd and the other is even.
				// Find ints m, n such that (2m+1)^2 + (2n)^2 = mod4
				mod4 = (mod4 - 1) / 4
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
				count := 0
				// First, let's consider possibility that the real component is the odd one.
				for {
					isFactor, quotient := modulo(z, [2]int{2 * m + 1, 2 * n})
					if isFactor {
						z = quotient
						count++
					} else {
						if count > 0 {
							gaussianFactors = append(gaussianFactors, [3]int{2 * m + 1, 2 * n, count})
						}
						break
					}
				}
				// For the remaining factors, the real component must be the even one.
				count2 := exponent - count
				if count2 > 0 {
					gaussianFactors = append(gaussianFactors, [3]int{2 * n, 2 * m + 1, count2})
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
	isPrime := len(gaussianFactors) == 1 && gaussianFactors[0][2] == 1
	// if isPrime {
		// for _, exponent := range factors {
			// if exponent == 1 {
				// isPrime = false
				// break
			// }
		// }
	// }
	return isPrime, n, gaussianFactors, ""
}

func main() {
	// port := os.Getenv("PORT")
	// if port == "" {
		// log.Fatal("$PORT must be set")
	// }
	// I opted not to use this version of router, for technical reasons.
	// router := gin.New()
	// router := gin.Default()
	// router.Use(gin.Logger())
	// router.LoadHTMLGlob("templates/*.tmpl.html")
	// router.Static("/static", "static")
	// router.GET("/", func(c *gin.Context) {
		// c.HTML(http.StatusOK, "index.tmpl.html", nil)
	// })
	// router.GET("/:input", func(c *gin.Context) {
		// inputStr := c.Param("input")
		// if inputStr[0:1] == "[" {
			// result, message := gcdParse(inputStr[1:])
			// c.HTML(http.StatusOK, "result.tmpl.html", gin.H{
				// "numbers": inputStr,
				// "result": result,
				// "message": message,
				// "type": "GCD",
				// "title": "GCD",
			// })
		// } else {
			// number, message := factorizeParse(inputStr)
			// results := factorize(number)
			// factors := [][2]string{}
			// for _, result := range results {
				// factors = append(factors, [2]string{strconv.Itoa(result[0]), strconv.Itoa(result[1])})
			// }
			// c.HTML(http.StatusOK, "result.tmpl.html", gin.H{
					// "number": inputStr,
					// "isPrime": len(factors) == 1,
					// "factors": factors,
					// "message": message,
					// "type": "integer",
					// "title": "prime factorization",
			// })
		// }
	// })
	// router.GET("/json/:input", func(c *gin.Context) {
		// inputStr := c.Param("input")
		// var resultStr string
		// if inputStr[0:1] == "[" {
			// result, message := gcdParse(inputStr[1:])
			// resultStr = "{\"numbers\": " + inputStr
			// if len(message) > 0 {
				// resultStr += ", \"message\": " + message
			// } else {
				// resultStr += ", \"gcd\": " + strconv.Itoa(result)
			// }
		// } else {
			// number, message := factorizeParse(inputStr)
			// result := factorize(number)
			// resultStr = "{\"number\": " + inputStr
			// if len(message) > 0 {
				// resultStr += ", \"message\": " + message
			// } else {
				// resultStr += ", \"isPrime\": " + strconv.FormatBool(len(factors) == 1)
				// if len(factors) != 1 {
					// factorStr, _ := json.Marshal(result)
					// resultStr += ", \"factors\": " + string(factorStr)
				// }
			// }
		// }
		// c.String(http.StatusOK, resultStr + "}")
	// })
//
	// router.GET("/complex/:input", func(c *gin.Context) {
		// inputStr := c.Param("input")
		// number, results, message := gaussianFactorize(inputStr)
		// PREFACTOR := [4]string{"", "i", "-", "-i"}
		// factors := [][2]string{}
		// for i, result := range results {
			// exponent := strconv.Itoa(result[2])
			// factor := ""
			// if i == 0 {
				// coef := PREFACTOR[number]
				// if number % 2 == 0 || result[1] != 0 {
					// factor += coef
				// } else {
					// factors = append(factors, [2]string{coef, "1"})
				// }
			// }
			// if result[1] == 0 {
				// factor += strconv.Itoa(result[0])
			// } else {
				// factor += "(" + strconv.Itoa(result[0])
				// if result[1] > 0 {
					// factor += " + "
				// } else {
					// factor += " - "
				// }
				// imCoef := math.Abs(float64(result[1]))
				// if imCoef != 1. {
					// factor += strconv.Itoa(int(imCoef))
				// }
				// factor += "i)"
			// }
			// factors = append(factors, [2]string{factor, exponent})
		// }
		// isPrime := len(results) == 1
		// c.HTML(http.StatusOK, "result.tmpl.html", gin.H{
				// "number": inputStr,
				// "factors": factors,
				// "message": message,
				// "isPrime": isPrime,
				// "type": "Gaussian",
				// "title": "Gaussian-prime factorization",
		// })
	// })
//
	// router.GET("/complex/json/:input", func(c *gin.Context) {
		// inputStr := c.Param("input")
		// var resultStr string
		// n, result, message := gaussianFactorize(inputStr)
		// resultStr = "{\"number\": " + inputStr
		// if len(message) > 0 {
			// resultStr += ", \"message\": " + message
		// } else {
			// resultStr += ", \"exponent\": " + strconv.Itoa(n)
			// factorStr, _ := json.Marshal(result)
			// resultStr += ", \"factors\": " + string(factorStr)
		// }
		// c.String(http.StatusOK, resultStr + "}")
	// })
//
	// router.Run(":" + port)
	number := 18
	fmt.Println(number)
	fmt.Println(factorize(number))
	input := "3+4i"
	fmt.Println(input)
	fmt.Println(gaussianFactorize(input))
	// Use the space below when testing the app as a CLI.
	// zStr := "-3"
	// fmt.Println(gaussianFactorize(zStr))
}
