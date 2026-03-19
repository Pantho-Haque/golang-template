package generic

func Value2Pointer[T any](value T) *T {
	return &value
}
