package cbpv

type Stack []Value

func (s *Stack) Push(v Value) { *s = append(*s, v) }
func (s *Stack) Pop() (Value, bool) {
	if len(*s) == 0 {
		return nil, false
	}
	old := *s
	n := len(old)
	res := old[n-1]
	*s = old[:n-1]
	return res, true
}
