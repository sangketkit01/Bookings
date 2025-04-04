package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sangketkit01/bookings/internal/config"
	"github.com/sangketkit01/bookings/internal/driver"
	"github.com/sangketkit01/bookings/internal/forms"
	"github.com/sangketkit01/bookings/internal/helpers"
	"github.com/sangketkit01/bookings/internal/models"
	"github.com/sangketkit01/bookings/internal/render"
	"github.com/sangketkit01/bookings/internal/repository"
	"github.com/sangketkit01/bookings/internal/repository/dbrepo"
)


// Repo is the handlers repository
var Repo *Repository

// Repository holds app config and able handlers to access DatabaseRepo functions
type Repository struct {
	App *config.AppConfig
	DB repository.DataBaseRepo
}

// NewRepository creates a new repository
func NewRepository(app *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: app,
		DB: dbrepo.NewPostgresRepo(db.SQL,app),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home renders home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	m.DB.AllUsers()
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About renders about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

// Reservation renders the make a reservation page and displays form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})

	reservation, ok := m.App.Session.Get(r.Context(),"reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(),"error","Can't get reservation from session")
		http.Redirect(w,r,"/search-availability",http.StatusTemporaryRedirect)
		return
	}

	room, err := m.DB.GetRoomByID(reservation.RoomID)
	if err != nil{
		m.App.Session.Put(r.Context(),"error",err.Error())
		http.Redirect(w,r,"/search-availability",http.StatusTemporaryRedirect)
		return
	}

	reservation.Room = room

	m.App.Session.Put(r.Context(),"reservation",reservation)

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data["reservation"] = reservation

	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
		StringMap: stringMap,
	})
}

// PostReservation handlers the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {

	res, ok := m.App.Session.Get(r.Context(),"reservation").(models.Reservation)
	if !ok{
		helpers.ServerError(w,errors.New("Cannot get reservation from session"))
		return
	}

	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w,err)
		return
	}


	roomID , err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil{
		helpers.ServerError(w,err)
		return
	}


	reservation := &models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		StartDate: res.StartDate,
		EndDate: res.EndDate,
		RoomID: roomID,
		Room: res.Room,
	}

	form := forms.New(r.PostForm)
	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}


	newReservationID, err := m.DB.InsertReservation(*reservation)
	if err != nil{
		helpers.ServerError(w,err)
		return
	}

	restriction := models.RoomRestriction{
		StartDate: res.StartDate,
		EndDate: res.EndDate,
		RoomID: roomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}


	stringMap := make(map[string]string)
	stringMap["start_date"] = res.StartDate.Format("2006-01-02")
	stringMap["end_date"] = res.EndDate.Format("2006-01-02")

	err = m.DB.InsertRoomRestriction(restriction)

	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// Generals renders generals-quarters page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// Majors renders major's suit page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// SearchAvailability renders search-availability page
func (m *Repository) SearchAvailability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostSearchAvailability renders the search availability page
func (m *Repository) PostSearchAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout,start)
	if err != nil{
		helpers.ServerError(w,err)
		return
	}

	endDate, err := time.Parse(layout,end)
	if err != nil{
		helpers.ServerError(w,err)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate,endDate)
	if err != nil{
		helpers.ServerError(w,err)
		return
	}

	for _, room := range rooms{
		m.App.InfoLog.Println("ROOM", room.ID, room.RoomName)
	}

	if len(rooms) == 0{
		m.App.Session.Put(r.Context(),"error","No available room")
		http.Redirect(w,r,"/search-availability",http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms
	
	res := models.Reservation{
		StartDate: startDate,
		EndDate: endDate,
	}

	m.App.Session.Put(r.Context(),"reservation",res)

	render.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
	RoomID string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate string `json:"end_date"`
}


// SearchAvailabilityJSON returns room availability in JSON format
func (m *Repository) SearchAvailabilityJSON(w http.ResponseWriter, r *http.Request) {

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")


	layout := "2006-01-02"
	startDate, _ := time.Parse(layout,sd)
	endDate, _ := time.Parse(layout,ed)

	roomID, _ := strconv.Atoi(r.Form.Get("room_id"))

	available, _ := m.DB.SearchAvailabilityByDatesByRoomID(startDate,endDate,roomID)
	resp := jsonResponse{
		OK:      available,
		Message: "",
		StartDate: sd,
		EndDate: ed,
		RoomID: strconv.Itoa(roomID),
	}

	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		helpers.ServerError(w,err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Contact renders the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}


// ReservationSummary renders reservation summary page after inserted a new reservation
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	m.App.InfoLog.Println(reservation.Room)
	if !ok{
		m.App.ErrorLog.Println("Can't get error from session")
		log.Println("Cannot get item from session")
		m.App.Session.Put(r.Context(),"error","Can't get reservation from session")
		http.Redirect(w,r,"/",http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation

	location, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	startDate := reservation.StartDate.In(location)
	endDate := reservation.EndDate.In(location)

	stringMap := make(map[string]string)
	stringMap["start_date"] = startDate.Format("2006-01-02")
	stringMap["end_date"] = endDate.Format("2006-01-02")

	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
		StringMap: stringMap,
	})
}

// ChooseRoom displays list of available rooms
func (m *Repository) ChooseRoom(w http.ResponseWriter,r *http.Request){
	roomID, err := strconv.Atoi(chi.URLParam(r,"id"))
	if err != nil{
		helpers.ServerError(w,err)
		return
	}

	res,ok := m.App.Session.Get(r.Context(),"reservation").(models.Reservation)
	if !ok{
		helpers.ServerError(w,err)
		return
	}

	res.RoomID = roomID

	m.App.Session.Put(r.Context(),"reservation",res)
	
	http.Redirect(w,r,"/reservation",http.StatusSeeOther)
}

// BookRoom take URL parameters, builds a sessional variable, and takes user to make reservation screen
func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request){
	roomID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout,sd)
	endDate, _ := time.Parse(layout,ed)


	var reservation models.Reservation
	reservation.RoomID = roomID
	reservation.StartDate = startDate
	reservation.EndDate = endDate

	room, err  := m.DB.GetRoomByID(roomID)
	if err != nil{
		helpers.ServerError(w,err)
		return
	}

	reservation.Room = room
	reservation.RoomID = room.ID

	m.App.Session.Put(r.Context(),"reservation",reservation)

	http.Redirect(w,r,"/reservation",http.StatusSeeOther)
}
