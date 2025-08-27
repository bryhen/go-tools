package utl

// Panics if error found, otherwise returns val
func Ok[T any](val T, err error) T {
	if err != nil {
		panic(err.Error())
	}

	return val
}
