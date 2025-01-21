package utils

type Nullable[T any] struct {
	V       T
	IsValid bool
}

func NewNullable[T any](value T, isValid bool) Nullable[T] {
	return Nullable[T]{
		V:       value,
		IsValid: isValid,
	}
}
