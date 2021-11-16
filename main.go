package main

import (
	// "fmt"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	// "reflect"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

func factorize(numberStr string) (bool, [][2]int, string) {
	j := 1
	factors := [][2]int{}
	isPrime := true
	var message string
	if numberStr[0:1] == "-" {
		numberStr = numberStr[1:]
	}
	if len(numberStr) > 19 {
		tooLarge := "Your number is too large."
		if len(numberStr) > 20 {
			return isPrime, factors, tooLarge
		}
		if len(numberStr) == 19 {
			numTrunc, err := strconv(numberStr[0:6])
			if err != nil {
				return isPrime, factors, "There is something wrong with your number."
			}
			if numTrunc > 922336 {
				return isPrime, factors, tooLarge
			}
		}
	}
	_, err := strconv.ParseFloat(numberStr, 64)
	if err != nil {
		return isPrime, factors, "There is something wrong with the number that you input."
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
	router.GET("/:number", func(c *gin.Context) {
		numberStr := c.Param("number")
		isPrime, factors, message := factorize(numberStr)
		c.HTML(http.StatusOK, "result.tmpl.html", gin.H{
				"number": numberStr,
				"isPrime": isPrime,
				"factors": factors,
				"message": message,
		})
	})
	router.GET("/json/:number", func(c *gin.Context) {
		numberStr := c.Param("number")
		// number, _ := strconv.Atoi(numberStr)
		isPrime, result, message := factorize(numberStr)
		resultStr := "{\"number\": " + numberStr
		if len(message) > 0 {
			resultStr += ", \"message\": " + message
		} else {
			resultStr += ", \"isPrime\": " + strconv.FormatBool(isPrime)
			if !isPrime {
				factorStr, _ := json.Marshal(result)
				resultStr += ", \"factors\": " + string(factorStr)
			}
		}
		c.String(http.StatusOK, resultStr + "}")
	})
	router.Run(":" + port)
	// Use the following when testing the app in a non-server configuration.
	// number := 1234567890123456789
	// bool, factors, message := factorize(number)
	// fmt.Println(number, bool, factors, message)
}
