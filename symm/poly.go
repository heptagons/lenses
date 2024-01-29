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

func (p *Polylines) New(vectors ...int) (*Polyline, error) {
	for v := 0; v < len(vectors); v++ {
		if vectors[v] < 1 {
			return nil, fmt.Errorf("Invalid vector %v at position %v", vectors[v], v)
		} else if vectors[v] > p.s.s {
			return nil, fmt.Errorf("Invalid vector %v at position %v", vectors[v], v)
		}
	}
	return &Polyline{
		vectors: vectors,
	}, nil
}

type Polyline struct {
	vectors []int
}