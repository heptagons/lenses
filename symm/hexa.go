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

func (hh *Hexagons) All2() {
	//all := make([]Gon, 0)
	min := hh.a.min
	max := hh.a.max
	sum := hh.a.sum
	i := 1
	g := ""
	for a := min; a <= max; a++ {
		for b := a; b <= max; b++ {
			ab := a + b
			for c := 1; c <= max; c++ {
				abc := ab + c
				for d := 1; d <= max; d++ {
					abcd := abc + d
					for e := 1; e <= max; e++ {
						abcde := abcd + e
						if f := sum - abcde; min <= f && f <= max {
							p := NewPolylineWithAngles(hh.p, 1, []int{
								a,b,c,d,e,
							})
							accums := p.Accums()
							last := accums[len(accums)-1]
							if ok, err := hh.p.s.Origin(last); err != nil { // last accum is at origin
								continue
							} else if !ok {
								continue
							}
							g = ""
							if a==b && b==c && c==d && d==e && e==f {
								g = "D_6"

							} else if a==c && c==e && b==d && b==f {
								g = "D_3"

							} else if a==d && b==e && c==f {
								if a==b || b==c {
									g = "D_2"
								} else {
									g = "C_2"
								}
								if a > b || b > c {
									continue
								}
							}
							fmt.Printf("%d) %s [%d,%d,%d,%d,%d,%d]\n", i, g, a,b,c,d,e,f)
							i++



						}
					}
				}
			}
		}
	}
}

func (hh *Hexagons) All() []Gon {
	all := make([]Gon, 0)
	max := hh.a.max
	for a := hh.a.min; a <= max; a++ {
		for b := a; b <= max; b++ {
			if a == b {
				t := hh.t1(a)
				if h, err := hh.New(t, 1, 1); err == nil {
					// regular hexagon
					all = append(all, h)
				}
			} else {
				t := hh.t2(a, b)
				if h, err := hh.New(t, 1, 1); err == nil {
					// triangular star
					all = append(all, h)
				}
			}
			for c := b; c <=max; c++ {
				if a == b && b == c {
					continue // regular hexagon already appended
				}
				t := hh.t3(a, b, c)
				if h, err := hh.New(t, 1, 1); err == nil {
					// lense C2/D2
					all = append(all, h)
				}
			}
		}
	}
	return all
}

func (hh *Hexagons) Transforms(angles []int) (*Transforms, error) {
	switch len(angles) {

	case 1:
		a := angles[0]
		return hh.t1(a), nil // regular hexagon D6

	case 2:
		a, b := angles[0], angles[1]
		if a == b {
			return hh.t1(a), nil // regular hexagon D6
		} else {
			return hh.t2(a,b), nil // triangular start D3
		}
	case 3:
		a, b, c := angles[0], angles[1], angles[2]
		if a == b && b == c {
			return hh.t1(a), nil // regular hexagon D6
		} else {
			return hh.t3(a,b,c), nil // lense C2/D2
		}

	default:
		return nil, fmt.Errorf("Invalid number of angles")
	}
}




// t1 returns a transformation with the symmetry group D6 (regular hexagon)
// shifts are only identity (all regular hexagon vertices are isogonal)
func (hh *Hexagons) t1(a int) *Transforms {
	angles := []int{ a }
	shifts :=  []int{ 1 }
	return NewTransforms(hh.p, angles, NewGroupD(6), shifts)
}

// t2 returns a transformation with the symmetry group D3 (equilateral triangle)
// shifts are two: star has two different vertices
func (hh *Hexagons) t2(a, b int) *Transforms {
	angles := []int{ a, b }
	shifts :=  []int{ 1, 2 }
	return NewTransforms(hh.p, angles, NewGroupD(3), shifts)
}

// t3 returns a transformation with symmetry groups D2 or C2.
// If at least two of the three angles a,b,c are equal then
// return with symmetry group D2 (letters N,S,Z) with 3 possible shifts.
// Otherwise return symmetry group C2 (letters H,I,X,O) with 6 possible shifts.
func (hh *Hexagons) t3(a, b, c int) *Transforms {
	angles := []int { a, b, c }
	if a == b || b == c {
		shifts :=  []int{ 1, 2, 3 }
		return NewTransforms(hh.p, angles, NewGroupD(2), shifts)
	} else {
		shifts :=  []int{ -3, -2, -1, 1, 2, 3 }
		return NewTransforms(hh.p, angles, NewGroupC(2), shifts)
	}
}

func (hh *Hexagons) New(t *Transforms, shift int, vector int) (Gon, error) {
	var five []int // the minimum five angles to form the equilateral hexagon
	sum := hh.a.sum
	n := len(t.angles)
	switch n {
	
	case 1: // group D6 expected
		a := t.angles[0]
		if 6*a != sum {
			return nil, fmt.Errorf("Angles sum error 6(%d) != %d", a, sum)
		}
		five = []int{ a,a,a,a,a } // shift=1 or default

	case 2: // group D3 expected
		a, b := t.angles[0], t.angles[1]
		if 3*(a+b) != sum {
			return nil, fmt.Errorf("Angles sum error 3(%d+%d) != %d", a, b, sum)
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
	if h, err := NewHexagon(hh.p, t, five, vector); err != nil {
		return nil, err
	} else {
		accums := h.Accums()
		last := accums[len(accums)-1]
		if ok, err := hh.p.s.Origin(last); err != nil { // last accum is at origin
			return nil, err
		} else if ok {
			return h, nil
		} else {
			return nil, fmt.Errorf("Hexagon is not equilateral")
		}
	}
}


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

