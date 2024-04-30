package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kortbyhotel/reservation/api"
	"github.com/kortbyhotel/reservation/data/fixtures"
	"github.com/kortbyhotel/reservation/middleware"
	"github.com/kortbyhotel/reservation/types"
)

func TestUserGetBooking(t *testing.T) {
	tdb := setup(t)
    defer tdb.teardown(t)

	noneAuthUser, _ := fixtures.AddUser(tdb.Store, "guest", "guest", false)
	user, _ := fixtures.AddUser(tdb.Store, "test", "test", false)
	hotel := fixtures.AddHotel(tdb.Store, "Marriot", "Dallas", 5, nil)
	room := fixtures.AddRoom(tdb.Store, "Small", false, 99, hotel.ID)
	booking := fixtures.AddBooking(tdb.Store, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 3))
	_ = booking

	app := fiber.New()
	route := app.Group("/:id", middleware.JWTAuthentication(tdb.User))
	bookingHandler := api.NewBookingHandler(tdb.Store)
	route.Get("/", bookingHandler.HandleGetBooking)
	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	token, err := api.CreateTokenFromUser(user)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("X-Api-Token", token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	var bookingResp *types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookingResp); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(bookingResp.ID, booking.ID) {
		t.Fatal("expected booking to be equals")
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("non 200 response %d", resp.StatusCode)
	}

	/// none auth user
	req = httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	token, err = api.CreateTokenFromUser(noneAuthUser)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("X-Api-Token", token)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Fatalf("%d shouldn't be 200 response", resp.StatusCode)
	}

}


func TestAdminGetBookings(t *testing.T) {
	tdb := setup(t)
    defer tdb.teardown(t)

	user, _ := fixtures.AddUser(tdb.Store, "test", "test", false)
	adminUser, _ := fixtures.AddUser(tdb.Store, "admin", "admin", true)
	hotel := fixtures.AddHotel(tdb.Store, "Marriot", "Dallas", 5, nil)
	room := fixtures.AddRoom(tdb.Store, "Small", false, 99, hotel.ID)
	booking := fixtures.AddBooking(tdb.Store, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 3))
	_ = booking

	app := fiber.New()
	admin := app.Group("/", middleware.JWTAuthentication(tdb.User), middleware.AdminAuth)
	bookingHandler := api.NewBookingHandler(tdb.Store)
	admin.Get("/", bookingHandler.HandleGetBookings)
	req := httptest.NewRequest("GET", "/", nil)
	token, err := api.CreateTokenFromUser(adminUser)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("X-Api-Token", token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("non 200 response %d", resp.StatusCode)
	}
	var bookings []*types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}
	if len(bookings) <= 1 {
		t.Fatalf("expected 1 or more booking  but we got %d", len(bookings))
	}
	if !reflect.DeepEqual(booking.ID, bookings[len(bookings) - 1].ID) {
		t.Fatal("expected booking to be equals")
	}


	/// none admin cannot access
	req = httptest.NewRequest("GET", "/", nil)
	token, err = api.CreateTokenFromUser(user)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("X-Api-Token", token)
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == http.StatusUnauthorized {
		t.Fatalf("expected status unauthorized be we got %d", resp.StatusCode)
	}
}