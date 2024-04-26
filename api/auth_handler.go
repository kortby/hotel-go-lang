package api

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kortbyhotel/reservation/data"
	"github.com/kortbyhotel/reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	userStore data.UserStore
}

func NewAuthHandler(userStore data.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

type AuthResponse struct {
	User *types.User `json:"user"`
	Token string	`json:"token"`
}
type AuthParams struct {
    Email   string `json:"email"` // Fixed JSON tags
    Password   string `json:"password"`
}

func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
    var params AuthParams
    if err := c.BodyParser(&params); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Error parsing request body"})
    }

    user, err := h.userStore.GetUserByEmail(c.Context(), params.Email)
    if err != nil {
        if errors.Is(err, mongo.ErrNoDocuments) {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
        }
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "server error"})
    }

    if !types.IsValidPassword(user.EncryptedPassword, params.Password) {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
    }
    
    token, err := createTokenFromUser(user)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create token"})
    }
    
    resp := AuthResponse{
        User: user,
        Token: token,
    }
    return c.JSON(resp)
}

func createTokenFromUser(user *types.User) (string, error) {
    claims := jwt.MapClaims{
        "userID": user.ID.Hex(), // Assuming user.ID is an ObjectID
        "email": user.Email,
        "exp": time.Now().Add(time.Hour * 4).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    secret := os.Getenv("JWT_SECRET")
    tokenStr, err := token.SignedString([]byte(secret))
    if err != nil {
        fmt.Println("Failed to sign token with secret:", err)
        return "", err
    }
    return tokenStr, nil
}
