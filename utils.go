package orm

// pointer returns a pointer of the provided value.
func pointer[T any](x T) *T {
	return &x
}
