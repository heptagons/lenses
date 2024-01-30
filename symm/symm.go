package symm

import (
	"fmt"
	"math"
)

type Symm struct {
	s int       // The symmetry number: 3,5,7,9,...
	t int       // (s+1)/2
	x []float64 // x coordinates (cos)
	y []float64 // y coordinates (sin)
	v [][]int   // vector pair of indices
	w [][]int   // vector inverted 
}

func (s *Symm) XY(accum *Accum) []float64 {
	ax := float64(0)
	ay := float64(0)
	for i := 0; i < s.t; i++ {
		if x := accum.x[i]; x != 0 {
			ax += float64(x) * s.x[i]
		}
		if y := accum.y[i]; y != 0 {
			ay += float64(y) * s.y[i]
		}
	}
	return []float64{ ax, ay }
}

func NewSymm(symm int) (*Symm, error) {
	if symm < 3 {
		return nil, fmt.Errorf("To small symmetry")
	}
	if symm % 2 == 0 {
		return nil, fmt.Errorf("Not an odd symmetry")
	}
	return newSymm(symm), nil
}

func newSymm(s int) *Symm {
	t := (s + 1) / 2
	theta := 2*math.Pi / float64(s)
	// [ x_i, y_i ] i = 1,2,...,t
	x := make([]float64, t)
	y := make([]float64, t)
	for i := 0; i < t; i++ {
		x[i] = math.Cos(float64(i)*theta)
		y[i] = math.Sin(float64(i)*theta)
	}
	// v_i = | (j,j)  for i < t j = i
	//       | (j,-j) for i > t j = s + 2 - i
	v := make([][]int, s)
	w := make([][]int, s)
	for i := 0; i < s; i++ {
		if i < t {
			j := i + 1
			v[i] = []int{ +j, +j }
			w[i] = []int{ -j, -j }
		} else {
			j := s + 2 - i - 1
			v[i] = []int{ +j,-j }
			w[i] = []int{ -j,+j }
		}
	}
	return &Symm{
		s: s,
		t: t,
		x: x,
		y: y,
		v: v,
		w: w,
	}
}
