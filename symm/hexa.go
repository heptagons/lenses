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

func (h *Hexagons) New(vertice int, angles []int) (*Hexagon, error) {
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

func (hh *Hexagons) All() []*Hexagon {
	all := make([]*Hexagon, 0)
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

type Octagons struct {
	p *Polylines
}

func NewOctagons(p *Polylines) *Octagons {
	return &Octagons{
		p: p,
	}
}

func (oo *Octagons) New(vector int, angles []int) (*Octagon, error) {
	n := len(angles)
	switch n {
	case 5:
		a, b, c, d, e := angles[0], angles[1], angles[2], angles[3], angles[4]
		return NewOctagon(oo.p, vector, []int{ a,b,c,d,e,d,c }, 1)
	default:
		return nil, fmt.Errorf("Invalid number of angles not [4]")
	}	
}

func (oo *Octagons) All(vector int) []*Octagon {
	all := make([]*Octagon, 0)
	min := 1            // minimal possible individual angle
	max := oo.p.s.s - 1 // maximum possible individual angle
	sum := 3*oo.p.s.s   // the sum of octagon internal angles
	for a := min; a <= max; a++ {
		for e := a; e <= max; e++ {
			if a + e + 6 > sum {
				continue
			}
			for b := min; b <= max; b++ {
				if a + e + 2*b + 4 > sum {
					continue
				}
				for c := min; c <= max; c++ {
					if a + e + 2*b + 2*c + 2 > sum {
						continue
					}
					for d := min; d <= max; d++ {
						if a + e + 2*b + 2*c + 2*d != sum {
							continue
						}
						if o, err := oo.New(vector, []int{ a,b,c,d,e }); err == nil {
							all = append(all, o)
						}
					}
				}
			}
		}		
	}	
	return all
}

type Octagon struct {
	p  *Polyline
	id string
}

func NewOctagon(pp *Polylines, vertice int, angles []int, size int) (*Octagon, error) {
	if p, err := pp.NewWithAngles(vertice, angles); err != nil {
		return nil, err
	} else {
		var ids []string
		for i := 0; i < size; i++ {
			ids = append(ids, strconv.Itoa(angles[i]))
		}
		return &Octagon{
			p:  p,
			id: strings.Join(ids, ","),
		}, nil
	}
}

func (o *Octagon) String() string {
	return fmt.Sprintf("id=%s angles=%v vectors=%v", o.id, o.p.Angles(), o.p.vectors)
}

func (o *Octagon) Accums() []*Accum {
	return o.p.Accums()
}

func (o *Octagon) Id() string {
	return o.id
}

func (o *Octagon) Angles() []int {
	return o.p.Angles()
}

func (o *Octagon) Vectors() []int {
	return o.p.vectors
}



type Hexagon struct {
	p  *Polyline
	id string // simplification of angles list
}

func NewHexagon(pp *Polylines, vertice int, angles []int, size int) (*Hexagon, error) {
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

func (h *Hexagon) SelfIntersecting() bool {
	return false
}

// gcd returns the greatest common divisor of two integers
func gcd(a, b int) int {
	if b == 0 {
		return a
	}
	return gcd(b, a % b)
}

func gcd2(a, b *int) {
	if g := gcd(*a, *b); g > 1 {
		*a /= g
		*b /= g
	}
}

func gcd3(a, b, c *int) {
	if g := gcd(gcd(*a, *b), *c); g > 1 {
		*a /= g
		*b /= g
		*c /= g
	}
}

func gcd4(a, b, c, d *int) {
	if g := gcd(gcd(gcd(*a, *b), *c), *d); g > 1 {
		*a /= g
		*b /= g
		*c /= g
		*d /= g		
	}
}


