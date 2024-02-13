package symm

import (
	"fmt"
)

// Gon has the common methods for a polygon such as:
// hexagon, octagon or star.
type Gon interface {
	
	// Polyline contains the edges and vertices of the polygon.
	Polyline() *Polyline
	
	// The transformations of the polygon to be placed in the plane.
	Transforms() *Transforms
	
	// Prime returns false when the polygon IS congruent with another one of of smaller symmetry.
	Prime() bool
	
	// Simple returns false when the polygon has intersecting edges.
	Simple() bool
}

type Gons interface {

	// All returns all the types of polygon of all symmetry groups possible.
	All() []Gon
	
	// Transforms validate the given minimal polygon angles and return
	// sanitized angles and possible shifts and vectors to set the polygon on the plane
	Transforms(angles []int) (*Transforms, error)
	
	// New returns a polygon given the angles in given transforms shifted and rotated
	New(t *Transforms, shift int, vector int) (Gon, error)
}

type Polygon struct {
	p      *Polyline
	t      *Transforms
	simple bool
}

func NewPolygon(pp *Polylines, t *Transforms, angles []int, vector int) (*Polygon, error) {
	if p, err := pp.NewWithAngles(vector, angles); err != nil {
		return nil, err
	} else {
		return &Polygon{
			p:      p,
			t:      t,
			//simple: Simple(pp, t),
		}, nil
	}
}

func (p *Polygon) Polyline() *Polyline {
	return p.p
}

func (p *Polygon) Transforms() *Transforms {
	return p.t
}

func (p *Polygon) Simple() bool {
	return p.simple
}

func (p *Polygon) String() string {
	return fmt.Sprintf("%v t=%v simple=%t", p.p, p.t, p.simple)
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
