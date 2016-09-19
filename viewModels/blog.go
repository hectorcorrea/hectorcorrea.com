package viewModels

import (
	"fmt"
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
	vm.Url = fmt.Sprintf("/blog/%s/%d", blog.Slug, blog.Id)
	vm.Html = template.HTML(blog.Content)
	vm.CreatedOn = blog.CreatedOn
	vm.PostedOn = blog.PostedOn
	vm.UpdatedOn = blog.UpdatedOn
	vm.IsDraft = (vm.PostedOn == "")
	vm.Session = session
	return vm
}

func FromBlogs(blogs []models.Blog, session Session) BlogList {
	var list []Blog
	for _, blog := range blogs {
		list = append(list, FromBlog(blog, session))
	}
	return BlogList{Blogs: list, Session: session}
}
