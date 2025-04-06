package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sangketkit01/bookings/internal/models"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"generals-quarters", "/generals-quarters", "GET", http.StatusOK},
	{"majors-suite", "/majors-suite", "GET", http.StatusOK},
	{"search-availability", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	{"make reservation pose","/reservation","POST",http.StatusOK},

	// {"post-search-availability","/search-availability","POST",[]postData{
	// 	{key:"start",value:"2025-04-01"},
	// 	{key: "end",value: "2025-04-01"},
	// },http.StatusOK},
	// {"post-search-availability-json","/search-availability-json","POST",[]postData{
	// 	{key:"start",value:"2025-04-01"},
	// 	{key: "end",value: "2025-04-01"},
	// },http.StatusOK},
	// {"post-make-reservation-without-data","/reservation","POST",nil,http.StatusOK},
	// {"post-make-reservation-missing-data","/reservation","POST",[]postData{
	// 	{key: "first_name",value: "John"},
	// },http.StatusOK},
}

func TestHandlers(t *testing.T){
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _,e := range theTests{
		resp , err := ts.Client().Get(ts.URL + e.url)
			if err != nil{
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode{
				t.Errorf("for %s expected %d but got %d",e.name,e.expectedStatusCode , resp.StatusCode)
			}
	}
}

func getSession() (*http.Request , error) {
	r , err := http.NewRequest("GET","/some-url",nil)
	if err != nil{
		return nil , err
	}

	ctx := r.Context()
	ctx , _ = session.Load(ctx,r.Header.Get("X-Session"))
	r = r.WithContext(ctx)

	return r,nil

}

func TestRepository_PostReservation(t *testing.T){
	reservation := models.Reservation{
		RoomID: 1,
		Room : models.Room{
			ID: 1,
			RoomName: "General's Quarters",
		},
	}

	reqBody := "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s",reqBody,"end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s",reqBody,"first_name=John")
	reqBody = fmt.Sprintf("%s&%s",reqBody,"last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s",reqBody,"email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s",reqBody,"phone=0123456789")
	reqBody = fmt.Sprintf("%s&%s",reqBody,"room_id=1")

	req, _ := http.NewRequest("POST","/reservation",strings.NewReader(reqBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	session.Put(ctx,"reservation",reservation)

	req.Header.Set("Content-Type","application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handlers := http.HandlerFunc(Repo.PostReservation)

	handlers.ServeHTTP(rr,req)
	if rr.Code != http.StatusSeeOther{
		t.Errorf("PostReservation handler returned wrong response code: got %d, wanted %d",rr.Code,http.StatusSeeOther)
	}

	// test for no session
	req, _ = http.NewRequest("POST","/reservation",strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type","application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handlers.ServeHTTP(rr,req)
	if rr.Code != http.StatusTemporaryRedirect{
		t.Errorf("PostReservation handler returned wrong response code: got %d, wanted %d",rr.Code,http.StatusTemporaryRedirect)
	}

	// test with no form
	req, _ = http.NewRequest("POST","/reservation",nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	session.Put(ctx,"reservation",reservation)

	handlers = http.HandlerFunc(Repo.PostReservation)

	req.Header.Set("Content-Type","application/x-www-form-urlencoded")

	handlers.ServeHTTP(rr,req)
	if rr.Code != http.StatusTemporaryRedirect{
		t.Errorf("PostReservation handler returned wrong response code for missing post body: got %d, wanted %d",rr.Code,http.StatusTemporaryRedirect)
	}

	// test for wrong terms of room id
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s",reqBody,"end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s",reqBody,"first_name=John")
	reqBody = fmt.Sprintf("%s&%s",reqBody,"last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s",reqBody,"email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s",reqBody,"phone=0123456789")
	reqBody = fmt.Sprintf("%s&%s",reqBody,"room_id=haha")

	req, _ = http.NewRequest("POST","/reservation",strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type","application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	session.Put(ctx,"reservation",reservation)

	handlers.ServeHTTP(rr,req)
	if rr.Code != http.StatusTemporaryRedirect{
		t.Errorf("PostReservation handler returned wrong response code: got %d, wanted %d",rr.Code,http.StatusTemporaryRedirect)
	}

	// test invalid from (not send last_name)
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s",reqBody,"end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s",reqBody,"first_name=J")
	reqBody = fmt.Sprintf("%s&%s",reqBody,"email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s",reqBody,"phone=0123456789")
	reqBody = fmt.Sprintf("%s&%s",reqBody,"room_id=1")

	req, _ = http.NewRequest("POST","/reservation",strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type","application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	session.Put(ctx,"reservation",reservation)

	handlers.ServeHTTP(rr,req)
	if rr.Code != http.StatusOK{
		t.Errorf("PostReservation handler returned wrong response code: got %d, wanted %d",rr.Code,http.StatusOK)
	}

}

func TestRepository_Reservation(t *testing.T){
	reservation := models.Reservation{
		RoomID: 1,
		Room : models.Room{
			ID: 1,
			RoomName: "General's Quarters",
		},
	}

	req, _ := http.NewRequest("GET","/reservation",nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation",reservation)

	handler := http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(rr,req)

	if rr.Code != http.StatusOK{
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d",rr.Code,http.StatusOK)
	}

	// test case where reservation is not session (reset everything)
	req, _ = http.NewRequest("GET","/reservation",nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr,req)
	if rr.Code != http.StatusTemporaryRedirect{
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d",rr.Code,http.StatusTemporaryRedirect)
	}

	// test with non-existent room
	req, _ = http.NewRequest("GET","/reservation",nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	reservation.RoomID = 100
	session.Put(ctx,"reservation",reservation)

	handler.ServeHTTP(rr,req)
	if rr.Code != http.StatusTemporaryRedirect{
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d",rr.Code,http.StatusTemporaryRedirect)
	}
}

func getCtx(req *http.Request) context.Context{
	ctx, err := session.Load(req.Context(),req.Header.Get("X-Session"))
	if err != nil{
		log.Println(err)
	}

	return ctx
}

func TestRepository_AvailabilityJSON(t *testing.T){
	// fisrt case - rooms are not available
	reqBody := "start=2025-01-01"
	reqBody = fmt.Sprintf("%s&%s",reqBody,"end=2050-01-01")
	reqBody= fmt.Sprintf("%s&%s",reqBody,"room_id=1")

	req, _ := http.NewRequest("POST","/search-availability-json",strings.NewReader(reqBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type","x-www-form-urlencoded")
	handler := http.HandlerFunc(Repo.SearchAvailabilityJSON)

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr,req)

	var j jsonResponse
	err := json.Unmarshal([]byte(rr.Body.String()),&j)
	if err != nil{
		t.Error("Failed to parse json")
	}
}