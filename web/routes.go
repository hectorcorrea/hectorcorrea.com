package web

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

type RouteHandler func(session, map[string]string)

type Route struct {
	method  string // GET or POST
	path    string // /blog/:title/:id
	handler RouteHandler
	re      *regexp.Regexp // /blog/(\w+)/(\w+)
	tokens  []string       // [:title, :id]
}

type Router struct {
	Routes []Route
}

func (r *Router) Add(method, path string, handler RouteHandler) {
	route := NewRoute(method, path, handler)
	r.Routes = append(r.Routes, route)
}

func (r Router) FindRoute(method, url string) (bool, Route) {
	for _, route := range r.Routes {
		if route.IsMatch(method, url) {
			return true, route
		}
	}
	return false, Route{}
}

// Path should be in the form /xxx/:title/:id
// Values preceded by a colon (e.g. :id) are considered
// named tokens.
func NewRoute(method, path string, handler RouteHandler) Route {
	route := Route{method: method, path: path, handler: handler}
	if !strings.Contains(path, "/:") {
		// Route without tokens
		route.re = regexp.MustCompile("^" + path + "/??$")
		return route
	}

	// Store the tokens indicated in the path (e.g. :title, :id)
	// and a regEx to match them
	tokenRe := regexp.MustCompile("/:([\\w\\-\\._]+)")
	pattern := path
	for _, token := range tokenRe.FindAllString(path, -1) {
		route.tokens = append(route.tokens, token)
		pattern = strings.Replace(pattern, token, "/([\\w\\-\\._]+)", 1)
	}
	route.re = regexp.MustCompile("^" + pattern + "/??$")
	return route
}

func (r Route) IsMatch(method, url string) bool {
	return r.method == method && r.re.MatchString(url)
}

func (r Route) UrlValues(url string) map[string]string {
	values := make(map[string]string)
	matches := r.re.FindStringSubmatch(url)
	if len(matches) == len(r.tokens)+1 {
		for i, token := range r.tokens {
			key := token[2:] // "/:title" becomes "title"
			values[key] = matches[i+1]
		}
		log.Printf("set values %s", values)
	} else {
		log.Printf("got NO values: %s %d %d\r\n", url, len(matches), len(r.tokens))
	}
	return values
}

func (r Route) String() string {
	return fmt.Sprintf("%s %s", r.method, r.path)
}
