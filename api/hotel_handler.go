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

type ResourceResp struct {
	Results int 	`json:"resutls"`
	Data 	any 	`json:"data"`
	Page 	int 	`json:"page"`
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
		return ErrResourceFound("room")
	}

	return c.JSON(rooms)
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	var pagination data.Pagination
	if err := c.QueryParser(&pagination); err != nil {
		return ErrBadRequest()
	}
	hotels, err := h.store.Hotel.GetHotels(c.Context(), nil, &pagination)
	if err != nil {
		return ErrInvalidID()
	}
	resp := ResourceResp{
		Results: len(hotels),
		Data: hotels,
		Page: int(pagination.Page),
	}
	return c.JSON(resp)
}

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {


	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	hotel, err := h.store.Hotel.GetHotelByID(c.Context(), oid)
	if err != nil {
		return ErrResourceFound("hotel")
	}
	return c.JSON(hotel)
}