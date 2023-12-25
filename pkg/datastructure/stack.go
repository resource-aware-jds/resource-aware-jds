package datastructure

// Stack implement the stack data structure.
// First come, Last out.
type Stack[T any] []T

// IsEmpty check if stack is empty
func (s *Stack[T]) IsEmpty() bool {
	return len(*s) == 0
}

// Push a new value onto the stack
func (s *Stack[T]) Push(str T) {
	*s = append(*s, str) // Simply append the new value to the end of the stack
}

// Pop Remove and return top element of stack. Return false if stack is empty.
func (s *Stack[T]) Pop() (*T, bool) {
	if s.IsEmpty() {
		return nil, false
	} else {
		index := len(*s) - 1   // Get the index of the top most element.
		element := (*s)[index] // Index into the slice and obtain the element.
		*s = (*s)[:index]      // Remove it from the stack by slicing it off.
		return &element, true
	}
}
