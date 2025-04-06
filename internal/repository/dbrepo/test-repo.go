package dbrepo

import (
	"errors"
	"time"

	"github.com/sangketkit01/bookings/internal/models"
)

// AllUsers() nothing
func (m *testDBRepo) AllUsers() bool{
	return true
}

// InsertReservation inserts a reservation into the database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int,error){

	return 1,nil
}

// InsertRoomRestriction inserts a room restriction into the database
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	return nil
}

// SearchAvailabilityByDatesByRoomID returns if availablity exists for roomID, and false if no availability exists
func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start , end time.Time, roomID int) (bool , error){
	return false, nil
}

// SearchAvailabilityForAllRooms return a slice of available rooms, if any, for given date range
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error){

	var rooms []models.Room
	return rooms,nil
}

// GetRoomByID gets the room by id
func (m *testDBRepo) GetRoomByID(id int) (models.Room, error){
	var room models.Room
	if id > 2 {
		return room, errors.New("Some error")
	}
	return room,nil
}


// GetUserByID returns a user by id
func (m *testDBRepo) GetUserByID(id int) (models.User, error){
	var user models.User

	return user,nil
}

// UpdateUser updates user in the database
func (m *testDBRepo) UpdateUser(u models.User) (error) {
	return nil
}

// Authenticate authenticates a user
func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {

	return 0, "", nil
}

// AllReservations returns a slice of all reservatins
func (m *testDBRepo) AllReservations() ([]models.Reservation, error){
	return []models.Reservation{},nil
}

// AllNewReservations returns a slice of all new reservatins
func (m *testDBRepo) AllNewReservations() ([]models.Reservation, error){
	return []models.Reservation{}, nil
}

// GetReservationByID returns one reservation by id
func (m *testDBRepo) GetReservationByID(id int) (models.Reservation, error){
	return models.Reservation{}, nil
}

// UpdateReservation updates a reservation in the database
func (m *testDBRepo) UpdateReservation(u models.Reservation) (error) {
	return nil
}

// DeleteReservation deletes a reservation by id from the database
func (m *testDBRepo) DeleteReservation(id int ) error{
	return nil
}

// UpdateProcessedForReservation updates processed for a reservation by id
func (m *testDBRepo) UpdateProcessedForReservation(id, processed int) error {
	return nil
}

func (m* testDBRepo) AllRooms() ([]models.Room, error){
	

	return []models.Room{}, nil
}


// GetRestrictionsForRoomByDate returns restrictions for a room by date range
func (m *testDBRepo) GetRestrictionsForRoomByDate(roomID int, start, end time.Time) ([]models.RoomRestriction, error){

	return []models.RoomRestriction{}, nil
}

// InsertBlockForRoom inserts  a room restriction
func (m *testDBRepo) InsertBlockForRoom(id int, startDate time.Time) error{
	return nil
}

// DeleteBlockByID deletes  a room restriction
func (m *testDBRepo) DeleteBlockByID(id int) error{

	return nil
}