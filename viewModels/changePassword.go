package viewModels

type ChangePassword struct {
	Message string
	Session
}

func NewChangePassword(message string, session Session) ChangePassword {
	return ChangePassword{Message: message, Session: session}
}
