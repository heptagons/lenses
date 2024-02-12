package symm

import (
	"fmt"
)

type Accum struct {
	x []int
	y []int
}

func NewAccum(t int) *Accum {
	return &Accum{
		x: make([]int, t),
		y: make([]int, t),
	}
}

func (a *Accum) Clone() *Accum {
	n := len(a.x)
	x := make([]int, n)
	y := make([]int, n)
	for i := 0; i < n; i++ {
		x[i] = a.x[i]
		y[i] = a.y[i]
	}
	return &Accum{
		x: x,
		y: y,
	}
}

func (a *Accum) X() []int {
	return a.x
}

func (a *Accum) Y() []int {
	return a.y
}

func (a *Accum) String() string {
	return fmt.Sprintf("xy=%v%v", a.x, a.y)
}

// originPrime checks if this accumulator represents the origin.
// All x,y values should be exactly zero
func (a *Accum) originPrime() bool {
	for i := 0; i < len(a.x); i++ {
		if a.x[i] != 0 {
			return false
		}
		if a.y[i] != 0 {
			return false
		}
	}
	return true
}

// origin9 checks if this accumulator for symmetry 3x3 represents the origin
func (a *Accum) origin9() bool {
	return accumOrigin9x(a.x) && accumOrigin9y(a.y)
}

// origin9 checks if this accumulator for symmetry 3x5 represents the origin
func (a *Accum) origin15() bool {
	return accumOrigin15x(a.x) && accumOrigin15y(a.y)
}

// symmetry 9 = 3x3 ω=40° X roots. Three rows: unit, primes, factor3
//	h1: x0           = cos(0)                      =  1
//	hP: x1 + x2 + x4 = cos(1ω) + cos(2ω) + cos(4ω) =  0
//	h3: x3           = cos(3ω)                     = -0.5
func accumOrigin9x(X []int) bool {
	// h1 is unity halfs equals the double of X[0]x[0] = X[0]*1 = 2X[0]
	h1 := 2*X[0]
	// hp is the primes half, first check X1 == X2 == X4 so
	// hp = X1x1 + X2x2 + X4x4 = X1(x1+x4+x4) = X1(0)
	// hp anyway should be zero so we don't use it
	if ok := accumAllEqual(X[1], X[2], X[4]); !ok {
		return false
	}
	// h3 is the factor-3 half X3x3 = (X3)(-1/2) so in halfs h3 = -X[3] 
	h3 := -X[3]
	// X is at origin when all halfs add up to zero
	return h1 + h3 == 0
}

// symmetry 9 = 3x3 ω=40° Y roots, three rows: unit, primes, factor-3
//	y0           = sin(0)                      = 0       : h1
//	y1 - y2 + y4 = sin(1ω) - sin(2ω) + sin(4ω) = 0       : hp
//	y3           = sin(3ω)                     = 0.866.. : h3
func accumOrigin9y(Y []int) bool {
	// Y[0] don't since product Y[0]*(y0) = Y[0](0) = 0
	// primes half, check primes: y1 == -y2 == y4
	if ok := accumAllEqual(Y[1], -Y[2], +Y[4]); !ok {
		return false
	} else {
		// since y1 - y2 + y4 = 0 we dont need to compute primes half
	}
	if Y[3] != 0 { // check (Y3)(y3) is zero
		return false
	}
	return true
}

// symmetry 15 = 3x5 ω=24° X roots:


// cos(0)                      = 1.000
// cos(1ω) + cos(4ω) + cos(6ω) = 0.913 - 0.104 - 0.809 = 0
// cos(2ω) + cos(3ω) + cos(7ω) = 0.669 + 0.309 - 0.978 = 0
// cos(5ω)                     = 0.500
func accumOrigin15x(X []int) bool {
	if accumAllEqual(X[1], X[4], X[6]) == false {
		return false
	}
	if accumAllEqual(X[2], X[3], X[7]) == false {
		return false
	}
	return 2*X[0] - X[5] == 0
}

// symmetry 15 = 3x5 ω=24° Y roots:
func accumOrigin15y(Y []int) bool {
	if Y[5] != 0 {
		return false
	}
	if accumAllEqual(Y[1], -Y[4], Y[6]) == false {
		return false
	}
	if accumAllEqual(Y[2], -Y[3], Y[7]) == false {
		return false
	}
	return true
}


func accumAllEqual(a ...int) bool {
	if n := len(a); n > 1 {
		for i := 1; i < n; i++ {
			if a[i-1] != a[i] {
				// at least one pair does not match
				return false
			}
		}
	}
	// all elements equal
	return true
}



/*
15x accumulators X exploration to get formulas
from localhost:8080/symm/15/hexagon/1,9,1,9,1,9?s=1&v={1...15}:

X0 X1 X2 X3 X4 X5 X6 X7
== == == == == == == ==              1,9,1,9,1,9
 1 -1  0  0 -1  2 -1  0  -> case +A
-1  1  0  0  1 -2  1  0  -> case -A
 0 -1  1  1 -1  0 -1  1  -> case +B
 0  0  0  0  0  0  0  0
 0  1 -1 -1  1  0  1 -1  -> case -B
 1 -1  0  0 -1  2 -1  0  -> case +A
-1  1  0  0  1 -2  1  0  -> case -A
 0 -1  1  1 -1  0 -1  1  -> case +B 
                                     2,8,2,8,2,8
 1  0 -1 -1  0  2  0 -1  -> case +C
 0  0  0  0  0  0  0  0
-1  0  1  1  0 -2  0  1  -> case -C
 0 -1  1  1 -1  0 -1  1  -> case +B
 0  1 -1 -1  1  0  1 -1  -> case -B
 1  0 -1 -1  0  2  0 -1  -> case +C
                                     3,7,3,7,3,7
 1  0 -1 -1  0  2  0 -1  -> case +C
 0  1 -1 -1  1  0  1 -1  -> case -B
 0 -1  1  1 -1  0 -1  1  -> case +B
-1  0  1  1  0 -2  0  1  -> case -C
                                     4,6,4,6,4,6
 1 -1  0  0 -1  2 -1  0  -> case +A
 0  1 -1 -1  1  0  1 -1  -> case -B

 X0 - X1           - X4 + 2X5 - X6      = 0   from A
     -X1 + X2 + X3 - X4       - X6 + X7 = 0   from B
 X0      - X2 - X3      + 2X5      - X7 = 0   from C

 B+C =>  X0 - X1 - X4 + 2X5 - X6
     =>   1             -1
     =   X1 + X4 + X6 = 0

 A-C => - X1 - X4  - X6 + X2 + X3 + X7
     = X2 + X3 + X7 = 0
               





*/