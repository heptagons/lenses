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
			r.Use(symmCtx)
			r.Get("/", getSymm)
			r.Get("/hexagons", getHexas)
			r.Get("/octagons", getOctas)
			r.Get("/stars", getStars)

			r.Get("/hexagon/{id}", getHexa)
			r.Get("/octagon/{id}", getOcta)
			r.Get("/star/{id}", getStar)
		})
	})
}

func getIndex(w http.ResponseWriter, r *http.Request) {
	h := getPage(func(h *dom.Html) {
		h.Div(h1, "Symmetries")
		for _, s := range[]int { 3, 5, 7, 9, 11, 13, 15, 17, 19, 21 } {
			link := fmt.Sprintf("/symm/%d", s)
			buttonLink(h, link, fmt.Sprintf("%d", s))
		}
	})
	h.Write(w)
}

func symmCtx(next http.Handler) http.Handler {
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

func getSymm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s, ok := ctx.Value("symm").(*S)
	if !ok {
		w.Write([]byte("symm context error"))
		return
	}
	h := getPage(func(h *dom.Html) {
		// back button to return to symmetries
		h.Elem(dom.Td, nil, func(h *dom.Html) {
			link := "/symm"
			buttonLink(h, link, "<")
		})
		// title
		h.Div(h1, fmt.Sprintf("Symmetry %d", s.S()))
		// button to go hexagons
		h.Div(nil, func(h *dom.Html) {
			link := fmt.Sprintf("/symm/%d/hexagons", s.S())
			buttonLink(h, link, fmt.Sprintf("Hexagons H<sub>%d</sub>", s.S()))
		})
		// button to go to octagons
		h.Div(nil, func(h *dom.Html) {
			link := fmt.Sprintf("/symm/%d/octagons", s.S())
			buttonLink(h, link, fmt.Sprintf("Octagons O<sub>%d</sub>", s.S()))
		})
		// button to go to stars
		h.Div(nil, func(h *dom.Html) {
			link := fmt.Sprintf("/symm/%d/stars", s.S())
			buttonLink(h, link, fmt.Sprintf("Stars S<sub>%d</sub>", s.S()))
		})
		// symmetry details tables
		s.getSymm(h)
	})
	h.Write(w)
}

func getHexas(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s, ok := ctx.Value("symm").(*S)
	if !ok {
		w.Write([]byte("symm context error"))
		return
	}
	h := getPage(func(h *dom.Html) {
		// back button to return to symmetry s
		h.Div(nil, func(h *dom.Html) {
			buttonLink(h, fmt.Sprintf("/symm/%d", s.S()), "<")
		})
		// title
		h.Div(h1, fmt.Sprintf("Hexagons H<sub>%d</sub>", s.S()))
		// hexagons table and links to go to particular hexagon
		s.getHexas(h, func(id string, h *dom.Html) {
			buttonLink(h, fmt.Sprintf("/symm/%d/hexagon/%s", s.S(), id), id)
		})
	})
	h.Write(w)
}

func getOctas(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s, ok := ctx.Value("symm").(*S)
	if !ok {
		w.Write([]byte("symm context error"))
		return
	}
	h := getPage(func(h *dom.Html) {
		// back button to return to symmetry s
		h.Div(nil, func(h *dom.Html) {
			buttonLink(h, fmt.Sprintf("/symm/%d", s.S()), "<")
		})
		// title
		h.Div(h1, fmt.Sprintf("Octagons O<sub>%d</sub>", s.S()))
		// octagons table and links to go to particular octagon
		s.getOctas(h, func(id string, h *dom.Html) {
			buttonLink(h, fmt.Sprintf("/symm/%d/octagon/%s", s.S(), id), id)
		})
	})
	h.Write(w)
}

func getStars(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s, ok := ctx.Value("symm").(*S)
	if !ok {
		w.Write([]byte("symm context error"))
		return
	}
	h := getPage(func(h *dom.Html) {
		// back button to return to symmetry s
		h.Div(nil, func(h *dom.Html) {
			buttonLink(h, fmt.Sprintf("/symm/%d", s.S()), "<")
		})
		// title
		h.Div(h1, fmt.Sprintf("Stars S<sub>%d</sub>", s.S()))
		// stars table and links to go to particular star
		s.getStars(h, func(id string, h *dom.Html) {
			buttonLink(h, fmt.Sprintf("/symm/%d/star/%s", s.S(), id), id)
		})
	})
	h.Write(w)
}


func getHexa(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s, ok := ctx.Value("symm").(*S)
	if !ok {
		w.Write([]byte("symm context error"))
		return
	}
	getPage(func(h *dom.Html) {
		defer h.Write(w)
		// back button
		h.Div(nil, func(h *dom.Html) {
			buttonLink(h, fmt.Sprintf("/symm/%d/hexagons", s.S()), "<")
		})
		sids, ids, err := idAngles(r)
		if err != nil {
			h.Div(domErr, err.Error())
			return
		}
		// title
		h.Div(h1, fmt.Sprintf("Hexagon H<sub>%d</sub>(%s)", s.S(), sids))

		// rows with buttons to change the shifts and vectors
		link := fmt.Sprintf("/symm/%d/hexagon/%s", s.S(), sids)
		shift, vector := buttonsShiftVector(r, s.S(), h, link)
		// particular hexagon controls (svg and tables)
		if err := s.getHexa(h, ids, shift, vector); err != nil {
			h.Div(domErr, fmt.Sprintf("%v", err))
			return
		}
	})
}

func getOcta(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s, ok := ctx.Value("symm").(*S)
	if !ok {
		w.Write([]byte("symm context error"))
		return
	}
	getPage(func(h *dom.Html) {
		defer h.Write(w)
		// back button
		h.Div(nil, func(h *dom.Html) {
			buttonLink(h, fmt.Sprintf("/symm/%d/octagons", s.S()), "<")
		})
		sids, ids, err := idAngles(r)
		if err != nil {
			h.Div(domErr, err.Error())
			return
		}
		// title including ids
		h.Div(h1, fmt.Sprintf("Octagon O<sub>%d</sub>(%s)", s.S(), sids))

		// rows with buttons to change the shifts and vectors
		link := fmt.Sprintf("/symm/%d/octagon/%s", s.S(), sids)
		shift, vector := buttonsShiftVector(r, s.S(), h, link)
		// particular octagon controls (svg and tables)
		if err := s.getOcta(h, ids, shift, vector); err != nil {
			h.Div(domErr, fmt.Sprintf("%v", err))
			return
		}
	})
}

func getStar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s, ok := ctx.Value("symm").(*S)
	if !ok {
		w.Write([]byte("symm context error"))
		return
	}
	getPage(func(h *dom.Html) {
		defer h.Write(w)
		// back button
		h.Div(nil, func(h *dom.Html) {
			buttonLink(h, fmt.Sprintf("/symm/%d/stars", s.S()), "<")
		})
		sids, ids, err := idAngles(r)
		if err != nil {
			h.Div(domErr, err.Error())
			return
		}
		// title including ids
		h.Div(h1, fmt.Sprintf("Star S<sub>%d</sub>(%s)", s.S(), sids))

		// rows with buttons to change the shifts and vectors
		link := fmt.Sprintf("/symm/%d/star/%s", s.S(), sids)
		shift, vector := buttonsShiftVector(r, s.S(), h, link)
		// particular star controls (svg and tables)
		if err := s.getStar(h, ids, shift, vector); err != nil {
			h.Div(domErr, fmt.Sprintf("%v", err))
			return
		}
	})
}


func buttonsShiftVector(r *http.Request, symm int, h *dom.Html, prelink string) (int, int) {
	return shiftVectorOpts(r, symm, h, func(shift,vector int, text string) {
		link := fmt.Sprintf("%s?shift=%d&vector=%d", prelink, shift, vector)
		buttonLink(h, link, text)
	})
}




func buttonLink(h *dom.Html, link, title string) {
	h.Elem(dom.Button, []dom.Attr{
		dom.NewOnclick("window.location.href= '%s'", link),
	}, title)
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



// idAngles return the ids integer array from URL param id. 
func idAngles(r *http.Request) (string, []int, error) {
	// read is the eight angles simplified example 1,1,7
	sids := chi.URLParam(r, "id")
	var ids []int
	for _, sid := range strings.Split(sids, ",") {
		if id, err := strconv.Atoi(sid); err != nil || id < 1 {
			return sids, nil, fmt.Errorf("Invalid angle %s", sid)
		} else {
			ids = append(ids, id)
		}
	}
	return sids, ids, nil
}

// vectorOptions
func shiftVectorOpts(r *http.Request, symm int, h *dom.Html, option func(s, v int,text string)) (int,int) {
	shift := 1
	vector := 1
	v := r.URL.Query().Get("vector")
	if v, err := strconv.Atoi(v); err == nil && v > 0 && v <= symm {
		vector = v
	}
	h.Div(nil, func(h *dom.Html) {
		h.Elem(dom.Span, nil, "Vector: ")
		for v := 1; v <= symm; v++ {
			if v == vector {
				h.Elem(dom.Span, h1, fmt.Sprintf("%d", v))
			} else {
				text := fmt.Sprintf("%d",v)
				option(shift, v, text)
			}
		}
	})
	return shift, vector
}

