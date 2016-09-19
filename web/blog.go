package web

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"hectorcorrea.com/models"
	"hectorcorrea.com/viewModels"
)

func blogFromForm(id int64, s session) models.Blog {
	var blog models.Blog
	blog.Id = id
	blog.Title = s.req.FormValue("title")
	blog.Summary = s.req.FormValue("summary")
	blog.Content = s.req.FormValue("content")
	return blog
}

func blogPages(resp http.ResponseWriter, req *http.Request) {
	session := newSession(resp, req)
	if req.Method == "GET" {
		blogView(session)
	} else if req.Method == "POST" {
		blogAction(session)
	} else {
		renderError(session, "Unknown HTTP Method", errors.New("HTTP method not supported"))
	}
}

func blogView(rr session) {
	if id, err := parseBlogViewUrl(rr.req.URL.Path); err != nil {
		renderError(rr, "Cannot parse Blog URL", err)
	} else if id != 0 {
		blogViewOne(rr, id)
	} else {
		blogViewAll(rr)
	}
}

func blogViewOne(rr session, id int64) {
	log.Printf("Loading %d", id)
	blog, err := models.BlogGetById(id)
	if err != nil {
		renderError(rr, "Fetching by ID", err)
	} else {
		t, err := loadTemplate("views/blogView.html")
		if err != nil {
			renderError(rr, "Loading template blogView", err)
		} else {
			vm := viewModels.FromBlog(blog, vmSession(rr))
			t.Execute(rr.resp, vm)
		}
	}
}

func vmSession(rr session) viewModels.Session {
	return viewModels.NewSession(rr.sessionId, rr.loginName)
}

func blogViewAll(rr session) {
	log.Printf("Loading all...")
	if blogs, err := models.BlogGetAll(); err != nil {
		renderError(rr, "Error fetching all", err)
	} else {
		vm := viewModels.FromBlogs(blogs, vmSession(rr))
		if t, err := loadTemplate("views/blogList.html"); err != nil {
			renderError(rr, "Loading template blogList", err)
		} else {
			t.Execute(rr.resp, vm)
		}
	}
}

func blogAction(rr session) {
	// TODO: make sure user is authenticated
	id, action, err := parseBlogEditUrl(rr.req.URL.Path)
	if err != nil {
		renderError(rr, "Cannot determine HTTP action", err)
	} else if action == "new" {
		blogNew(rr)
	} else if action == "edit" {
		blogEdit(rr, id)
	} else if action == "save" {
		blogSave(rr, id)
	} else if action == "post" {
		blogPost(rr, id)
	} else if action == "draft" {
		blogDraft(rr, id)
	} else {
		renderError(rr, "Unknown action", nil)
	}
}

func blogSave(rr session, id int64) {
	blog := blogFromForm(id, rr)
	if err := blog.Save(); err != nil {
		renderError(rr, fmt.Sprintf("Saving blog ID: %d"), err)
	} else {
		url := fmt.Sprintf("/blog/%s/%d", blog.Slug, id)
		log.Printf("Redirect to %s", url)
		http.Redirect(rr.resp, rr.req, url, 301)
	}
}

func blogNew(rr session) {
	if newId, err := models.SaveNew(); err != nil {
		renderError(rr, fmt.Sprintf("Error creating new blog"), err)
	} else {
		log.Printf("Redirect to (edit for new) %d", newId)
		blogEdit(rr, newId)
	}
}

func blogDraft(rr session, id int64) {
	// AJAX request
}

func blogPost(rr session, id int64) {
	// AJAX request
}

func blogEdit(rr session, id int64) {
	log.Printf("Loading %d", id)
	if blog, err := models.BlogGetById(id); err != nil {
		renderError(rr, fmt.Sprintf("Loading ID: %d", id), err)
	} else {
		t, err := loadTemplate("views/blogEdit.html")
		if err != nil {
			renderError(rr, "Loading template blogEdit", err)
		} else {
			vm := viewModels.FromBlog(blog, vmSession(rr))
			t.Execute(rr.resp, vm)
		}
	}
}

func idFromString(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}

func parseBlogViewUrl(url string) (id int64, err error) {
	if url == "/blog/" {
		return 0, nil
	}
	// url /blog/:title/:id
	// parts[0] empty
	// parts[1] blog
	// parts[2] title
	// parts[3] id
	parts := strings.Split(url, "/")
	if len(parts) == 4 && parts[0] == "" && parts[1] == "blog" {
		return idFromString(parts[3])
	}
	return 0, errors.New("Could not parse (view) blog URL")
}

func parseBlogEditUrl(url string) (id int64, action string, err error) {
	if url == "/blog/new" {
		return 0, "new", nil
	}
	// url /blog/:title/:id/:action
	// parts[0] empty
	// parts[1] blog
	// parts[2] title
	// parts[3] id
	// parts[4] action (edit, post, draft)
	parts := strings.Split(url, "/")
	if len(parts) == 5 && parts[0] == "" && parts[1] == "blog" {
		if id, err := idFromString(parts[3]); err != nil {
			return 0, "", err
		} else {
			action := parts[4]
			if action == "edit" || action == "save" ||
				action == "post" || action == "draft" {
				return id, action, nil
			}
			return 0, "", errors.New("Invalid action")
		}
	}
	return 0, "", errors.New("Could not parse (edit) blog URL")
}
