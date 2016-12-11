package viewModels

type Error struct {
	Title   string
	Details string
	Session
}

func NewError(title string, err error, session Session) Error {
	details := ""
	if err != nil {
		details = err.Error()
	}
	return Error{Title: title, Details: details, Session: session}
}

func NewErrorFromStr(title string, err string, session Session) Error {
	return Error{Title: title, Details: err, Session: session}
}
