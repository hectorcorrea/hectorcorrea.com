package viewModels

import (
	"html/template"

	"hectorcorrea.com/models"
)

type Blog struct {
	Id        int64
	Title     string
	Summary   string
	Slug      string
	Url       string
	CreatedOn string
	PostedOn  string
	UpdatedOn string
	IsDraft   bool
	Html      template.HTML
	Markdown  string
	Session
}

type BlogList struct {
	Blogs []Blog
	Session
}

func FromBlog(blog models.Blog, session Session) Blog {
	var vm Blog
	vm.Id = blog.Id
	vm.Title = blog.Title
	vm.Summary = blog.Summary
	vm.Slug = blog.Slug
	vm.Url = blog.URL("")
	vm.Html = template.HTML(blog.ContentHtml)
	vm.Markdown = blog.ContentMarkdown
	vm.CreatedOn = blog.CreatedOn
	vm.PostedOn = blog.PostedOn
	vm.UpdatedOn = blog.UpdatedOn
	vm.IsDraft = (vm.PostedOn == "")
	vm.Session = session

	vm.Session.TwitterCard = true
	vm.Session.TwitterTitle = blog.Title
	vm.Session.TwitterDescription = blog.Summary

	return vm
}

func FromBlogs(blogs []models.Blog, session Session) BlogList {
	var list []Blog
	for _, blog := range blogs {
		list = append(list, FromBlog(blog, session))
	}
	return BlogList{Blogs: list, Session: session}
}
