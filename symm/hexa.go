package symm

import (
	"fmt"
	"strconv"
	"strings"
)

type Hexagons struct {
	p *Polylines
}

func NewHexagons(p *Polylines) *Hexagons {
	return &Hexagons{
		p: p,
	}
}

func (h *Hexagons) New(vertice int, angles []int) (Gon, error) {
	s := h.p.s.s
	n := len(angles)
	switch n {
	case 1:
		a := angles[0]
		if 6*a == 2*s {
			// try hexagon with angles a,a,a,a,a,a. Rotational symmetry D_6
			return NewHexagon(h.p, vertice, []int{ a,a,a,a,a }, 1)
		} else {
			return nil, fmt.Errorf("6(%d) != 2(%d)", a, s)
		}

	case 2:
		a, b := angles[0], angles[1]
		if 3*(a+b) == 2*s {
			if a == b {
				return NewHexagon(h.p, vertice, []int{ a,a,a,a,a }, 1)
			} else {
				// try hexagon with angles a,b,a,b,a,b. Rotational symmetry D_3
				return NewHexagon(h.p, vertice, []int{ a,b,a,b,a }, 2)
			}
		} else {
			return nil, fmt.Errorf("3(%d+%d) != 2(%d)", a,b,s)
		}

	case 3:
		a, b, c := angles[0], angles[1], angles[2]
		if 2*(a+b+c) == 2*s {
			// try hexagon with angles:a,b,c,a,b,c. Rotational symmetry C_2
			return NewHexagon(h.p, vertice, []int{ a,b,c,a,b }, 3)
		} else {
			return nil, fmt.Errorf("2(%d+%d+%d) != 2(%d)", a,b,c,s)
		}

	default:
		return nil, fmt.Errorf("Invalid number of angles not [1,2,3]")
	}
}

func (hh *Hexagons) All() []Gon {
	all := make([]Gon, 0)
	min := 1
	max := hh.p.s.s - 1
	for a := min; a <= max; a++ {
		for b := a; b <= max; b++ {
			if h, err := hh.New(1, []int{ a, b }); err == nil {
				all = append(all, h)
			}
			for c := b; c <=max; c++ {
				if a == b && b == c {
					// case 1
				} else if h, err := hh.New(1, []int{ a, b, c }); err == nil {
					all = append(all, h)
				}
			}
		}
	}
	return all
}


type Hexagon struct {
	p  *Polyline
	id string // simplification of angles list
}

func NewHexagon(pp *Polylines, vertice int, angles []int, size int) (Gon, error) {
	if p, err := pp.NewWithAngles(vertice, angles); err != nil {
		return nil, err
	} else {
		var ids []string
		for i := 0; i < size; i++ {
			ids = append(ids, strconv.Itoa(angles[i]))
		}
		return &Hexagon{
			p:  p,
			id: strings.Join(ids, ","),
		}, nil
	}
}

func (h *Hexagon) Accums() []*Accum {
	return h.p.Accums()
}

func (h *Hexagon) Id() string {
	return h.id
}

func (h *Hexagon) Angles() []int {
	return h.p.Angles()
}

func (h *Hexagon) Vectors() []int {
	return h.p.vectors
}

func (h *Hexagon) Prime() bool {
	s := h.p.pp.s.s // symmetry
	angles := h.p.Angles()
	a := s
	b := angles[0]
	c := angles[1]
	d := angles[2]
	gcd4(&s, &b, &c, &d)
	if s == a {
		return true
	} else {
		return false
	}
}

func (h *Hexagon) Intersecting() bool {
	return false
}

