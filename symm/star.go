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

func (ss *Stars) All(shift, vector int) []Gon {
	symm := ss.p.s.s
	all := make([]Gon, 0)
	min := 1
	max := (symm-1) / 2
	fmt.Println("Stars.All", vector, min, max)
	for i := min; i <= max; i++ {
		if i == max {
			if s, err := ss.New([]int { i }, shift, vector); err == nil {
				all = append(all, s)
			}
		} else {
			if s, err := ss.New([]int { i, symm-1-i }, shift, vector); err == nil {
				all = append(all, s)
			}
		}
	}
	return all
}

func (ss *Stars) New(angles []int, shift, vector int) (Gon, error) {
	n := len(angles)
	symm := ss.p.s.s
	all := make([]int, 2*symm-1)
	id := ss.p.IdFromAngles(angles)
	switch n {
	case 1: 
		// regular polygon of size 2*symm
		a := angles[0]
		for i := range all {
			all[i] = a
		}
		return NewStar(ss.p, id, vector, all, n, NewGroupD(2*symm))

	case 2:
		// star
		a,b := angles[0], angles[1]
		for i := range all {
			if i % 2 == 0 {
				all[i] = a
			} else {
				all[i] = b
			}
		}
		return NewStar(ss.p, id, vector, all, n, NewGroupD(symm))

	default:
		return nil, fmt.Errorf("Number of angles out of range [1,2]")
	}
}


type Star struct {
	*Polygon
}

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
}

func (s *Star) Prime() bool {
	return true
}

func (s *Star) Intersecting() bool {
	return false
}




