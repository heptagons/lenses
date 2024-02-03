package dom

import (
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
	Path    elem = "path"
	Polygon elem = "polygon"
	Script  elem = "script"
	Select  elem = "select"
	Span    elem = "span"
	Style   elem = "style"
	Svg     elem = "svg"
	Table   elem = "table"
	Td      elem = "td"
	Th      elem = "th"
	Title_  elem = "title"
	Tr      elem = "tr"
)


