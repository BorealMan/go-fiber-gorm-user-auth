package auth

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"

	"app/config"
)

var JWTSecretKey = []byte(config.JWT_SECRET)

func IssueJWT(userId uint, userRole string) (string, error) {

	claims := jwt.MapClaims{
		"user_id": fmt.Sprintf("%d", userId),
		"role":    userRole,
		"exp":     time.Now().Add(time.Minute * 3600).Unix(), // 1 Day
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS384, claims)

	t, err := token.SignedString(JWTSecretKey)

	return t, err
}

func ValidateJWT(c *fiber.Ctx) error {
	headers := c.GetReqHeaders()
	authHeader := headers["Authorization"]
	authToken := strings.Split(authHeader, " ")

	if len(authToken) != 2 {
		return c.SendStatus(401)
	}

	token, err := jwt.Parse(authToken[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		return JWTSecretKey, nil
	})

	if err != nil {
		return c.SendStatus(401)
	}

	if token.Valid {
		claims, ok := token.Claims.(jwt.MapClaims)
		// fmt.Printf("Claims:\nRole:%s\nUserId:%s\nExp: %f\n", claims["role"], claims["user_id"], claims["exp"])
		if ok {
			// Check Expiration
			err = TokenExpired(claims)
			if err != nil {
				return c.SendStatus(401)
			}
			// Add Values To Locals
			c.Locals("user_id", fmt.Sprintf("%s", claims["user_id"]))
			c.Locals("role", fmt.Sprintf("%s", claims["role"]))
		}
		return c.Next()
	}

	return c.SendStatus(500)
}

// Use Validate JWT First
func ValidateAdmin(c *fiber.Ctx) error {
	if fmt.Sprintf("%s", c.Locals("role")) != "admin" {
		return c.SendStatus(403)
	}
	return c.Next()
}

func TokenExpired(claims jwt.MapClaims) error {
	tokenExp, err := strconv.ParseFloat(fmt.Sprintf("%f", claims["exp"]), 64)
	if err != nil {
		return fmt.Errorf("Invalid Token")
	}
	now := time.Now().Add(time.Minute * 3600).Unix()
	diff := now - int64(tokenExp)
	b := diff > config.JWT_EXPIRES
	// Token Is Expired
	if b {
		return fmt.Errorf("Token Is Expired")
	}
	return nil
}
