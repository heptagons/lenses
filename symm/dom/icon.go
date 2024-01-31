package dom

import (
	"fmt"
	"strings"
)

type Icon struct {
	pos  int
	path string
}

func NewIcon(pos int, path string) *Icon {
	return &Icon{
		pos:  pos,
		path: path,
	}
}

// Svg returns the icon path of size 15 and without fill set.
// The fill color should be set on parent DOM element as in a button class.
func (i *Icon) Svg() string {
	return i.svg(15, nil)
}

// path returns the svg DOM text to represent the icon.
// size is commonly set to 15. fill is nil for Path
func (i *Icon) svg(size int, fill []byte) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d"`, size))
	sb.WriteString(fmt.Sprintf(` height="%d"`, size))
	sb.WriteString(fmt.Sprintf(` viewBox="0 0 %d %d">`, size, size))
	sb.WriteString(fmt.Sprintf(`<path d="%s"`, i.path))
	if len(fill) == 3 {
		sb.WriteString(fmt.Sprintf(` fill="rgb(%d,%d,%d)"`, fill[0], fill[1], fill[2]))
	}
	sb.WriteString(`/>`)
	sb.WriteString(`</svg>`)
	s := sb.String()
	return s
}

func (i *Icon) style(size int, h *Html, fill []byte) {
	const image = "\tbackground-image:url('data:image/svg+xml;charset=utf-8,%s');\n"
	h.WriteF(".icon-%d {\n", i.pos)
	h.WriteF(image, i.svg(size, fill))
	h.WriteF("}\n")
}

func (i *Icon) Html(h *Html) {
	h.WriteF(`<div class='icon0'><div class="icon icon-%d"></div></div>`, i.pos)
}

type Icons map[int]*Icon

func NewIcons() Icons {
	return make(map[int]*Icon)
}

func (i *Icons) Add(icon *Icon) {
	if icon == nil {
		return
	}
	(*i)[icon.pos] = icon
}

// Style returns all icons styles as background images
// use min-width instead width in order this icon inside expands properlly:
// <td style="text-align:center;"><div class=".icon"/>...
func (i *Icons) Style(h *Html, fill []byte) {
	h.WriteF(`
.icon {
	background-position: center;
	background-repeat: no-repeat;
	min-height: 24px;
	min-width: 24px;
}
.icon0 {
	display: inline-block;
	padding: 0px;
	margin: 2px;
}
`,	
	)
	for _, icon := range (*i) {
		icon.style(15, h, fill)
	}
}
