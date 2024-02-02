package dom

import (
	"fmt"
	"io"
	"strings"
)

type F func(h *Html)

type Html struct {
	Dark bool	
	sb   strings.Builder
}

func NewHtml(parent *Html, t string, args ...interface{}) *Html {
	h := &Html{}
	if parent != nil {
		h.Dark = parent.Dark
	}
	h.sb.WriteString(fmt.Sprintf(t, args...))
	return h
}

func (h *Html) Elem(e elem, attrs []Attr, inner interface{}) {
	if e != "" {
		h.write("<", string(e))
		for _, a := range attrs {
			a.WriteAttr(h)
		}
		h.write(">")
	}
	h.inner(inner)
	if e != "" {
		h.write("</", string(e), ">")
	}
}

func (h *Html) Div(attrs []Attr, inner interface{}) {
	h.Elem(Div, attrs, inner)
}

func (h *Html) Button(attrs []Attr, inner interface{}) {
	h.Elem(Button, attrs, inner)
}

func (h *Html) Svg(width, height int, viewBox []int, inner interface{}) {
	h.Elem(Svg, []Attr{
		&Asvg{
			width:   width,
			height:  height,
			viewBox: viewBox,
		},
	}, inner)
}

func (h *Html) Polygon(points [][]float64, fill, stroke string) {
	h.Elem(Polygon, []Attr{
		&Apolygon{
			points: points,
			fill:   fill,
			stroke: stroke,
		},
	}, nil)
}

func (h *Html) WriteF(f string, args ...interface{}) {
	h.sb.WriteString(fmt.Sprintf(f, args...))
}

func (h *Html) Size() int {
	return len(h.sb.String())
}

func (h *Html) Write(w io.Writer) {
	w.Write([]byte(h.sb.String()))
}

func (h *Html) write(args ...string) {
	for _, arg := range args {
		h.sb.WriteString(arg)
	}
}

func (h *Html) inner(inner interface{}) {
	if inner != nil {
		switch inner.(type) {

		case []F:
			for _, f := range inner.([]F) { f(h) }

		case func(*Html):
			inner.(func(*Html))(h)

		case string:
			h.sb.WriteString(inner.(string))

		default:
			panic("inner invalid")
		}
	}
}




