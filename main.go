package main

import (
	"fmt"
	// "io"
	"encoding/json"
		"log"
	"net/http"
	"os"
	"strconv"
	"reflect"
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
		numberStr := c.Param("number")
		// Eventually, I'll need to error-handle the following.
		number, _ := strconv.Atoi(numberStr)
		isPrime, result := factorize(number)
		resultStr := ""
		for prime, power := range result {
			resultStr += `&nbsp;` + prime
			if power > 1 {
				resultStr += `<SUP>` + strconv.Itoa(power) + `</SUP>`
			}
		}
		fmt.Println(resultStr)
		c.HTML(http.StatusOK, "result.tmpl.html", gin.H{
				// "expressionText": expressionText,
				// "expressionValue": expression,
				// "resultText": resultText,
				"numberString": numberStr,
				"resultString": resultStr,
				"isPrime": isPrime,
		})
	})
	router.GET("/json/:number", func(c *gin.Context) {
		numberSt := c.Param("number")
		number, _ := strconv.Atoi(numberStr)
		isPrime, result := factorize(number)
		factorStr, _ := json.Marshal(result)
		isPrimeStr := strconv.FormatBool(isPrime)
		fmt.Println(reflect.TypeOf(numberStr), reflect.TypeOf(isPrimeStr), reflect.TypeOf(string(factorStr)))
		resultStr := "string" //"{\"number\": " + numberString + ", \"isPrime\": " + strconv.FormatBool(isPrime) + ", \"factors\": ", string(factorStr) + "}"
		// resultString := "{\"" + expressionText + "\": " + expression + ", \"" + resultText + "\": " + handler(expression) + "}"
		c.String(http.StatusOK, resultStr)
	})
	router.Run(":" + port)
	// Use the following when testing the app in a non-server configuration.
	// number := 1234567890123456789
	// factoredString := factorize(number)
	// fmt.Println(number, factoredString)
}
