package main

import (
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Do some actions here

	// Then run the test and exit
	os.Exit(m.Run())
}

type myHandler struct{}

func (mh *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}
