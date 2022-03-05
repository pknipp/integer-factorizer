(helper) simplify ((*num, *den) (no return) {
    fac = gcd(n0, n1)
    n0 /= g
    n1 /= g
}

split input on decimal point
if len(slice) = 2
    n0 := parse(0th element)
    d := len(1st element)
    declare num, den
    if 1st element doesn't contain an r:
        num = parse(1st element)
        den = 10 ** d
    else:
        arr := 1st element split on r
        num0 := parse(arr[0])
        den0 := 10 ** len(arr[0])
        num1 := parse(arr[1])
        den1 := (10 ** len(arr[1]) - 1) * den0
        num = num0 * den1 + num1 * den0
        den = den0 * den1
    simplify(&num, &den)
    return n0, num, den, properly formatted


If the number contains a decimal point, return its fraction equivalent (in both mixed and improper forms).
If the number contains a "repeat", "R". "r", comma (recommended) or space (discouraged) after the decimal point, assume the subsequent digits to repeat.
