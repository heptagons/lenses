package dom

import (
	"fmt"
	"strings"
)

type Dom struct {
	i Id   // id
	h *Html // innerHTML
}

func NewIH(i Id, h *Html) *Dom {
	return &Dom{
		i: i,
		h: h,
	}
}

type Doms []*Dom

func (doms *Doms) Json() string {
	var sb strings.Builder
	sb.WriteString(`{"doms":[`)
	comma := ""
	for _, d := range *doms {
		sb.WriteString(comma)
		sb.WriteString(`{"i":"`)
		sb.WriteString(string(d.i))
		sb.WriteString(`","h":"`)
		if d.h == nil {
			sb.WriteString("&nbsp;")
		} else {
			h := strings.ReplaceAll(d.h.sb.String(), `"`, `\"`)
			// h must not contain \n !!!
			sb.WriteString(h)
		}
		sb.WriteString(`"}`)
		comma = `,`
	}
	sb.WriteString(`]}`)
	return sb.String()
}

type elem string

const (
	Body    elem = "body"
	Br      elem = "br"
	Button  elem = "button"
	Caption elem = "caption"
	Div     elem = "div"
	Header  elem = "header"
	Html_   elem = "html"
	Hr      elem = "hr"
	Head    elem = "head"
	Input   elem = "input"
	Label   elem = "label"
	Option  elem = "option"
	P       elem = "p"
	Script  elem = "script"
	Select  elem = "select"
	Span    elem = "span"
	Style   elem = "style"
	Table   elem = "table"
	Td      elem = "td"
	Th      elem = "th"
	Title_  elem = "title"
	Tr      elem = "tr"
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




