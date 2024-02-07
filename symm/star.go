package symm

import (
	"fmt"
)

type Stars struct {
	p *Polylines
}

func NewStars(p *Polylines) *Stars {
	return &Stars{
		p: p,
	}
}

func (ss *Stars) All() []Gon {
	symm := ss.p.s.s
	all := make([]Gon, 0)
	min := 1            // todo create Angles in ss first
	max := (symm-1) / 2 
	for i := min; i <= max; i++ {
		var t *Transforms
		if i == max {
			t = ss.transforms1([]int{ i })
			//if s, err := ss.New([]int { i }, shift, vector); err == nil {
			//	all = append(all, s)
			//}
		} else {
			t = ss.transforms2([]int{ i, symm - 1 - i })
			//if s, err := ss.New([]int { i, symm-1-i }, shift, vector); err == nil {
			//	all = append(all, s)
			//}
		}
		if s, err := ss.New(t, 1, 1); err == nil {
			all = append(all, s)
		}
	}
	return all
}

func (ss *Stars) Transforms(angles []int) (*Transforms, error) {
	switch len(angles) {

	case 1:
		return ss.transforms1(angles), nil

	case 2:
		if angles[0] == angles[1] {
			return ss.transforms1([]int{ angles[0] }), nil
		} else {
			return ss.transforms2(angles), nil
		}

	default:
		return nil, fmt.Errorf("Invalid number of angles")
	}
}

func (ss *Stars) transforms1(angles []int) *Transforms {
	// group is D_(2s) the regular 2s-gon
	// shifts are only identity
	// vectors is list 1,2,3,...,symm
	return &Transforms{
		id:      ss.p.IdFromAngles(angles),
		angles:  angles,
		group:   NewGroupD(2*ss.p.s.s),
		shifts:  []int{ 1 }, 
		vectors: ss.p.vectors,             
	}
}

func (ss *Stars) transforms2(angles []int) *Transforms {
	// group is D_s for a start
	// shifts are 2: internal and external vertices.
	// vectors is list 1,2,3,...,symm
	return &Transforms{
		id:      ss.p.IdFromAngles(angles),
		angles:  angles,
		group:   NewGroupD(ss.p.s.s),
		shifts:  []int{ 1,2 }, 
		vectors: ss.p.vectors,             
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


type Star struct {
	*Polygon
}

func NewStar(pp *Polylines, t *Transforms, angles []int, vector int) (Gon, error) {
	if p, err := NewPolygonT(pp, t, angles, vector); err != nil {
		return nil, err
	} else {
		return &Star{
			Polygon: p,
		}, nil
	}
}

/*
func NewStar(pp *Polylines, id string, vertice int, angles []int, size int, group *Group) (Gon, error) {
	t := &Transforms{
		group: group,
	}
	if p, err := NewPolygon(pp, id, vertice, angles, size, t); err != nil {
		return nil, err
	} else {
		return &Star{
			Polygon: p,
		}, nil
	}
}*/

func (s *Star) Prime() bool {
	return true
}

func (s *Star) Intersecting() bool {
	return false
}




