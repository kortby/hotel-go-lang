package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kortbyhotel/reservation/types"
)

func HandleGetUsers(c *fiber.Ctx) error {
    u := types.User{
        FirstName: "John",
        LastName: "Doe",
    }
    return c.JSON(u)
} 

func HandleGetUser(c *fiber.Ctx) error {
    return c.JSON(map[string]string{"user" : "one user John Doe! Thanks"})
} 
