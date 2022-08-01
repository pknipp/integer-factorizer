package main

import (
	"strconv"
	"strings"
	"math"
	"regexp"
	"math/big"
)

func simplify(num *big.Int, den *big.Int) {
    fac := gcd2(num, den)
	num.Div(num, fac)
	den.Div(den, fac)
}

func pow10(n int) int {
	if n == 0 { return 1 }
	if n == 1 { return 10 }
	y := pow10(n/2)
	if n % 2 == 0 { return y*y }
	return 10 * y * y
 }

type fraction struct {
	whole int
	num *big.Int
	den *big.Int
	nonrepeating string
	repeating string
}

func decimal(inputStr string) (fraction, string) {
	parts := strings.Split(inputStr, ".")
	var result fraction
	whole, err := strconv.Atoi(parts[0])
	if err != nil {
		return result, "Part (" + parts[0] + ") which is left of decimal cannot be parsed as an integer."
	}
	var num, den *big.Int
    decimalPart := regexp.MustCompile("repeat").ReplaceAllString(parts[1], "r")
	decimalPart = regexp.MustCompile("R").ReplaceAllString(decimalPart, "r")
	decimalPart = regexp.MustCompile(",").ReplaceAllString(decimalPart, "r")
	n_rs := strings.Count(decimalPart, "r")
	var ok bool
	if n_rs > 1 {
		return result, "Part (" + decimalPart + ") which is right of the decimal has more than one character which signals the start of the repeating part of the decimal."
	} else {
		decimalParts := strings.Split(decimalPart, "r")
		if decimalParts[0] == "" {
			num = big.NewInt(0)
			den = big.NewInt(1)
		} else {
			num, ok = new(big.Int).SetString(decimalParts[0], 10)
			if !ok {
				return result, "Terminating part of decimal (" + decimalParts[0] + ") cannot be parsed as an integer."
			}
			den.Exp(big.NewInt(10), big.NewInt(int64(len(decimalParts[0]))), nil)
		}
		if len(decimalParts) > 1 {
			num2, ok := new(big.Int).SetString(decimalParts[1], 10)
			if !ok {
				return result, "Repeating part of decimal (" + decimalParts[1] + ") cannot be parsed as an integer."
			}
			big10 := big.NewInt(10)
			big10.Exp(big10, big.NewInt(int64(len(decimalParts[1]) - 1)), nil)
			den2 := big.NewInt(0).Mul(den, big10)
			num.Add(num.Mul(num, den2), den.Mul(den, num2))
			den.Mul(den, den2)
		}
		simplify(num, den)
		repeating := ""
		if len(decimalParts) > 1 {
			repeating = decimalParts[1]
		}
		return fraction{whole, num, den, parts[0] + "." + decimalParts[0], repeating}, ""
	}
}

func modulus(z [2]int) int {
	return z[0] * z[0] + z[1] * z[1]
}

// determines whether 2nd gaussian integer is a factor of the 1st
func modulo(z0, z1 [2]int) (bool, [2]int) {
	var isFactor bool
	quotient := [2]int{0, 0}
	den := z1[0] * z1[0] + z1[1] * z1[1]
	numR:= z0[0] * z1[0] + z0[1] * z1[1]
	numI:= z0[1] * z1[0] - z0[0] * z1[1]
	if numR % den == 0 && numI % den == 0 {
		quotient = [2]int{numR / den, numI / den}
		isFactor = true
	}
	return isFactor, quotient
}

// ensure that string may be converted to integer, for complex case
func partParse(str, part string) (int, string) {
	if str == "" && part == "imaginary" {
		return 1, ""
	}
	if _, err := strconv.ParseFloat(str, 64); err == nil {
		if integer, err := strconv.Atoi(str); err == nil {
			return integer, ""
		} else {
			return 0, "Your " + part + " part (" + str + ") does not seem to be an integer"
		}
	} else {
		return 0, "Your " + part + " part (" + str + ") does not seem to be a number."
	}
}

func factorizeParse(nStr string) (*big.Int, string) {
	number, message := checkIntStr(nStr)
	if number.Uint64() == 1 {
		message =   "1" + UNITY
	}
	return number, message
}

//Euclid's algorithm, which is extremely efficient.
func gcd2(n1, n2 *big.Int) *big.Int {
	for {
	  t := n2;
	  if n1, n2 = t, n1.Mod(n1, n2); n2 == big.NewInt(0) {
		return n1
	  }
	}
}

func gcd(ns []*big.Int) *big.Int {
	// base case
	if len(ns) == 1 {
		return ns[0]
	}
	// recursive call
	return gcd(append(ns[2:], gcd2(ns[0], ns[1])))
}

func gcdParse(nStr string) (*big.Int, string) {
	var message string
	// slice'll hold all integers whose gcd is needed
	var ns []*big.Int
	var n *big.Int
	nsStr := strings.Split(nStr, ",")
	for _, nStr := range nsStr {
		if n, message = checkIntStr(nStr); len(message) == 0 {
			ns = append(ns, n)
		} else {
			message = ""
			n, _ = new(big.Int).SetString("0", 10)
			return n, ""
		}
	}
	return gcd(ns), ""
}

func gcdComplex(gausss []map[string][2]int) map[string][2]int {
	// base case
	if len(gausss) == 1 {
		return gausss[0]
	}
	// recursive call
	return gcdComplex(append(gausss[2:], gcd2Complex(gausss[0], gausss[1])))
}

func gcd2Complex(gaussa, gaussb map[string][2]int) map[string][2]int {
	// In 2-component array, zeroth element is squared modulus of factor, and last element is power.
	gauss := map[string][2]int{}
	if len(gaussa) > len(gaussb) {
		gaussa, gaussb = gaussb, gaussa // better to iterate over a shorter map
	}
	for prime, paira := range gaussa {
		mod2, exponenta := paira[0], paira[1]
		if pairb, found := gaussb[prime]; found {
			gauss[prime] = [2]int{mod2, int(math.Min(float64(exponenta), float64(pairb[1])))}
		}
	}
	return gauss
}
