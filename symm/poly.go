package symm

import (
	"fmt"
)

type Polylines struct {
	s *Symm
}

func NewPolylines(s *Symm) *Polylines {
	return &Polylines{
		s: s,
	}
}

func (pp *Polylines) New(vectors ...int) (*Polyline, error) {
	for v := 0; v < len(vectors); v++ {
		if vectors[v] < 1 {
			return nil, fmt.Errorf("Invalid vector %v at position %v", vectors[v], v)
		} else if vectors[v] > pp.s.s {
			return nil, fmt.Errorf("Invalid vector %v at position %v", vectors[v], v)
		}
	}
	return NewPolyline(pp, vectors), nil
}

type Polyline struct {
	pp      *Polylines
	vectors []int
}

func NewPolyline(pp *Polylines, vectors []int) *Polyline {
	return &Polyline{
		pp:      pp,
		vectors: vectors,
	}
}

func (p *Polyline) Angles() []int {
	if n :=  len(p.vectors); n < 2 {
		return nil
	} else {
		angles := make([]int, n-1)
		s := p.pp.s.s
		for i := 1; i < n; i++ {
			m := p.vectors[i-1]
			n := p.vectors[i]
			u := (s + m - n) % s
			angles[i-1] = u
		}
		return angles
	}
}

func (p *Polyline) Accums() []*Accum {
	n := len(p.vectors)
	t := p.pp.s.t
	base := NewAccum(t)
	pos := 0
	accum := 0
	accums := make([]*Accum, n)
	var indices []int	
	for i := 0; i < n; i++ {
		vindex := p.vectors[i] 
		if i % 2 == 0 {
			// v normal
			indices = p.pp.s.v[vindex-1]
		} else {
			// v overline
			indices = p.pp.s.w[vindex-1]
		}
		x := indices[0]
		y := indices[1]
		fmt.Printf("vindex=%d indices=%v\n", vindex, indices)
		if x < 0 {
			pos = -x - 1
			accum = -1
		} else {
			pos = +x - 1
			accum = +1
		}
		if pos < t {
			base.x[pos] += accum
		}
		// y row
		if y < 0 {
			pos = -y - 1
			accum = -1
		} else {
			pos = +y - 1
			accum = +1
		}
		if pos < t {
			base.y[pos] += accum
		}
		accums[i] = base.Clone()
	}
	return accums
}

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

func (a *Accum) String() string {
	return fmt.Sprintf("[x=%v y=%v", a.x, a.y)
}