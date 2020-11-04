package stack

// Stackuint64 a stack of uint64
type Stackuint64 struct {
	data []uint64
	size int
}

// Push pushes a uint64 on the Stack
func (s *Stackuint64) Push(v uint64) {
	if s.size == len(s.data) {
		s.data = append(s.data, v)
	} else {
		s.data[s.size] = v
	}
	s.size++
}

// Pop pops a uint64 from the stack
func (s *Stackuint64) Pop() (r uint64) {
	if s.size == 0 {
		return
	}
    s.size--
    r = s.data[s.size]
	return
}

// Size returns current size of the Stack
func (s *Stackuint64) Size() int {
	return s.size
}

// IsEmpty returns true if stack is empty
func (s *Stackuint64) IsEmpty() bool {
	return s.size == 0
}