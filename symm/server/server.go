package server

import (
	"context"
	"fmt"
	"strconv"

	"net/http"
	"github.com/go-chi/chi/v5"

	"github.com/heptagons/lenses/symm/dom"
)

func style(h *dom.Html) {
	h.WriteF(`
* { font-family:Arial; font-size:14px; }
.h1 { font-size:20px; color:#08f; }
.err { font-size:10px: color:#f00; }
table { border-collapse: collapse; }
table,td,th { border:1px solid #888888; }
td,th { padding:0px 5px; }
`	)
}

var (
	h1  = []dom.Attr{ dom.Class("h1") }
	err = []dom.Attr{ dom.Class("err") }
)

func New(r *chi.Mux) {
	r.Get("/", getIndex)
	r.Route("/symm", func(r chi.Router) {
		r.Get("/", getIndex)
		r.Route("/{symm}", func(r chi.Router) {
			r.Use(SymmCtx)
			r.Get("/", getSymm)
			r.Get("/hexagons", getHexas)
			r.Get("/stars", getStars)
		})
	})
}

func buttonLink(h *dom.Html, link, title string) {
	h.Elem(dom.Button, []dom.Attr{
		dom.NewOnclick("window.location.href= '%s'", link),
	}, title)
}


func SymmCtx(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  	if s, err := strconv.Atoi(chi.URLParam(r, "symm")); err != nil {
		http.Error(w, http.StatusText(404), 404)
    	w.Write([]byte("Invalid symm number"))
    } else if s, err := newS(s); err == nil {
	    ctx := context.WithValue(r.Context(), "symm", s)
    	next.ServeHTTP(w, r.WithContext(ctx))
    } else {
		http.Error(w, http.StatusText(404), 404)
    	w.Write([]byte(err.Error()))
    }
  })
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
				for _, s := range[]int { 3, 5, 7, 9, 11, 13, 15 } {
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

func getSymm(w http.ResponseWriter, r *http.Request) {
	h := dom.NewHtml(nil, "<!DOCTYPE html>\n")
	h.Elem(dom.Html_, nil, func(h *dom.Html) {
		h.Elem(dom.Head, nil, func(h *dom.Html) {
			h.Elem(dom.Script, nil, nil)
			h.Elem(dom.Style, nil, style)
		})
		h.Elem(dom.Body, nil, func(h *dom.Html) {
			ctx := r.Context()
  			if s, ok := ctx.Value("symm").(*S); ok {
				h.Div(h1, fmt.Sprintf("Symmetry %d", s.S()))
				h.Div(nil, func(h *dom.Html) {
					buttonLink(h, fmt.Sprintf("/symm/%d/hexagons", s.S()), "Hexagons")
				})
				s.getSymm(h)
  			} else {
  				h.Div(err, "Symmetry value error")
  			}
		})
	})
	h.Write(w)
}

func getHexas(w http.ResponseWriter, r *http.Request) {
	h := dom.NewHtml(nil, "<!DOCTYPE html>\n")
	h.Elem(dom.Html_, nil, func(h *dom.Html) {
		h.Elem(dom.Head, nil, func(h *dom.Html) {
			h.Elem(dom.Script, nil, nil)
			h.Elem(dom.Style, nil, style)
		})
		h.Elem(dom.Body, nil, func(h *dom.Html) {
			ctx := r.Context()
  			if symm, ok := ctx.Value("symm").(*S); ok {
				h.Div(h1, fmt.Sprintf("Symmetry %d hexagons", symm.S()))
			}
		})
	})
	h.Write(w)
}

func getStars(w http.ResponseWriter, r *http.Request) {

}



