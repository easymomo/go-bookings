package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields are missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")

	r, _ = http.NewRequest("POST", "/whatever", nil)
	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("shows 'does not have required fields' when it does")
	}
}

func TestForm_Has(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	has := form.Has("whatever")
	if has {
		t.Error("form shows has field when it does not")
	}

	// Set up some values
	postedData = url.Values{} // Created a variable of url.Values{}, which is what PostForm is
	// Put a value in the form
	postedData.Add("a", "a")
	// Reinitialise the form variable, past it the form data
	form = New(postedData)

	has = form.Has("a")
	if !has {
		t.Error("shows form does not have field when it should")
	}
}

func TestForm_MinLength(t *testing.T) {
	postedValues := url.Values{}
	form := New(postedValues)

	form.MinLength("x", 10)
	if form.Valid() {
		t.Error("form shows min length for non-existing field")
	}

	// Testing the Get method of the errors.go
	// We look for a place where there should be an error in the form because of the data processed
	// and test for the presence of an error
	isError := form.Errors.Get("x")
	if isError == "" {
		t.Error("should have an error, but did not get one")
	}

	// Set up some values
	postedValues = url.Values{} // Created a variable of url.Values{}, which is what PostForm is
	// Put a value in the form
	postedValues.Add("some_field", "some value")
	// Reinitialise the form variable, past it the form data
	form = New(postedValues)

	form.MinLength("some_field", 100)
	if form.Valid() {
		t.Error("shows min length of 100 met when data is shorter")
	}

	postedValues = url.Values{}
	postedValues.Add("another_field", "some other value")
	form = New(postedValues)

	form.MinLength("another_field", 1)
	if !form.Valid() {
		t.Error("shows min length of 1 is not met when it is")
	}

	// Testing the Get method of the errors.go
	// We look for a place where there should be no error in the form because of the data processed
	// and test for the absence of an error
	isError = form.Errors.Get("x")
	if isError != "" {
		t.Error("should have no error, but did get one")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postedValues := url.Values{}
	form := New(postedValues)

	form.IsEmail("x")
	if form.Valid() {
		t.Error("form shows valid email for non existing field")
	}

	postedValues = url.Values{}
	postedValues.Add("email", "my@email.com")
	form = New(postedValues)

	form.IsEmail("email")
	if !form.Valid() {
		t.Error("got an invalid email when we should not have been")
	}

	postedValues = url.Values{}
	postedValues.Add("email", "g")
	form = New(postedValues)

	form.IsEmail("email")
	if form.Valid() {
		t.Error("got a valid email when we should not have been")
	}
}
