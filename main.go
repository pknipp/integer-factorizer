package main

import (
	"fmt"
	// "io"
	"encode/json"
		"log"
	"net/http"
	"os"
	"strconv"
	// "strings"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

func addFactor(j string, factors map[string]int) {
	if _, ok := factors[j]; ok {
		factors[j]++
	} else {
		factors[j] = 1
	}
}

func factorize(number int) (bool, map[string]int) {
	j := 2
	var factors = make(map[string]int)
	isPrime := true
	for j * j <= number {
		if number % j == 0 {
			addFactor(strconv.Itoa(j), factors)
			number /= j
			isPrime = false
		} else {
			if j == 2 {
				j++
			} else {
				j += 2
			}
		}
	}
	addFactor(strconv.Itoa(number), factors)
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
	// expressionText := "your expression"
	// resultText := "numerical value"
	router.GET("/:number", func(c *gin.Context) {
		numberString := c.Param("number")
		// Eventually, I'll need to error-handle the following.
		number, _ := strconv.Atoi(numberString)
		isPrime, result := factorize(number)
		resultString := ""
		for prime, power := range result {
			resultString += `&nbsp;` + prime
			if power > 1 {
				resultString += `<SUP>` + strconv.Itoa(power) + `</SUP>`
			}
		}
		fmt.Println(resultString)
		c.HTML(http.StatusOK, "result.tmpl.html", gin.H{
				// "expressionText": expressionText,
				// "expressionValue": expression,
				// "resultText": resultText,
				"numberString": numberString,
				"resultString": resultString,
				"isPrime": isPrime,
		})
	})
	router.GET("/json/:number", func(c *gin.Context) {
		numberString := c.Param("number")
		number, _ := strconv.Atoi(numberString)
		isPrime, result := factorize(number)
		factorStr, err := json.Marshal(result)
		resultString := "{\"number\": " + numberString + ", \"isPrime\": " + strconv.FormatBool(isPrime) + ", \"factors\": ", factorStr + "}"
		// resultString := "{\"" + expressionText + "\": " + expression + ", \"" + resultText + "\": " + handler(expression) + "}"
		c.String(http.StatusOK, numberString)
	})
	router.Run(":" + port)
	// Use the following when testing the app in a non-server configuration.
	// number := 1234567890123456789
	// factoredString := factorize(number)
	// fmt.Println(number, factoredString)
}
