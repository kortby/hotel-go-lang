package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/kortbyhotel/reservation/api"
	"github.com/kortbyhotel/reservation/data"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
    ErrorHandler: func(ctx *fiber.Ctx, err error) error {
        return ctx.JSON(map[string]string{"error": err.Error()})
    },
}

func main() {
    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(data.DBURI))
	if err != nil {
		log.Fatal(err)
	}
    userStore := data.NewMongoUserStore(client)
    hotelStore := data.NewMongoHotelStore(client)
    roomStore := data.NewMongoRoomStore(client, hotelStore)
    store := &data.Store{
        Room: roomStore,
        Hotel: hotelStore,
        User: userStore,
    }
    userHandler := api.NewUserHandler(userStore)
    hotelhandler := api.NewHotelHandler(store)

    // listenAddr := flag.String("listenAddr", ":5000", "The listen address of the API server")
    // flag.Parse()
    app := fiber.New(config)
    apiv1 := app.Group("/api/v1") // /api

    // users
    apiv1.Get("/users", userHandler.HandleGetUsers)
    apiv1.Post("/users", userHandler.HandlePostUser)
    apiv1.Get("/users/:id", userHandler.HandleGetUser)
    apiv1.Put("/users/:id", userHandler.HandlePutUser)
    apiv1.Delete("/users/:id", userHandler.HandleDeleteUser)

    // hotels
    apiv1.Get("/hotels", hotelhandler.HandleGetHotels)
    apiv1.Get("/hotels/:id", hotelhandler.HandleGetHotel)
    apiv1.Get("/hotels/:id/rooms", hotelhandler.HandleGetRooms)

    app.Listen(":3000")
}
