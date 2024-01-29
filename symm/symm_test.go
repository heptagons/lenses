package symm

import (
	"testing"
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
	} else {
		t.Logf("9[1,6,4] -> %v", p.vectors)
	}
}