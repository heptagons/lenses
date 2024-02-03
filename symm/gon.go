package symm

type Gon interface {
	Accums() []*Accum
	Id() string
	Angles() []int
	Vectors() []int
	Prime() bool
	Intersecting() bool
}

// gcd returns the greatest common divisor of two integers
func gcd(a, b int) int {
	if b == 0 {
		return a
	}
	return gcd(b, a % b)
}

func gcd2(a, b *int) {
	if g := gcd(*a, *b); g > 1 {
		*a /= g
		*b /= g
	}
}

func gcd3(a, b, c *int) {
	if g := gcd(gcd(*a, *b), *c); g > 1 {
		*a /= g
		*b /= g
		*c /= g
	}
}

func gcd4(a, b, c, d *int) {
	if g := gcd(gcd(gcd(*a, *b), *c), *d); g > 1 {
		*a /= g
		*b /= g
		*c /= g
		*d /= g		
	}
}

func gcd5(a, b, c, d, e *int) {
	if g := gcd(gcd(gcd(gcd(*a, *b), *c), *d), *e); g > 1 {
		*a /= g
		*b /= g
		*c /= g
		*d /= g		
		*e /= g
	}
}

func gcd6(a, b, c, d, e, f *int) {
	if g := gcd(gcd(gcd(gcd(gcd(*a, *b), *c), *d), *e), *f); g > 1 {
		*a /= g
		*b /= g
		*c /= g
		*d /= g		
		*e /= g
		*f /= g
	}
}
