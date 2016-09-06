package web

import (
	"log"
	"net/http"
	"time"

	"hectorcorrea.com/models"
)

type requestSession struct {
	resp      http.ResponseWriter
	req       *http.Request
	cookie    *http.Cookie
	loginName string
}

func sessionFromRequest(resp http.ResponseWriter, req *http.Request) requestSession {
	login := ""
	cookie, err := req.Cookie("sessionId")
	if err == nil {
		sessionId := cookie.Value
		userSession, err := models.ValidateUserSession(sessionId)
		if err == nil {
			login = userSession.Login
			log.Printf("Session is valid (%s), user: %s", cookie.Value, login)
		} else {
			log.Printf("Session was not valid (%s), %s", cookie.Value, err)
		}
	} else {
		cookie = nil
		log.Printf("No cookie: %s", err)
	}
	return requestSession{resp: resp, req: req, cookie: cookie, loginName: login}
}

func (s requestSession) IsAuth() bool {
	return s.loginName != ""
}

func (s requestSession) logout() {
	s.loginName = ""
	if s.cookie != nil {
		sessionId := s.cookie.Value
		models.DeleteUserSession(sessionId)
		s.cookie.Value = ""
		s.cookie.Expires = time.Unix(0, 0)
		s.cookie.Path = "/"
		s.cookie.HttpOnly = true
		http.SetCookie(s.resp, s.cookie)
	}
}

func (s requestSession) login(loginName string) {
	if s.cookie == nil {
		s.cookie = &http.Cookie{Name: "sessionId"}
	}
	// TODO: validate login/password
	userSession, err := models.AddUserSession(loginName)
	if err != nil {
		log.Printf("ERROR saving session %s", err)
	} else {
		s.loginName = userSession.Login
		s.cookie.Value = userSession.SessionId
		s.cookie.Expires = userSession.ExpiresOn
		s.cookie.Path = "/"
		s.cookie.HttpOnly = true
		http.SetCookie(s.resp, s.cookie)
	}
}
