package web

import (
	"log"
	"net/http"
	"time"

	"hectorcorrea.com/models"
)

type session struct {
	resp   http.ResponseWriter
	req    *http.Request
	cookie *http.Cookie
	isAuth bool
}

func sessionFromRequest(resp http.ResponseWriter, req *http.Request) session {
	isAuth := false
	cookie, err := req.Cookie("sessionId")
	if err == nil {
		if models.IsValidSession(cookie.Value) {
			isAuth = true
			log.Printf("Session is valid (%s)", cookie.Value)
		} else {
			log.Printf("Session was not valid (%s)", cookie.Value)
		}
	} else {
		cookie = nil
		log.Printf("No cookie: %s", err)
	}
	return session{resp: resp, req: req, cookie: cookie, isAuth: isAuth}
}

func (s session) logout() {
	s.isAuth = false
	if s.cookie != nil {
		models.DeleteSession(s.cookie.Value)
		s.cookie.Value = ""
		s.cookie.Expires = time.Unix(0, 0)
		s.cookie.Path = "/"
		s.cookie.HttpOnly = true
		http.SetCookie(s.resp, s.cookie)
	}
}

func (s session) login() {
	s.isAuth = true
	if s.cookie == nil {
		s.cookie = &http.Cookie{Name: "sessionId"}
	}
	x, err := models.AddSession()
	if err != nil {
		log.Printf("ERROR saving session %s", err)
	} else {
		s.cookie.Value = x.Id
		s.cookie.Expires = x.ExpiresOn
		s.cookie.Path = "/"
		s.cookie.HttpOnly = true
		http.SetCookie(s.resp, s.cookie)
	}
}
