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

func (pp *Polylines) String() string {
	return fmt.Sprintf("{s=%d v=%v}", pp.s.s, pp.vectors)
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
	p := NewPolylineN(pp, len(angles)+1)
	p.SetAngles(vector, angles)
	return p, nil
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
	pp    *Polylines
	edges []int
}

func NewPolyline(pp *Polylines, edges []int) *Polyline {
	return &Polyline{
		pp:      pp,
		edges: edges,
	}
}

func NewPolylineN(pp *Polylines, numEdges int) *Polyline {
	return &Polyline{
		pp:    pp,
		edges: make([]int, numEdges),
	}
}

func (p *Polyline) String() string {
	return fmt.Sprintf("p=%v e=%v", p.pp, p.edges)
}

func (p *Polyline) Edges() []int {
	return p.edges
}

// SetAngles updates this polyline edges. First edge is set as given vector
// The rest of edges is computed according the given angles
// Returns an error if number of edges + 1 don't equal the number of angles.
func (p *Polyline) SetAngles(vector int, angles []int) error {
	n := len(angles) + 1
	if n != len(p.edges) {
		return fmt.Errorf("number of angles + 1 != number of edges")
	}
	//edges := make([]int, n)
	s := p.pp.s.s
	p.edges[0] = vector
	for i := 1; i < n; i++ {
		m := p.edges[i-1]
		a := angles[i-1]
		n := (s + m - a) % s
		if n == 0 {
			n = s // TODO document
		}
		p.edges[i] = n
	}
	return nil
}

// Angles return the angles of this polyline reading the internal edges
func (p *Polyline) Angles() []int {
	if n :=  len(p.edges); n < 2 {
		return nil
	} else {
		angles := make([]int, n-1)
		s := p.pp.s.s
		for i := 1; i < n; i++ {
			m := p.edges[i-1]
			n := p.edges[i]
			u := (s + m - n) % s
			angles[i-1] = u
		}
		return angles
	}
}

// Accums return the accumulators of the vertices of this polyline
// reading the internal edges.
func (p *Polyline) Accums() []*Accum {
	n := len(p.edges)
	t := p.pp.s.t
	base := NewAccum(t)
	pos := 0
	accum := 0
	accums := make([]*Accum, n)
	var indices []int	
	for i := 0; i < n; i++ {
		vindex := p.edges[i] 
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

