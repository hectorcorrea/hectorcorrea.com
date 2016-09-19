package web

import (
	"fmt"
	"log"
	"net/http"

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
		log.Printf("Login in user %s, password %s", login, password)
		err := session.login(login, password)
		if err != nil {
			vm := viewModels.NewLogin("Sorry, not sorry", session.toViewModel())
			renderAuth(session, "views/login.html", vm)
		} else {
			http.Redirect(resp, req, "/", 302)
		}
		return
	}

	if req.Method == "GET" && req.URL.Path == "/auth/logout" {
		vm := viewModels.NewLogin("Logged out", session.toViewModel())
		session.logout()
		renderAuth(session, "views/logout.html", vm)
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
