package stack

// Stackstring a stack of string
type Stackstring struct {
	data []string
	size int
}

// Push pushes a string on the Stack
func (s *Stackstring) Push(v string) {
	if s.size == len(s.data) {
		s.data = append(s.data, v)
	} else {
		s.data[s.size] = v
	}
	s.size++
}

// Pop pops a string from the stack
func (s *Stackstring) Pop() (r string) {
	if s.size == 0 {
		return
	}
    s.size--
    r = s.data[s.size]
	return
}

// Size returns current size of the Stack
func (s *Stackstring) Size() int {
	return s.size
}

// IsEmpty returns true if stack is empty
func (s *Stackstring) IsEmpty() bool {
	return s.size == 0
}