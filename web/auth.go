package web

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"hectorcorrea.com/models"
	"hectorcorrea.com/viewModels"
)

var authRoutes = []Route{
	Route{method: "GET", path: "/auth/login", handler: handleLogin},
	Route{method: "POST", path: "/auth/login", handler: handleLoginPost},
	Route{method: "GET", path: "/auth/logout", handler: handleLogout},
	Route{method: "GET", path: "/auth/changepassword", handler: handleChangePass},
	Route{method: "POST", path: "/auth/changepassword", handler: handleChangePassPost},
}

func authPages(resp http.ResponseWriter, req *http.Request) {
	session := newSession(resp, req)
	handled := false
	for _, r := range authRoutes {
		if r.method == req.Method && req.URL.Path == r.path {
			r.handler(resp, req, session)
			handled = true
			break
		}
	}

	if !handled {
		renderNotFound(session)
	}
}

func handleLogin(resp http.ResponseWriter, req *http.Request, s session) {
	vmSession := s.toViewModel()
	vm := viewModels.NewLogin("", vmSession)
	renderTemplate(s, "views/login.html", vm)
}

func handleLoginPost(resp http.ResponseWriter, req *http.Request, s session) {
	login := s.req.FormValue("user")
	password := s.req.FormValue("password")
	err := s.login(login, password)
	if err != nil {
		log.Printf("Login FAILED for user: %s", login)
		vmSession := s.toViewModel()
		vm := viewModels.NewLogin("Sorry, not sorry", vmSession)
		renderTemplate(s, "views/login.html", vm)
	} else {
		log.Printf("Login OK for user: %s", login)
		http.Redirect(resp, req, "/", 302)
	}
}

func handleLogout(resp http.ResponseWriter, req *http.Request, s session) {
	s.logout()
	homeUrl := fmt.Sprintf("/?cb?=%s", cacheBuster())
	http.Redirect(resp, req, homeUrl, 302)
}

func handleChangePass(resp http.ResponseWriter, req *http.Request, s session) {
	vmSession := s.toViewModel()
	if !vmSession.IsAuth {
		renderNotAuthorized(s)
		return
	}

	vm := viewModels.NewChangePassword("", vmSession)
	renderTemplate(s, "views/changePassword.html", vm)
}

func handleChangePassPost(resp http.ResponseWriter, req *http.Request, s session) {
	vmSession := s.toViewModel()
	if !vmSession.IsAuth {
		renderNotAuthorized(s)
		return
	}

	if vmSession.LoginName != s.req.FormValue("user") {
		renderNotAuthorized(s)
		return
	}

	login := vmSession.LoginName
	password := s.req.FormValue("oldPassword")
	newPassword := s.req.FormValue("newPassword")
	repeatPassword := s.req.FormValue("repeatPassword")
	message := ""

	err := s.login(login, password)
	if err != nil {
		message += "Invalid password."
	}

	if len(newPassword) == 0 {
		message += "New password cannot be empty."
	}
	if newPassword != repeatPassword {
		message += "Password and Repeat Password must match."
	}

	if len(message) > 0 {
		vm := viewModels.NewChangePassword(message, vmSession)
		renderTemplate(s, "views/changePassword.html", vm)
	} else {
		err := models.SetPassword(login, newPassword)
		if err != nil {
			renderError(s, "Could not change passowrd", err)
		} else {
			http.Redirect(resp, req, "/", 302)
		}
	}
}

func cacheBuster() string {
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	return fmt.Sprintf("%d", r.Int())
}
