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

func TestPolylines(t *testing.T) {

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

	p, _ = pp.NewWithAngles(1, []int{ 4, 2, 3, 4, 2 })
	t.Logf("3) vectors:%v", p.vectors)
	t.Logf("3) angles:%v", p.Angles()) 
	for _, accum := range p.Accums() {
		t.Logf("3) accum:%v", accum)
	}
}

func TestPolyStars(t *testing.T) {
	s3, _ := NewSymm(3)
	p3 := NewPolylines(s3)
	for s, star := range [][]int { // arrays of size 10-1=9
		[]int{ 1,1,1,1,1 }, // S_3(1) = "A" = regular hexagon
	} {
		testPolyStars(t, p3, s, star)
	}
	s5, _ := NewSymm(5)
	p5 := NewPolylines(s5)
	for s, star := range [][]int { // arrays of size 10-1=9
		[]int{ 3,1,3,1,3,1,3,1,3 }, // S_5(1) = "B"
		[]int{ 2,2,2,2,2,2,2,2,2 }, // S_5(2) = "C" = regular 10-gon 
	} {
		testPolyStars(t, p5, s, star)
	}
	s7, _ := NewSymm(7)
	p7 := NewPolylines(s7)
	for s, star := range [][]int { // arrays of size 14-1=13
		[]int{ 5,1,5,1,5,1,5,1,5,1,5,1,5 }, // S_7(1) = "D"
		[]int{ 4,2,4,2,4,2,4,2,4,2,4,2,4 }, // S_7(2) = "E"
		[]int{ 3,3,3,3,3,3,3,3,3,3,3,3,3 }, // S_7(3) = "F" = regular 14-gon 
	} {
		testPolyStars(t, p7, s, star)
	}
	s9, _ := NewSymm(9)
	p9 := NewPolylines(s9)
	for s, star := range [][]int{ // arrays size 18-1=17
		[]int{ 7,1,7,1,7,1,7,1,7,1,7,1,7,1,7,1,7 }, // S_9(1) = "G"
		[]int{ 6,2,6,2,6,2,6,2,6,2,6,2,6,2,6,2,6 }, // S_9(2) = "H"
		[]int{ 5,3,5,3,5,3,5,3,5,3,5,3,5,3,5,3,5 }, // S_9(3) = "I"
		[]int{ 4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4 }, // S_9(4) = "J" = regular 18-gon 
	} {
		testPolyStars(t, p9, s, star)
	}
}

func testPolyStars(t *testing.T, pp *Polylines, s int, star []int) {
	if p, err := pp.NewWithAngles(1, star); err != nil {
		t.Fatalf("star %v error:%v", star, err)
	} else if angles := p.Angles(); !reflect.DeepEqual(star, angles) {
		t.Fatalf("star expected %v got angles %v", star, angles)
	} else {
		t.Logf("star %d) vectors:%v", s, p.vectors)
		t.Logf(" angles:%v", p.Angles()) 
		accums := p.Accums()
		for _, accum := range accums {
			t.Logf(" accum:%v", accum)
		}
		last := accums[len(accums)-1]
		for pos, x := range last.x {
			if x != 0 {
				t.Fatalf("x[%d] not zero: %d", pos, x)
			}
		}
		for pos, y := range last.y {
			if y != 0 {
				t.Fatalf("y[%d] not zero: %d", pos, y)
			}
		}
	}
}

func TestHexa(t *testing.T) {
	s, _ := NewSymm(9)
	p := NewPolylines(s)
	hh := NewHexagons(p)
	if _, err := hh.New(1, nil); err == nil {
		t.Fatalf("Accepted no angles")
	} else if _, err := hh.New(1, []int{1,1,1,1}); err == nil {
		t.Fatalf("Accepted four angles")
	} else if _, err := hh.New(1, []int{1,1,1}); err == nil {
		t.Fatalf("Accepted not closed hexagon")
	}

	min := 1
	max := s.s - 1
	//for a := min; a <= max; a++ {
	//	if h, err := hh.New(1, []int { a }); err == nil {
	//		t.Logf("Hexagon a angles=%v", h.p.Angles())
	//	}
	//}
	for a := min; a <= max; a++ {
		for b := a; b <= max; b++ {
			if h, err := hh.New(1, []int{ a, b }); err == nil {
				t.Logf("Hexagon a,b v=%v, a=%v", h.p.vectors, h.p.Angles())
			} else {
				//t.Logf("error %v %v", ab, err)
			}
		}
	}
	for a := 1; a <= max; a++ {
		for b := a; b <= max; b++ {
			for c := b; c <=max; c++ {
				if a == b && b == c {
					// case 1
				} else if h, err := hh.New(1, []int{ a, b, c }); err == nil {
					t.Logf("Hexagon a,b,c v=%v a=%v", h.p.vectors, h.p.Angles())
				} else {
					//t.Logf("error %v %v", ab, err)
				}
			}
		}
	}
}