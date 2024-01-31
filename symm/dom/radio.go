package dom

type Radio struct {
	class   Class
	display bool
	text    string
}

func NewRadio(class Class, display bool, text string) *Radio {
	return &Radio{
		class:   class,
		display: display,
		text:    text,
	}
}

func (r *Radio) Content(h *Html, inner func(h *Html)) {
	h.Div([]Attr{
		r.class,
		NewStyleDisplay(r.display),
	}, inner)
}

func RadioButtons(class Class, radios []*Radio, h *Html) {
	if len(radios) == 1 {
		return
	}
	h.Elem(Table, nil, func (h *Html) {
		h.Elem(Tr, nil, func (h *Html) {
			for i, r := range radios {
				ons := []Class { r.class }
				offs := make([]Class, 0)
				for j, r := range radios {
					if i != j {
						offs = append(offs, r.class)
					}
				}
				h.Elem(Td, []Attr{ class }, func (h *Html) {
					h.Button([]Attr{
						NewOnclickSets(ons, offs, nil),
					}, r.text)
					h.Elem(Hr, []Attr{
						r.class,
						NewStyleDisplay(r.display),			
					}, nil)
				})
			}
		})
	})
}





