package api

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/kortbyhotel/reservation/data"
	"github.com/kortbyhotel/reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
        return ErrBadRequest()
    }
    validationErrors := params.Validate()
    if len(validationErrors) > 0 {
        return c.Status(fiber.StatusBadRequest).JSON(validationErrors)
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

func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {
    var (
        // values bson.M
        params types.UpdateUserParams
        userID = c.Params("id")
    )
    oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return ErrInvalidID()
	}
    if err := c.BodyParser(&params); err != nil {
        return ErrBadRequest()
    }
    filter := bson.M{"_id": oid}
    
    if err := h.userStore.UpdateUser(c.Context(), filter, params); err != nil {
        return err
    }
    return c.JSON(map[string]string{"updated": userID})
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
    var (
        id = c.Params("id")
    )
    user, err := h.userStore.GetUserByID(c.Context(), id)
    if err != nil {
        if errors.Is(err, mongo.ErrNoDocuments) {
            return ErrResourceFound("user")
        }
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


func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
    userId := c.Params("id")
    if err := h.userStore.DeleteUser(c.Context(), userId); err != nil {
        return err
    }
    return c.JSON(map[string]string{"deleted": userId})
}