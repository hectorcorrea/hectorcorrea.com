package viewModels

type Session struct {
	Id        string
	LoginName string
	IsAuth    bool
}

func NewSession(id, loginName string) Session {
	return Session{
		Id:        id,
		LoginName: loginName,
		IsAuth:    loginName != "",
	}
}
