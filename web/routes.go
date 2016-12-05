package web

import (
	"regexp"
	"strings"
)

type RouteHandler func(session)

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

func NewRoute(method, path string, handler RouteHandler) Route {
	route := Route{method: method, path: path, handler: handler}
	if !strings.Contains(path, "/:") {
		// Route without tokens
		route.re = regexp.MustCompile(path)
		return route
	}

	// Store the tokens indicated in the path (e.g. :title, :id)
	// and a regEx to match them
	tokenRe := regexp.MustCompile("/:(\\w+)")
	pattern := path
	for _, token := range tokenRe.FindAllString(path, -1) {
		route.tokens = append(route.tokens, token)
		pattern = strings.Replace(pattern, token, "/(\\w+)", 1)
	}
	route.re = regexp.MustCompile(pattern)
	return route
}

func (r Route) IsMatch(method, url string) bool {
	return r.method == method && r.re.MatchString(url)
}

func (r Route) UrlValues(url string) map[string]string {
	values := make(map[string]string)
	matches := r.re.FindStringSubmatch(url)
	if len(matches)+1 == len(r.tokens) {
		for i, token := range r.tokens {
			values[token] = matches[i+1]
		}
	}
	return values
}
