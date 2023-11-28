package main

import (
	"fmt"
	"github.com/easymomo/go-bookings/internal/config"
	"github.com/go-chi/chi/v5"
	"testing"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
		// Do nothing, test passed
	default:
		t.Error(fmt.Sprintf("Test failed, type is not *chi.Mux, it is %T", v))
	}
}
