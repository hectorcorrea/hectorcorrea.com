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

func authPages(resp http.ResponseWriter, req *http.Request) {
	session := newSession(resp, req)
	vmSession := session.toViewModel()
	path := req.URL.Path

	if req.Method == "GET" {
		if path == "/auth/login" {
			vm := viewModels.NewLogin("", vmSession)
			renderTemplate(session, "views/login.html", vm)
			return
		} else if path == "/auth/changepassword" {

			if !vmSession.IsAuth {
				renderNotAuthorized(session)
				return
			}

			vm := viewModels.NewChangePassword("", vmSession)
			renderTemplate(session, "views/changePassword.html", vm)
			return
		} else if path == "/auth/logout" {
			session.logout()
			homeUrl := fmt.Sprintf("/?cb?=%s", cacheBuster())
			http.Redirect(resp, req, homeUrl, 302)
			return
		}
	}

	if req.Method == "POST" {
		if path == "/auth/login" {
			login := session.req.FormValue("user")
			password := session.req.FormValue("password")
			err := session.login(login, password)
			if err != nil {
				log.Printf("Login FAILED for user: %s", login)
				vm := viewModels.NewLogin("Sorry, not sorry", vmSession)
				renderTemplate(session, "views/login.html", vm)
			} else {
				log.Printf("Login OK for user: %s", login)
				http.Redirect(resp, req, "/", 302)
			}
			return
		} else if path == "/auth/changepassword" {

			if !vmSession.IsAuth {
				renderNotAuthorized(session)
				return
			}

			if vmSession.LoginName != session.req.FormValue("user") {
				renderNotAuthorized(session)
				return
			}

			login := vmSession.LoginName
			password := session.req.FormValue("oldPassword")
			newPassword := session.req.FormValue("newPassword")
			repeatPassword := session.req.FormValue("repeatPassword")
			message := ""

			err := session.login(login, password)
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
				renderTemplate(session, "views/changePassword.html", vm)
			} else {
				err := models.SetPassword(login, newPassword)
				if err != nil {
					renderError(session, "Could not change passowrd", err)
				} else {
					http.Redirect(resp, req, "/", 302)
				}
			}
			return
		}
	}

	renderNotFound(session)
}

func cacheBuster() string {
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	return fmt.Sprintf("%d", r.Int())
}
