package main

import (
	"math/big"
)

// type factor struct {
	// prime *big.Int
	// exponent int
// }

type modExp struct {
	modulus *big.Int
	exponent int
}

// returns: isPrime, squared modulus, factor map (each w/value equalling squared modulus and power)
func gaussian(z [2]*big.Int) (bool, int, map[string]modExp) {
	BIG := big.NewInt(0)
	gaussianFactors := map[string]modExp{}
	// Factoring a gaussian is facilitated by finding the (real) factors of its (squared) modulus.
	_, factors := factorize(modulus(z))
	for _, pair := range factors {
		prime := pair.prime
		exponent := pair.exponent
		// Here are the factors of 1 + i
		if prime.Cmp(big.NewInt(2)) == 0 {
			gaussianFactors["1+i"] = modExp{big.NewInt(2), exponent}
			for count := 0; count < exponent; count++ {
				_, z = modulo(z, [2]*big.Int{big.NewInt(1), big.NewInt(1)})
			}
		} else {
			// Here are the (irreducible) real prime factors, which occur in pairs.
			if BIG.Mod(prime, big.NewInt(4)).Cmp(big.NewInt(3)) == 0 {
				gaussianFactors[prime.Text(10)] = modExp{prime, exponent / 2}
				for count := 0; count < exponent / 2; count++ {
					for i := range z {
						z[i].Div(z[i], prime)
					}
				}
			} else {
				// Here are Gaussian integers for which one component is odd and the other is even.
				// Find ints m, n such that (2m+1)^2 + (2n)^2 = prime
				mod4 := BIG.Div(BIG.Add(prime, BIG.Neg(big.NewInt(1))), big.NewInt(4))
				// Now this becomes m*(m+1) + n^2 = mod4, which is solved via a while loop.
				m := big.NewInt(0)
				var n *big.Int
				for {
					nm := BIG.Sqrt(BIG.Add(mod4, BIG.Neg(BIG.Mul(m, (BIG.Add(m, big.NewInt(1)))))))
					np := BIG.Add(nm, big.NewInt(1))
					if BIG.Add(BIG.Mul(m, BIG.Add(m, big.NewInt(1))), BIG.Mul(nm, nm)) == mod4 {
						n = nm
						break
					} else if BIG.Add(BIG.Mul(m, BIG.Add(m, big.NewInt(1))), BIG.Mul(np, np)) == mod4 {
						n = np
						break
					}
					m.Add(m, big.NewInt(1))
				}
				odd := BIG.Add(BIG.Mul(big.NewInt(2), m), big.NewInt(1))
				even := BIG.Mul(big.NewInt(2), n)
				// First, let's consider possibility that the real component is the odd one.
				count := 0
				for {
					isFactor, quotient := modulo(z, [2]*big.Int{odd, even})
					if isFactor {
						z = quotient
						count++
					} else {
						if count > 0 {
							im := even.Text(10)
							gaussianFactors[odd.Text(10) + "+" + im + "i"] = modExp{prime, count}
						}
						break
					}
				}
				// For the remaining factors, the real component must be the even one.
				count2 := exponent - count
				if count2 > 0 {
					im := odd.Text(10)
					if im == "1" {
						im = ""
					}
					gaussianFactors[even.Text(10) + "+" + im + "i"] = modExp{prime, count2}
				}
				for count = 0; count < count2; count++ {
					_, z = modulo(z, [2]*big.Int{BIG.Mul(big.NewInt(2), n), BIG.Add(BIG.Mul(big.NewInt(2), m), big.NewInt(1))})
				}
			}
		}
	}
	// The following logic is a bit obtuse, but it determines exponent of i, based upon what is left after dividing by all Gaussian primes.
	var n *big.Int
	if BIG.Abs(z[0]) == big.NewInt(1) {
		n = BIG.Add(big.NewInt(1), BIG.Neg(z[0]))
	} else {
		n = BIG.Add(big.NewInt(2), BIG.Neg(z[1]))
	}
	// Below is a necessary - but not sufficient - condition.
	isPrime := len(gaussianFactors) == 1
	// The next condition is required to make it "sufficient"
	// An example of this would be 9 (= 3^2).
	if isPrime {
		for _, pair := range gaussianFactors {
			if pair.exponent > 1 {
				isPrime = false
				break
			}
		}
	}
	return isPrime, int(n.Int64()), gaussianFactors
}
