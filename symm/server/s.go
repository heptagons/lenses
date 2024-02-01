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

func (s *S) getHexas(h *dom.Html) {

}