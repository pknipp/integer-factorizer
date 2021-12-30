package main

import (
	"strconv"
	"regexp"
)

var TOOLARGE string = "Your number is too large."

func checkIntStr(nStr string) (int, string) {
	nStr = regexp.MustCompile(" ").ReplaceAllString(nStr, "")
	var number int
	missingNumber := "Number is missing."
	if nStr == "" {
		return number, missingNumber
	} else {
		if nStr[0:1] == "-" {
			nStr = nStr[1:]
		}
	}
	if nStr == "" {
		return number, missingNumber
	}
	badNumber := "There is something wrong with your number."
	if len(nStr) > 18 {
		if len(nStr) > 19 {
			return number, TOOLARGE
		}
		if len(nStr) == 19 {
			if numTrunc, err := strconv.Atoi(nStr[0:6]); err != nil {
				return number, badNumber
			} else if numTrunc > 922336 {
				return number, TOOLARGE
			}
		}
	}
	if _, err := strconv.ParseFloat(nStr, 64); err != nil {
		return number, badNumber
	}
	if number, err := strconv.Atoi(nStr); err != nil {
		return number, "Note that the number may not be a decimal."
	} else {
		return number, ""
	}
}
