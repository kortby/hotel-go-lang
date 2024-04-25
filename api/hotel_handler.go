package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kortbyhotel/reservation/data"
	"go.mongodb.org/mongo-driver/bson"
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

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {

	hotels, err := h.hotelStore.GetHotels(c.Context(), nil)
	if err != nil {
		return err
	}
	return c.JSON(hotels)
}