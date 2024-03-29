package server

import (
	"fmt"

	"github.com/heptagons/lenses/symm"
	"github.com/heptagons/lenses/symm/dom"
)

type S struct {
	*symm.Symm
}

func newS(s int) (*S, error) {
	if s, err := symm.NewSymm(s); err != nil {
		return nil, err
	} else {
		return &S{
			Symm: s,
		}, nil
	}
}

func (s *S) getSymm(h *dom.Html) {
	xs, ys := s.XYs()
	h.Elem(dom.Table, nil, func(h *dom.Html) {
		h.Elem(dom.Caption, nil, "Points")
		h.Elem(dom.Tr, nil, func(h *dom.Html) {
			for _, c := range []string { "i", "x<sub>i</sub>", "y<sub>i</sub>" } {
				h.Elem(dom.Th, nil, c)
			}
		})
		for c := 0; c < len(xs); c++ {
			h.Elem(dom.Tr, nil, func(h *dom.Html) {
				h.Elem(dom.Th, nil, fmt.Sprintf("%d", c))
				h.Elem(dom.Td, nil, fmt.Sprintf("%f", xs[c]))	
				h.Elem(dom.Td, nil, fmt.Sprintf("%f", ys[c]))	
			})
		}
	})

	vs, ws := s.Vectors()
	h.Elem(dom.Table, nil, func(h *dom.Html) {
		h.Elem(dom.Caption, nil, "Vectors")
		h.Elem(dom.Tr, nil, func(h *dom.Html) {
			h.Elem(dom.Th, nil, "&nbsp;")
			h.Elem(dom.Th, nil, "Direct")
			h.Elem(dom.Th, nil, "Inverse")
		})
		for i := range vs {
			h.Elem(dom.Tr, nil, func(h *dom.Html) {
				h.Elem(dom.Th, nil, fmt.Sprintf("v<sub>%d</sub>", i+1))
				if v := vs[i]; len(v) == 2 {
					h.Elem(dom.Td, nil, fmt.Sprintf("(%d,%d)",v[0],v[1]))
				}
				if w := ws[i]; len(w) == 2 {
					h.Elem(dom.Td, nil, fmt.Sprintf("(%d,%d)",w[0],w[1]))
				}
			})
		}
	})
}

func (s *S) getHexas() *symm.Hexagons {
	p := symm.NewPolylines(s.Symm)
	return symm.NewHexagons(p)
}

func (s *S) getOctas() symm.Gons {
	p := symm.NewPolylines(s.Symm)
	return symm.NewOctagons(p)
}

func (s *S) getStars() symm.Gons {
	p := symm.NewPolylines(s.Symm)
	return symm.NewStars(p)
}

func (s *S) getGon(h *dom.Html, gon symm.Gon) {
	s.gonSvg(h, gon, 250)
	s.gonTables(h, gon)
}

// ---

func (s *S) gonTableHeader(h *dom.Html, gon string) {
	h.Elem(dom.Tr, nil, func(h *dom.Html) {
		h.Elem(dom.Td, nil, "&nbsp;")
		h.Elem(dom.Th, nil, "Group")
		h.Elem(dom.Th, nil, "Angles")
		h.Elem(dom.Th, nil, "Simple")
	})
}

func (s *S) gonTableRow(h *dom.Html, c,simple int, gon symm.Gon, call func(id string, h *dom.Html)) {
	h.Elem(dom.Tr, nil, func(h *dom.Html) {
		h.Elem(dom.Th, nil, fmt.Sprintf("%d", c))
		t := gon.Transforms()
		g := t.Group()
		h.Elem(dom.Td, nil, fmt.Sprintf("%s<sub>%d</sub>", g.Letter, g.Number))
		h.Elem(dom.Td, nil, func(h *dom.Html) {
			call(t.Id(), h)
			if !gon.Prime() {
				h.Elem(dom.Span, nil, fmt.Sprintf("Not prime"))
			}
		})
		h.Elem(dom.Th, nil, func(h *dom.Html) {
			if gon.Simple() {
				h.Elem(dom.Span, nil, fmt.Sprintf("%d", simple))
			}
		})
	})
}


func (s *S) gonSvg(h *dom.Html, gon symm.Gon, size int) {
	accums := gon.Polyline().Accums()
	var points [][]float64
	max := float64(0)
	for _, accum := range accums {
		xy := s.Symm.XY(accum)
		x, y := 20*xy[0], 20*xy[1]
		if x > 0 {
			if max < x { max = x }
		} else {
			if max < -x { max = -x }
		}
		if y > 0 {
			if max < y { max = y }
		} else {
			if max < -y { max = -y }
		}
		// Invert y to make we page y increases to top (like latex/tikz)
		points = append(points, []float64{ x, -y })
	}
	max *= 1.2
	a, b := int(-max), int(2*max)
	viewBox := []int{ a, a, b, b }
	fill0, stroke0 := "none", "#ddd"
	fill1, stroke1 := "#0ff8", "blue"
	lines := dom.NewAPath(fill0, stroke0)
	lines.M(int(-max), 0); lines.H(int(2*max)) // horizontal axis
	lines.M(0,int(-max));  lines.V(int(2*max)) // vertical axis
	h.Svg(size, size, viewBox, func(h *dom.Html) {
		h.Elem(dom.Path, []dom.Attr{ lines }, nil)
		h.Polygon(points, fill1, stroke1)
	})
}

func (s *S) gonTables(h *dom.Html, gon symm.Gon) {
	p := gon.Polyline()
	h.Elem(dom.Table, nil, func(h *dom.Html) {
		h.Elem(dom.Tr, nil, func(h *dom.Html) {
			h.Elem(dom.Th, nil, "Angles")
			h.Elem(dom.Td, nil, fmt.Sprintf("%v", p.Angles()))
		})
		h.Elem(dom.Tr, nil, func(h *dom.Html) {
			h.Elem(dom.Th, nil, "Edges")
			h.Elem(dom.Td, nil, fmt.Sprintf("%v", p.Edges()))
		})
	})
	h.Elem(dom.Table, nil, func(h *dom.Html) {
		h.Elem(dom.Caption, nil, "Accumulators")
		h.Elem(dom.Tr, nil, func(h *dom.Html) {
			h.Elem(dom.Td, nil, "")
			h.Elem(dom.Th, nil, "X")
			h.Elem(dom.Th, nil, "Y")
		})
		for c, accum := range p.Accums() {
			h.Elem(dom.Tr, nil, func(h *dom.Html) {
				h.Elem(dom.Th, nil, fmt.Sprintf("%d", c+1))
				h.Elem(dom.Td, nil, fmt.Sprintf("%v", accum.X()))
				h.Elem(dom.Td, nil, fmt.Sprintf("%v", accum.Y()))
			})
		}

	})
}

