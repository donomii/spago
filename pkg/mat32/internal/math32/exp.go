package math32

// Exp returns e**x, the base-e exponential of x.
func Exp(float32) float32

func exp(x float32) float32 {
	const (
		Ln2Hi = float32(6.9313812256e-01)
		Ln2Lo = float32(9.0580006145e-06)
		Log2e = float32(1.4426950216e+00)

		Overflow  = 7.09782712893383973096e+02
		Underflow = -7.45133219101941108420e+02
		NearZero  = 1.0 / (1 << 28) // 2**-28
	)
	// hx := Float32bits(x) & uint32(0x7fffffff)

	// special cases
	switch {
	case IsNaN(x) || IsInf(x, 1):
		return x
	case IsInf(x, -1):
		return 0
	case x > Overflow:
		return Inf(1)
	case x < Underflow:
		return 0
		// case hx > LogMax:
		// 	return Inf(1)
		// case x < 0 && hx > LogMin:
		//  return 0
	case -NearZero < x && x < NearZero:
		return 1 + x
	}

	// reduce; computed as r = hi - lo for extra precision.
	var k int
	switch {
	case x < 0:
		k = int(Log2e*x - 0.5)
	case x > 0:
		k = int(Log2e*x + 0.5)
	}
	hi := x - float32(k)*Ln2Hi
	lo := float32(k) * Ln2Lo

	// compute
	return expmulti(hi, lo, k)
}

// exp1 returns e**r × 2**k where r = hi - lo and |r| ≤ ln(2)/2.
func expmulti(hi, lo float32, k int) float32 {
	const (
		P1 = float32(1.6666667163e-01)  /* 0x3e2aaaab */
		P2 = float32(-2.7777778450e-03) /* 0xbb360b61 */
		P3 = float32(6.6137559770e-05)  /* 0x388ab355 */
		P4 = float32(-1.6533901999e-06) /* 0xb5ddea0e */
		P5 = float32(4.1381369442e-08)  /* 0x3331bb4c */
	)

	r := hi - lo
	t := r * r
	c := r - t*(P1+t*(P2+t*(P3+t*(P4+t*P5))))
	y := 1 - ((lo - (r*c)/(2-c)) - hi)
	// TODO(rsc): make sure Ldexp can handle boundary k
	return Ldexp(y, k)
}
