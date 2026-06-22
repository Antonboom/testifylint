package checkers

type stack[T any] []T

func (s stack[T]) Len() int {
	return len(s)
}

func (s *stack[T]) Push(v T) {
	*s = append(*s, v)
}

func (s *stack[T]) Pop() T {
	n := len(*s)
	if n == 0 {
		var zero T
		return zero
	}

	last := (*s)[n-1]
	*s = (*s)[:n-1]
	return last
}

func (s stack[T]) Last() T {
	if len(s) == 0 {
		var zero T
		return zero
	}
	return s[len(s)-1]
}
