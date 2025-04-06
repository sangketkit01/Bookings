package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/sangketkit01/bookings/internal/config"
)

var app *config.AppConfig

// NewHelpers creates new Helpers
func NewHelpers(a *config.AppConfig){
	app = a
}

func ClientError(w http.ResponseWriter, status int){
	app.InfoLog.Println("Client error with the status of",status)
	http.Error(w, http.StatusText(status),status)
}

func ServerError(w http.ResponseWriter, err error){
	trace := fmt.Sprintf("%s\n%s",err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)

	http.Error(w,
		fmt.Sprintf("%s\n%s\n%s",http.StatusText(http.StatusInternalServerError),err.Error(),debug.Stack()),
		http.StatusInternalServerError,
	)
}

func IsAuthenticated(r *http.Request) bool {
	exists := app.Session.Exists(r.Context(),"user_id")
	return exists
}