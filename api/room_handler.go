package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kortbyhotel/reservation/data"
	"github.com/kortbyhotel/reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookRoomParams struct {
	FromDate 	time.Time 	`json:"fromDate"`
	UntilDate 	time.Time 	`json:"untilDate"`
	NumPersons 	int			`json:"numPersons"`
}

type RoomHandler struct {
	store *data.Store
}

func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	rooms, err := h.store.Room.GetRooms(c.Context(), bson.M{})
	if err != nil {
		return ErrResourceFound("room")
	}

	return c.JSON(rooms)
}

func (p BookRoomParams) validate() error {
	now := time.Now()
	if now.After(p.FromDate) || now.After(p.UntilDate) {
		return fmt.Errorf("cannot book a room in the past")
	}
	if p.FromDate.After(p.UntilDate) {
		return fmt.Errorf("the checkout date connnot be before check in date")
	}
	return nil
}

func NewRoomHandler(store *data.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params BookRoomParams
	if err:= c.BodyParser(&params); err != nil {
		return err
	}
	if err := params.validate(); err != nil {
		return err
	}
	roomID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON("Internal error")
	}

	ok, err = h.isRoomAvailableForBooking(c.Context(), roomID, params) 
	if err != nil {
		return err
	}
	if !ok {
		return c.Status(http.StatusBadRequest).JSON("Room Already booked")
	}

	booking := types.Booking{
		UserID: user.ID,
		RoomID: roomID,
		FromDate: params.FromDate,
		UntilDate: params.UntilDate,
		NumPersons: params.NumPersons,
	}
	inserted, err := h.store.Booking.InsertBooking(c.Context(), &booking)
	if err != nil {
		return err
	}
	return c.JSON(inserted)
}

func (h *RoomHandler) isRoomAvailableForBooking(ctx context.Context,roomID primitive.ObjectID, params BookRoomParams) (bool, error) {
	filter := bson.M{
		"roomID": roomID,
		"fromDate": bson.M{
			"$gte":params.FromDate,
		},
		"untilDate": bson.M{
			"$lte": params.UntilDate,
		},
	}
	bookings, err := h.store.Booking.GetBookings(ctx, filter)
	if err != nil {
		return false, ErrResourceFound("bookings")
	} 

	return len(bookings) == 0, nil
}