package symm

import (
	"fmt"
)

func NewHexaAngles(symm int) *GonAngles {
	return &GonAngles {
		min: 1,        // minimal possible individual angle
		max: symm - 1, // maximum possible individual angle
		num: 6,        // number of angles
		sum: 2*symm,   // the sum of octagon internal angles
	}
}


type Hexagons struct {
	p *Polylines
	a *GonAngles
}

func NewHexagons(p *Polylines) *Hexagons {
	return &Hexagons{
		p: p,
		a: NewHexaAngles(p.s.s),
	}
}

// All returns all hexagons
func (hh *Hexagons) All() []Gon {
	return hh.all6()
}

// Transforms validate the given minimal hexagon angles and return
// sanitized angles and possible shifts and vectors to transform the hexagon.
func (hh *Hexagons) Transforms(angles []int) (*Transforms, error) {
	if err := hh.a.Valid(angles); err != nil {
		return nil, err
	}
	return hh.t6(
		angles[0], angles[1], angles[2],
		angles[3], angles[4], angles[5],
	), nil
}

func (hh *Hexagons) New(t *Transforms, shift, vector int) (Gon, error) {
	if err := hh.a.Valid(t.angles); err != nil {
		return nil, err
	}
	return hh.new6(t, shift, vector)
}

func (hh *Hexagons) all6() []Gon {
	min := hh.a.min
	max := hh.a.max
	sum := hh.a.sum
	gons := make([]Gon, 0)
	N := uint64(0)
	p := NewPolylineN(hh.p, 6)
	for a := min; a <= max; a++ {
		//N = 0x000000000000
		N = uint64(a) << 40
		for b := min; b <= max; b++ { ab := a + b
			N &= 0xFF0000000000
			N |= uint64(b) << 32
			for c := min; c <= max; c++ { abc := ab + c
				N &= 0xFFFF00000000
				N |= uint64(c) << 24
				for d := min; d <= max; d++ { abcd := abc + d
					N &= 0xFFFFFF000000
					N |= uint64(d) << 16
					for e := min; e <= max; e++ { abcde := abcd + e
						N &= 0xFFFFFFFF0000
						N |= uint64(e) << 8
						if f := sum - abcde; f >= min && f <= max {
							N &= 0xFFFFFFFFFF00
							N |= uint64(f)
							//fmt.Println("all6 * ", N)
							if out := hexaReduce(N); out != N {
								// rotation/reflection already appended
								continue
							}
							//fmt.Println("all6", N)
							five := []int{ a,b,c,d,e }
							if err := p.SetAngles(1, five); err != nil {
								panic(err)
							}
							accums := p.Accums()
							last := accums[len(accums)-1]
							if ok, err := hh.p.s.Origin(last); err != nil {
								// last accum is at origin
								continue
							} else if !ok {
								continue
							}
							t := hh.t6(a,b,c,d,e,f)
							//fmt.Printf("%016x %d %v\n", N, C, t.Group()); C++
							if h, err := hh.New(t, 1, 1); err == nil {
								gons = append(gons, h)
							}
						}
					}
				}
			}
		}
	}
	return gons
}

// hexaReduce receive six 8-bit numbers (6 bytes) disposed side by side in a uint64.
// The bytes represent the six consecutive angles an hexagon (closed or not) has.
// When the bytes are A,B,C,D,E,F then n is:
//	n = A<<40 + B<<32 + C<<24 + D<<16 + E<<8 + F
// The bytes set is shifted and then reflected and shifted to locate the transormed n lower
// which is returned. Sixteen combinations are tested:
//	ABCDEF, ABCDEF, ABCDEF, FABCDE
//	EFABCD, DEFABC, CDEFAB, BCDEFA
//	FEDCBA, AFEDCB, BAFEDC, CBAFED
//	DCBAFE, EDCBAF, FEDCBA, FEDCBA
func hexaReduce(n uint64) uint64 {
	m0 := uint64(0xFFFFFFFFFFFF) // the biggest (the first)
	m1 := n
	m2 := uint64(0)
	for i := 0; i < 6; i++ {
		if m0 > m1 {
			m0 = m1
		}
		//fmt.Printf("\tm1 %016x\n", m1)
		low := m1 & 0xFF
		m2 |= low << (40 - 8*i)
		m1 >>= 8
		m1 |= low << 40
	}
	for i := 0; i < 6; i++ {
		if m0 > m2 {
			m0 = m2
		}
		//fmt.Printf("\tm2 %016x\n", m2)
		low := m2 & 0xFF
		m2 >>= 8
		m2 |= low << 40
	}
	//fmt.Printf("\tm0 %016x\n", m0)
	return m0
}


// t6 returns the hexagons Trasforms of six valid angles
func (hh *Hexagons) t6(a,b,c,d,e,f int) *Transforms {

	t := func(group *Group, shifts ...int) *Transforms {
		all := []int{ a,b,c,d,e,f }
		return NewTransforms(hh.p, all, group, shifts)
	}
	// D6
	// a b c d e f
	// 1 1 1 1 1 1 : case A
	if (a==b && b==c && c==d && d==e && e==f) { // (1)
		return t(NewGroupD(6), 1)
	}
	// D3
	// a b c d e f
	// 1 2 1 2 1 2 : case A
	if (a==c && c==e) && (b==d && d==f) { // (1) && (2)
		return t(NewGroupD(3), 1,2)
	}
	// D2
	// a b c d e f
	// 1 2 2 1 2 2 : case A
	// 2 1 2 2 1 2 : case B
	// 2 2 1 2 2 1 : case C
	if a==d && (b==c && c==e && e==f) { // 1 && (2) case A
		return t(NewGroupD(2), 1,2,3)
	}
	if b==e && (a==c && c==d && d==e) { // 1 && (2) case B
		return t(NewGroupD(2), 1,2,3)
	}
	if c==f && (a==b && b==d && d==e) { // 1 && (2) case C
		return t(NewGroupD(2), 1,2,3)
	}
	// D1
	// a b c d e f
	// - 1 2 - 2 1 : case A
	// 1 - 1 2 - 2 : case B
	// 2 1 - 1 2 - : case C
	if b==f && c==e { // 1 && 2 case A
		return t(NewGroupD(1), 1,2,3,4,5,6)
	}
	if a==c && d==f { // 1 && 2 case B
		return t(NewGroupD(1), 1,2,3,4,5,6)
	}
	if b==d && a==e { // 1 && 2 case C
		return t(NewGroupD(1), 1,2,3,4,5,6)
	}
	// C2
	// a b c d e f
	// 1 2 3 1 2 3 : case A
	// 3 1 2 3 1 2 : case B
	// 2 3 1 2 3 1 : case C
	if a==d && b==e && c==f { // cases A,B,C
		return t(NewGroupC(2), -3,-2,-1, 1,2,3)
	}
	// C1
	return t(NewGroupC(1), -6,-5,-4,-3,-2,-1, 1,2,3,4,5,6)
}

func (hh *Hexagons) new6(t *Transforms, shift, vector int) (Gon, error) {
	var five []int
	a,b,c,d,e,f := t.angles[0], t.angles[1], t.angles[2],
		t.angles[3], t.angles[4], t.angles[5]
	
	switch shift {
	default: five = []int{ a,b,c,d,e } // C1,C2,D1,D2,D3,D6
	case +2: five = []int{ b,c,d,e,f } // C1,C2,D1,D2,D3
	case +3: five = []int{ c,d,e,f,a } // C1,C2,D1,D2
	case +4: five = []int{ d,e,f,a,b } // C1,D1
	case +5: five = []int{ e,f,a,b,c } // C1,D1
	case +6: five = []int{ f,a,b,c,d } // C1,D1
	// rotate
	case -1: five = []int{ f,e,d,c,b } // C1,C2
	case -2: five = []int{ e,d,c,b,a } // C1,C2
	case -3: five = []int{ d,c,b,a,f } // C1,C2
	case -4: five = []int{ c,b,a,f,e } // C1
	case -5: five = []int{ b,a,f,e,d } // C1
	case -6: five = []int{ a,f,e,d,c } // C1
	}
	return NewHexagon(hh.p, t, five, vector)
}

type Hexagon struct {
	*Polygon
}

func NewHexagon(pp *Polylines, t *Transforms, angles []int, vector int) (Gon, error) {
	if p, err := NewPolygonT(pp, t, angles, vector); err != nil {
		return nil, err
	} else {
		return &Hexagon{
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



// Deprecate later:

func (hh *Hexagons) NewFast(t *Transforms, shift int, vector int) (Gon, error) {
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


// All2 returns all possible hexagons including possible existing symmetry group C_1
func (hh *Hexagons) _All2() {
	min := hh.a.min
	max := hh.a.max
	sum := hh.a.sum
	i := 1
	g := ""
	p := NewPolylineN(hh.p, 6)
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
							_ = p.SetAngles(1, []int{ a,b,c,d,e })
							accums := p.Accums()
							last := accums[len(accums)-1]
							if ok, err := hh.p.s.Origin(last); err != nil { // last accum is at origin
								continue
							} else if !ok {
								continue
							}
							g = "???"
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
							if g == "???" {
								fmt.Printf("\t%v\n", accums)
							}
						}
					}
				}
			}
		}
	}
}

// All return all the hexagons of symmetries groups D6,D3,D2 and C2
// ignoring possible (inexisting not prove available yet) C1
func (hh *Hexagons) allFast() []Gon {
	gons := make([]Gon, 0)
	max := hh.a.max
	for a := hh.a.min; a <= max; a++ {
		for b := a; b <= max; b++ {
			if a == b {
				t := hh.t1(a)
				if h, err := hh.New(t, 1, 1); err == nil {
					// regular hexagon
					gons = append(gons, h)
				}
			} else {
				t := hh.t2(a, b)
				if h, err := hh.New(t, 1, 1); err == nil {
					// triangular star
					gons = append(gons, h)
				}
			}
			for c := b; c <=max; c++ {
				if a == b && b == c {
					continue // regular hexagon already appended
				}
				t := hh.t3(a, b, c)
				if h, err := hh.New(t, 1, 1); err == nil {
					// lense C2/D2
					gons = append(gons, h)
				}
			}
		}
	}
	return gons
}

func (hh *Hexagons) transformsFast(angles []int) (*Transforms, error) {
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

