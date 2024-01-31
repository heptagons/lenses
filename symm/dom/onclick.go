package dom

import (
	"fmt"
	"sort"
	"strings"
)

func NewOnclickAjax(uri string, id Id, args map[string]string) Onclick {
	var sb strings.Builder
	sb.WriteString(`a0(`)
	sb.WriteString(`'` + uri + `',`)
	sb.WriteString(`'` + string(id) + `',{`)
	first := true
	for k, v := range args {
		if !first {
			sb.WriteString(`,`)
		}
		sb.WriteString(`'`)
		sb.WriteString(k)
		sb.WriteString(`':'`)
		sb.WriteString(v)
		sb.WriteString(`'`)
		first = false
	}
	sb.WriteString(`})`)
	return Onclick(sb.String())
}

func NewOnclickAjaxValue(uri string, id, value Id, ons, offs []Class) Onclick {
	var sb strings.Builder
	sb.WriteString(`a1(`)
	sb.WriteString(`'` + uri + `',`)
	sb.WriteString(`'` + string(id) + `',`)
	sb.WriteString(`'` + string(value) + `',`)
	onclickSets(ons, offs, nil, &sb)
	sb.WriteString(`)`)
	return Onclick(sb.String())
}

func NewOnclickAjaxValues(uri string, id, base Id, n int) Onclick {
	return NewOnclick("a2('%s','%s','%s',%d)", uri, id, base, n)
}


type Onclick string

func NewOnclick(f string, args ...interface{}) Onclick {
	return Onclick(fmt.Sprintf(f, args...))
}

type OnclickSets struct {
	// Inner is the address of the text to update in element innerHTML
	html  *string
	// Value is the address of the value to update into data-value attribute
	value *string
	// Class is the address of the classname to update into class attribute
	class *Class

	input *Id
}

type OnclickSetsMap map[Id]OnclickSets

var empty = ""

func NewOnclickValsInnerEmpty() OnclickSets {
	return OnclickSets{
		html: &empty,
	}
}

func NewOnclickValsDirect(html string, value string, class Class) OnclickSets {
	return OnclickSets{
		html:  &html,
		value: &value,
		class: &class,
	}
}

func NewOnclickValsInput(input Id, class Class) OnclickSets {
	return OnclickSets{
		input: &input,
		class: &class,
	}
}

func (o *OnclickSets) write(sb *strings.Builder) {
	comma := false
	pair := func(k, v string) {
		if comma {
			sb.WriteString(`,`)
		}
		comma = true
		sb.WriteString(`'`)
		sb.WriteString(k)
		sb.WriteString(`':'`)
		sb.WriteString(v)
		sb.WriteString(`'`)
	}
	if o.html != nil {
		pair("h", *o.html)
	}
	if o.value != nil {
		pair("v", *o.value)
	}
	if o.class != nil {
		pair("c", string(*o.class))
	}
	if o.input != nil {
		pair("i", string(*o.input))
	}
}

// calls javascript s0()
func NewOnclickSet(id Id, f string, args ...interface{}) Onclick {
	return Onclick(fmt.Sprintf("s0('%s','%s')", string(id), fmt.Sprintf(f, args...)))
}

// calls javascript s1()
func NewOnclickSetValue(id Id, f string, args ...interface{}) Onclick {
	var sb strings.Builder
	sb.WriteString("s1(")
	sb.WriteString(`'`)
	sb.WriteString(string(id))
	sb.WriteString(`','`)
	sb.WriteString(fmt.Sprintf(f, args...))	
	sb.WriteString(`'`)
	sb.WriteString(`)`)
	return Onclick(sb.String())
}

// calls javascript s2()
func NewOnclickSets(ons, offs []Class, set map[Id]OnclickSets) Onclick {
	var sb strings.Builder
	sb.WriteString("s2(")
	onclickSets(ons, offs, set, &sb)
	sb.WriteString(`)`) // pass the element with the "onclick"
	return Onclick(sb.String())
}

func onclickSets(ons, offs []Class, set map[Id]OnclickSets, sb *strings.Builder) {
	onclickClasses(ons, sb)
	sb.WriteString(",")
	onclickClasses(offs, sb)
	sb.WriteString(",{")
	comma1 := ""
	ids := make([]string, 0, len(set))
	for i := range set {
		ids = append(ids, string(i))
	}
	sort.Strings(ids)
	for _, id := range ids {
		ss := set[Id(id)]
		sb.WriteString(comma1)
		sb.WriteString("'")
		sb.WriteString(id)
		sb.WriteString("':{")
		ss.write(sb)
		sb.WriteString(`}`)
		comma1 = `,`
	}
	sb.WriteString(`}`)
}

func onclickClasses(cs []Class, sb *strings.Builder) {
	sb.WriteString(`[`)
	for i, c := range cs {
		if i > 0 {
			sb.WriteString(`,`)
		}
		sb.WriteString(`'`)
		sb.WriteString(string(c))
		sb.WriteString(`'`)
	}
	sb.WriteString(`]`)
}

func (o Onclick) WriteAttr(h *Html) {
	h.write(` onclick="`, string(o), `"`)
}





