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

func (oo *Octagons) All(shift, vector int) []Gon {
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
						if o, err := oo.New([]int{ a,b,c,d,e }, shift, vector); err == nil {
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


func (oo *Octagons) New(angles []int, shift, vector int) (Gon, error) {
	n := len(angles)
	switch n {
	case 5:
		a, b, c, d, e := angles[0], angles[1], angles[2], angles[3], angles[4]
		return NewOctagon(oo.p, vector, []int{ a,b,c,d,e,d,c }, n, NewGroupD(1))
	default:
		return nil, fmt.Errorf("Invalid number of angles not [4]")
	}	
}

// NewD1 returns and octagon with symmetry dihedral 1
// angles array must include five angles valid respect to octagons symmetry angles.
// Angles a,b,c,d,e must have this conditions:
//	min <= a,b,c,d,e <= max
//	a <= e
//  a + e + 2b + 2c + 2d == sum
//  last accumulators must be zero (at origin)
func (oo *Octagons) NewD1(angles []int, shift int, vector int) (Gon, error) {
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
	} else {
		seven := []int{ a,b,c,d,e,d,c }
		switch shift {
		case +1: seven = []int{ a,b,c,d,e,d,c }
		case +2: seven = []int{ b,c,d,e,d,c,b }
		case +3: seven = []int{ c,d,e,d,c,b,a }
		case +4: seven = []int{ d,e,d,c,b,a,b }
		case +5: seven = []int{ e,d,c,b,a,b,c }
		case +6: seven = []int{ d,c,b,a,b,c,d }
		case +7: seven = []int{ c,b,a,b,c,d,e }
		case +8: seven = []int{ b,a,b,c,d,e,a }
		}
		//id := oo.p.IdFromAngles(angles)
		return NewOctagon(oo.p, vector, seven, 5, NewGroupD(1))
	}
}



type Octagon struct {
	*Polygon
}

func NewOctagon(pp *Polylines, vector int, angles []int, size int, group *Group) (Gon, error) {
	if p, err := NewPolygon(pp, vector, angles, size, group); err != nil {
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

