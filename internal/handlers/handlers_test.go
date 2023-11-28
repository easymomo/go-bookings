package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"gq", "/generals-quarters", "GET", []postData{}, http.StatusOK},
	{"ms", "/majors-suite", "GET", []postData{}, http.StatusOK},
	{"sa", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"mr", "/make-reservation", "GET", []postData{}, http.StatusOK},

	{"post-search-availability", "/search-availability", "POST", []postData{
		{key: "start", value: "2020-01-01"},
		{key: "end", value: "2020-01-01"},
	}, http.StatusOK},
	{"post-search-availability", "/search-availability", "POST", []postData{
		{key: "start", value: "2020-01-01"},
		{key: "end", value: "2020-01-01"},
	}, http.StatusOK},
	{"post-search-availability", "/search-availability", "POST", []postData{
		{key: "first_name", value: "John"},
		{key: "last_name", value: "Smith"},
		{key: "email", value: "me@here.com"},
		{key: "phone", value: "555-555-5555"},
	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes() // returns an http handler

	// Create a web server that returns status code, something we can post to
	// We need to create a server, and we need to create a client that can post to that server
	ts := httptest.NewTLSServer(routes)
	// Close the server after we are done with the test
	// Whatever is after the keyword defer is not executed until after the current function
	// (TestHandlers in this case) is finished running
	defer ts.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			// We need to make a request as a client, as if we were a web browser
			// ts.Client() creates a web browser
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		} else {
			// Construct a variable that is in a format our server expects to receive
			// Create an empty variable in the form that is required by the method we are going to
			// call on our test server
			// url.Values{} builtin type part of the standard library that holds information as a POST request
			values := url.Values{}
			// Populate this with our entries
			for _, x := range e.params {
				values.Add(x.key, x.value)
			}
			// We need to make a request as a client, as if we were a web browser
			// ts.Client() creates a web browser
			resp, err := ts.Client().PostForm(ts.URL+e.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}
