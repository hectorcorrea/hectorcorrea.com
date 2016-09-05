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

func blogFromForm(id int, rr webRR) models.Blog {
	var blog models.Blog
	blog.Id = id
	blog.Title = rr.req.FormValue("title")
	blog.Summary = rr.req.FormValue("summary")
	blog.Content = rr.req.FormValue("content")
	return blog
}

func blogPages(resp http.ResponseWriter, req *http.Request) {
	rr := newWebRR(resp, req)
	if req.Method == "GET" {
		blogView(rr)
	} else if req.Method == "POST" {
		blogAction(rr)
	} else {
		renderError(rr, "Unknown HTTP Method", errors.New("HTTP method not supported"))
	}
}

func blogView(rr webRR) {
	if id, err := parseBlogViewUrl(rr.req.URL.Path); err != nil {
		renderError(rr, "Cannot parse Blog URL", err)
	} else if id != 0 {
		blogViewOne(rr, id)
	} else {
		blogViewAll(rr)
	}
}

func blogViewOne(rr webRR, id int) {
	log.Printf("Loading %d", id)
	blog, err := models.BlogGetById(id)
	if err != nil {
		renderError(rr, "Fetching by ID", err)
	} else {
		t, err := loadTemplate("views/blogView.html")
		if err != nil {
			renderError(rr, "Loading template blogView", err)
		} else {
			vm := viewModels.FromBlog(blog)
			t.Execute(rr.resp, vm)
		}
	}
}

func blogViewAll(rr webRR) {
	log.Printf("Loading all...")
	if blogs, err := models.BlogGetAll(); err != nil {
		renderError(rr, "Error fetching all", err)
	} else {
		vm := viewModels.FromBlogs(blogs)
		if t, err := loadTemplate("views/blogList.html"); err != nil {
			renderError(rr, "Loading template blogList", err)
		} else {
			t.Execute(rr.resp, vm)
		}
	}
}

func blogAction(rr webRR) {
	id, action, err := parseBlogEditUrl(rr.req.URL.Path)
	if err != nil {
		renderError(rr, "Cannot determine HTTP action", err)
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

func blogSave(rr webRR, id int) {
	blog := blogFromForm(id, rr)
	if err := blog.Save(); err != nil {
		renderError(rr, fmt.Sprintf("Saving blog ID: %d"), err)
	} else {
		url := fmt.Sprintf("/blog/%s/%d", blog.Slug, id)
		log.Printf("Redirect to %s", url)
		http.Redirect(rr.resp, rr.req, url, 301)
	}
}

func blogDraft(rr webRR, id int) {
	// AJAX request
}

func blogPost(rr webRR, id int) {
	// AJAX request
}

func blogEdit(rr webRR, id int) {
	log.Printf("Loading %d", id)
	if blog, err := models.BlogGetById(id); err != nil {
		renderError(rr, fmt.Sprintf("Loading ID: %d", id), err)
	} else {
		t, err := loadTemplate("views/blogEdit.html")
		if err != nil {
			renderError(rr, "Loading template blogEdit", err)
		} else {
			vm := viewModels.FromBlog(blog)
			t.Execute(rr.resp, vm)
		}
	}
}

func parseBlogViewUrl(url string) (id int, err error) {
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
		return strconv.Atoi(parts[3])
	}
	return 0, errors.New("Could not parse (view) blog URL")
}

func parseBlogEditUrl(url string) (id int, action string, err error) {
	// url /blog/:title/:id/:action
	// parts[0] empty
	// parts[1] blog
	// parts[2] title
	// parts[3] id
	// parts[4] action (edit, post, draft)
	parts := strings.Split(url, "/")
	if len(parts) == 5 && parts[0] == "" && parts[1] == "blog" {
		if id, err := strconv.Atoi(parts[3]); err != nil {
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
