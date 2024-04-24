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

const dburi = "mongodb://localhost:27017"
const dbname = "hotel-reservation"
// const userColl = "users"

var config = fiber.Config{
    ErrorHandler: func(ctx *fiber.Ctx, err error) error {
        return ctx.JSON(map[string]string{"error": err.Error()})
    },
}

func main() {
    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}
    userhandler := api.NewUserHandler(data.NewMongoUserStore(client, data.DBNAME))

    // listenAddr := flag.String("listenAddr", ":5000", "The listen address of the API server")
    // flag.Parse()
    app := fiber.New(config)
    apiv1 := app.Group("/api/v1") // /api

    apiv1.Get("/users", userhandler.HandleGetUsers)
    apiv1.Post("/users", userhandler.HandlePostUser)
    apiv1.Get("/users/:id", userhandler.HandleGetUser)
    apiv1.Put("/users/:id", userhandler.HandlePutUser)
    apiv1.Delete("/users/:id", userhandler.HandleDeleteUser)

    app.Listen(":3000")
}
