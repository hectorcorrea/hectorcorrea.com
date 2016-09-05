package viewModels

import (
	"fmt"
	"html/template"

	"hectorcorrea.com/models"
)

type Blog struct {
	Id        int
	Title     string
	Summary   string
	Slug      string
	Url       string
	CreatedOn string
	PostedOn  string
	IsDraft   bool
	Html      template.HTML
}

type BlogList struct {
	Whatever string
	Blogs    []Blog
}

func FromBlog(blog models.Blog) Blog {
	var vm Blog
	vm.Id = blog.Id
	vm.Title = blog.Title
	vm.Summary = blog.Summary
	vm.Slug = blog.Slug
	vm.Url = fmt.Sprintf("/blog/%s/%d", blog.Slug, blog.Id)
	vm.Html = template.HTML(blog.Content)
	vm.CreatedOn = blog.CreatedOn
	vm.PostedOn = blog.PostedOn
	vm.IsDraft = (vm.PostedOn == "")
	return vm
}

func FromBlogs(blogs []models.Blog) BlogList {
	var list []Blog
	for _, blog := range blogs {
		list = append(list, FromBlog(blog))
	}
	return BlogList{Whatever: "whatever", Blogs: list}
}
