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

func (a *Accum) AtOrigin() bool {
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

func (a *Accum) X() []int {
	return a.x
}

func (a *Accum) Y() []int {
	return a.y
}

func (a *Accum) String() string {
	return fmt.Sprintf("xy=%v%v", a.x, a.y)
}