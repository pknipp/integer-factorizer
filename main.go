package main

import (
	"fmt"
	// "io"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	// "reflect"
	// "strings"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

// type Factor struct {
	// Prime   int
	// Power   int
// }

func factorize(number int) (bool, [][2]int) {
	j := 1
	factors := [][2]int{}
	isPrime := true
	var factor [2]int
	var facFound bool
	// One only needs to search up until the square root of number.
	for j * j < number {
		// After 3, primes skip by at least two.
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
	// router.GET("/:number", func(c *gin.Context) {
		// numberStr := c.Param("number")
		// Eventually, I'll need to error-handle the following.
		// number, _ := strconv.Atoi(numberStr)
		// isPrime, result := factorize(number)
		// resultStr := ""
		// for prime, power := range result {
			// resultStr += `&nbsp;` + prime
			// if power > 1 {
				// resultStr += `<SUP>` + strconv.Itoa(power) + `</SUP>`
			// }
		// }
		// fmt.Println(resultStr)
		// c.HTML(http.StatusOK, "result.tmpl.html", gin.H{
				// "numberString": numberStr,
				// "resultString": resultStr,
				// "isPrime": isPrime,
		// })
	// })
	router.GET("/json/:number", func(c *gin.Context) {
		numberStr := c.Param("number")
		number, _ := strconv.Atoi(numberStr)
		isPrime, result := factorize(number)
		resultStr := "{\"number\": " + numberStr + ", \"isPrime\": " + strconv.FormatBool(isPrime)
		if !isPrime {
			factorStr, _ := json.Marshal(result)
			resultStr += ", \"factors\": " + string(factorStr)
		}
		c.String(http.StatusOK, resultStr + "}")
	})
	router.Run(":" + port)
	// Use the following when testing the app in a non-server configuration.
	number := 1234567890123456789
	bool, factoredString := factorize(number)
	fmt.Println(number, bool, factoredString)
}
