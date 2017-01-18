package viewModels

type Login struct {
	Message string
	Session
}

func NewLogin(message string, session Session) Login {
	return Login{Message: message, Session: session}
}
