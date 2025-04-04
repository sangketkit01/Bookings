package main

import (
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/sangketkit01/bookings/internal/config"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig
	mux := route(&app)

	switch mux.(type){
	case *chi.Mux : // do nothing
	default: 
		t.Error("type is not *chi.Mux")
	}
}