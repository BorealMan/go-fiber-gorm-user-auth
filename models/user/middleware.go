package user

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"app/database"
)

/*
Verifies That Account Is Enabled - JWT Tokens Are Valid Until Expiration
This Additional Middleware Will Add A Way To Disable An Account Immediately
*/
func VerifyAccountEnabled(c *fiber.Ctx) error {
	// Get UserID From Locals
	user_id, err := strconv.ParseUint(fmt.Sprintf("%s", c.Locals("user_id")), 10, 32)

	if err != nil {
		if DEBUG {
			log.Printf("Verify Account Enabled Error: %s\n", err.Error())
		}
		return c.SendStatus(500)
	}

	user := new(User)
	user.ID = uint(user_id)

	err = database.DB.First(&user).Error

	if err != nil {
		if DEBUG {
			log.Printf("Invalid JWT Token For User: %d\n", user_id)
		}
		return c.SendStatus(401)
	}

	// Check If Account Is Enabled
	if !*user.AccountEnabled {
		if DEBUG {
			log.Printf("Blocked Disabled Account: %s\n", user.Username)
		}
		return c.Status(403).JSON(fiber.Map{"error": "User Account Disabled"})
	}
	return c.Next()
}
