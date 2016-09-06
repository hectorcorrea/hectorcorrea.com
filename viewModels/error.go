package viewModels

type Error struct {
	Title   string
	Details string
}

func NewError(title string, err error) Error {
	details := ""
	if err != nil {
		details = err.Error()
	}
	return Error{Title: title, Details: details}
}
