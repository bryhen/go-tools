package utl

// Panics if error found
func Check(err error) {
	if err != nil {
		panic(err.Error())
	}
}
