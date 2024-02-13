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

func Simple(pp *Polylines, t *Transforms) bool {
	return true
}

func simpleStrategy(p *Polyline, t *Transforms) error {
	edges := p.Edges()
	n := len(edges) 
	if n < 3 {
		return fmt.Errorf("Invalid polygon number of edges: %d", n)
	}
	for p1, e1 := range edges {
		for p2, e2 := range edges {
			if p2 < p1 {
				continue
			}
			next     := (p1 == p2 - 1)         // consecutive letter edges.
			extremes := (p1 == 0 && p2 == n-1) // first/last edges are adjacent (edges ring).
			if next || extremes {
				// skip because edges are adjacent.
				// they share a vertice and cannot intersect.
			} else if e1 == e2 {
				// skip because edges are parallel.
				// having the same vece2r cannot intersect.
			} else {
				fmt.Println("Strategy", p1, p2, e1, e2)
			}
		}
	}
	return nil
}
