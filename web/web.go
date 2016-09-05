package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"hectorcorrea.com/models"
	"hectorcorrea.com/viewModels"
)

// Web Request and Response wrapper
type webRR struct {
	resp http.ResponseWriter
	req  *http.Request
}

func newWebRR(resp http.ResponseWriter, req *http.Request) webRR {
	return webRR{resp: resp, req: req}
}

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

func homePage(resp http.ResponseWriter, req *http.Request) {
	rr := newWebRR(resp, req)
	viewName := viewForPath(req.URL.Path)
	if viewName != "" {
		t, err := loadTemplate(viewName)
		if err != nil {
			renderError(rr, fmt.Sprintf("Loading view %s", viewName), err)
		} else {
			t.Execute(resp, nil)
		}
	} else {
		renderError(rr, "Unknown path", nil)
	}
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

func renderError(rr webRR, title string, err error) {
	// TODO: log more about the Request
	// TODO: http 404 vs 500
	log.Printf("ERROR: %s - %s", title, err)
	vm := viewModels.NewError(title, err)
	t, err := template.New("layout").ParseFiles("views/layout.html", "views/error.html")
	if err != nil {
		log.Printf("Error rendering error page :(")
		// perhaps render a hard coded string?
	} else {
		t.Execute(rr.resp, vm)
	}
}

func loadTemplate(viewName string) (*template.Template, error) {
	return template.New("layout").ParseFiles("views/layout.html", viewName)
}
