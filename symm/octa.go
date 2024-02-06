package symm

import (
	"fmt"
)

type Octagons struct {
	p *Polylines
	a *Angles
}

func NewOctagons(p *Polylines) *Octagons {
	return &Octagons{
		p: p,
		a: &Angles {
			min: 1,         // minimal possible individual angle
			max: p.s.s - 1, // maximum possible individual angle
			sum: 3*p.s.s,   // the sum of octagon internal angles
		},
	}
}

// All returns all the types of octagons (of group D1)
func (oo *Octagons) All() []Gon {
	all := make([]Gon, 0)
	min := oo.a.min
	max := oo.a.max
	sum := oo.a.sum
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
						t := oo.transforms([]int{ a,b,c,d,e })
						if o, err := oo.New(t, 1, 1); err == nil {
							all = append(all, o)
						}
					}
				}
			}
		}		
	}	
	return all
}

// Transforms validate the given minimal octagon angles and return
// sanitized angles and possible shifts and vectors to transform the octagon.
func (oo *Octagons) Transforms(angles []int) (*Transforms, error) {
	if len(angles) != 5 {
		return nil, fmt.Errorf("Invalid number of angles")
	}
	for pos, angle := range angles {
		if !oo.a.ValidAngle(angle) {
			return nil, fmt.Errorf("Invalid angle at position %d", pos)
		}
	}
	a, b, c, d, e := angles[0], angles[1], angles[2], angles[3], angles[4]
	if a > e {
		return nil, fmt.Errorf("Invalid angles: a > e")	
	} else if !oo.a.ValidSum(a + 2*b + 2*c + 2*d + e) {
		return nil, fmt.Errorf("Invalid angles: a + 2b + 2c + 2d + e != sum")
	}
	return oo.transforms(angles), nil
}

func (oo *Octagons) transforms(angles []int) *Transforms {
	// group is mirror symmetry like letters A,B,C,D,E,K...
	// shifts are eight vertices (no negative since no rotations)             
	// vectors is list 1,2,3,...,symm
	return &Transforms{
		id:      oo.p.IdFromAngles(angles),
		angles:  angles,
		group:   NewGroupD(1),
		shifts:  []int{ 1,2,3,4,5,6,7,8 }, 
		vectors: oo.p.vectors,             
	}
}

// New returns and octagon with symmetry dihedral 1
// angles array must include five angles valid respect to octagons symmetry angles.
// Angles a,b,c,d,e must have this conditions:
//	min <= a,b,c,d,e <= max
//	a <= e
//  a + e + 2b + 2c + 2d == sum
//  last accumulators must be zero (at origin)
func (oo *Octagons) New(t *Transforms, shift int, vector int) (Gon, error) {
	a, b, c, d, e := t.angles[0], t.angles[1], t.angles[2], t.angles[3], t.angles[4]
	seven := []int{ a,b,c,d,e,d,c }
	switch shift {
	case +1: seven = []int{ a,b,c,d,e,d,c }
	case +2: seven = []int{ b,c,d,e,d,c,b }
	case +3: seven = []int{ c,d,e,d,c,b,a }
	case +4: seven = []int{ d,e,d,c,b,a,b }
	case +5: seven = []int{ e,d,c,b,a,b,c }
	case +6: seven = []int{ d,c,b,a,b,c,d }
	case +7: seven = []int{ c,b,a,b,c,d,e }
	case +8: seven = []int{ b,a,b,c,d,e,d }
	}
	if o, err := NewOctagon(oo.p, t, seven, vector); err != nil {
		return nil, err
	} else {
		accums := o.Accums()
		if accums[len(accums)-1].AtOrigin() {
			return o, nil
		} else {
			return nil, fmt.Errorf("Not equilateral")
		}
	}
}




type Octagon struct {
	*Polygon
}

func NewOctagon(pp *Polylines, t *Transforms, angles []int, vector int) (Gon, error) {
	if p, err := NewPolygonT(pp, t, angles, vector); err != nil {
		return nil, err
	} else {
		return &Octagon{
			Polygon: p,
		}, nil
	}
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

