package web

import "testing"

var router Router

func init() {
	router.Add("GET", "/auth/login", dummyHandler)
	router.Add("GET", "/auth/logout", dummyHandler)
	router.Add("GET", "/auth/changepassword", dummyHandler)
	router.Add("POST", "/auth/login", dummyHandler)
	router.Add("POST", "/auth/changepassword", dummyHandler)
	router.Add("GET", "/blog/:title/:id", dummyHandler)
	router.Add("GET", "/blog", dummyHandler)
	router.Add("POST", "/blog/:title/:id/edit", dummyHandler)
	router.Add("POST", "/blog/new", dummyHandler)
}

func dummyHandler(s session, values map[string]string) {
}

func TestRoutes(t *testing.T) {
	// GET valid URLs
	tests := []string{
		"/blog/title/123",
		"/blog/title/123/",
		"/blog/title.xyz/123/",
		"/blog/",
		"/blog",
		"/auth/login",
		"/auth/logout",
		"/auth/changepassword",
	}
	for _, testUrl := range tests {
		if found, _ := router.FindRoute("GET", testUrl); !found {
			t.Errorf("Failed to find route for GET %s", testUrl)
		}
	}

	// POST valid URLs
	tests = []string{"/blog/t1/1/edit", "/blog/new",
		"/auth/login", "/auth/changepassword"}
	for _, testUrl := range tests {
		if found, _ := router.FindRoute("POST", testUrl); !found {
			t.Errorf("Failed to find route for POST %s", testUrl)
		}
	}

	// Invalid GET URLs
	tests = []string{"/blog/title/1/edit", "/blog/title/123/whatever",
		"/blog-", "/xyz/blog", "blog/title/1",
		"/auth/bad", "/auth/login/bad", "auth/login"}
	for _, testUrl := range tests {
		if found, _ := router.FindRoute("GET", testUrl); found {
			t.Errorf("Found an incorrect route for GET %s", testUrl)
		}
	}

	// Invalid POST URLs
	tests = []string{"/blog/title/123/edit/xxx", "/blog/newX",
		"/auth/bad", "/auth/login/auth", "auth/login"}
	for _, testUrl := range tests {
		if found, r := router.FindRoute("POST", testUrl); found {
			t.Errorf("Found an incorrect route for %s, %s", testUrl, r)
		}
	}

}

func TestValuesInRoutes(t *testing.T) {
	url := "/blog/the-title/123"
	_, route := router.FindRoute("GET", url)
	values := route.UrlValues(url)
	if values["title"] != "the-title" {
		t.Errorf(":title was not parsed correctly from URL %s", url)
	}
	if values["id"] != "123" {
		t.Errorf(":id was not parsed correctly from URL %s", url)
	}

	url = "/blog/the-title/123/edit"
	_, route = router.FindRoute("POST", url)
	values = route.UrlValues(url)
	if values["title"] != "the-title" {
		t.Errorf(":title was not parsed correctly from URL %s", url)
	}
	if values["id"] != "123" {
		t.Errorf(":id was not parsed correctly from URL %s", url)
	}
}
