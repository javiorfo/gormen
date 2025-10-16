package error

type PageError struct {
}

func (p PageError) Error() string {
	return ""
}
