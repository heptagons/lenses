package symm

type Transforms struct {
	id      string
	angles  []int
	group   *Group
	shifts  []int
	vectors []int
}

func (t *Transforms) Group() *Group {
	return t.group
}

