package symm

import (
	//"fmt"
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

	edges := []int{1,6,4}
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
		t.Fatalf("edges 9[1,6,4] error:%v", err)
	} else if !reflect.DeepEqual(edges, p.edges){
		t.Fatalf("edges expected [1,6,4] got %v", p.edges)
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
	t.Logf("2) edges:%v", p.edges)
	t.Logf("2) angles:%v", p.Angles()) // [4 2 3 4 2 3]

	p, _ = pp.NewWithAngles(1, []int{ 4, 2, 3, 4, 2 })
	t.Logf("3) edges:%v", p.edges)
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
		t.Logf("star %d) edges:%v", s, p.edges)
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

/*func TestHexa(t *testing.T) {

	for _, symm := range []int{ 3, 5, 7, 9, 11, 13, 15 } {
		s, _ := NewSymm(symm)
		p := NewPolylines(s)
		hh := NewHexagons(p)
		t.Logf("Symmetry %d", symm)
		min := 1
		max := s.s - 1
		for a := min; a <= max; a++ {
			for b := a; b <= max; b++ {

				if h, err := hh.New(1, []int{ a, b }); err == nil {
					t.Logf("a,b   a=%v v=%v", h.p.Angles(), h.p.vectors)
				} else {
					//t.Logf("error %v %v", ab, err)
				}

				for c := b; c <=max; c++ {
					if a == b && b == c {
						// case 1
					} else if h, err := hh.New(1, []int{ a, b, c }); err == nil {
						t.Logf("a,b,c a=%v v=%v", h.p.Angles(), h.p.vectors)
					} else {
						//t.Logf("error %v %v", ab, err)
					}
				}
			}
		}
	}
}*/

/*func TestOcta(t *testing.T) {
	symm := 7
	s, _ := NewSymm(symm)
	p := NewPolylines(s)
	oo := NewOctagons(p)
	t.Logf("Symmetry %d", symm)
	for c, o := range oo.All() {
		t.Logf("%d %s", c+1, o.String())
	}
}*/

/*func TestHexaAll2(t *testing.T) {
	symm := 37
	s, _ := NewSymm(symm)
	p := NewPolylines(s)
	hh := NewHexagons(p)
	t.Logf("Symmetry %d", symm)
	hh.All2() // prints...
}*/

/*func TestOctaAll5(t *testing.T) {
	symm := 9
	s, _ := NewSymm(symm)
	p := NewPolylines(s)
	oo := NewOctagons(p)
	t.Logf("Symmetry %d", symm)
	for c, o := range oo.all5() {
		t.Logf("%d %v", c+1, o.Angles())
	}
}*/

func TestHexaReduce(t *testing.T) {
	for _, io := range [][]uint64 {
		[]uint64{ 0x0123456789AB, 0x0123456789AB },
		[]uint64{ 0xBA9876543210, 0x1032547698BA },
		[]uint64{ 0x000002010000, 0x000000000102 },
		[]uint64{ 0x000002010000, 0x000000000102 },
		[]uint64{ 0x010101010000, 0x000001010101 },
	} {
		in := io[0]
		exp := io[1]
		if got := hexaReduce(in); got != exp {
			t.Fatalf("in: %016x out: exp %016x got %016x", in, exp, got)
		}
	}
}

func TestHexaAll6(t *testing.T) {
	symm := 19
	s, _ := NewSymm(symm)
	p := NewPolylines(s)
	hh := NewHexagons(p)
	t.Logf("Symmetry %d", symm)
	for c, gon := range hh.all6() {
		t.Logf("%d) %v", c+1, gon.Transforms())
	}
}


func TestOctaReduce(t *testing.T) {
	for _, io := range [][]uint64 {
		[]uint64{ 0x0123456789ABCDEF, 0x0123456789ABCDEF },
		[]uint64{ 0xFEDCBA9876543210, 0x1032547698BADCFE },
		[]uint64{ 0x0000000002010000, 0x0000000000000102 },
		[]uint64{ 0x0000000002010000, 0x0000000000000102 },
		[]uint64{ 0x0101010101010100, 0x0001010101010101 },
	} {
		in := io[0]
		exp := io[1]
		if got := octaReduce(in); got != exp {
			t.Fatalf("in: %016x out: exp %016x got %016x", in, exp, got)
		}
	}
}

func TestOctaAll8(t *testing.T) {
	symm := 7
	s, _ := NewSymm(symm)
	p := NewPolylines(s)
	oo := NewOctagons(p)
	t.Logf("Symmetry %d", symm)
	for c, gon := range oo.all8() {
		t.Logf("%d) %v", c+1, gon)
	}
}


func TestAccumOrigin15(t *testing.T)  {
	// H_15(1,9,1,9,1,9)
	//         0   1  2  3   4  5   6  7
	h_1_9 := [][]int{
		[]int{ 1, -1, 0, 0, -1, 2, -1, 0 },
		[]int{ 1,  1, 0, 0, -1, 0,  1, 0 },
	}
	t.Logf("H15(1,9,1,9,1,9) x=%t", accumOrigin15(h_1_9[0], h_1_9[1]))
}


func TestSimple(t *testing.T) {
	type T struct {
		gons   Gons
		angles []int
		simple bool
	}
	s7, _  := NewSymm(7);   p7 := NewPolylines(s7)
	s9, _  := NewSymm(9);   p9 := NewPolylines(s9)
	s11, _ := NewSymm(11); p11 := NewPolylines(s11)
	s13, _ := NewSymm(13); p13 := NewPolylines(s13)
	s15, _ := NewSymm(15); p15 := NewPolylines(s15)


	h7  := NewHexagons(p7)
	h9  := NewHexagons(p9)
	h11 := NewHexagons(p11)
	h13 := NewHexagons(p13)
	h15 := NewHexagons(p15)

	for _, test := range []*T {
		&T{ gons:h7, angles:[]int { 1,1,5, 1,1,5 }, simple:false },
		&T{ gons:h7, angles:[]int { 1,2,4, 1,2,4 }, simple:true  },

		&T{ gons:h9, angles:[]int { 1,1,7, 1,1,7 }, simple:false },
		&T{ gons:h9, angles:[]int { 1,2,6, 1,2,6 }, simple:true  },

		&T{ gons:h11, angles:[]int { 1,1,9, 1,1,9 }, simple:false },
		&T{ gons:h11, angles:[]int { 1,2,8, 1,2,8 }, simple:false },
		&T{ gons:h11, angles:[]int { 1,3,7, 1,3,7 }, simple:true  },

		&T{ gons:h13, angles:[]int { 1,1,11, 1,1,11 }, simple:false },
		&T{ gons:h13, angles:[]int { 1,2,10, 1,2,10 }, simple:false },
		&T{ gons:h13, angles:[]int { 1,3, 9, 1,3, 9 }, simple:true  },

		&T{ gons:h15, angles:[]int { 1,1,13, 1,1,13 }, simple:false },
		&T{ gons:h15, angles:[]int { 1,2,12, 1,2,12 }, simple:false },
		&T{ gons:h15, angles:[]int { 1,3,11, 1,3,11 }, simple:false },
		&T{ gons:h15, angles:[]int { 1,4,10, 1,4,10 }, simple:true  },
		&T{ gons:h15, angles:[]int { 2,2,11, 2,2,11 }, simple:false },
		&T{ gons:h15, angles:[]int { 2,3,10, 2,3,10 }, simple:true  },

	} {
		if trs, err := test.gons.Transforms(test.angles); err != nil {
			t.Fatalf("transform error: %v", err)
		} else if gon, err := test.gons.New(trs, 1, 1); err != nil {
			t.Fatalf("hexagon error: %v", err)
		} else {
			if simple, err := Simple(gon); err != nil {
				t.Fatalf("simple error: %v", err)
			} else if test.simple != simple {
				t.Fatalf("simple exp: %t got %t", test.simple, simple)
			} else {
				t.Logf("%v simple=%t", gon.Transforms(), simple)
			}
		}
	}
}
