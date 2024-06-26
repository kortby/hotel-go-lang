package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/kortbyhotel/reservation/api"
	"github.com/kortbyhotel/reservation/data"
	"github.com/kortbyhotel/reservation/config"
	"github.com/kortbyhotel/reservation/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
    _ "github.com/joho/godotenv/autoload"
)

var fiberConfig = fiber.Config{
    ErrorHandler: func(ctx *fiber.Ctx, err error) error {
        if apiError, ok := err.(api.Error); ok {
            return ctx.Status(apiError.Code).JSON(apiError)
        }
        return api.NewError(http.StatusInternalServerError, err.Error())
    },
}

func main() {
    cfg := config.New()

    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.DBURI))
	if err != nil {
		log.Fatal(err)
	}
    userStore := data.NewMongoUserStore(client)
    hotelStore := data.NewMongoHotelStore(client)
    roomStore := data.NewMongoRoomStore(client, hotelStore)
    bookingStore := data.NewMongoBookingStore(client)
    store := &data.Store{
        Room: roomStore,
        Hotel: hotelStore,
        User: userStore,
        Booking: bookingStore,
    }
    authHandler := api.NewAuthHandler(userStore)
    userHandler := api.NewUserHandler(userStore)
    hotelHandler := api.NewHotelHandler(store)
    roomHandler := api.NewRoomHandler(store)
    bookingHandler := api.NewBookingHandler(store)

    // listenAddr := flag.String("listenAddr", ":5000", "The listen address of the API server")
    // flag.Parse()
    app := fiber.New(fiberConfig)
    api := app.Group("/api") // /api
    apiv1 := app.Group("/api/v1", middleware.JWTAuthentication(userStore)) // /api
    admin := apiv1.Group("/admin", middleware.AdminAuth) // /api

    // auth
    api.Post("/auth", authHandler.HandleAuthenticate)


    // users
    apiv1.Get("/users", userHandler.HandleGetUsers)
    apiv1.Post("/users", userHandler.HandlePostUser)
    apiv1.Get("/users/:id", userHandler.HandleGetUser)
    apiv1.Put("/users/:id", userHandler.HandlePutUser)
    apiv1.Delete("/users/:id", userHandler.HandleDeleteUser)

    // hotels
    apiv1.Get("/hotels", hotelHandler.HandleGetHotels)
    apiv1.Get("/hotels/:id", hotelHandler.HandleGetHotel)
    apiv1.Get("/hotels/:id/rooms", hotelHandler.HandleGetRooms)

    // booking
    apiv1.Get("/rooms", roomHandler.HandleGetRooms)
    apiv1.Post("/rooms/:id/book", roomHandler.HandleBookRoom)

    apiv1.Get("/bookings/:id", bookingHandler.HandleGetBooking)
    apiv1.Get("/bookings/:id/cancel", bookingHandler.HandleCancelBooking)
    // admin root booking handlers
    admin.Get("/bookings", bookingHandler.HandleGetBookings)

    app.Listen(":3000")
}
