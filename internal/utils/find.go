package utils

// FindBy is a generic function that takes a slice of any type and a predicate function.
// It returns the first element in the slice that satisfies the predicate function.
func FindBy[T any](items *[]T, predicate func(T) bool) (T, bool) {
	for _, item := range *items {
		if predicate(item) {
			return item, true
		}
	}
	var zeroValue T
	return zeroValue, false
}
