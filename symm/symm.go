package symm

import (
	"fmt"
	"math"
)

const (
	MAX_SYMM = 37
)

var (
	SYMMS = []int{
		3, 5, 7, 9,
		11, 13, 15, 17, 19,
		21, 23, 25, 27, 29,
		31, 33, 35, 37,
	}
)

type Symm struct {
	s int       // The symmetry number: 3,5,7,9,...
	t int       // (s+1)/2
	x []float64 // x coordinates (cos)
	y []float64 // y coordinates (sin)
	v [][]int   // vector pair of indices
	w [][]int   // vector inverted 
}

func NewSymm(symm int) (*Symm, error) {
	if symm % 2 == 0 {
		return nil, fmt.Errorf("Not an odd symmetry")
	}
	if symm < 3 {
		return nil, fmt.Errorf("To small symmetry")
	}
	if symm > MAX_SYMM {
		return nil, fmt.Errorf("To big symmetry")
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

func (s *Symm) S() int {
	return s.s
}

func (s *Symm) XYs() ([]float64, []float64) {
	return s.x, s.y
}

func (s *Symm) Vectors() ([][]int, [][]int) {
	return s.v, s.w
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

func (s *Symm) Origin(a *Accum) (bool, error) {
	// confirm accumulator size matches with this symmetry
	if t := len(a.x); t != (s.s + 1) / 2 {
		return false, fmt.Errorf("Wrong accums size %d for symmetry %d", t, s.s)
	}
	switch s.s {

	case 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37:
		// for primes we expect all accumulators are zero
		return a.originPrime(), nil

	case 9: // 3x3
		return a.origin9(), nil

	case 15: // 3x5
		return a.origin15(), nil

	case 21, 25, 27, 33, 35:
		// provisional
		// TODO create origin21, origina25, etc...
		return a.originPrime(), nil

	default:
		return false, fmt.Errorf("Wrong symmetry %d", s.s) // should not happen
	}
}
