package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kortbyhotel/reservation/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelHandler struct {
	store *data.Store
}

func NewHotelHandler(store *data.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

type HotelQueryPrams struct {
	Rooms bool 
	Rating int
}

func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"hotelID": oid}
	rooms, err := h.store.Room.GetRooms(c.Context(), filter)
	if err != nil {
		return err
	}

	return c.JSON(rooms)
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {

	hotels, err := h.store.Hotel.GetHotels(c.Context(), nil)
	if err != nil {
		return err
	}
	return c.JSON(hotels)
}

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {


	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	hotel, err := h.store.Hotel.GetHotelByID(c.Context(), oid)
	if err != nil {
		return err
	}
	return c.JSON(hotel)
}