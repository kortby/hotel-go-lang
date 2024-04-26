package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)


func JWTAuthentication(c *fiber.Ctx) error {
    fmt.Println("--- JWT auth ----")

    token := c.Get("X-Api-Token")
    if token == "" {
        fmt.Println("token not present in the header")
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
    }

    claims, err := validateToken(token)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
    }

    expires := int64(claims["exp"].(float64))
    if time.Now().Unix() > expires {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "token expired"})
    }

    fmt.Println("claims:", claims)
    return c.Next()
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
        fmt.Println("Invalid token")
        return nil, fmt.Errorf("unauthorized")
    }

    claims, ok := token.Claims.(*jwt.MapClaims)
    if !ok {
        fmt.Println("Invalid claims structure")
        return nil, fmt.Errorf("unauthorized")
    }

    return *claims, nil
}
