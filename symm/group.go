package symm

// https://mathstat.slu.edu/escher/index.php/Introduction_to_Symmetry
type Group struct {
	Letter string
	Number int
}

// NewGroupC builds a cyclic groups
// C1 is the symmetry with identity only, like letters F,G,J,L,P,Q,R
// C2 is the symmetry of 180° like letters N,S,Z
// C3 is the symmetry of 120° triskelion
// C4 is the symmetry of 90° swastika
func NewGroupC(number int) *Group {
	return &Group{
		Letter: "C",
		Number: number,
	}
}

// NewGroupD builds diedral groups
// D1 is the single mirror symmetry (like of letters A,B,C,D,E,K,M,T,U,V,W,Y)
// D2 is the symmetry of the rectangle (like of letters H,I,X,O)
// D3 is the symmetry of the equilateral triangle 
// D5 is the symmetry of the regular pentagon
// D6 is the symmetry of the regular hexagon
// D7 is the symmetry of the regular heptagon
// D9 is the symmetry of the regular 9-gon
// DN N >= 10 is the symmetry of the regular N-gon
func NewGroupD(number int) *Group {
	return &Group{
		Letter: "D",
		Number: number,
	}
}

var (
	C1 = NewGroupC(1) // only identity symmetry
	C2 = NewGroupC(2) // 180°
	D1 = NewGroupD(1) // mirror
	D2 = NewGroupD(2) // rectangle symmetry
	D3 = NewGroupD(3) // equilateral triangle
	D6 = NewGroupD(6) // regular hexagon
)
