package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"hectorcorrea.com/models"
	"hectorcorrea.com/viewModels"
)

func StartWebServer(address string) {
	log.Printf("Listening for requests at %s\n", "http://"+address)
	models.InitDB()
	log.Printf("Database: %s", models.DbConnStringSafe())

	fs := http.FileServer(http.Dir("/Users/hector/dev/src/hectorcorrea.com/public"))
	http.Handle("/favicon.ico", fs)
	http.Handle("/robots.txt", fs)
	http.Handle("/public/", http.StripPrefix("/public/", fs))
	http.HandleFunc("/auth/", authPages)
	http.HandleFunc("/blog/", blogPages)
	http.HandleFunc("/", staticPages)

	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal("Failed to start the web server: ", err)
	}
}

func staticPages(resp http.ResponseWriter, req *http.Request) {
	session := newSession(resp, req)
	vm := session.toViewModel()
	viewName := viewForPath(req.URL.Path)
	if viewName != "" {
		t, err := loadTemplate(viewName)
		if err != nil {
			renderError(session, fmt.Sprintf("Loading view %s", viewName), err)
		} else {
			cacheResponse(resp)
			t.Execute(resp, vm)
		}
	} else {
		cacheResponse(resp)
		renderNotFound(session)
	}
}

func cacheResponse(resp http.ResponseWriter) {
	fiveMinutes := time.Minute * 5
	later := time.Now().Add(fiveMinutes)
	cacheControl := fmt.Sprintf("public, max-age=%.f", time.Duration(fiveMinutes).Seconds())
	resp.Header().Add("Cache-Control", cacheControl)
	resp.Header().Add("Expires", later.UTC().String())
}

func viewForPath(path string) string {
	var viewName string
	if path == "/about" || path == "/credits" {
		viewName = fmt.Sprintf("views%s.html", path)
	} else if path == "/" {
		viewName = "views/index.html"
	}
	return viewName
}

func renderNotFound(s session) {
	// TODO: log more about the Request
	log.Printf("Not found")
	t, err := template.New("layout").ParseFiles("views/layout.html", "views/notFound.html")
	if err != nil {
		log.Printf("Error rendering not found page :(")
		// perhaps render a hard coded string?
	} else {
		s.resp.WriteHeader(http.StatusNotFound)
		t.Execute(s.resp, s.toViewModel())
	}
}

func renderError(s session, title string, err error) {
	// TODO: log more about the Request
	log.Printf("ERROR: %s - %s", title, err)
	vm := viewModels.NewError(title, err, s.toViewModel())
	t, err := template.New("layout").ParseFiles("views/layout.html", "views/error.html")
	if err != nil {
		log.Printf("Error rendering error page :(")
		// perhaps render a hard coded string?
	} else {
		s.resp.WriteHeader(http.StatusInternalServerError)
		t.Execute(s.resp, vm)
	}
}

func loadTemplate(viewName string) (*template.Template, error) {
	return template.New("layout").ParseFiles("views/layout.html", viewName)
}
