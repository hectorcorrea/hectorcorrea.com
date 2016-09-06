package viewModels

type Session struct {
	Id        string
	LoginName string
	IsAuth    bool
}

func NewSession(sessionId, loginName string) Session {
	return Session{
		Id:        sessionId,
		LoginName: loginName,
		IsAuth:    loginName != "",
	}
}
