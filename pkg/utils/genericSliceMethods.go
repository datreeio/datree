package utils

// MapSlice should be replaced with an external package once we find a reliable one
func MapSlice[Input any, Output any](inputs []Input, f func(Input) Output) []Output {
	n := make([]Output, len(inputs))
	for index, input := range inputs {
		n[index] = f(input)
	}
	return n
}
