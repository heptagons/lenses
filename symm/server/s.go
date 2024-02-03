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

func (s *S) getHexas(h *dom.Html, call func(id string, h *dom.Html)) {
	p := symm.NewPolylines(s.Symm)
	hh := symm.NewHexagons(p)
	h.Elem(dom.Table, nil, func(h *dom.Html) {
		s.gonTableHeader(h, "Hexagon")
		for c, gon := range hh.All() {
			s.gonTableRow(h, c, gon, call)
		}
	})
}

func (s *S) getOctas(h *dom.Html, call func(id string, h *dom.Html)) {
	p := symm.NewPolylines(s.Symm)
	oo := symm.NewOctagons(p)
	h.Elem(dom.Table, nil, func(h *dom.Html) {
		s.gonTableHeader(h, "Octagon")
		for c, gon := range oo.All(1) {
			s.gonTableRow(h, c, gon, call)
		}
	})
}

func (s *S) getStars(h *dom.Html, call func(id string, h *dom.Html)) {
	p := symm.NewPolylines(s.Symm)
	ss := symm.NewStars(p)
	h.Elem(dom.Table, nil, func(h *dom.Html) {
		s.gonTableHeader(h, "Star")
		for c, gon := range ss.All(1) {
			s.gonTableRow(h, c, gon, call)
		}
	})
}

func (s *S) getHexa(h *dom.Html, vector int, angles []int) error {
	p := symm.NewPolylines(s.Symm)
	g := symm.NewHexagons(p)
	if gon, err := g.New(vector, angles); err != nil {
		return err
	} else {
		s.gonSvg(h, gon, 200)
		s.gonTables(h, gon)
		return nil
	}
}

func (s *S) getOcta(h *dom.Html, vector int, angles []int) error {
	p := symm.NewPolylines(s.Symm)
	g := symm.NewOctagons(p)
	if gon, err := g.New(vector, angles); err != nil {
		return err
	} else {
		s.gonSvg(h, gon, 250)
		s.gonTables(h, gon)
		return nil
	}
}

func (s *S) getStar(h *dom.Html, vector int, angles []int) error {
	p := symm.NewPolylines(s.Symm)
	g := symm.NewStars(p)
	if gon, err := g.New(vector, angles); err != nil {
		return err
	} else {
		s.gonSvg(h, gon, 300)
		s.gonTables(h, gon)
		return nil
	}
}

// ---

func (s *S) gonTableHeader(h *dom.Html, gon string) {
	h.Elem(dom.Tr, nil, func(h *dom.Html) {
		h.Elem(dom.Td, nil, "&nbsp;")
		h.Elem(dom.Th, nil, gon)
		h.Elem(dom.Th, nil, "Group")
		h.Elem(dom.Th, nil, "Angles")
		//h.Elem(dom.Th, nil, "Vectors")
	})
}

func (s *S) gonTableRow(h *dom.Html, c int, gon symm.Gon, call func(id string, h *dom.Html)) {
	h.Elem(dom.Tr, nil, func(h *dom.Html) {
		h.Elem(dom.Th, nil, fmt.Sprintf("%d", c))
		if !gon.Prime() {
			h.Elem(dom.Td, nil, fmt.Sprintf("Not prime"))
		} else if gon.Intersecting() {
			h.Elem(dom.Td, nil, fmt.Sprintf("Self intersecting"))
		} else {
			// button/link for particular valid hexagon
			h.Elem(dom.Td, nil, func(h *dom.Html) {
				call(gon.Id(), h)
			})
		}
		g := gon.Group()
		h.Elem(dom.Td, nil, fmt.Sprintf("%s<sub>%d</sub>", g.Letter, g.Number))
		h.Elem(dom.Td, nil, fmt.Sprintf("%v", gon.Angles()))
		//h.Elem(dom.Td, nil, fmt.Sprintf("%v", gon.Vectors()))
	})
}


func (s *S) gonSvg(h *dom.Html, gon symm.Gon, size int) {
	accums := gon.Accums()
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
	h.Elem(dom.Table, nil, func(h *dom.Html) {
		h.Elem(dom.Tr, nil, func(h *dom.Html) {
			h.Elem(dom.Th, nil, "Angles")
			h.Elem(dom.Td, nil, fmt.Sprintf("%v", gon.Angles()))
		})
		h.Elem(dom.Tr, nil, func(h *dom.Html) {
			h.Elem(dom.Th, nil, "Vectors")
			h.Elem(dom.Td, nil, fmt.Sprintf("%v", gon.Vectors()))
		})
	})
	h.Elem(dom.Table, nil, func(h *dom.Html) {
		h.Elem(dom.Caption, nil, "Accumulators")
		h.Elem(dom.Tr, nil, func(h *dom.Html) {
			h.Elem(dom.Td, nil, "")
			h.Elem(dom.Th, nil, "X")
			h.Elem(dom.Th, nil, "Y")
		})
		for c, accum := range gon.Accums() {
			h.Elem(dom.Tr, nil, func(h *dom.Html) {
				h.Elem(dom.Th, nil, fmt.Sprintf("%d", c+1))
				h.Elem(dom.Td, nil, fmt.Sprintf("%v", accum.X()))
				h.Elem(dom.Td, nil, fmt.Sprintf("%v", accum.Y()))
			})
		}

	})
}

