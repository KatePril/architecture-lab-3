package lang

type InvalidArgumentsError struct {
}

func (e InvalidArgumentsError) Error() string {
	return "Given arguments are invalid"
}
