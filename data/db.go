package data

const (
	DBNAME = "hotel-reservation"
	TESTDBNAME = "hotel-reservation-test"
	DBURI = "mongodb://localhost:27017"
)

type Store struct {
	User UserStore
	Hotel HotelStore
	Room RoomStore
}
