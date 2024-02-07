package symm

import (
	"fmt"
)

type Hexagons struct {
	p *Polylines
	a *Angles
}

func NewHexagons(p *Polylines) *Hexagons {
	return &Hexagons{
		p: p,
		a: &Angles {
			min: 1,         // minimal possible individual angle
			max: p.s.s - 1, // maximum possible individual angle
			sum: 2*p.s.s,   // the sum of octagon internal angles
		},
	}
}

func (hh *Hexagons) All() []Gon {
	all := make([]Gon, 0)
	min := 1
	max := hh.p.s.s - 1
	shift := 1
	vector := 1
	for a := min; a <= max; a++ {
		for b := a; b <= max; b++ {
			if h, err := hh.New([]int{ a, b }, shift, vector); err == nil {
				all = append(all, h)
			}
			for c := b; c <=max; c++ {
				if a == b && b == c {
					// case 1
				} else if h, err := hh.New([]int{ a, b, c }, shift, vector); err == nil {
					all = append(all, h)
				}
			}
		}
	}
	return all
}


func (hh *Hexagons) New(angles []int, shift, vector int) (Gon, error) {
	s := hh.p.s.s
	n := len(angles)
	id := hh.p.IdFromAngles(angles)
	switch n {
	case 1:
		a := angles[0]
		if 6*a == 2*s {
			// try hexagon with angles a,a,a,a,a,a. Rotational symmetry D_6
			return hh.new(id, vector, []int{ a,a,a,a,a }, 1, D6)
		} else {
			return nil, fmt.Errorf("6(%d) != 2(%d)", a, s)
		}

	case 2:
		a, b := angles[0], angles[1]
		if 3*(a+b) == 2*s {
			if a == b {
				return hh.new(id, vector, []int{ a,a,a,a,a }, 1, D6)
			} else {
				// try hexagon with angles a,b,a,b,a,b. Rotational symmetry D_3
				return hh.new(id, vector, []int{ a,b,a,b,a }, 2, D3)
			}
		} else {
			return nil, fmt.Errorf("3(%d+%d) != 2(%d)", a,b,s)
		}

	case 3:
		a, b, c := angles[0], angles[1], angles[2]
		if 2*(a+b+c) == 2*s {
			// try hexagon with angles:a,b,c,a,b,c. Rotational symmetry C_2
			return hh.new(id, vector, []int{ a,b,c,a,b }, 3, C2)
		} else {
			return nil, fmt.Errorf("2(%d+%d+%d) != 2(%d)", a,b,c,s)
		}

	default:
		return nil, fmt.Errorf("Invalid number of angles not [1,2,3]")
	}
}

func (hh *Hexagons) new(id string, vertice int, angles []int, size int, group *Group) (Gon, error) {
	
	t := &Transforms{
		group: group,
	}

	if p, err := NewPolygon(hh.p, id, vertice, angles, size, t); err != nil {
		return nil, err
	} else {
		return &Hexagon{
			Polygon: p,
		}, nil
	}
}





type Hexagon struct {
	*Polygon
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

