package main

import (
	"strconv"
	"regexp"
)

var TOOLARGE string = "Your number is too large."

func checkIntStr(nStr string) (int, string) {
	// Remove whitespace to facilitate parsing.
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
	badNumber := "There is something wrong with your number"
	// max integer = 2^63 = 9.22...E18
	if len(nStr) > 18 {
		// coarse check of integer size
		if len(nStr) > 19 {
			return number, TOOLARGE
		}
		if len(nStr) == 19 {
			// finer check of integer size
			if numTrunc, err := strconv.Atoi(nStr[0:6]); err != nil {
				return number, badNumber + " (" + nStr + ")."
			} else if numTrunc > 922336 {
				return number, TOOLARGE
			}
		}
	}
	if _, err := strconv.ParseFloat(nStr, 64); err != nil {
		return number, badNumber + " (" + nStr + ")."
	}
	if number, err := strconv.Atoi(nStr); err == nil {
		return number, ""
	} else {
		return number, "Note that the number may not be a decimal."
	}
}
