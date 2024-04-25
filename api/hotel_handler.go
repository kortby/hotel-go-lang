package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kortbyhotel/reservation/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelHandler struct {
	hotelStore data.HotelStore
	roomStore data.RoomStore
}

func NewHotelHandler(hs data.HotelStore, rs data.RoomStore) *HotelHandler {
	return &HotelHandler{
		hotelStore: hs,
		roomStore: rs,
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
	rooms, err := h.roomStore.GetRooms(c.Context(), filter)
	if err != nil {
		return err
	}

	return c.JSON(rooms)
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {

	hotels, err := h.hotelStore.GetHotels(c.Context(), nil)
	if err != nil {
		return err
	}
	return c.JSON(hotels)
}