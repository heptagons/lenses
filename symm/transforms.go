package symm

import (
	"fmt"
)

type Transforms struct {
	p       *Polylines
	angles  []int
	group   *Group
	shifts  []int
}

func NewTransforms(p *Polylines, angles []int, group *Group, shifts []int) *Transforms {
	return &Transforms{
		p:       p,
		angles:  angles,
		group:   group,
		shifts:  shifts, 
	}
}

func (t *Transforms) Id() string {
	return t.p.IdFromAngles(t.angles)
}

func (t *Transforms) Group() *Group {
	return t.group
}

func (t *Transforms) Shifts() []int {
	return t.shifts
}

func (t *Transforms) Vectors() []int {
	return t.p.vectors
}

func (t *Transforms) String() string {
	return fmt.Sprintf("{symm=%d a=%v g=%s s=%v}",
		t.p.s.s, t.angles, t.group, t.shifts)
}
