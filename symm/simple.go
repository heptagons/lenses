package symm

/* 

H7(1,1,5,1,1,5) is not a simple hexagon
because edges A(a,f) and C(b,c) AT LEAST intersect at (X)
                  
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
                  
Strategy, compare edges A..F by pairs with these results:

  A   B   C   D   E   F        skip    Test
=========================    ========  ====
  |-->|   |   |   |   |      adjacent
  |------>|   |   |   |                 AC (intersects)
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
              |------>|                 DF (intersects)  
                  |-->|      adjacent


*/

func Simple(pp *Polylines, t *Transforms) bool {
	return true
}

func simple(pp *Polylines, t *Transforms) {
	
}
