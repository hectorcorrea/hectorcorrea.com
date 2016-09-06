package web

import (
	"log"
	"net/http"
	"time"

	"hectorcorrea.com/models"
)

type session struct {
	resp      http.ResponseWriter
	req       *http.Request
	cookie    *http.Cookie
	loginName string
	sessionId string
}

func newSession(resp http.ResponseWriter, req *http.Request) session {
	sessionId := ""
	login := ""
	cookie, err := req.Cookie("sessionId")
	if err == nil {
		sessionId = cookie.Value
		userSession, err := models.GetUserSession(sessionId)
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
	s := session{
		resp:      resp,
		req:       req,
		cookie:    cookie,
		loginName: login,
		sessionId: sessionId,
	}
	return s
}

func (s session) IsAuth() bool {
	return s.loginName != ""
}

func (s session) logout() {
	models.DeleteUserSession(s.sessionId)
	s.loginName = ""
	s.sessionId = ""
	if s.cookie != nil {
		s.cookie.Value = ""
		s.cookie.Expires = time.Unix(0, 0)
		s.cookie.Path = "/"
		s.cookie.HttpOnly = true
		http.SetCookie(s.resp, s.cookie)
	}
}

func (s session) login(loginName string) {
	if s.cookie == nil {
		s.cookie = &http.Cookie{Name: "sessionId"}
	}
	// TODO: validate login/password
	userSession, err := models.NewUserSession(loginName)
	if err != nil {
		log.Printf("ERROR creating new session: %s", err)
	} else {
		s.loginName = userSession.Login
		s.sessionId = userSession.SessionId
		s.cookie.Value = s.sessionId
		s.cookie.Expires = userSession.ExpiresOn
		s.cookie.Path = "/"
		s.cookie.HttpOnly = true
		http.SetCookie(s.resp, s.cookie)
	}
}
