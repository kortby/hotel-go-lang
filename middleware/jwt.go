package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kortbyhotel/reservation/api"
	"github.com/kortbyhotel/reservation/data"
)


func JWTAuthentication(userStore data.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {

		token := c.Get("X-Api-Token")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}
	
		claims, err := validateToken(token)
		if err != nil {
			return api.ErrUnauthorized()
		}
	
		expires := int64(claims["exp"].(float64))
		if time.Now().Unix() > expires {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "token expired"})
		}
	
		userID, ok := claims["userID"].(string)
		if !ok || userID == "" {
			fmt.Println("userID claim is missing or not a string")
			return api.ErrUnauthorized()
		}

		user, err := userStore.GetUserByID(c.Context(), userID)
		if err != nil {
			fmt.Println("error retrieving user:", err)
			return api.ErrUnauthorized()
		}
		c.Context().SetUserValue("user", user)
		return c.Next()
	}

}

func validateToken(tokenString string) (jwt.MapClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            fmt.Printf("Unexpected signing method: %v\n", token.Header["alg"])
            return nil, fmt.Errorf("unauthorized - wrong signing method")
        }
        secret := os.Getenv("JWT_SECRET")
        return []byte(secret), nil
    })
    if err != nil {
        fmt.Println("failed to parse JWT token: ", err)
        return nil, fmt.Errorf("unauthorized - token parse error")
    }
    if !token.Valid {
        return nil, api.ErrUnauthorized()
    }

    claims, ok := token.Claims.(*jwt.MapClaims)
    if !ok {
        return nil, api.ErrUnauthorized()
    }

    return *claims, nil
}
