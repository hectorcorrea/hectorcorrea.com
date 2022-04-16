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

	if err := models.InitDB(); err != nil {
		log.Print("ERROR: Failed to initialize database: ", err)
	}
	log.Printf("Database: %s", models.DbInfo())

	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/favicon.ico", fs)
	http.Handle("/robots.txt", fs)
	http.Handle("/public/", http.StripPrefix("/public/", fs))
	http.HandleFunc("/auth/", authPages)
	http.HandleFunc("/blog/", blogRoutes)
	http.HandleFunc("/Blog/", blogRoutes)
	http.HandleFunc("/draft/", blogRoutes)
	http.HandleFunc("/page/", blogRoutes)
	http.HandleFunc("/about", blogRoutes)
	http.HandleFunc("/", homePage)

	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal("Failed to start the web server: ", err)
	}
}

func homePage(resp http.ResponseWriter, req *http.Request) {
	s := newSession(resp, req)
	if req.URL.Path != "/" {
		cacheResponse(resp)
		renderNotFound(s)
		return
	}

	blog, err := models.BlogGetBySlug("home")
	vm := viewModels.FromBlog(blog, s.toViewModel())
	if err != nil {
		renderError(s, "Home page not found", nil)
		return
	}
	renderTemplate(s, "views/home.html", vm)
}

func cacheResponse(resp http.ResponseWriter) {
	fiveMinutes := time.Minute * 5
	later := time.Now().Add(fiveMinutes)
	cacheControl := fmt.Sprintf("public, max-age=%.f", time.Duration(fiveMinutes).Seconds())
	resp.Header().Add("Cache-Control", cacheControl)
	resp.Header().Add("Expires", later.UTC().String())
}

func renderNotFound(s session) {
	path := s.req.URL.Path
	log.Printf(fmt.Sprintf("Not found (%s)", path))
	t, err := template.New("layout").ParseFiles("views/layout.html", "views/notFound.html")
	if err != nil {
		log.Printf("Error rendering not found page :(")
		// perhaps render a hard coded string?
	} else {
		s.resp.WriteHeader(http.StatusNotFound)
		t.Execute(s.resp, s.toViewModel())
	}
}

func renderNotAuthorized(s session) {
	title := "Not Authorized"
	details := fmt.Sprintf("You are not authorized to perform this action: %s %s",
		s.req.Method, s.req.URL.Path)
	log.Printf(fmt.Sprintf("%s: %s", title, details))
	vm := viewModels.NewErrorFromStr(title, details, s.toViewModel())
	t, err := template.New("layout").ParseFiles("views/layout.html", "views/error.html")
	if err != nil {
		log.Printf("Error rendering not authorized page :(")
		// perhaps render a hard coded string?
	} else {
		s.resp.WriteHeader(http.StatusUnauthorized)
		t.Execute(s.resp, vm)
	}
}

func renderError(s session, title string, err error) {
	log.Printf("ERROR: %s - %s (%s)", title, err, s.req.URL.Path)
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

func loadTemplate(s session, viewName string) (*template.Template, error) {
	t, err := template.New("layout").ParseFiles("views/layout.html", viewName)
	if err != nil {
		renderError(s, fmt.Sprintf("Loading view %s", viewName), err)
		return nil, err
	} else {
		log.Printf("Loaded template %s (%s)", viewName, s.req.URL.Path)
		return t, nil
	}
}

func renderTemplate(s session, viewName string, viewModel interface{}) {
	t, err := loadTemplate(s, viewName)
	if err != nil {
		log.Printf("Error loading: %s, %s ", viewName, err)
	} else {
		err = t.Execute(s.resp, viewModel)
		if err != nil {
			log.Printf("Error rendering: %s, %s ", viewName, err)
		}
	}
}
