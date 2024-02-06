package symm

import (
	"fmt"
	"strconv"
	"strings"
)

// https://mathstat.slu.edu/escher/index.php/Introduction_to_Symmetry
type Group struct {
	Letter string
	Number int
}

// NewGroupC builds a rotational group
// C2 is the symmetry of letters N,S,Z
// C3 is the symmetry of triskelion
// C4 is the symmetry of swastika
func NewGroupC(number int) *Group {
	return &Group{
		Letter: "C",
		Number: number,
	}
}

// NewGroupD builds a diedral group
// D1 is the single mirror symmetry (like of letters A,B,C,D,E,K,M,T,U,V,W,Y)
// D2 is the symmetry of the rectangle (like of letters H,I,X,O)
// D3 is the symmetry of the equilateral triangle 
// D5 is the symmetry of the regular pentagon
// D6 is the symmetry of the regular hexagon
// D7 is the symmetry of the regular heptagon
// D9 is the symmetry of the regular 9-gon
// DN N >= 10 is the symmetry of the regular N-gon
func NewGroupD(number int) *Group {
	return &Group{
		Letter: "D",
		Number: number,
	}
}

var (
	C1 = NewGroupC(1) // only identity symmetry
	C2 = NewGroupC(2) // 180°
	D1 = NewGroupD(1) // mirror
	D2 = NewGroupD(2) // rectangle symmetry
	D3 = NewGroupD(3) // equilateral triangle
	D6 = NewGroupD(6) // regular hexagon
)



// Gon has the common methods for a polygon such as:
// hexagon, octagon or star.
type Gon interface {
	// Id is the identifier of the polygon.
	// Is the fewest number of angles separated by commas identifying the polygon.
	Id() string
	// The rotational group of the polygon
	Group() *Group
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
	group *Group
}

func NewPolygon(pp *Polylines, vertice int, angles []int, size int, group *Group) (*Polygon, error) {
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

func (p *Polygon) Group() *Group {
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
