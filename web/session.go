package web

import (
	"errors"
	"log"
	"net/http"
	"time"

	"hectorcorrea.com/models"
	"hectorcorrea.com/viewModels"
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
		} else {
			log.Printf("Session was not valid (%s), %s", cookie.Value, err)
			cookie = nil
			sessionId = ""
		}
	} else {
		cookie = nil
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

func (s *session) logout() {
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

func (s *session) login(loginName, password string) error {
	if s.cookie == nil {
		s.cookie = &http.Cookie{Name: "sessionId"}
	}

	logged, err := models.LoginUser(loginName, password)
	if err != nil {
		return err
	}

	if logged {
		userSession, err := models.NewUserSession(loginName)
		if err != nil {
			log.Printf("ERROR creating new session: %s", err)
			return err
		}

		s.loginName = userSession.Login
		s.sessionId = userSession.SessionId
		s.cookie.Value = s.sessionId
		s.cookie.Path = "/"
		s.cookie.HttpOnly = true

		// Give the cookie a long expiration date (we control
		// the session expiration via the date on the database.)
		daysFromToday := time.Now().UTC().Add(time.Hour * 24 * 90)
		s.cookie.Expires = daysFromToday

		http.SetCookie(s.resp, s.cookie)
		return nil
	}

	log.Printf("ERROR invalid user/password received: %s/***", loginName)
	return errors.New("Invalid user/password received")
}

func (s session) isAuth() bool {
	return s.loginName != ""
}

// Provide toViewModel() here since this type does not have
// a model per-se.
func (s session) toViewModel() viewModels.Session {
	return viewModels.NewSession(s.sessionId, s.loginName, s.isAuth())
}
