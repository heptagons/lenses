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
	return accumOrigin9x(a.x) || accumOrigin9y(a.y)
}

// origin9 checks if this accumulator for symmetry 3x5 represents the origin
func (a *Accum) origin15() bool {
	return accumOrigin15x(a.x) || accumOrigin15y(a.y)
}

// symmetry 9 = 3x3 ω=40° X roots:
//	h1: x0           = cos(0)                      =  1
//	hP: x1 + x2 + x4 = cos(1ω) + cos(2ω) + cos(4ω) =  0
//	h3: x3           = cos(3ω)                     = -1/2
func accumOrigin9x(X []int) bool {
	// h1 is unity halfs equals the double of X[0]x[0] = X[0]*1 = 2X[0]
	h1 := 2*X[0]
	// hp is the primes half, first check X1 == X2 == X4 so
	// hp = X1x1 + X2x2 + X4x4 = X1(x1+x4+x4) = X1(0)
	// hp anyway should be zero so we don't use it
	if ok := accumGroupEqual(X[1], X[2], X[4]); !ok {
		return false
	}
	// h3 is the factor-3 half X3x3 = (X3)(-1/2) so in halfs h3 = -X[3] 
	h3 := -X[3]
	// X is at origin when all halfs add up to zero
	return h1 + h3 == 0
}

// symmetry 9 = 3x3 ω=40° Y roots:
//	h1: y0           = sin(0)                      = 1
//	hP: y1 - y2 + y4 = sin(1ω) - sin(2ω) + sin(4ω) = 0
//	h3: y3           = sin(3ω)                     = 0.866.. 
func accumOrigin9y(Y []int) bool {
	// primes half, check primes: y1 == -y2 == y4
	if ok := accumGroupEqual(Y[1], -Y[2], +Y[4]); !ok {
		return false
	} else {
		// since y1 - y2 + y4 = 0 we dont need to compute primes half
	}
	h3 := Y[3]
	if h3 != 0 {
		return false
	}
	return true
}

// symmetry 15 = 3x5 ω=24° X roots:
//	h1: x0                     = +1     cos(0)                         ax = 2*x_0
//	hP: x1 + x2 + x4 + x6 + x7 = cos(1ω) + cos(2ω) + cos(4ω)   +0.5
//	h3: x3 + x6                = -0.5
//	h5: x5                     = -0.5
func accumOrigin15x(X []int) bool {
	return true
}

func accumOrigin15y(Y []int) bool {
	return true
}


func accumGroupEqual(ints ...int) bool {
	for i := 1; i < len(ints); i++ {
		if ints[i-1] != ints[i] {
			return false
		}
	}
	return true
}



