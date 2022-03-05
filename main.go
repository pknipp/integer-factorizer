package main

import (
	"fmt"
	"sort"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"regexp"
	"strings"
	"github.com/gin-gonic/gin"
)

const UNITY string = " is neither prime nor composite."
// This allows for easy toggling between cli and web versions of this app.
var isWebVersion bool = true

func main() {
	if isWebVersion {
		port := os.Getenv("PORT")
		if port == "" {
			log.Fatal("$PORT must be set")
		}
		// I opted not to use this version of router, for technical reasons: router := gin.New()
		router := gin.Default()
		router.Use(gin.Logger())
		router.LoadHTMLGlob("templates/*.tmpl.html")
		router.Static("/static", "static")
		router.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.tmpl.html", nil)
		})
		// gcd and factorization of real integers
		router.GET("/:input", func(c *gin.Context) {
			inputStr := c.Param("input")
			if strings.Count(inputStr, ".") > 1 {
				c.HTML(http.StatusOK, "result.tmpl.html", gin.H{
					// "input": inputStr,
					"message": "input (" + inputStr + ") has too many decimal points.",
				})
			// if len(strings.Split(inputStr, ".")) > 1 {
				// inputStr = regexp.MustCompile("repeat").ReplaceAllString(inputStr, "r")
				// inputStr = regexp.MustCompile("R").Copy().ReplaceAllString(inputStr, "r")
			// twoParts := strings.split(input, )
				// fmt.Println(inputStr)
			// }

			// real gcd
			} else if len(strings.Split(inputStr, ",")) > 1  {
				// Reduce white-space, to facilitate parsing.
				inputStr = regexp.MustCompile(" ").ReplaceAllString(inputStr, "")
				// Reinsert a space, so that rendered input is easy to read.
				inputStr = strings.Join(strings.Split(inputStr, ","), ", ")
				// parse input to ensure that it is well-formed
				result, message := gcdParse(inputStr)
				_, results := factorize(result)
				factors := [][2]string{}
				var isPrime bool
				if result > 1 {
					isPrime = false
					for _, pair := range results {
						prime, exponent := pair[0], pair[1]
						// Create slice of 2-component arrays of strings, for use in template.
						factors = append(factors, [2]string{strconv.Itoa(prime), strconv.Itoa(exponent)})
					}
				} else {
					isPrime = true
					factors = append(factors, [2]string{"1", "1"})
				}
				c.HTML(http.StatusOK, "result.tmpl.html", gin.H{
					"input": inputStr,
					"factors": factors,
					"message": message,
					"isPrime": isPrime,
					"type": "GCD",
					"title": "Real GCD",
				})
			} else {
				number, message := factorizeParse(inputStr)
				isPrime, results := factorize(number)
				// Convert from map to slice of 2-component arrays so that 0-th element can be handled separately in results.html.
				factors := [][2]string{}
				for _, pair := range results {
					prime, exponent := pair[0], pair[1]
					// Create slice of 2-component arrays of strings, for use in template.
					factors = append(factors, [2]string{strconv.Itoa(prime), strconv.Itoa(exponent)})
				}
				c.HTML(http.StatusOK, "result.tmpl.html", gin.H{
					"input": inputStr,
					"isPrime": isPrime,
					"factors": factors,
					"message": message,
					"type": "integer",
					"title": "Real factorization",
				})
			}
		})
		router.GET("/json/:input", func(c *gin.Context) {
			inputStr := c.Param("input")
			var resultStr string
			if len(strings.Split(inputStr, ",")) > 1 {
				result, message := gcdParse(inputStr)
				resultStr = "{\"input\": \"" + strings.Join(strings.Split(inputStr, ","), ", ") + "\""
				if len(message) > 0 {
					resultStr += ", \"message\": " + message
				} else {
					areRelativelyPrime := result == 1
					resultStr += ", \"areRelativelyPrime\": " + strconv.FormatBool(areRelativelyPrime)
					if !areRelativelyPrime {
						resultStr += ", \"gcd\": " + strconv.Itoa(result)
					}
				}
			} else {
				number, message := factorizeParse(inputStr)
				isPrime, result := factorize(number)
				resultStr = "{\"input\": \"" + inputStr + "\""
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
			if inputStr := c.Param("input"); len(strings.Split(inputStr, ",")) > 1 {
				results, message := gcdComplexParse(inputStr);
				factors := [][2]string{}
				var isPrime bool
				if len(message) == 0 {
					if len(results) > 0 {
						for _, singleFactor := range results {
							prime, exponent := singleFactor.prime, strconv.Itoa(singleFactor.exponent)
							factors = append(factors, [2]string{"(" + prime + ")", exponent})
						}
					} else {
						isPrime = true
						factors = append(factors, [2]string{"1", "1"})
					}
				}
				c.HTML(http.StatusOK, "result.tmpl.html", gin.H{
					"input": strings.Join(strings.Split(inputStr, ","), ", "),
					"factors": factors,
					"message": message,
					"isPrime": isPrime,
					"type": "GCD",
					"title": "Complex GCD",
				})
			} else {
				z, message := gaussianParse(inputStr)
				factors := [][2]string{}
				var isPrime bool
				if len(message) == 0 {
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
		})
		router.GET("/complex/json/:input", func(c *gin.Context) {
			inputStr := c.Param("input")
			var resultStr string
			if len(strings.Split(inputStr, ",")) > 1 {
				results, message := gcdComplexParse(inputStr)
				twoFields := [][2]string{}
				for _, result := range results {
					twoFields = append(twoFields, [2]string{result.prime, strconv.Itoa(result.exponent)})
				}
				resultStr = "{\"input\": \"" + inputStr + "\""
				if len(message) > 0 {
					resultStr += ", \"message\": " + message
				} else {
					areRelativelyPrime := len(twoFields) == 0
					resultStr += ", \"areRelativelyPrime\": " + strconv.FormatBool(areRelativelyPrime)
					if !areRelativelyPrime {
						gcdResult, _ := json.Marshal(twoFields)
						resultStr += ", \"gcd\": " + string(gcdResult)
					}
				}
			} else {
				resultStr = "{\"input\": \"" + inputStr + "\""
				if z, message := gaussianParse(inputStr); len(message) > 0 {
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
					if !isPrime {
						factorStr, _ := json.Marshal(twoFields)
						resultStr += ", \"factors\": " + string(factorStr)
					}
				}
			}
			c.String(http.StatusOK, resultStr + "}")
		})
		router.Run(":" + port)
	} else {
		// Use this block when testing app as CLI./
		inputStr := "1234"
		number, message := factorizeParse(inputStr)
		isPrime, results := factorize(number)
		// results, message := gcdComplexParse(inputStr)
		// results, message := gcdParse(inputStr)
		// fmt.Println(results, message)
		// _, result := factorize(results)
		fmt.Println(number, message, isPrime, results)
	}
}
