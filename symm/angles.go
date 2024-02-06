package symm

type Angles struct {
	min int
	max int
	sum int
}

func (a *Angles) ValidAngle(angle int) bool {
	if a.min > angle || a.max < angle {
		return false
	}
	return true
}

func (a *Angles) ValidSum(sum int) bool {
	return a.sum == sum
}


