package main

import (
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var mh myHandler
	h := NoSurf(&mh)

	switch h.(type) {
	case http.Handler:
		// do nothing
	default:
		t.Error("type is not http.Handler")
	}
}

func TestSessionLoad(t *testing.T) {
	var mh myHandler
	h := SessionLoad(&mh)

	switch h.(type) {
	case http.Handler:
		// do nothing
	default:
		t.Error("type is not http.Handler")
	}
}