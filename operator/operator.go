package operator

// Go Generics for Ternary function that returns the first value if the condition is true, and returns the second value if the condition is false
func Ternary[T comparable](condition bool, a, b T) T {
	if condition {
		return a
	}

	return b
}
