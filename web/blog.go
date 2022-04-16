package web

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"hectorcorrea.com/models"
	"hectorcorrea.com/viewModels"
)

var blogRouter Router

func init() {
	blogRouter.Add("GET", "/about", blogAbout)
	blogRouter.Add("GET", "/page", blogPages)
	blogRouter.Add("GET", "/draft", blogDrafts)
	blogRouter.Add("GET", "/blog/rss", blogRss)
	blogRouter.Add("GET", "/blog/:title/:id", blogViewOne)
	blogRouter.Add("GET", "/Blog/:title", blogLegacyOne)
	blogRouter.Add("GET", "/blog/:title", blogLegacyOne)
	blogRouter.Add("GET", "/blog", blogPosts)
	blogRouter.Add("POST", "/blog/:title/:id/edit", blogEdit)
	blogRouter.Add("POST", "/blog/:title/:id/save", blogSave)
	blogRouter.Add("POST", "/blog/:title/:id/post", blogMarkPosted)
	blogRouter.Add("POST", "/blog/:title/:id/draft", blogMarkDraft)
	blogRouter.Add("POST", "/blog/new", blogNew)
}

func blogRoutes(resp http.ResponseWriter, req *http.Request) {
	session := newSession(resp, req)
	found, route := blogRouter.FindRoute(req.Method, req.URL.Path)
	if found {
		values := route.UrlValues(req.URL.Path)
		route.handler(session, values)
	} else {
		renderNotFound(session)
	}
}

func blogViewOne(s session, values map[string]string) {
	id := values["id"]
	log.Printf("Loading %s", id)
	blog, err := models.BlogGetById(id)
	if err != nil {
		blog, err = models.BlogGetByOldId(id)
		if err != nil {
			renderError(s, "Fetching by ID", err)
			return
		}
		newUrl := fmt.Sprintf("/blog/%s/%s", blog.Slug, blog.Id)
		log.Printf("Legacy blog Redirected to %s", newUrl)
		http.Redirect(s.resp, s.req, newUrl, http.StatusMovedPermanently)
		return
	}

	if blog.IsDraft() && !s.isAuth() {
		renderNotAuthorized(s)
		return
	}

	vm := viewModels.FromBlog(blog, s.toViewModel())
	renderTemplate(s, "views/blogView.html", vm)
}

// The About page
func blogAbout(s session, values map[string]string) {
	blog, err := models.BlogGetBySlug("about")
	if err == nil {
		vm := viewModels.FromBlog(blog, s.toViewModel())
		renderTemplate(s, "views/about.html", vm)
	} else {
		renderError(s, "About page not found", nil)
	}
}

func blogPosts(s session, values map[string]string) {
	log.Printf("Loading posts...")
	blogs := models.BlogGetPosts()
	vm := viewModels.FromBlogs(blogs, s.toViewModel())
	renderTemplate(s, "views/blogList.html", vm)
}

func blogDrafts(s session, values map[string]string) {
	if !s.isAuth() {
		renderNotAuthorized(s)
		return
	}
	log.Printf("Loading drafts...")
	blogs := models.BlogGetDrafts()
	vm := viewModels.FromBlogs(blogs, s.toViewModel())
	renderTemplate(s, "views/blogList.html", vm)
}

func blogPages(s session, values map[string]string) {
	if !s.isAuth() {
		renderNotAuthorized(s)
		return
	}
	log.Printf("Loading pages...")
	blogs := models.BlogGetPages()
	vm := viewModels.FromBlogs(blogs, s.toViewModel())
	renderTemplate(s, "views/blogList.html", vm)
}

func blogSave(s session, values map[string]string) {
	if !s.isAuth() {
		renderNotAuthorized(s)
		return
	}

	id := values["id"]
	blog := blogFromForm(id, s)
	blog, err := blog.Save()
	if err != nil {
		renderError(s, fmt.Sprintf("Saving blog ID: %s", id), err)
		return
	}
	url := fmt.Sprintf("/blog/%s/%s", blog.Slug, blog.Id)
	log.Printf("Redirect to %s", url)
	http.Redirect(s.resp, s.req, url, 301)
}

func blogNew(s session, values map[string]string) {
	if !s.isAuth() {
		renderNotAuthorized(s)
		return
	}
	id, err := models.SaveNew()
	if err != nil {
		renderError(s, fmt.Sprintf("Error creating new blog"), err)
		return
	}
	log.Printf("Redirect to (edit for new) %s", id)
	values["id"] = id
	blogEdit(s, values)
}

func blogMarkDraft(s session, values map[string]string) {
	if !s.isAuth() {
		renderNotAuthorized(s)
		return
	}
	id := values["id"]
	blog, err := models.MarkAsDraft(id)
	if err != nil {
		renderError(s, fmt.Sprintf("Mark as draft: %s", id), err)
		return
	}

	url := fmt.Sprintf("/blog/%s/%s", blog.Slug, id)
	log.Printf("Marked as draft: %s", url)
	http.Redirect(s.resp, s.req, url, 301)
}

func blogMarkPosted(s session, values map[string]string) {
	if !s.isAuth() {
		renderNotAuthorized(s)
		return
	}
	id := values["id"]
	blog, err := models.MarkAsPosted(id)
	if err != nil {
		renderError(s, fmt.Sprintf("Mark as posted: %s", id), err)
		return
	}

	url := fmt.Sprintf("/blog/%s/%s", blog.Slug, id)
	log.Printf("Mark as posted: %s", url)
	http.Redirect(s.resp, s.req, url, 301)
}

func blogEdit(s session, values map[string]string) {
	if !s.isAuth() {
		renderNotAuthorized(s)
		return
	}
	id := values["id"]
	if id == "" {
		renderError(s, "No blog ID was received", nil)
		return
	}

	log.Printf("Loading %s", id)
	blog, err := models.BlogGetById(id)
	if err != nil {
		renderError(s, fmt.Sprintf("Loading ID: %s", id), err)
		return
	}

	vm := viewModels.FromBlog(blog, s.toViewModel())
	renderTemplate(s, "views/blogEdit.html", vm)
}

func blogRss(s session, values map[string]string) {
	title := "Hector Correa"
	desc := "Hector Correa's blog"
	url := "http://hectorcorrea.com"
	rss := models.NewRss(title, desc, url)

	blogs := models.BlogGetPosts()
	for _, blog := range blogs {
		rss.Add(blog.Title, blog.Summary, blog.URL(url), blog.PostedOn)
	}

	xml, err := rss.ToXml()
	if err != nil {
		renderError(s, "Error generating RSS feed", err)
		return
	}
	fmt.Fprint(s.resp, xml)
}

func blogLegacyOne(s session, values map[string]string) {
	oldSlug := values["title"]
	if oldSlug == "" {
		renderError(s, "No slug was received in legacy URL", nil)
		return
	}

	log.Printf("Handling legacy URL: %s", oldSlug)
	slug := strings.ToLower(oldSlug)
	if strings.HasSuffix(slug, ".aspx") {
		slug = slug[0 : len(slug)-5]
	}

	slug = strings.Replace(slug, ".", "-", -1)

	blog, err := models.BlogGetBySlug(slug)
	if err != nil {
		renderError(s, "Fetching legacy URL by slug", err)
		return
	}

	newUrl := fmt.Sprintf("/blog/%s/%s", blog.Slug, blog.Id)
	log.Printf("Legacy blog Redirected to %s", newUrl)
	http.Redirect(s.resp, s.req, newUrl, http.StatusMovedPermanently)
}

func blogFromForm(id string, s session) models.Blog {
	var blog models.Blog
	blog.Id = id
	blog.Title = s.req.FormValue("title")
	blog.Summary = s.req.FormValue("summary")
	blog.ContentMarkdown = s.req.FormValue("content")
	blog.Type = s.req.FormValue("type")
	return blog
}
