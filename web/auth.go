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

var authRouter Router

func authPages(resp http.ResponseWriter, req *http.Request) {

	// This should be initialized only once, not on every call.
	authRouter.Add("GET", "/auth/login", handleLogin)
	authRouter.Add("POST", "/auth/login", handleLoginPost)
	authRouter.Add("GET", "/auth/logout", handleLogout)
	authRouter.Add("GET", "/auth/changepassword", handleChangePass)
	authRouter.Add("POST", "/auth/changepassword", handleChangePassPost)

	session := newSession(resp, req)
	found, route := authRouter.FindRoute(req.Method, req.URL.Path)
	if found {
		route.handler(session, nil)
	} else {
		renderNotFound(session)
	}
}

func handleLogin(s session, values map[string]string) {
	vmSession := s.toViewModel()
	vm := viewModels.NewLogin("", vmSession)
	renderTemplate(s, "views/login.html", vm)
}

func handleLoginPost(s session, values map[string]string) {
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
		http.Redirect(s.resp, s.req, "/", 302)
	}
}

func handleLogout(s session, values map[string]string) {
	s.logout()
	homeUrl := fmt.Sprintf("/?cb?=%s", cacheBuster())
	http.Redirect(s.resp, s.req, homeUrl, 302)
}

func handleChangePass(s session, values map[string]string) {
	if !s.isAuth() {
		renderNotAuthorized(s)
		return
	}

	vmSession := s.toViewModel()
	vm := viewModels.NewChangePassword("", vmSession)
	renderTemplate(s, "views/changePassword.html", vm)
}

func handleChangePassPost(s session, values map[string]string) {
	if !s.isAuth() || (s.loginName != s.req.FormValue("user")) {
		renderNotAuthorized(s)
		return
	}

	login := s.loginName
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
		vmSession := s.toViewModel()
		vm := viewModels.NewChangePassword(message, vmSession)
		renderTemplate(s, "views/changePassword.html", vm)
	} else {
		err := models.SetPassword(login, newPassword)
		if err != nil {
			renderError(s, "Could not change passowrd", err)
		} else {
			http.Redirect(s.resp, s.req, "/", 302)
		}
	}
}

func cacheBuster() string {
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	return fmt.Sprintf("%d", r.Int())
}
