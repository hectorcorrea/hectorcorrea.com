package web

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"hectorcorrea.com/models"
	"hectorcorrea.com/viewModels"
)

func StartWebServer(address string) {
	log.Printf("Listening for requests at %s\n", "http://"+address)
	models.InitDB()
	log.Printf("Database: %s", models.SafeConnString())

	fs := http.FileServer(http.Dir("/Users/hector/dev/src/hectorcorrea.com/public"))
	http.Handle("/favicon.ico", fs)
	http.Handle("/robots.txt", fs)
	http.Handle("/public/", http.StripPrefix("/public/", fs))
	http.HandleFunc("/blog/", blogPages)
	http.HandleFunc("/", homePage)

	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal("Failed to start the web server: ", err)
	}
}

func blogPages(resp http.ResponseWriter, req *http.Request) {
	id, err := parseBlogUrl(req.URL.Path)
	if err != nil {
		t, _ := template.ParseFiles("views/error.html")
		t.Execute(resp, nil)
		log.Printf("\t ERROR: %s", err)
		return
	}

	if id == 0 {
		blogs, err := models.BlogGetAll()
		if err != nil {
			log.Printf("ERROR %s", err)
		}
		vm := viewModels.FromBlogs(blogs)
		t, _ := template.ParseFiles("views/blogList.html")
		t.Execute(resp, vm)
		return
	}

	log.Printf("Loading %d", id)
	blog, err := models.BlogGetById(id)
	if err != nil {
		log.Printf("ERROR %s", err)
	}
	t, _ := template.ParseFiles("views/blogView.html")
	vm := viewModels.FromBlog(blog)
	t.Execute(resp, vm)
}

func homePage(resp http.ResponseWriter, req *http.Request) {
	viewName := viewForPath(req.URL.Path)
	t, _ := template.ParseFiles(viewName)
	t.Execute(resp, nil)
}

func viewForPath(path string) string {
	var viewName string
	if path == "/about" || path == "/credits" {
		viewName = fmt.Sprintf("views%s.html", path)
	} else if path == "/" {
		viewName = "views/index.html"
	} else {
		viewName = "views/error.html"
	}
	return viewName
}

func parseBlogUrl(url string) (id int, err error) {
	if url == "/blog/" {
		return 0, nil
	}

	// url /blog/:title/:id
	// parts[0] empty
	// parts[1] blog
	// parts[2] title
	// parts[3] id
	parts := strings.Split(url, "/")
	if len(parts) == 4 && parts[1] == "blog" {
		return strconv.Atoi(parts[3])
	}

	return 0, errors.New("Could not parse blog URL")
}
