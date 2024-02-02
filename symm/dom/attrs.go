package dom

import (
	"fmt"
)

type Attr interface {
	WriteAttr(h *Html)
}

type Class string

func NewClass(f string, args ...interface{}) Class {
	return Class(fmt.Sprintf(f, args...))
}

func (c *Class) Append(next Class) {
	// Use pointer to update this!
	*c = NewClass("%s %s", *c, next)
}

func (c Class) WriteAttr(h *Html) {
	h.write(` class='`, string(c), `'`)
}

type Id string

func NewId(f string, args ...interface{}) Id {
	return Id(fmt.Sprintf(f, args...))
}

func (i Id) WriteAttr(h *Html) {
	h.write(` id='`, string(i), `'`)
}


// A DataValue is a DOM attribute of the type data-value='XXX'
type DataValue string

// NewDataValue returns a DataValue attribute
func NewDataValue(v string) DataValue {
	return DataValue(v)
}

func (d DataValue) WriteAttr(h *Html) {
	h.write(` data-value='`, string(d), `'`)
}


type StyleDisplay string

func NewStyleDisplay(display bool) StyleDisplay {
	if display {
		return StyleDisplay("block")
	} else {
		return StyleDisplay("none")
	}
}

func (s StyleDisplay) WriteAttr(h *Html) {
	h.write(` style='display:`, string(s), `;'`)
}


type Colspan string

func NewColspan(colspan int) Colspan {
	return Colspan(fmt.Sprintf("%d", colspan))
}

func (c Colspan) WriteAttr(h *Html) {
	h.write(` colspan="`, string(c), `"`)
}

type Number struct {
	Min   string
	Max   string
	Step  string
	Value string
	Width string
}

func (n Number) WriteAttr(h *Html) {
	h.write(` type="number"`)
	if n.Min != "" {
		h.write(` min="`, n.Min, `"`)
	}
	if n.Max != "" {
		h.write(` max="`, n.Max, `"`)
	}
	if n.Step != "" {
		h.write(` step="`, n.Step, `"`)
	}
	if n.Value != "" {
		h.write(` value="`, n.Value, `"`)
	}
	if n.Width != "" {
		h.write(` style="width:`, n.Width, `"`)
	}
}

type Checkbox bool

// implements dom.Attr
func (c Checkbox) WriteAttr(h *Html) {
	h.WriteF(` type="checkbox"`)
	if c {
		h.WriteF(` checked="checked"`)
	}
}

type Title string

// implements dom.Attr
func (t Title) WriteAttr(h *Html) {
	h.WriteF(` title="%s"`, t)
}

type Value string

// implements dom.Attr
func (v Value) WriteAttr(h *Html) {
	h.WriteF(` value="%s"`, v)
}

// Svg are attributes used within elem "svg"
type Asvg struct {
	width   int
	height  int
	viewBox []int // viewBox
}

// implements dom.Attr
func (a Asvg) WriteAttr(h *Html) {
	h.WriteF(` width="%dpx"`, a.width)
	h.WriteF(` height="%dpx"`, a.height)
	if v := a.viewBox; len(v) == 4 {
		h.WriteF(` viewBox="%d %d %d %d"`, v[0], v[1], v[2], v[3])
	}
	h.WriteF(` xmlns="http://www.w3.org/2000/svg"`)
}

// Polygon are attributes used within elem "polygon"
type Apolygon struct {
	points [][]float64
	fill   string
	stroke string
}

// implements dom.Attr
func (a Apolygon) WriteAttr(h *Html) {
	if p := a.points; len(p) > 0 {
		h.WriteF(` points="`)
		for c, pair := range p {
			if c > 0 {
				h.WriteF(` `)
			}
			if len(pair) == 2 {
				h.WriteF(`%.3f,%.3f`, pair[0], pair[1])
			}
		}
		h.WriteF(`"`)
	}
	if f := a.fill; f != "" {
		h.WriteF(` fill="%s"`, f)
	}
	if f := a.stroke; f != "" {
		h.WriteF(` stroke="%s"`, f)
	}
}


