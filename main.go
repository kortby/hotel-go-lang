package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kortbyhotel/reservation/api"
)

func main() {
    // listenAddr := flag.String("listenAddr", ":5000", "The listen address of the API server")
    // flag.Parse()
    app := fiber.New()
    apiv1 := app.Group("/api/v1") // /api

    apiv1.Get("/users", api.HandleGetUsers)
    apiv1.Get("/users/:id", api.HandleGetUsers)

    app.Listen(":3000")
}
