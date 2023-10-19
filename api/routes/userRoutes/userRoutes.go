package userRoutes

import (
	"github.com/gofiber/fiber/v2"

	"app/api/auth"
	"app/models/user"
)

func SetUserRoutes(api fiber.Router) {
	userGroup := api.Group("/user")
	userGroup.Post("/login", user.Login)
	userGroup.Post("/create", user.CreateUser)
	userGroup.Get("/", auth.ValidateJWT, user.VerifyAccountEnabled, user.GetUser)
	userGroup.Put("/update-user", auth.ValidateJWT, user.VerifyAccountEnabled, user.UpdateUser)
	userGroup.Put("/update-password", auth.ValidateJWT, user.VerifyAccountEnabled, user.UpdatePassword)
	userGroup.Delete("/", auth.ValidateJWT, user.VerifyAccountEnabled, user.DeleteUser)

	// Admin Functions
	userGroup.Put("/admin-user-update", auth.ValidateJWT, auth.ValidateAdmin, user.VerifyAccountEnabled, user.AdminUpdateUser)
	userGroup.Get("/getall", auth.ValidateJWT, auth.ValidateAdmin, user.VerifyAccountEnabled, user.GetAll)
	userGroup.Get("/get-user-roles", auth.ValidateJWT, auth.ValidateAdmin, user.VerifyAccountEnabled, user.GetUserRoles)
}
