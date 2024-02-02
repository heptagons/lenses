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
		h.Elem(dom.Tr, nil, func(h *dom.Html) {
			for _, c := range []string {
				"i", 
				"x<sub>i</sub>",
				"y<sub>i</sub>",
			} {
				h.Elem(dom.Th, nil, c)
			}
		})
		for c := 0; c < len(xs); c++ {
			h.Elem(dom.Tr, nil, func(h *dom.Html) {
				h.Elem(dom.Td, nil, fmt.Sprintf("%d", c))
				h.Elem(dom.Td, nil, fmt.Sprintf("%f", xs[c]))	
				h.Elem(dom.Td, nil, fmt.Sprintf("%f", ys[c]))	
			})
		}
	})

	vs, ws := s.Vectors()
	h.Elem(dom.Table, nil, func(h *dom.Html) {
		h.Elem(dom.Tr, nil, func(h *dom.Html) {
			h.Elem(dom.Th, nil, "&nbsp;")
			for i := range vs {
				h.Elem(dom.Th, nil, fmt.Sprintf("v<sub>%d</sub>",i+1))
			}
		})
		h.Elem(dom.Tr, nil, func(h *dom.Html) {
			h.Elem(dom.Th, nil, "direct")
			for _, v := range vs {
				if len(v) == 2 {
					h.Elem(dom.Td, nil, fmt.Sprintf("(%d,%d)",v[0],v[1]))
				}
			}
		})
		h.Elem(dom.Tr, nil, func(h *dom.Html) {
			h.Elem(dom.Th, nil, "inverse")
			for _, v := range ws {
				if len(v) == 2 {
					h.Elem(dom.Td, nil, fmt.Sprintf("(%d,%d)",v[0],v[1]))	
				}
			}
		})
	})
}

func (s *S) getHexas(h *dom.Html, hexagon func(id string, h *dom.Html)) {
	p := symm.NewPolylines(s.Symm)
	hh := symm.NewHexagons(p)
	all := hh.All()
	h.Elem(dom.Table, nil, func(h *dom.Html) {
		h.Elem(dom.Tr, nil, func(h *dom.Html) {
			h.Elem(dom.Th, nil, "&nbsp;")
			h.Elem(dom.Th, nil, "Hexagon") // not prime, intersecting, etc
			h.Elem(dom.Th, nil, "Angles")
			h.Elem(dom.Th, nil, "Vectors")
		})
		for c, hexa := range all {
			h.Elem(dom.Tr, nil, func(h *dom.Html) {
				h.Elem(dom.Th, nil, fmt.Sprintf("%d", c+1))
				if !hexa.Prime() {
					h.Elem(dom.Td, nil, fmt.Sprintf("Not prime"))
				} else if hexa.SelfIntersecting() {
					h.Elem(dom.Td, nil, fmt.Sprintf("Self intersecting"))
				} else {
					h.Elem(dom.Td, nil, func(h *dom.Html) {
						// button/link for particular valid hexagon
						hexagon(hexa.Id(), h)
					})
				}
				h.Elem(dom.Td, nil, fmt.Sprintf("%v", hexa.Angles()))
				h.Elem(dom.Td, nil, fmt.Sprintf("%v", hexa.Vectors()))
			})
		}
	})
}

func (s *S) getHexa(h *dom.Html, vector int, angles []int) error {
	p := symm.NewPolylines(s.Symm)
	hh := symm.NewHexagons(p)
	hexa, err := hh.New(vector, angles)
	if err != nil {
		return err
	}
	accums := hexa.Accums()
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
		points = append(points, []float64{ x, y })
	}
	max *= 1.4
	a, b := int(-max), int(2*max)
	viewBox := []int{ a, a, b, b }
	// first the svg of the hexagon
	h.Svg(250, 250, viewBox, func(h *dom.Html) {
		fill := "cyan"
		stroke := "blue"
		h.Polygon(points, fill, stroke)
	})
	// later the tables
	h.Elem(dom.Table, nil, func(h *dom.Html) {
		h.Elem(dom.Tr, nil, func(h *dom.Html) {
			h.Elem(dom.Th, nil, "Angles")
			h.Elem(dom.Td, nil, fmt.Sprintf("%v", hexa.Angles()))
		})
		h.Elem(dom.Tr, nil, func(h *dom.Html) {
			h.Elem(dom.Th, nil, "Vectors")
			h.Elem(dom.Td, nil, fmt.Sprintf("%v", hexa.Vectors()))
		})
		for c, accum := range accums {
			h.Elem(dom.Tr, nil, func(h *dom.Html) {
				h.Elem(dom.Th, nil, fmt.Sprintf("Accum %d", c+1))
				h.Elem(dom.Td, nil, fmt.Sprintf("%v", accum))
			})
		}
	})
	return nil
}


func (s *S) getOctas(h *dom.Html, gon func(id string, h *dom.Html)) {
	p := symm.NewPolylines(s.Symm)
	hh := symm.NewOctagons(p)
	all := hh.All(1)
	h.Elem(dom.Table, nil, func(h *dom.Html) {
		h.Elem(dom.Tr, nil, func(h *dom.Html) {
			h.Elem(dom.Th, nil, "&nbsp;")
			h.Elem(dom.Th, nil, "Octagon") // not prime, intersecting, etc
			h.Elem(dom.Th, nil, "Angles")
			h.Elem(dom.Th, nil, "Vectors")
		})
		for c, p := range all {
			h.Elem(dom.Tr, nil, func(h *dom.Html) {
				h.Elem(dom.Th, nil, fmt.Sprintf("%d", c+1))
				//if !p.Prime() {
				//	h.Elem(dom.Td, nil, fmt.Sprintf("Not prime"))
				//} else if p.SelfIntersecting() {
				//	h.Elem(dom.Td, nil, fmt.Sprintf("Self intersecting"))
				//} else {
					h.Elem(dom.Td, nil, func(h *dom.Html) {
						// button/link for particular valid hexagon
						gon(p.Id(), h)
					})
				//}
				h.Elem(dom.Td, nil, fmt.Sprintf("%v", p.Angles()))
				h.Elem(dom.Td, nil, fmt.Sprintf("%v", p.Vectors()))
			})
		}
	})
}

func (s *S) getOcta(h *dom.Html, vector int, angles []int) error {
	p := symm.NewPolylines(s.Symm)
	oo := symm.NewOctagons(p)
	gon, err := oo.New(vector, angles)
	if err != nil {
		return err
	}
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
		points = append(points, []float64{ x, y })
	}
	max *= 1.4
	a, b := int(-max), int(2*max)
	viewBox := []int{ a, a, b, b }
	// first the svg of the hexagon
	h.Svg(250, 250, viewBox, func(h *dom.Html) {
		fill := "cyan"
		stroke := "blue"
		h.Polygon(points, fill, stroke)
	})
	// later the tables
	h.Elem(dom.Table, nil, func(h *dom.Html) {
		h.Elem(dom.Tr, nil, func(h *dom.Html) {
			h.Elem(dom.Th, nil, "Angles")
			h.Elem(dom.Td, nil, fmt.Sprintf("%v", gon.Angles()))
		})
		h.Elem(dom.Tr, nil, func(h *dom.Html) {
			h.Elem(dom.Th, nil, "Vectors")
			h.Elem(dom.Td, nil, fmt.Sprintf("%v", gon.Vectors()))
		})
		for c, accum := range accums {
			h.Elem(dom.Tr, nil, func(h *dom.Html) {
				h.Elem(dom.Th, nil, fmt.Sprintf("Accum %d", c+1))
				h.Elem(dom.Td, nil, fmt.Sprintf("%v", accum))
			})
		}
	})
	return nil
}
