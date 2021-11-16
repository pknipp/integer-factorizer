package main

import (
	// "fmt"
	// "io"
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

func factorize(number int) map[string]int {
	j := 2
	var factors = make(map[string]int)
	for j * j <= number {
		if number % j == 0 {
			addFactor(strconv.Itoa(j), factors)
			number /= j
		} else {
			if j == 2 {
				j++
			} else {
				j += 2
			}
		}
	}
	addFactor(strconv.Itoa(number), factors)
	return factors
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
		result := factorize(number)
		resultString := ""
		for prime, power := range result {
			resultString += prime
			if power > 1 {
				resultString += "<sup>" + strconv.Itoa(power) + "</sup>"
			}
		}
		c.HTML(http.StatusOK, "result.tmpl.html", gin.H{
				// "expressionText": expressionText,
				// "expressionValue": expression,
				// "resultText": resultText,
				"number": number,
				"resultString": resultString,
		})
	})
	// router.GET("/json/:expression", func(c *gin.Context) {
		// expression := doRegExp(c.Param("expression"))
		// resultString := "{\"" + expressionText + "\": " + expression + ", \"" + resultText + "\": " + handler(expression) + "}"
		// c.String(http.StatusOK, resultString)
	// })
	router.Run(":" + port)
	// Use the following when testing the app in a non-server configuration.
	// number := 1234567890123456789
	// factoredString := factorize(number)
	// fmt.Println(number, factoredString)
}
