package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sangketkit01/bookings/internal/config"
	"github.com/sangketkit01/bookings/internal/handlers"
	"net/http"
)

func route(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/generals-quarters", handlers.Repo.Generals)
	mux.Get("/majors-suite", handlers.Repo.Majors)

	mux.Get("/user/login",handlers.Repo.ShowLogin)
	mux.Post("/user/login",handlers.Repo.PostLogin)
	mux.Get("/user/logout",handlers.Repo.Logout)

	mux.Get("/search-availability", handlers.Repo.SearchAvailability)
	mux.Post("/search-availability", handlers.Repo.PostSearchAvailability)
	mux.Post("/search-availability-json", handlers.Repo.SearchAvailabilityJSON)
	mux.Get("/choose-room/{id}",handlers.Repo.ChooseRoom)
	mux.Get("/book-room",handlers.Repo.BookRoom)

	mux.Get("/reservation", handlers.Repo.Reservation)
	mux.Post("/reservation", handlers.Repo.PostReservation)
	mux.Get("/reservation-summary",handlers.Repo.ReservationSummary)

	mux.Get("/contact", handlers.Repo.Contact)

	fileServer := http.FileServer(http.Dir("../../stat/"))
	mux.Handle("/stat/*", http.StripPrefix("/stat", fileServer))

	mux.Route("/admin",func(mux chi.Router){
		//mux.Use(Auth)

		mux.Get("/dashboard",handlers.Repo.AdminDashBoard)
		mux.Get("/reservations-new",handlers.Repo.AdminNewReservations)
		mux.Get("/reservations-all",handlers.Repo.AdminAllReservations)
		mux.Get("/reservations-calendar",handlers.Repo.AdminReservationsCalendar)
		mux.Post("/reservations-calendar",handlers.Repo.AdminPostReservationsCalendar)
		mux.Get("/process-reservation/{src}/{id}",handlers.Repo.AdminProcessReservation)
		mux.Get("/delete-reservation/{src}/{id}",handlers.Repo.AdminDeleteReservation)

		mux.Get("/reservations/{src}/{id}",handlers.Repo.AdminShowReservation)
		mux.Post("/reservations/{src}/{id}",handlers.Repo.AdminPostShowReservation)
	})	
	
	return mux
}
