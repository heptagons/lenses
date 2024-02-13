package symm

import (
	"fmt"
)

/* 

H7(1,1,5,1,1,5) is not a simple hexagon
because edges A(a,f) and C(b,c) AT LEAST intersect at point X:
                  
                 |         b
                          / \
                 |       /   \
                      C /     \ B
                 |     /       \
                      /         \
 - - - - - - - - f - X - - - - - a - - 
         D      /   /      A
   d-----------/-|-c
    \         /   
     \       /   |
      \     / F     
     E \   /     |
        \ /       
         e       |
                  
For shift=1, vece2r=1
Angles = 1 1 5 1 1
Edges  = 1 7 6 1 7 6

Strategy, compare edges A..F by pairs with these results:

  A   B   C   D   E   F        skip    Test   Result
=========================    ========  ====  =========
  |-->|   |   |   |   |      adjacent
  |------>|   |   |   |                 AC   intersect
  |---------->|   |   |      parallel
  |-------------->|   |                 AE 
  |------------------>|      adjacent
      |-->|                  adjacent
      |------>|                         BD
      |---------->|          parallel
      |-------------->|                 BF
          |-->|              adjacent
          |------>|                     CE
          |---------->|      parallel
              |-->|          adjacent
              |------>|                 DF   intersect
                  |-->|      adjacent
*/
func Simple(gon Gon) (bool, error) {
    p := gon.Polyline()
    //t := gon.Transforms()
	edges := p.Edges()
	n := len(edges) 
	if n < 3 {
		return false, fmt.Errorf("Invalid polygon number of edges: %d", n)
	}
    accums := p.Accums()
	for p1, e1 := range edges {
		for p2, e2 := range edges {
			if p2 < p1 {
				continue
			}
			if (p1 == p2 - 1) {
                // skip consecutive letter edges adjacent
                // they share a vertice and cannot intersect.
                continue
            }
			if p1 == 0 && p2 == n-1 {
                // first/last edges are adjacent (edges ring).
                // they share a vertice and cannot intersect.
                continue
            }
			if e1 == e2 {
				// skip because edges are parallel.
				// having the same vector cannot intersect.
                continue
			}
            m1 := p.pp.s.XY(accums[(p1+n-1) % n])
            m2 := p.pp.s.XY(accums[      p1 % n])
            m4 := p.pp.s.XY(accums[(p2+n-1) % n])
            m3 := p.pp.s.XY(accums[      p2 % n])
            //fmt.Printf("simple e[%d]=%d e[%d]=%d\n", p1, e1, p2, e2)
            if simpleBezier(m1, m2, m3, m4) == false {
                // intersect => not simple
                return false, nil
            }
		}
	}
	return true, nil
}

func simpleBezier(m1, m2, m3, m4 []float64) bool {
    x1, y1 :=  m1[0], m1[1]
    x2, y2 :=  m2[0], m2[1]
    x3, y3 :=  m3[0], m3[1]
    x4, y4 :=  m4[0], m4[1]
    // https://en.wikipedia.org/wiki/Line%E2%80%93line_intersection#Given_two_points_on_each_line_segment
    //      | x1 |    | x2 - x1 |
    // L1 = |    | + t|  
    //      | y1 |    | y2 - y1 |
    //
    //      | x3 |    | x4 - x3 |
    // L2 = |    | + u|         |
    //      | y3 |    | y4 - y3 |
    //
    // d = (x1 - x2)(y3 - y4) - (y1 - y2)(x3 - x4)
    //
    //       (x1 - x3)(y3 - y4) - (y1 - y3)(x3 - x4)
    // t = + ---------------------------------------
    //                          d
    //       (x1 - x2)(y1 - y3) - (y1 - y2)(x1 - x3)
    // u = - ---------------------------------------
    //                          d
    d  :=   (x1 - x2)*(y3 - y4) - (y1 - y2)*(x3 - x4)
    t  :=  ((x1 - x3)*(y3 - y4) - (y1 - y3)*(x3 - x4)) / d
    u  := -((x1 - x2)*(y1 - y3) - (y1 - y2)*(x1 - x3)) / d
    tt := 0 <= t && t <= 1
    uu := 0 <= u && u <= 1
    if tt && uu {
        return false
    }
    return true
} 



func simpleFloat(m1, m2, m3, m4 []float64) bool {
    x1, y1 :=  m1[0], m1[1]
    x2, y2 :=  m2[0], m2[1]
    x3, y3 :=  m3[0], m3[1]
    x4, y4 :=  m4[0], m4[1]
    // https://en.wikipedia.org/wiki/Intersection_(geometry)#Two_line_segments

    // Xs, Ys =  x1 + (x2 - x1)s,   y1 + (y2 - y1)s

    // (x2 - x1)s - (x4 - x3)t = x3 - x1   ->   a1X + b1Y = c1
    // (y2 - y1)s - (y4 - y3)t = y3 - y1   ->   a2X + b2Y = c2
    // D  = (a1b2 - a2b1)
    // Xs = (c1b2 - c2b1) / D
    // Ys = (a1c2 - a2c1) / D
    a1, b1, c1 := x2 - x1, -(x4 - x3), x3 - x1
    a2, b2, c2 := y2 - y1, -(y4 - y3), y3 - y1
    D := a1*b2 - a2*b1
    if D == 0 {
        return true
    }
    x0 := (c1*b2 - c2*b1) / D
    y0 := (a1*c2 - a2*c1) / D
    fmt.Printf("\tx0=%f y0=%f \n", x0, y0)
    if 0 < x0 && y0 < 1 {
        return false
    }
    return true
}

