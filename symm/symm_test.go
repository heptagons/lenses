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
	s, _ := NewSymm(9)
	pp := NewPolylines(s)
	if _, err := pp.New(-1); err == nil {
		t.Fatalf("Accepted negative vector: %v", err)
	} else if _, err := pp.New(10); err == nil {
		t.Fatalf("Accepted out of range vector: %v", err)
	} else if p, err := pp.New(1,6,4); err != nil {
		t.Fatalf("vectors 9[1,6,4] error:%v", err)
	} else if !reflect.DeepEqual([]int{1,6,4}, p.vectors){
		t.Fatalf("vectors expected [1,6,4] got %v", p.vectors)
	} else if angles := p.Angles(); !reflect.DeepEqual([]int{4,2}, angles) {
		t.Fatalf("angles expected [4,2] got %v", angles)
	} else {
		accums := p.Accums()
		for _, a := range accums {
			t.Logf("accum: %v", a)
		}
		x, y := s.XY(accums[2])
		t.Logf("Third point x=%f y=%f", x, y)
	}
}