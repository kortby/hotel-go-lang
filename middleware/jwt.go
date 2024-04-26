package middleware

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(c *fiber.Ctx) error {
	fmt.Println("--- JWT auth ----")

	tokens, ok := c.GetReqHeaders()["X-Api-Token"]
	if !ok {
		return fmt.Errorf("unauthorized")
	}

	
	token := tokens[0]
	if err := parseToken(token); err != nil {
		return err
	}

	fmt.Println("token:", token)
	return nil
}

func parseToken(tokenString string) error {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Printf("Unexpected signing method: %v\n", token.Header["alg"])
			return nil, fmt.Errorf("Unauthrized")
		}

		secret := os.Getenv("JWT_SECRET")

		// fmt.Println("never print this but ---:", secret)

		return []byte(secret), nil
	})

	if err != nil {
		fmt.Println("failed to parse JWT claims ", err)
		return fmt.Errorf("Unauthorized")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims)
	} 
	
	return fmt.Errorf("Unauthorized")

}