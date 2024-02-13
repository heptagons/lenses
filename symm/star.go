package symm

import (
	"fmt"
)

type Stars struct {
	p *Polylines
	a *GonAngles
}

func NewStars(p *Polylines) *Stars {
	return &Stars{
		p: p,
		a: &GonAngles {
			min: 1,               // minimal possible individual angle
			max: (p.s.s - 1) / 2, // maximum possible individual angle
		},
	}
}

func (ss *Stars) All() []Gon {
	symm := ss.p.s.s
	all := make([]Gon, 0)
	// stars
	for i := ss.a.min; i < ss.a.max; i++ {
		t := ss.tDsymm([]int{ i, symm - 1 - i })
		if s, err := ss.New(t, 1, 1); err == nil {
			all = append(all, s)
		}
	}
	// single regular polygon
	t := ss.tD2symm([]int{ ss.a.max })
	if s, err := ss.New(t, 1, 1); err == nil {
		all = append(all, s)
	}
	return all
}

func (ss *Stars) Transforms(angles []int) (*Transforms, error) {
	switch len(angles) {

	case 1:
		return ss.tD2symm(angles), nil

	case 2:
		if angles[0] == angles[1] {
			return ss.tD2symm([]int{ angles[0] }), nil
		} else {
			return ss.tDsymm(angles), nil
		}

	default:
		return nil, fmt.Errorf("Invalid number of angles")
	}
}

func (ss *Stars) New(t *Transforms, shift int, vector int) (Gon, error) {
	symm := ss.p.s.s
	all := make([]int, 2*symm-1)
	switch len(t.angles) {
	
	case 1:
		// regular polygon of size 2*symm
		a := t.angles[0]
		for i := range all {
			all[i] = a
		}
		return NewStar(ss.p, t, all, vector)
	
	case 2:
		// star
		a,b := t.angles[0], t.angles[1]
		if shift == 2 {
			a,b = t.angles[1], t.angles[0]
		}
		for i := range all {
			if i % 2 == 0 {
				all[i] = a
			} else {
				all[i] = b
			}
		}
		return NewStar(ss.p, t, all, vector)

	default:
		return nil, fmt.Errorf("Number of angles out of range [1,2]")
	}
}

// tD2symm returns a transformation with the group of the regular 2symm-gon
// shifts are only identity (all regular polygon vertices are isogonal)
// vectors is list 1,2,3,...,symm
func (ss *Stars) tD2symm(angles []int) *Transforms {
	return NewTransforms(ss.p, angles, NewGroupD(2*ss.p.s.s), []int{ 1 })
}

// tDsymm returns a transformation with the group of the star with symm points.
// shifts are 2: internal and external vertices of stars are different.
// vectors is list 1,2,3,...,symm
func (ss *Stars) tDsymm(angles []int) *Transforms {
	return NewTransforms(ss.p, angles, NewGroupD(ss.p.s.s), []int{ 1,2 })
}



type Star struct {
	*Polygon
}

func NewStar(pp *Polylines, t *Transforms, angles []int, vector int) (Gon, error) {
	if p, err := NewPolygon(pp, t, angles, vector); err != nil {
		return nil, err
	} else {
		return &Star{
			Polygon: p,
		}, nil
	}
}

func (s *Star) Prime() bool {
	return true
}



