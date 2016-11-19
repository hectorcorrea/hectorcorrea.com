package web

import (
	"fmt"
	"log"
	"net/http"
	"math/rand"
	"time"
	"hectorcorrea.com/viewModels"
)

func authPages(resp http.ResponseWriter, req *http.Request) {
	session := newSession(resp, req)

	if req.Method == "GET" && req.URL.Path == "/auth/login" {
		vm := viewModels.NewLogin("", session.toViewModel())
		renderAuth(session, "views/login.html", vm)
		return
	}

	if req.Method == "POST" && req.URL.Path == "/auth/login" {
		login := session.req.FormValue("user")
		password := session.req.FormValue("password")
		err := session.login(login, password)
		if err != nil {
			log.Printf("Login FAILED for user: %s", login)
			vm := viewModels.NewLogin("Sorry, not sorry", session.toViewModel())
			renderAuth(session, "views/login.html", vm)
		} else {
			log.Printf("Login OK for user: %s", login)
			http.Redirect(resp, req, "/", 302)
		}
		return
	}

	if req.Method == "GET" && req.URL.Path == "/auth/logout" {
		session.logout()
		homeUrl := fmt.Sprintf("/?cb?=%s", cacheBuster())
		http.Redirect(resp, req, homeUrl, 302)
		return
	}

	renderNotFound(session)
}

func renderAuth(s session, viewName string, vm viewModels.Login) {
	t, err := loadTemplate(viewName)
	if err != nil {
		renderError(s, fmt.Sprintf("Loading view %s", viewName), err)
	} else {
		t.Execute(s.resp, vm)
	}
}

func cacheBuster() string {
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	return fmt.Sprintf("%d", r.Int())
}
