package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/kortbyhotel/reservation/data/fixtures"
)


func TestAdminGetBookings(t *testing.T) {
	tdb := setup(t)
    defer tdb.teardown(t)

	user, _ := fixtures.AddUser(tdb.Store, "admin", "admin", false)
	hotel := fixtures.AddHotel(tdb.Store, "Marriot", "Dallas", 5, nil)
	room := fixtures.AddRoom(tdb.Store, "Small", false, 99, hotel.ID)
	booking := fixtures.AddBooking(tdb.Store, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 3))
	fmt.Println(booking)
}