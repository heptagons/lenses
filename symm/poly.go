package symm

import (
	"fmt"
	"strconv"
	"strings"
)

type Polylines struct {
	s       *Symm
	vectors []int
}

func NewPolylines(s *Symm) *Polylines {
	vectors := make([]int, s.s)
	for pos := range vectors {
		vectors[pos] = pos + 1 // 1,2,3,...symm
	}
	return &Polylines{
		s:       s,
		vectors: vectors, 
	}
}

func (pp *Polylines) New(edges ...int) (*Polyline, error) {
	for v := 0; v < len(edges); v++ {
		if edges[v] < 1 {
			return nil, fmt.Errorf("Invalid edge vector %v at position %v", edges[v], v)
		} else if edges[v] > pp.s.s {
			return nil, fmt.Errorf("Invalid edge vector %v at position %v", edges[v], v)
		}
	}
	return NewPolyline(pp, edges), nil
}

func (pp *Polylines) NewWithAngles(vector int, angles []int) (*Polyline, error) {
	s := pp.s.s
	if vector < 1 {
		return nil, fmt.Errorf("Invalid vector %v", vector)
	} else if vector > s {
		return nil, fmt.Errorf("Invalid vector %v", vector)
	}
	for pos, angle := range angles {
		if angle < 1 || angle > s {
			return nil, fmt.Errorf("Invalid angle %v at position %v", angle, pos)
		}
	}
	return NewPolylineWithAngles(pp, vector, angles), nil
}

// IdFromAngles returns a string of angles array values separated by commas
func (pp *Polylines) IdFromAngles(angles []int) string {
	var ids []string
	for _, angle := range angles {
		ids = append(ids, strconv.Itoa(angle))
	}
	return strings.Join(ids, ",")
}

type Polyline struct {
	pp      *Polylines
	vectors []int
}

func NewPolyline(pp *Polylines, vectors []int) *Polyline {
	return &Polyline{
		pp:      pp,
		vectors: vectors,
	}
}

func NewPolylineWithAngles(pp *Polylines, vector int, angles []int) *Polyline {
	s := pp.s.s
	n := len(angles) + 1
	vectors := make([]int, n)
	vectors[0] = vector
	for i := 1; i < n; i++ {
		m := vectors[i-1]
		a := angles[i-1]
		n := (s + m - a) % s
		if n == 0 {
			n = s // TODO document
		}

		//fmt.Println("n", n)
		vectors[i] = n
	}
	return &Polyline{
		pp:      pp,
		vectors: vectors,
	}
}

func (p *Polyline) Angles() []int {
	if n :=  len(p.vectors); n < 2 {
		return nil
	} else {
		angles := make([]int, n-1)
		s := p.pp.s.s
		for i := 1; i < n; i++ {
			m := p.vectors[i-1]
			n := p.vectors[i]
			u := (s + m - n) % s
			angles[i-1] = u
		}
		return angles
	}
}

func (p *Polyline) Accums() []*Accum {
	n := len(p.vectors)
	t := p.pp.s.t
	base := NewAccum(t)
	pos := 0
	accum := 0
	accums := make([]*Accum, n)
	var indices []int	
	for i := 0; i < n; i++ {
		vindex := p.vectors[i] 
		if i % 2 == 0 {
			// even elements are taken as normal vector (v)
			indices = p.pp.s.v[vindex-1]
		} else {
			// odd elements are taken as vector rotated 180° (w)
			indices = p.pp.s.w[vindex-1]
		}
		//fmt.Printf("vindex=%d indices=%v\n", vindex, indices)
		// X array
		if x := indices[0]; x < 0 {
			pos = -x - 1
			accum = -1
		} else {
			pos = +x - 1
			accum = +1
		}
		if pos < t {
			base.x[pos] += accum
		}
		// Y array		
		if y := indices[1]; y < 0 {
			pos = -y - 1
			accum = -1
		} else {
			pos = +y - 1
			accum = +1
		}
		if pos < t {
			base.y[pos] += accum
		}
		accums[i] = base.Clone()
	}
	return accums
}

