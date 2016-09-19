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
		renderAuth(session, "views/login.html", viewModels.Login{})
		return
	}

	if req.Method == "POST" && req.URL.Path == "/auth/login" {
		login := session.req.FormValue("user")
		password := session.req.FormValue("password")
		log.Printf("Login in user %s, password %s", login, password)
		err := session.login(login, password)
		if err != nil {
			vm := viewModels.Login{LoginName: login, Message: "Sorry, not sorry"}
			renderAuth(session, "views/login.html", vm)
		} else {
			// vm := viewModels.Login{LoginName: login, Message: "Welcome"}
			http.Redirect(resp, req, "/", 302)
			// renderAuth(session, "views/welcomeback.html", vm)
		}
		return
	}

	if req.Method == "GET" && req.URL.Path == "/auth/logout" {
		session.logout()
		renderAuth(session, "views/logout.html", viewModels.Login{})
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
