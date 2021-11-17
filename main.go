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
	// "math"
	// "reflect"
	// "github.com/gin-gonic/gin"
	// _ "github.com/heroku/x/hmetrics/onload"
)

func gcd2(n1, n2 int) int {
	for {
	  t := n2;
	//   n2 = n1 % n2;
	//   n1 = t;
	  n1, n2 = t, n1 % n2
	  if n2 == 0 {
		//   break
		return n1
	  }
	}
	// return n1;
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
	// j := 1
	// factors := [][2]int{}
	// isPrime := true
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

func factorize(number int) (bool,[][2]int) {
	j := 1
	var isPrime bool
	var factor [2]int
	var factors [][2]int
	var facFound bool
	// One only needs to search up until the square root of number.
	for j * j < number {
		// After 3, only odd numbers can be prime.
		if j < 3 {
			j++
		} else {
			j += 2
		}
		facFound = false
		for {
			if number % j == 0 {
				if !facFound {
					factor[0], factor[1] = j, 1
					facFound = true
					if number > 2 {
						isPrime = false
					}
				} else {
					factor[1]++
				}
				number /= j
			} else {
				if facFound {
					factors = append(factors, factor)
					facFound = false
				}
				break
			}
		}
	}
	// The last factor is needed if the largest factor occurs by itself.
	if !facFound && number != 1 {
		factors = append(factors, [2]int{number, 1})
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

// func gaussFactorize(z [2]int) [][3]int {
	// iMin := 1
	// k := [2]int{iMin, 0}
	// sizeK := modulus(k)
	// sizeZ := modulus(z)
	// factors := [][3]int{}
	// isPrime := true
	// var factor [3]int
	// var facFound bool
	// var i, j int
	// One only needs to search up until the square root of number.
		// for sizeK * sizeK <= sizeZ {
		// i = int(math.Sqrt(float64(sizeK)))
		// j = 0
		// for {
			// for {
				// if i * i + j * j <= sizeK {
					// j++
				// } else {
					// break
				// }
			// }
			// sizeK = i * i + j * j
			// facFound = false
			// for {
				// isFactor, quotient := modulo(z, [2]int{i, j})
				// if isFactor {
					// if !facFound {
						// factor[0], factor[1], factor[2] = i, j, 1
						// facFound = true
					// } else {
						// factor[2]++
					// }
					// z = quotient
				// } else {
					// if facFound {
						// factors = append(factors, factor)
						// facFound = false
					// }
					// break
				// }
			// }
		// }
	// }
	// The last factor is needed if the largest factor occurs by itself.
	// if !facFound {
		// factors = append(factors, [3]int{i, j, 1})
	// }
	// return factors
// }

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
				// "isGCD": true,
				// "title": "GCD",
			// })
		// } else {
			// isPrime, factors, message := factorize(inputStr)
			// c.HTML(http.StatusOK, "result.tmpl.html", gin.H{
					// "number": inputStr,
					// "isPrime": isPrime,
					// "factors": factors,
					// "message": message,
					// "isGCD": false,
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
			// isPrime, result, message := factorize(inputStr)
			// resultStr = "{\"number\": " + inputStr
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
	// router.Run(":" + port)
	// Use the following when testing the app as a CLI.
	// fmt.Println(modulo([2]int{5, 5}, [2]int{2, 1}))
	// fmt.Println(gaussFactorize([2]int{1, 3}))
	// input := "[16, 18]"
	// var isPrime bool
	// var result [][2]int
	input := "1234567890123456789"
	number, message := factorizeParse(input)
	fmt.Println(number, message)
	fmt.Println(factorize(number))
}
