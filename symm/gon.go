package symm

import (
	"fmt"
)

// Gon has the common methods for a polygon such as:
// hexagon, octagon or star.
type Gon interface {
	// Id is the identifier of the polygon.
	// Is the fewest number of angles separated by commas identifying the polygon.
	Id() string
	// The transformations of the polygon
	Transforms() *Transforms
	// The accumulators of the polygon to locate the vertices
	Accums() []*Accum
	// Angles is the complete list of angles of the polygon in sort order.
	Angles() []int
	// Edges is the complete list of edges vectors of the polygon.
	Edges() []int
	// Prime prints true if this polygon is congruent with another of smaller 
	// symmetry
	Prime() bool
	// Intersecting prints true if at least 
	Intersecting() bool
}

type Gons interface {
	// All returns all the types of polygon of all symmetry groups possible.
	All() []Gon
	// Transforms validate the given minimal polygon angles and return
	// sanitized angles and possible shifts and vectors to transform the polygon
	Transforms(angles []int) (*Transforms, error)
	// New returns a polygon given the angles in given transforms shifted and rotated
	New(t *Transforms, shift int, vector int) (Gon, error)
}

type Polygon struct {
	p  *Polyline
	id string // deprecate
	t  *Transforms
}

func NewPolygonT(pp *Polylines, t *Transforms, angles []int, vector int) (*Polygon, error) {
	if p, err := pp.NewWithAngles(vector, angles); err != nil {
		return nil, err
	} else {
		return &Polygon{
			p:  p,
			id: t.id,
			t:  t,
		}, nil
	}
}

// deprecate
func NewPolygon(pp *Polylines, id string, vector int, angles []int, size int, t *Transforms) (*Polygon, error) {
	if p, err := pp.NewWithAngles(vector, angles); err != nil {
		return nil, err
	} else {
		return &Polygon{
			p:  p,
			id: id,
			t:  t,
		}, nil
	}
}

func (p *Polygon) Transforms() *Transforms {
	return p.t
}

func (p *Polygon) String() string {
	return fmt.Sprintf("id=%s a=%v e=%v t=%v",
		p.id, p.p.Angles(), p.p.edges, p.t)
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

func (p *Polygon) Edges() []int {
	return p.p.edges
}


type GonAngles struct {
	min int
	max int
	num int
	sum int
}

func (a *GonAngles) Valid(angles []int) error {
	if a.num != len(angles) {
		return fmt.Errorf("Angles error, num=%d != %d", a.num, len(angles))
	}
	sum := 0
	for _, angle := range angles {
		if a.min > angle {
			return fmt.Errorf("Angle too small, min=%d > %d", a.min, angle)
		}
		if a.max < angle {
			return fmt.Errorf("Angle too big, max=%d < %d", a.max, angle)
		}
		sum += angle
	}
	if a.sum != sum {
		return fmt.Errorf("Angles sum error: sum=%d != %d", a.sum, sum)
	}
	return nil
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
