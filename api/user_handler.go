package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kortbyhotel/reservation/data"
	"github.com/kortbyhotel/reservation/types"
)

type UserHandler struct {
    userStore data.UserStore
}

func NewUserHandler(userStore data.UserStore) *UserHandler {
    return &UserHandler{
        userStore: userStore,
    }
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
    var params types.CreateUserParams
    if err := c.BodyParser(&params); err != nil {
        return err
    }
    if err := params.Validate(); err != nil {
        return c.JSON(err)
    }
    user, err := types.NewUserFromParams(params)
    if err != nil {
        return err
    }
    insertedUser, err := h.userStore.CreateUser(c.Context(), user)
    if err != nil {
        return err
    }
    return c.JSON(insertedUser)
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
    var (
        id = c.Params("id")
    )
    user, err := h.userStore.GetUserByID(c.Context(), id)
    if err != nil {
        return err
    }
    return c.JSON(user)
} 

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
    users, err := h.userStore.GetUsers(c.Context())
    if err != nil {
        return err
    }
    return c.JSON(users)
} 
