package main

import (
	"net/http"

	"github.com/justinas/nosurf"
	"github.com/sangketkit01/bookings/internal/helpers"
)

// NoSurf adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// SessionLoad loads and save the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

// Auth checks if user has authenticate will serve incoming request, otherwise redirect to login page
func Auth(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
		if !helpers.IsAuthenticated(r){
			session.Put(r.Context(),"error","Login first")
			http.Redirect(w,r,"/user/login",http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w,r)
	})
}
