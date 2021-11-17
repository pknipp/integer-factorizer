package main

import (
	"fmt"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"regexp"
	"strings"
	// "reflect"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
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

func factorize(numberStr string) (bool, [][2]int, string) {
	j := 1
	factors := [][2]int{}
	isPrime := true
	var message string
	if numberStr[0:1] == "-" {
		numberStr = numberStr[1:]
	}
	if numberStr == "1" {
		return isPrime, factors, "This number is neither prime nor composite."
	}
	badNumber := "There is something wrong with your number."
	if len(numberStr) > 18 {
		tooLarge := "Your number is too large."
		if len(numberStr) > 20 {
			return isPrime, factors, tooLarge
		}
		if len(numberStr) == 19 {
			numTrunc, err := strconv.Atoi(numberStr[0:6])
			fmt.Println(numTrunc)
			if err != nil {
				return isPrime, factors, badNumber
			}
			if numTrunc > 922336 {
				return isPrime, factors, tooLarge
			}
		}
	}
	_, err := strconv.ParseFloat(numberStr, 64)
	if err != nil {
		return isPrime, factors, badNumber
	}
	var number int
	number, err = strconv.Atoi(numberStr)
	if err != nil {
		return isPrime, factors, "Note that the number may not be a decimal."
	}
	var factor [2]int
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
	return isPrime, factors, message
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
				"isGCD": true,
				"title": "GCD",
			})
		} else {
			isPrime, factors, message := factorize(inputStr)
			c.HTML(http.StatusOK, "result.tmpl.html", gin.H{
					"number": inputStr,
					"isPrime": isPrime,
					"factors": factors,
					"message": message,
					"isGCD": false,
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
			isPrime, result, message := factorize(inputStr)
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
	router.Run(":" + port)
	// Use the following when testing the app as a CLI.
	// input := "[16, 18]"
	// var isPrime bool
	// var result [][2]int
	// var result int
	// var message string
	// if input[0:1] == "[" {
		// result, message = gcdParse(input[1:])
	// } else {
		// isPrime, result, message = factorize(input)
	// }
	// fmt.Println(input, isPrime, result, message)
	// xs := []int{48, 52, 54}
	// fmt.Println(gcd(xs))
}
