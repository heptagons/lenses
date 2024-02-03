package symm

import (
	"fmt"
	"strconv"
	"strings"
)

type Octagons struct {
	p *Polylines
}

func NewOctagons(p *Polylines) *Octagons {
	return &Octagons{
		p: p,
	}
}

func (oo *Octagons) New(vector int, angles []int) (Gon, error) {
	n := len(angles)
	switch n {
	case 5:
		a, b, c, d, e := angles[0], angles[1], angles[2], angles[3], angles[4]
		return NewOctagon(oo.p, vector, []int{ a,b,c,d,e,d,c }, n)
	default:
		return nil, fmt.Errorf("Invalid number of angles not [4]")
	}	
}

func (oo *Octagons) All(vector int) []Gon {
	all := make([]Gon, 0)
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
							accums := o.Accums()
							if accums[len(accums)-1].AtOrigin() {
								all = append(all, o)
							}
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

func NewOctagon(pp *Polylines, vertice int, angles []int, size int) (Gon, error) {
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

func (o *Octagon) Prime() bool {
	s := o.p.pp.s.s // symmetry
	angles := o.p.Angles()
	a := s
	b := angles[0]
	c := angles[1]
	d := angles[2]
	e := angles[3]
	f := angles[4]
	gcd6(&s, &b, &c, &d, &e, &f)
	if s == a {
		return true
	} else {
		return false
	}
}

func (o *Octagon) Intersecting() bool {
	return false
}
