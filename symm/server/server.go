package server

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"net/http"
	"github.com/go-chi/chi/v5"

	"github.com/heptagons/lenses/symm/dom"
)

func style(h *dom.Html) {
	h.WriteF(`
* { font-family:Arial; font-size:14px; }
.h1 { font-size:20px; color:#08f; padding:0px 5px;}
.err { font-size:10px: color:#f00; }
table { border-collapse: collapse; }
table,td,th { border:1px solid #888888; }
td,th { padding:0px 5px; }
`	)
}

var (
	h1  = []dom.Attr{ dom.Class("h1") }
	domErr = []dom.Attr{ dom.Class("err") }
)

func New(r *chi.Mux) {
	r.Get("/", getIndex)
	r.Route("/symm", func(r chi.Router) {
		r.Get("/", getIndex)
		r.Route("/{symm}", func(r chi.Router) {
			r.Use(SymmCtx)
			r.Get("/", getSymm)
			r.Get("/hexagons", getHexas)
			r.Get("/hexagon/{id}", getHexa)
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

func getPage(body func(h *dom.Html)) *dom.Html {
	h := dom.NewHtml(nil, "<!DOCTYPE html>\n")
	h.Elem(dom.Html_, nil, func(h *dom.Html) {
		h.Elem(dom.Head, nil, func(h *dom.Html) {
			h.Elem(dom.Script, nil, nil)
			h.Elem(dom.Style, nil, style)
		})
		h.Elem(dom.Body, nil, func(h *dom.Html) {
			body(h)
		})
	})
	return h
}

func getIndex(w http.ResponseWriter, r *http.Request) {
	h := getPage(func(h *dom.Html) {
		h.Div(h1, "Symmetries")
		h.Elem(dom.Table, nil, func(h *dom.Html) {
			h.Elem(dom.Tr, nil, func(h *dom.Html) {
				h.Elem(dom.Th, nil, "Symm")
			})
			for _, s := range[]int { 3, 5, 7, 9, 11, 13, 15, 17, 19, 21 } {
				h.Elem(dom.Tr, nil, func(h *dom.Html) {
					h.Elem(dom.Td, nil, func(h *dom.Html) {
						buttonLink(h, fmt.Sprintf("/symm/%d", s), fmt.Sprintf("%d", s))
					})
				})
			}
		})
	})
	h.Write(w)
}

func getSymm(w http.ResponseWriter, r *http.Request) {
	h := getPage(func(h *dom.Html) {
		ctx := r.Context()
		if s, ok := ctx.Value("symm").(*S); ok {
			h.Elem(dom.Td, nil, func(h *dom.Html) {
				buttonLink(h, "/symm", "< Symmetries")
			})
			h.Div(h1, fmt.Sprintf("Symmetry %d", s.S()))
			s.getSymm(h)
			h.Div(nil, func(h *dom.Html) {
				buttonLink(h, fmt.Sprintf("/symm/%d/hexagons", s.S()), fmt.Sprintf("Hexagons H<sub>%d</sub>", s.S()))
			})
		} else {
			h.Div(domErr, "Symmetry value error")
		}
	})
	h.Write(w)
}

func getHexas(w http.ResponseWriter, r *http.Request) {
	h := getPage(func(h *dom.Html) {
		ctx := r.Context()
		if s, ok := ctx.Value("symm").(*S); ok {
			h.Div(nil, func(h *dom.Html) {
				buttonLink(h, fmt.Sprintf("/symm/%d", s.S()), fmt.Sprintf("< Symmetry %d", s.S()))
			})
			h.Div(h1, fmt.Sprintf("Hexagons H<sub>%d</sub>", s.S()))
			// todo add buttons
			s.getHexas(h, func(id string, h *dom.Html) {
				buttonLink(h, fmt.Sprintf("/symm/%d/hexagon/%s", s.S(), id), id)
			})
		} else {
			h.Div(domErr, "Symmetry value error")
		}
	})
	h.Write(w)
}

func getHexa(w http.ResponseWriter, r *http.Request) {
	h := getPage(func(h *dom.Html) {
		ctx := r.Context()
		if s, ok := ctx.Value("symm").(*S); ok {
			h.Div(nil, func(h *dom.Html) {
				buttonLink(h, fmt.Sprintf("/symm/%d/hexagons", s.S()), fmt.Sprintf("< Hexagons H<sub>%d</sub>", s.S()))
			})
			sids := chi.URLParam(r, "id") // is the six angles simplified example 1,1,7
			var ids []int
			for _, sid := range strings.Split(sids, ",") {
				if id, err := strconv.Atoi(sid); err != nil || id < 1 {
					h.Div(domErr, fmt.Sprintf("Invalid angle %s", sid))
					h.Write(w)
					return
				} else {
					ids = append(ids, id)
				}
			}
			h.Div(h1, fmt.Sprintf("Hexagon H<sub>%d</sub>(%s)", s.S(), sids))

			// row with buttons to change the first vector
			vector := 1
			h.Div(nil, func(h *dom.Html) {
				h.Elem(dom.Span, nil, "First vector: ")
				if v, err := strconv.Atoi(r.URL.Query().Get("vector")); err == nil && v > 0 && v <= s.S() {
					vector = v
				}
				for v := 1; v < s.S(); v++ {
					if v == vector {
						h.Elem(dom.Span, h1, fmt.Sprintf("%d", v))
					} else {
						buttonLink(h, fmt.Sprintf("/symm/%d/hexagon/%s?vector=%d", s.S(), sids, v), fmt.Sprintf("%d",v))
					}
				}
			})
			if err := s.getHexa(h, vector, ids); err != nil {
				h.Div(domErr, fmt.Sprintf("%v", err))
				h.Write(w)
				return
			}
		}
	})
	h.Write(w)
}

func getStars(w http.ResponseWriter, r *http.Request) {

}



