package symm

import (
	"fmt"
)

func NewOctagonsAngles(symm int) *Angles {
	return &Angles {
		min: 1,        // minimal possible individual angle
		max: symm - 1, // maximum possible individual angle
		sum: 3*symm,   // the sum of octagon internal angles
	}
}

type Octagons struct {
	p *Polylines
	a *Angles
}

func NewOctagons(p *Polylines) *Octagons {
	return &Octagons{
		p: p,
		a: NewOctagonsAngles(p.s.s),
	}
}

// All returns all the types of octagons (of symmetry group D1)
func (oo *Octagons) All() []Gon {
	return oo.all5()
}

func (oo *Octagons) all5() []Gon {
	min := oo.a.min
	max := oo.a.max
	sum := oo.a.sum
	all := make([]Gon, 0)
	for a := min; a <= max; a++ {
		for e := a; e <= max; e++ {
			for b := min; b <= max; b++ {
				for c := min; c <= max; c++ {
					for d := min; d <= max; d++ {
						if a + e + 2*b + 2*c + 2*d != sum {
							continue
						}
						t := oo.t5([]int{ a,b,c,d,e })
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
	return oo.t5(angles), nil
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
		last := accums[len(accums)-1]
		if ok, err := oo.p.s.Origin(last); err != nil { // last accum is at origin
			return nil, err
		} else if ok {
			return o, nil
		} else {
			return nil, fmt.Errorf("Octagon D1 is not equilateral")
		}
	}
}


// t5 returns transforms with symmetry group of mirror symmetry like letters A,B,C,D,E,K...
// shifts are eight positives: for eight vertices (no negative since no rotations)             
func (oo *Octagons) t5(angles []int) *Transforms {
	shifts :=  []int{ 1,2,3,4,5,6,7,8 }
	return NewTransforms(oo.p, angles, NewGroupD(1), shifts)
}

func (oo *Octagons) all7() {
	min := oo.a.min
	max := oo.a.max
	sum := oo.a.sum
	//all := make([]Gon, 0)
	z := 0
	p := NewPolylineN(oo.p, 8)
	for a := min; a <= max; a++ {
		for b := a; b <= max; b++ {
			for c := min; c <= max; c++ {
				for d := min; d <= max; d++ {
					for e := min; e <= max; e++ {
						for f := min; f <= max; f++ {
							for g := min; g <= max; g++ {
								for h := min; h < sum - (a+b+c+d+e+f+g); h++ {
									_ = p.SetAngles(1, []int{ a,b,c,d,e,f,g })
									accums := p.Accums()
									last := accums[len(accums)-1]
									if ok, err := oo.p.s.Origin(last); err != nil { // last accum is at origin
										continue
									} else if !ok {
										continue
									}
									z++
									fmt.Println(z,  a,b,c,d,e,f,g,h)
								}
							}
						}
					}
				}
			}
		}
	}	
}

func (oo *Octagons) AllAngles() {
	min := oo.a.min
	max := oo.a.max
	sum := oo.a.sum
	N := uint64(0) // octaAngles
	C := 1
	p := NewPolylineN(oo.p, 8)
	for a := min; a <= max; a++ { 
		N = uint64(a) << 56
		for b := min; b <= max; b++ { ab := a+b
			N &= 0xFF00000000000000
			N |= uint64(b) << 48
			for c := min; c <= max; c++ { abc := ab+c
				N &= 0xFFFF000000000000
				N |= uint64(c) << 40
				for d := min; d <= max; d++ { abcd := abc + d
					N &= 0xFFFFFF0000000000
					N |= uint64(d) << 32
					for e := min; e <= max; e++ { abcde := abcd + e
						N &= 0xFFFFFFFF00000000
						N |= uint64(e) << 24
						for f := min; f <= max; f++ { abcdef := abcde + f
							N &= 0xFFFFFFFFFF000000
							N |= uint64(f) << 16
							for g := min; g <= max; g++ { abcdefg := abcdef + g
								N &= 0xFFFFFFFFFFFF0000
								N |= uint64(g) << 8
								if h := sum - abcdefg; h >= min && h <= max {
									N &= 0xFFFFFFFFFFFFFF00
									N |= uint64(h)
									out := octaAnglesReduce(N)
									if out != N {
										// rotation/reflection already appeared
										continue
									}
									_ = p.SetAngles(1, []int{ a,b,c,d,e,f,g })
									accums := p.Accums()
									last := accums[len(accums)-1]
									if ok, err := oo.p.s.Origin(last); err != nil {
										// last accum is at origin
										continue
									} else if !ok {
										continue
									}
									fmt.Printf("%016x %d\n", N, C); C++
								}
							}
						}
					}
				}
			}
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

// octaReduce receive eight 8-bit numbers (8 bytes) disposed side by side in a uint64.
// The bytes represent the eight consecutive angles an octagon (closed or not) has.
// When the eight bytes are A,B,C,D,E,F,G and H then n is:
//	n = A<<56 + B<<48 + C<<40 + D<<32 + E<<24 + F<<16 + G<<8 + H
// The bytes set is shifted and then reflected and shifted to locate the transormed n lower
// which is returned. Sixteen combinations are tested:
//	ABCDEFGH, HABCDEFG, GHABCDEF, FGHABCDE
//	EFGHABCD, DEFGHABC, CDEFGHAB, BCDEFGHA
//	HGFEDCBA, AHGFEDCB, BAHGFEDC, CBAHGFED
//	DCBAHGFE, EDCBAHGF, FEDCBAHG, GFEDCBAH
func octaAnglesReduce(n uint64) uint64 {
	m0 := uint64(0xFFFFFFFFFFFFFFFF) // the biggest (the first)
	m1 := n
	m2 := uint64(0)
	for i := 0; i < 8; i++ {
		if m0 > m1 {
			m0 = m1
		}
		//fmt.Printf("\tm1 %016x\n", m1)
		low := m1 & 0xFF
		m2 |= low << (56-8*i)
		m1 >>= 8
		m1 |= low << 56
	}
	for i := 0; i < 8; i++ {
		if m0 > m2 {
			m0 = m2
		}
		//fmt.Printf("\tm2 %016x\n", m2)
		low := m2 & 0xFF
		m2 >>= 8
		m2 |= low << 56
	}
	//fmt.Printf("\tm0 %016x\n", m0)
	return m0
}

// Octagon symmetry groups

// D_4 symmetry
// p = 12345678
// n = ABABABAB
//   = BABABABA

// D_2 symmetry
// p = 12345678
// n = A


