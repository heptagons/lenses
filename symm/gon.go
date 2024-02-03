package symm

import (
	"fmt"
	"strconv"
	"strings"
)

// Group is a rotational group
type Group int

const (
	// C2 is the symmetry of letters N,S,Z
	C2 Group = iota
	// C3 is the symmetry of triskelion
	C3 
	// D2 is the symmetry of the rectangle or letters H,I,X,O
	D2
	// D3 is the symmetry of the equilateral triangle 
	D3
	// D5 is the symmetry of the regular pentagon
	D5
	// D6 is the symmetry of the regular hexagon
	D6
	// D7 is the symmetry of the regular heptagon
	D7
	// D9 is the symmetry of the regular 9-gon
	D9 
	// DN N >= 10 is the symmetry of the regular N-gon
	D10
	D14
	D18
	// M1 is the mirror symmetry letters: A,B,C,D,E,K,M,T,U,V,W,Y
	M1
)

func (g Group) Name() (string, int) {
	switch g {
	case C2:  return "C", 2
	case C3:  return "C", 3
	case D2:  return "D", 2
	case D3:  return "D", 3
	case D5:  return "D", 5
	case D6:  return "D", 6
	case D7:  return "D", 7
	case D10: return "D", 10
	case D14: return "D", 14
	case D18: return "D", 18
	case M1:  return "M", 1
	default:  return "", 0
	}
}


// Gon has the common methods for a polygon such as:
// hexagon, octagon or star.
type Gon interface {
	// Id is the identifier of the polygon.
	// Is the fewest number of angles separated by commas identifying the polygon.
	Id() string
	// The rotational group of the polygon
	Group() Group
	// The accumulators of the polygon to locate the vertices
	Accums() []*Accum
	// Angles is the complete list of angles of the polygon in sort order.
	Angles() []int
	// Vectors is the complete list of vectors of the polygon (the edges).
	Vectors() []int
	// Prime prints true if this polygon is congruent with another of smaller 
	// symmetry
	Prime() bool
	// Intersecting prints true if at least 
	Intersecting() bool
}

type Polygon struct {
	p     *Polyline
	id    string
	group Group
}

func NewPolygon(pp *Polylines, vertice int, angles []int, size int, group Group) (*Polygon, error) {
	if p, err := pp.NewWithAngles(vertice, angles); err != nil {
		return nil, err
	} else {
		var ids []string
		for i := 0; i < size; i++ {
			ids = append(ids, strconv.Itoa(angles[i]))
		}
		return &Polygon{
			p:     p,
			id:    strings.Join(ids, ","),
			group: group,
		}, nil
	}
}

func (p *Polygon) Group() Group {
	return p.group
}

func (p *Polygon) String() string {
	return fmt.Sprintf("id=%s angles=%v vectors=%v", p.id, p.p.Angles(), p.p.vectors)
}

func (p *Polygon) Accums() []*Accum {
	return p.p.Accums()
}

func (p *Polygon) Id() string {
	return p.id
}

func (p *Polygon) Angles() []int {
	return p.p.Angles()
}

func (p *Polygon) Vectors() []int {
	return p.p.vectors
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

func gcd5(a, b, c, d, e *int) {
	if g := gcd(gcd(gcd(gcd(*a, *b), *c), *d), *e); g > 1 {
		*a /= g
		*b /= g
		*c /= g
		*d /= g		
		*e /= g
	}
}

func gcd6(a, b, c, d, e, f *int) {
	if g := gcd(gcd(gcd(gcd(gcd(*a, *b), *c), *d), *e), *f); g > 1 {
		*a /= g
		*b /= g
		*c /= g
		*d /= g		
		*e /= g
		*f /= g
	}
}
