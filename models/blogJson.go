package models

//
// import (
// 	"encoding/json"
// 	"html/template"
// 	"io/ioutil"
// )
//
// type Blog struct {
// 	Title     string
// 	CreatedOn string
// 	PostedOn  string
// 	Html      template.HTML
// }
//
// // blog - model
// // blogViewModel - template.html
// //
//
// func (blog *Blog) fromJson(json BlogJson, body []byte) {
// 	blog.Title = json.Title
// 	blog.Html = template.HTML(body)
// 	blog.CreatedOn = json.CreatedOn
// 	blog.PostedOn = json.PostedOn
// }
//
// type BlogJson struct {
// 	Key       int    `json:"key"`
// 	Title     string `json:"title"`
// 	Summary   string `json:"summary"`
// 	Url       string `json:"url"`
// 	CreatedOn string
// 	PostedOn  string `json:"postedOn"`
// }
//
// func Load(key string) (*Blog, error) {
// 	contentFile := "./data/" + key + ".html"
// 	body, err := ioutil.ReadFile(contentFile)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	metaFile := "./data/" + key + ".json"
// 	meta, err := ioutil.ReadFile(metaFile)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	var blogJson BlogJson
// 	json.Unmarshal(meta, &blogJson)
//
// 	var blog Blog
// 	blog.fromJson(blogJson, body)
// 	return &blog, nil
// }
