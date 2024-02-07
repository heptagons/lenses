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
	max := hh.a.max
	for a := hh.a.min; a <= max; a++ {
		for b := a; b <= max; b++ {
			if a == b {
				t := hh.tD6([]int{ a })
				if h, err := hh.New(t, 1, 1); err == nil {
					// equilateral hexagon
					all = append(all, h)
				}
			} else {
				t := hh.tD3([]int{ a,b })
				if h, err := hh.New(t, 1, 1); err == nil {
					// triangular star
					all = append(all, h)
				}
			}
			for c := b; c <=max; c++ {
				if a == b && b == c {
					// append already above
				} else {
					t := hh.tC2([]int{ a,b,c })
					if s, err := hh.New(t, 1, 1); err == nil {
						// lense
						all = append(all, s)
					}
				}
			}
		}
	}
	return all
}

func (hh *Hexagons) Transforms(angles []int) (*Transforms, error) {
	switch len(angles) {

	case 1:
		return hh.tD6(angles), nil

	case 2:
		if angles[0] == angles[1] {
			return hh.tD6([]int{ angles[0] }), nil
		} else {
			return hh.tD3(angles), nil
		}
	case 3:
		return hh.tC2(angles), nil

	default:
		return nil, fmt.Errorf("Invalid number of angles")
	}
}




// tD6 returns a transformation with the symmetry group of the regular hexagon
// shifts are only identity (all regular hexagon vertices are isogonal)
func (hh *Hexagons) tD6(angles []int) *Transforms {
	shifts :=  []int{ 1 }
	return NewTransforms(hh.p, angles, NewGroupD(6), shifts)
}

// tD6 returns a transformation with the symmetry group of the equilateral triangle
// shifts are two: star has two different vertices
func (hh *Hexagons) tD3(angles []int) *Transforms {
	shifts :=  []int{ 1, 2 }
	return NewTransforms(hh.p, angles, NewGroupD(3), shifts)
}

// tD6 returns a transformation with the symmetry group of the isoscelles triangle
// shifts are six: three for each different angle and other three after mirror reflection
func (hh *Hexagons) tC2(angles []int) *Transforms {
	shifts :=  []int{ -3, -2, -1, 1, 2, 3 }
	return NewTransforms(hh.p, angles, NewGroupC(2), shifts)
}

func (hh *Hexagons) New(t *Transforms, shift int, vector int) (Gon, error) {
	var five []int // the minimum five angles to form the equilateral hexagon
	sum := hh.a.sum
	switch len(t.angles) {
	
	case 1: // group D6 expected
		a := t.angles[0]
		if 6*a != sum {
			return nil, fmt.Errorf("Angles sum error 6(%d) != %d", a, sum)
		}
		five = []int{ a,a,a,a,a } // shift=1 or default

	case 2: // group D3 expected
		a, b := t.angles[0], t.angles[1]
		if 3*(a+b) != hh.a.sum {
			return nil, fmt.Errorf("Angles sum error 3(%d+%d) != sum", a, b, sum)
		}
		switch shift {
		default: five = []int{ a,b,a,b,a }
		case 2:  five = []int{ b,a,b,a,b }
		}

	case 3: // group C2 expected
		a, b, c := t.angles[0], t.angles[1], t.angles[2]
		if 2*(a+b+c) != sum {
			return nil, fmt.Errorf("Angles sum error 2(%d+%d+%d) != %d", a, b, c, sum)
		}
		switch shift {
		default: five = []int { a,b,c,a,b } // +1
		case +2: five = []int { b,c,a,b,c }
		case +3: five = []int { c,a,b,c,a }
		case -1: five = []int { b,a,c,b,a }
		case -2: five = []int { c,b,a,c,b }
		case -3: five = []int { a,c,b,a,c }
		}

	default:
		return nil, fmt.Errorf("Number of angles out of range [1,2,3]")
	}
	return NewHexagon(hh.p, t, five, vector)
}

// to deprecate
/*
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
}*/





type Hexagon struct {
	*Polygon
}

func NewHexagon(pp *Polylines, t *Transforms, angles []int, vector int) (Gon, error) {
	if p, err := NewPolygonT(pp, t, angles, vector); err != nil {
		return nil, err
	} else {
		return &Star{
			Polygon: p,
		}, nil
	}
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

