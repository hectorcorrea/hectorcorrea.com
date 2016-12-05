package web

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"hectorcorrea.com/models"
	"hectorcorrea.com/viewModels"
)

var blogRouter Router

func blogPages(resp http.ResponseWriter, req *http.Request) {

	// This should be initialized only once, not on every call.
	blogRouter.Add("GET", "/blog/:title/:id", blogViewOne)
	blogRouter.Add("GET", "/blog/", blogViewAll)
	blogRouter.Add("POST", "/blog/new", blogNew)
	blogRouter.Add("POST", "/blog/:title/:id/edit", blogEdit)
	blogRouter.Add("POST", "/blog/:title/:id/save", blogSave)
	blogRouter.Add("POST", "/blog/:title/:id/post", blogPost)
	blogRouter.Add("POST", "/blog/:title/:id/draft", blogDraft)

	session := newSession(resp, req)
	found, route := blogRouter.FindRoute(req.Method, req.URL.Path)
	if found {
		values := route.UrlValues(req.URL.Path)
		route.handler(session, values)
	} else {
		renderNotFound(session)
	}
}

// func blogPagesOld(resp http.ResponseWriter, req *http.Request) {
//
// 	session := newSession(resp, req)
// 	if req.Method == "GET" {
// 		blogView(session)
// 	} else if req.Method == "POST" {
// 		blogAction(session)
// 	} else {
// 		renderError(session, "Unknown HTTP Method", errors.New("HTTP method not supported"))
// 	}
// }

// func blogView(s session) {
// 	if id, err := parseBlogViewUrl(s.req.URL.Path); err != nil {
// 		renderError(s, "Cannot parse Blog URL", err)
// 	} else if id != 0 {
// 		blogViewOne(s, id)
// 	} else {
// 		blogViewAll(s)
// 	}
// }

func blogViewOne(s session, values map[string]string) {
	id := idFromString(values["id"])
	log.Println(values)
	if id == 0 {
		renderError(s, "No Blog ID was received", nil)
		return
	}

	log.Printf("Loading %d", id)
	blog, err := models.BlogGetById(id)
	if err != nil {
		renderError(s, "Fetching by ID", err)
		return
	}

	vm := viewModels.FromBlog(blog, s.toViewModel())
	renderTemplate(s, "views/blogView.html", vm)
}

func blogViewAll(s session, values map[string]string) {
	log.Printf("Loading all...")
	if blogs, err := models.BlogGetAll(); err != nil {
		renderError(s, "Error fetching all", err)
	} else {
		vm := viewModels.FromBlogs(blogs, s.toViewModel())
		renderTemplate(s, "views/blogList.html", vm)
	}
}

// func blogAction(s session) {
// 	id, action, err := parseBlogEditUrl(s.req.URL.Path)
// 	if err != nil {
// 		renderError(s, "Cannot determine HTTP action", err)
// 		return
// 	}
//
// 	if !s.isAuth() {
// 		renderNotAuthorized(s)
// 		return
// 	}
//
// 	if action == "new" {
// 		blogNew(s)
// 	} else if action == "edit" {
// 		blogEdit(s, id)
// 	} else if action == "save" {
// 		blogSave(s, id)
// 	} else if action == "post" {
// 		blogPost(s, id)
// 	} else if action == "draft" {
// 		blogDraft(s, id)
// 	} else {
// 		renderError(s, "Unknown action", nil)
// 	}
// }

func blogSave(s session, values map[string]string) {
	id := idFromString(values["id"])
	blog := blogFromForm(id, s)
	if err := blog.Save(); err != nil {
		renderError(s, fmt.Sprintf("Saving blog ID: %d"), err)
	} else {
		url := fmt.Sprintf("/blog/%s/%d", blog.Slug, id)
		log.Printf("Redirect to %s", url)
		http.Redirect(s.resp, s.req, url, 301)
	}
}

func blogNew(s session, values map[string]string) {
	newId, err := models.SaveNew()
	if err != nil {
		renderError(s, fmt.Sprintf("Error creating new blog"), err)
		return
	}
	log.Printf("Redirect to (edit for new) %d", newId)
	values["id"] = fmt.Sprintf("%d", newId)
	blogEdit(s, values)
}

func blogDraft(s session, values map[string]string) {
	id := idFromString(values["id"])
	if id == 0 {
		renderError(s, "No blog ID was received", nil)
		return
	}

	blog, err := models.MarkAsDraft(id)
	if err != nil {
		renderError(s, fmt.Sprintf("Mark as draft: %d", id), err)
		return
	}

	url := fmt.Sprintf("/blog/%s/%d", blog.Slug, id)
	log.Printf("Marked as draft: %s", url)
	http.Redirect(s.resp, s.req, url, 301)
}

func blogPost(s session, values map[string]string) {
	id := idFromString(values["id"])
	if id == 0 {
		renderError(s, "No blog ID was received", nil)
		return
	}

	blog, err := models.MarkAsPosted(id)
	if err != nil {
		renderError(s, fmt.Sprintf("Mark as posted: %d", id), err)
		return
	}

	url := fmt.Sprintf("/blog/%s/%d", blog.Slug, id)
	log.Printf("Mark as posted: %s", url)
	http.Redirect(s.resp, s.req, url, 301)
}

func blogEdit(s session, values map[string]string) {
	id := idFromString(values["id"])
	if id == 0 {
		renderError(s, "No blog ID was received", nil)
		return
	}

	log.Printf("Loading %d", id)
	blog, err := models.BlogGetById(id)
	if err != nil {
		renderError(s, fmt.Sprintf("Loading ID: %d", id), err)
		return
	}

	vm := viewModels.FromBlog(blog, s.toViewModel())
	renderTemplate(s, "views/blogEdit.html", vm)
}

func idFromString(str string) int64 {
	id, _ := strconv.ParseInt(str, 10, 64)
	return id
}

// func parseBlogViewUrl(url string) (id int64, err error) {
// 	if url == "/blog/" {
// 		return 0, nil
// 	}
// 	// url /blog/:title/:id
// 	// parts[0] empty
// 	// parts[1] blog
// 	// parts[2] title
// 	// parts[3] id
// 	parts := strings.Split(url, "/")
// 	if len(parts) == 4 && parts[0] == "" && parts[1] == "blog" {
// 		return idFromString(parts[3])
// 	}
// 	return 0, errors.New("Could not parse (view) blog URL")
// }
//
// func parseBlogEditUrl(url string) (id int64, action string, err error) {
// 	if url == "/blog/new" {
// 		return 0, "new", nil
// 	}
// 	// url /blog/:title/:id/:action
// 	// parts[0] empty
// 	// parts[1] blog
// 	// parts[2] title
// 	// parts[3] id
// 	// parts[4] action (edit, post, draft)
// 	parts := strings.Split(url, "/")
// 	if len(parts) == 5 && parts[0] == "" && parts[1] == "blog" {
// 		if id, err := idFromString(parts[3]); err != nil {
// 			return 0, "", err
// 		} else {
// 			action := parts[4]
// 			if action == "edit" || action == "save" ||
// 				action == "post" || action == "draft" {
// 				return id, action, nil
// 			}
// 			return 0, "", errors.New("Invalid action")
// 		}
// 	}
// 	return 0, "", errors.New("Could not parse (edit) blog URL")
// }

func blogFromForm(id int64, s session) models.Blog {
	var blog models.Blog
	blog.Id = id
	blog.Title = s.req.FormValue("title")
	blog.Summary = s.req.FormValue("summary")
	blog.Content = s.req.FormValue("content")
	return blog
}
