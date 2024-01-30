package symm

import (
	"testing"
	"reflect"
)

func TestSymm(t *testing.T) {
	if _, err := NewSymm(-3); err == nil {
		t.Fatalf("Accepted symm -3")
	} else if _, err := NewSymm(1); err == nil {
		t.Fatalf("Accepted symm 1")
	} else if _, err := NewSymm(4); err == nil {
		t.Fatalf("Accepted symm 4")
	} else if s, err := NewSymm(9); err != nil {
		t.Fatalf("err:%v", err)
	} else if s.s != 9 {
		t.Fatalf("Expected symmetry 9 got %d", s.s)
	} else {
		t.Logf("cos %v", s.x)
		t.Logf("sin %v", s.y)
		t.Logf("+v  %v", s.v)
		t.Logf("-v  %v", s.w)
	}
}

func TestPoly(t *testing.T) {

	vectors := []int{1,6,4}
	angles := []int{4,2}
	accums := []*Accum{
		&Accum{ x:[]int{ 1,0,0,0,0  }, y:[]int{ 1,0,0,0,0 } },
		&Accum{ x:[]int{ 1,0,0,0,-1 }, y:[]int{ 1,0,0,0,1 } },
		&Accum{ x:[]int{ 1,0,0,1,-1 }, y:[]int{ 1,0,0,1,1 } },
	}
	xys := [][]float64 {
		[]float64{ 1, 0 },
		[]float64{ 1.9396926207859084, 0.3420201433256689 },
		[]float64{ 1.4396926207859084, 1.2080455471101077 },
	}
	
	s, _ := NewSymm(9)
	pp := NewPolylines(s)
	if _, err := pp.New(-1); err == nil {
		t.Fatalf("Accepted negative vector: %v", err)
	} else if _, err := pp.New(10); err == nil {
		t.Fatalf("Accepted out of range vector: %v", err)
	} else if p, err := pp.New(1,6,4); err != nil {
		t.Fatalf("vectors 9[1,6,4] error:%v", err)
	} else if !reflect.DeepEqual(vectors, p.vectors){
		t.Fatalf("vectors expected [1,6,4] got %v", p.vectors)
	} else if pa := p.Angles(); !reflect.DeepEqual(angles, pa) {
		t.Fatalf("angles: expected %v got %v", angles, pa)
	} else if aa := p.Accums(); !reflect.DeepEqual(accums, aa) {
		t.Fatalf("accums expected %v got %v", accums, aa)
	} else {
		for i := 0; i < 3; i++ {
			if xy := s.XY(aa[i]); !reflect.DeepEqual(xys[i], xy) {
				t.Fatalf("xy expected %v got %v", xys[i], xy)
			}
		}
	}


	p, _ := pp.New(1,6,4,1,6,4,1)
	t.Logf("2) vectors:%v", p.vectors)
	t.Logf("2) angles:%v", p.Angles()) // [4 2 3 4 2 3]

	p, _ = pp.NewWithAngles(1, []int{ 4, 2, 3, 4, 2, 3 })
	t.Logf("3) vectors:%v", p.vectors)
	t.Logf("3) angles:%v", p.Angles()) 
	for _, accum := range p.Accums() {
		t.Logf("3) accum:%v", accum)
	}

	p, _ = pp.NewWithAngles(1, []int{ 4,4,4,4,4,4,4,4,4,4,4 })
	t.Logf("4) vectors:%v", p.vectors)
	t.Logf("4) angles:%v", p.Angles()) 
	for _, accum := range p.Accums() {
		t.Logf("4) accum:%v", accum)
	}
}