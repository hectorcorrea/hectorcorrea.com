package viewModels

import (
	"html/template"

	"hectorcorrea.com/models"
)

type Blog struct {
	Title     string
	Summary   string
	Url       string
	CreatedOn string
	PostedOn  string
	Html      template.HTML
}

type BlogList struct {
	Whatever string
	Blogs    []Blog
}

func FromBlog(blog models.Blog) Blog {
	var vm Blog
	vm.Title = blog.Title
	vm.Summary = blog.Summary
	vm.Url = blog.Url
	vm.Html = template.HTML(blog.Content)
	vm.CreatedOn = blog.CreatedOn
	vm.PostedOn = blog.PostedOn
	return vm
}

func FromBlogs(blogs []models.Blog) BlogList {
	var list []Blog
	for _, blog := range blogs {
		list = append(list, FromBlog(blog))
	}
	return BlogList{Whatever: "whatever", Blogs: list}
}
