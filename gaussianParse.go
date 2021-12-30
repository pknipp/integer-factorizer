package main

import (
	"regexp"
	"strings"
	"math"
)

var MAXINT float64 = 0.999999 * math.Pow(2., 63.)

func gaussianParse(zStr string) ([2]int, string) {
	z := [2]int{1., 0.}
	noNumber := "You need to input a Gaussian integer."
	neither := "This number is neither prime nor composite."
	zStr = regexp.MustCompile(" ").ReplaceAllString(zStr, "")
	zStr = regexp.MustCompile("j").ReplaceAllString(zStr, "i")
	if zStr[0:1] == "+" {
		zStr = zStr[1:]
	}
	if len(zStr) == 0 {
		return z, noNumber
	}

	if zStr[len(zStr) - 1:] == "i" {
		// Number has an imaginary part
		zStr = zStr[0: len(zStr) - 1]
		if len(zStr) == 0 {
			return z, neither
		} else if zStr == "-" {
			return z, neither
		}
		zSlice := strings.Split(zStr, "+")
		if len(zSlice) == 2 {
			// Number's real part is nonzero and imaginary part is positive.
			int, message := partParse(zSlice[0], "real")
			if len(message) > 0 {
				return z, message
			} else {
				z[0] = int
			}
			int, message = partParse(zSlice[1], "imaginary")
			if len(message) > 0 {
				return z, message
			} else {
				z[1] = int
			}
		} else {
			zSlice = strings.Split(zStr, "-")
			if zStr[0:1] != "-" && len(zSlice) == 2 {
				// Numbers real part is nonzero and imaginary part is negative.
				int, message := partParse(zSlice[0], "real")
				if len(message) > 0 {
					return z, message
				} else {
					z[0] = int
				}
				int, message = partParse(zSlice[1], "imaginary")
				if len(message) > 0 {
					return z, message
				} else {
					z[1] = -int
				}
			} else {
				// Number's real part is zero.
				z[0] = 0
				int, message := partParse(zStr, "imaginary")
				if len(message) > 0 {
					return z, message
				} else {
					if math.Abs(float64(int)) == 1. {
						return z, neither
					}
					z[1] = int
				}
			}
		}
	} else {
		// Number is purely real.
		z[1] = 0
		int, message := partParse(zStr, "real")
		if len(message) > 0 {
			return z, message
		} else {
			if math.Abs(float64(int)) == 1. {
				return z, neither
			}
			z[0] = int
		}
	}
	x := float64(z[0])
	if x > math.Sqrt(MAXINT) || float64(z[1]) > math.Sqrt(MAXINT - x * x) {
		return [2]int{0, 0}, TOOLARGE
	}
	return z, ""
}
