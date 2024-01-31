package symm

import (
	"fmt"
)

type Hexagons struct {
	p *Polylines
}

func NewHexagons(p *Polylines) *Hexagons {
	return &Hexagons{
		p: p,
	}
}

func (h *Hexagons) New(vertice int, angles []int) (*Hexagon, error) {
	s := h.p.s.s
	n := len(angles)
	switch n {
	case 1:
		a := angles[0]
		if 6*a == 2*s {
			// try hexagon with angles a,a,a,a,a,a. Rotational symmetry D_6
			return NewHexagon(h.p, vertice, []int{ a,a,a,a,a })
		} else {
			return nil, fmt.Errorf("6(%d) != 2(%d)", a, s)
		}

	case 2:
		a, b := angles[0], angles[1]
		if 3*(a+b) == 2*s {
			// try hexagon with angles a,b,a,b,a,b. Rotational symmetry D_3
			return NewHexagon(h.p, vertice, []int{ a,b,a,b,a })
		} else {
			return nil, fmt.Errorf("3(%d+%d) != 2(%d)", a,b,s)
		}

	case 3:
		a, b, c := angles[0], angles[1], angles[2]
		if 2*(a+b+c) == 2*s {
			// try hexagon with angles:a,b,c,a,b,c. Rotational symmetry C_2
			return NewHexagon(h.p, vertice, []int{ a,b,c,a,b })
		} else {
			return nil, fmt.Errorf("2(%d+%d+%d) != 2(%d)", a,b,c,s)
		}

	default:
		return nil, fmt.Errorf("Invalid number of angles not [1,2,3]")
	}
}

type Hexagon struct {
	p      *Polyline
	accums []*Accum
}

func NewHexagon(pp *Polylines, vertice int, angles []int) (*Hexagon, error) {
	if p, err := pp.NewWithAngles(vertice, angles); err != nil {
		return nil, err
	} else {
		accums := p.Accums()
		if last := accums[len(accums)-1]; !last.Zero() {
			fmt.Println("NOT CLOSED", angles, p.vectors, accums)
			fmt.Println("LAST", pp.s.XY(last))
			return nil, fmt.Errorf("Not closed, last accum %v", last)
		}
		return &Hexagon{
			p:      p,
			accums: accums,
		}, nil
	}
}