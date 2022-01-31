If the number contains a decimal point, return its fraction equivalent (in both mixed and improper forms).
If the number contains a comma (recommended) or space (discouraged) after the decimal point, assume the subsequent digits to repeat.

intPart = part to left of decimal = n0

parse decimal places to left of comma:
    convert string to integer (= n1)
    construct appropriate power of 10: n2 = 10^p0
    simplify(n1, n2), so fraction = n1 / n2

parse decimal places to right of comma
    parse string to integer (= n3)
    construct appropriate denominator: n4 = 10^p1 - 1
    simplify(n3, n4), so fraction = n3 / n4

combine fractions:
    num = n1 * n4 + n2 * n3, den = n2 * n4
    simplify(num, den), so fraction = num / den
    mixed number: n0 + num / den
    improper fraction: (n0 * den + num) / den

helper:
simplify(*n0, *n1) (no return) {
    let g = gcd(n0, n1)
    n0 /= g
    n1 /= g
}
