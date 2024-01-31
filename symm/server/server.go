package server

import (
	"context"
	"fmt"

	"net/http"
	"github.com/go-chi/chi/v5"

	"github.com/heptagons/lenses/symm/dom"
)

func style(h *dom.Html) {
	h.WriteF(`
* { font-family:Arial; font-size:14px; }
.h1 { font-size:20px; color:#08f; }
.err { font-size:10px: color:#f00; }
`	)
}

var (
	h1  = []dom.Attr{ dom.Class("h1") }
	err = []dom.Attr{ dom.Class("err") }
)

func New(r *chi.Mux) {
	r.Get("/", getIndex)
	r.Route("/symm", func(r chi.Router) {
		r.Route("/{symm}", func(r chi.Router) {
			r.Use(SymmCtx)
			r.Get("/", getSymm)
		})
	})
}

func buttonLink(h *dom.Html, link, title string) {
	h.Elem(dom.Button, []dom.Attr{
		dom.NewOnclick("window.location.href= '%s'", link),
	}, title)
}

func getIndex(w http.ResponseWriter, r *http.Request) {
	h := dom.NewHtml(nil, "<!DOCTYPE html>\n")
	h.Elem(dom.Html_, nil, func(h *dom.Html) {
		h.Elem(dom.Head, nil, func(h *dom.Html) {
			h.Elem(dom.Script, nil, nil)
			h.Elem(dom.Style, nil, style)
		})
		h.Elem(dom.Body, nil, func(h *dom.Html) {
			h.Div(h1, "Symmetries")
			h.Elem(dom.Table, nil, func(h *dom.Html) {
				h.Elem(dom.Tr, nil, func(h *dom.Html) {
					h.Elem(dom.Th, nil, "Symm")
				})
				for _, s := range[]int { 3, 5, 7, 9 } {
					h.Elem(dom.Tr, nil, func(h *dom.Html) {
						h.Elem(dom.Td, nil, func(h *dom.Html) {
							buttonLink(h, fmt.Sprintf("/symm/%d", s), fmt.Sprintf("%d", s))
						})
					})
				}
			})
		})
	})
	h.Write(w)
}

func SymmCtx(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if symm, err := parseSymm(chi.URLParam(r, "symm")); err == nil {
	    ctx := context.WithValue(r.Context(), "symm", symm)
    	next.ServeHTTP(w, r.WithContext(ctx))
    } else {
		http.Error(w, http.StatusText(404), 404)
    	w.Write([]byte(err.Error()))
    }
  })
}

func parseSymm(symm string) (*Symm, error) {
	switch symm {
	case "3": return &Symm{ s:3 }, nil
	case "5": return &Symm{ s:5 }, nil
	case "7": return &Symm{ s:7 }, nil
	case "9": return &Symm{ s:9 }, nil
	default: return nil, fmt.Errorf("Invalid symmetry: %s", symm)
	}
}

func getSymm(w http.ResponseWriter, r *http.Request) {



	h := dom.NewHtml(nil, "<!DOCTYPE html>\n")
	h.Elem(dom.Html_, nil, func(h *dom.Html) {
		h.Elem(dom.Head, nil, func(h *dom.Html) {
			h.Elem(dom.Script, nil, nil)
			h.Elem(dom.Style, nil, style)
		})
		h.Elem(dom.Body, nil, func(h *dom.Html) {
			ctx := r.Context()
  			if symm, ok := ctx.Value("symm").(*Symm); ok {
				h.Div(h1, fmt.Sprintf("Symmetry %d", symm.s))
  			} else {
  				h.Div(err, "Symmetry value error")
  			}
		})
	})
	h.Write(w)
}

type Symm struct {
	s int
}
