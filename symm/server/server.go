package server

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"net/http"
	"github.com/go-chi/chi/v5"

	"github.com/heptagons/lenses/symm"
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
	getPage(w, func(h *dom.Html) error {
		h.Div(h1, "Symmetries")
		for _, s := range symm.SYMMS {
			link := fmt.Sprintf("/symm/%d", s)
			buttonLink(h, link, fmt.Sprintf("%d", s))
		}
		return nil
	})
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

func symmOk(w http.ResponseWriter, r *http.Request) *S {
	ctx := r.Context()
	if s, ok := ctx.Value("symm").(*S); !ok {
		w.Write([]byte("symm context error"))
		return nil
	} else {
		return s
	}
}

func getSymm(w http.ResponseWriter, r *http.Request) {
	s := symmOk(w, r)
	if s == nil {
		return
	}
	getPage(w, func(h *dom.Html) error {
		// back button to return to symmetries
		h.Elem(dom.Td, nil, func(h *dom.Html) {
			buttonLink(h, "/symm", "<")
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
		return nil
	})
}

func getHexas(w http.ResponseWriter, r *http.Request) {
	if s := symmOk(w, r); s == nil {
		return
	} else {
		getPage(w, func(h *dom.Html) error {
			return getGons(s, h, s.getHexas(),
				"Hexagons H<sub>%d</sub>",
				"Hexagon",
				"/symm/%d/hexagon/%s",
			)
		})
	}
}

func getOctas(w http.ResponseWriter, r *http.Request) {
	if s := symmOk(w, r); s == nil {
		return
	} else {
		getPage(w, func(h *dom.Html) error {
			return getGons(s, h, s.getOctas(),
				"Octagons O<sub>%d</sub>",
				"Octagon",
				"/symm/%d/octagon/%s",
			)
		})
	}
}

func getStars(w http.ResponseWriter, r *http.Request) {
	if s := symmOk(w, r); s == nil {
		return
	} else {
		getPage(w, func(h *dom.Html) error {
			return getGons(s, h, s.getStars(),
				"Stars S<sub>%d</sub>",
				"Star",
				"/symm/%d/star/%s",
			)
		})
	}
}


func getHexa(w http.ResponseWriter, r *http.Request) {
	if s := symmOk(w, r); s == nil {
		return
	} else {
		getPage(w, func(h *dom.Html) error {
			return getGon(r, s, h, s.getHexas(),
				"/symm/%d/hexagons",
				"Hexagon H<sub>%d</sub>(%s)",
				"/symm/%d/hexagon/%s",
			)
		})
	}
}


func getOcta(w http.ResponseWriter, r *http.Request) {
	if s := symmOk(w, r); s == nil {
		return
	} else {
		getPage(w, func(h *dom.Html) error {
			return getGon(r, s, h, s.getOctas(),
				"/symm/%d/octagons",
				"Octagon O<sub>%d</sub>(%s)",
				"/symm/%d/octagon/%s",
			)
		})
	}
}

func getStar(w http.ResponseWriter, r *http.Request) {
	if s := symmOk(w, r); s == nil {
		return
	} else {
		getPage(w, func(h *dom.Html) error {
			return getGon(r, s, h, s.getStars(),
				"/symm/%d/stars",
				"Star S<sub>%d</sub>(%s)",
				"/symm/%d/star/%s",
			)
		})
	}
}



func getPage(w http.ResponseWriter, body func(h *dom.Html) error) {
	h := dom.NewHtml(nil, "<!DOCTYPE html>\n")
	h.Elem(dom.Html_, nil, func(h *dom.Html) {
		h.Elem(dom.Head, nil, func(h *dom.Html) {
			h.Elem(dom.Script, nil, nil)
			h.Elem(dom.Style, nil, style)
		})
		h.Elem(dom.Body, nil, func(h *dom.Html) {
			if err := body(h); err != nil {
				h.Div(domErr, err.Error())
			}
		})
	})
	h.Write(w)
}

func buttonLink(h *dom.Html, link, title string) {
	h.Elem(dom.Button, []dom.Attr{
		dom.NewOnclick("window.location.href= '%s'", link),
	}, title)
}


func getGons(s *S, h *dom.Html, gg symm.Gons, title, name, link string) error {
	// back button to return to symmetry s
	h.Div(nil, func(h *dom.Html) {
		buttonLink(h, fmt.Sprintf("/symm/%d", s.S()), "<")
	})
	// title
	h.Div(h1, fmt.Sprintf(title, s.S()))
	// octagons table and links to go to particular octagon
	h.Elem(dom.Table, nil, func(h *dom.Html) {
		s.gonTableHeader(h, name)
		for c, gon := range gg.All() {
			s.gonTableRow(h, c+1, gon, func(id string, h *dom.Html) {
				buttonLink(h, fmt.Sprintf(link, s.S(), id), id)
			})
		}
	})
	return nil
}

func getGon(r *http.Request, s *S, h *dom.Html, gg symm.Gons, back, title, link string) error {
	// back button
	h.Div(nil, func(h *dom.Html) {
		buttonLink(h, fmt.Sprintf(back, s.S()), "<")
	})
	sids, ids, err := idAngles(r) // read URL param "id"
	if err != nil {
		return err
	}
	// title including ids
	h.Div(h1, fmt.Sprintf(title, s.S(), sids))
	t, err := gg.Transforms(ids)
	if err != nil {
		return err
	}
	shift, vector := btnsShiftVector(r, h, t, fmt.Sprintf(link, s.S(), sids))
	if gon, err := gg.New(t, shift, vector); err != nil {
		return err
	} else {
		s.getGon(h, gon)
	}
	return nil
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

func btnsShiftVector(r *http.Request, h *dom.Html, t *symm.Transforms, link string) (int,int) {
	// 
	button := func(shift, vector int, text string) {
		link := fmt.Sprintf("%s?s=%d&v=%d", link, shift, vector)
		buttonLink(h, link, text)
	}
	shift := 1
	s := r.URL.Query().Get("s")
	if s, err := strconv.Atoi(s); err == nil {
		shift = s
	}
	vector := 1
	v := r.URL.Query().Get("v")
	if v, err := strconv.Atoi(v); err == nil {
		vector = v
	}
	h.Div(nil, func(h *dom.Html) {
		h.Elem(dom.Span, nil, "Shift: ")
		for _, s := range t.Shifts() {
			if s == shift {
				h.Elem(dom.Span, h1, fmt.Sprintf("%d", s))
			} else {
				button(s, vector, fmt.Sprintf("%d",s))
			}
		}
	})
	h.Div(nil, func(h *dom.Html) {
		h.Elem(dom.Span, nil, "Vector: ")
		for _, v := range t.Vectors() {
			if v == vector {
				h.Elem(dom.Span, h1, fmt.Sprintf("%d", v))
			} else {
				button(shift, v, fmt.Sprintf("%d",v))
			}
		}
	})
	return shift, vector
}


// deprecate
func buttonsShiftVector(r *http.Request, symm int, h *dom.Html, prelink string) (int, int) {
	return shiftVectorOpts(r, symm, h, func(shift,vector int, text string) {
		link := fmt.Sprintf("%s?shift=%d&vector=%d", prelink, shift, vector)
		buttonLink(h, link, text)
	})
}

// deprecate
func shiftVectorOpts(r *http.Request, symm int, h *dom.Html, option func(s, v int,text string)) (int,int) {
	shift := 1
	s := r.URL.Query().Get("shift")
	if s, err := strconv.Atoi(s); err == nil {
		shift = s
	}
	vector := 1
	v := r.URL.Query().Get("vector")
	if v, err := strconv.Atoi(v); err == nil && v > 0 && v <= symm {
		vector = v
	}
	h.Div(nil, func(h *dom.Html) {
		h.Elem(dom.Span, nil, "Shift: ")
		for v := 1; v <= symm; v++ {
			if v == vector {
				h.Elem(dom.Span, h1, fmt.Sprintf("%d", v))
			} else {
				text := fmt.Sprintf("%d",v)
				option(shift, v, text)
			}
		}
	})
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

