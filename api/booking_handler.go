package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/kortbyhotel/reservation/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookingHandler struct {
	store *data.Store
}

func NewBookingHandler(store *data.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}


func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {

	bookings, err := h.store.Booking.GetBookings(c.Context(), nil)
	if err != nil {
		return err
	}
	return c.JSON(bookings)
}

func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	booking, err := h.store.Booking.GetBookingByID(c.Context(), oid)
	if err != nil {
		return err
	}
	user, err := GetAuthUser(c)
	if err != nil {
		return err
	}
	if booking.UserID != user.ID {
		return c.Status(http.StatusUnauthorized).JSON("Unauthorized")
	}
	return c.JSON(booking)
}

func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	booking, err := h.store.Booking.GetBookingByID(c.Context(), oid)
	if err != nil {
		return err
	}
	user, err := GetAuthUser(c)
	if err != nil {
		return err
	}
	if booking.UserID != user.ID {
		return c.Status(http.StatusUnauthorized).JSON("Unauthorized")
	}
	if err := h.store.Booking.UpdateBooking(c.Context(), oid, bson.M{"$set": bson.M{"cancelled": true}}); err != nil {
		return err
	}
	return c.JSON(map[string]string{"msg": "Successfully updated"})
}

