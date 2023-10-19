package user

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"app/api/auth"
	"app/config"
	"app/database"

	"app/util"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Override By Reassignment
var DEBUG = config.DEBUG

// var DEBUG = false

// Every User Must Be Assigned A Role
type UserRole struct {
	gorm.Model
	Role        string `json:"role" gorm:"type:VARCHAR(32);unique;not null" validate:"omitempty"`
	Description string `json:"description" gorm:"type:VARCHAR(100);" validate:"omitempty"`
}

type User struct {
	gorm.Model
	Username       string   `json:"username" gorm:"type:VARCHAR(16);not null;uniqueIndex:idx_username;" validate:"required,min=1,max=16"`
	Password       string   `json:"password" gorm:"type:VARCHAR(64);not null" validate:"omitempty,min=1,max=32"`
	Email          string   `json:"email" gorm:"type:VARCHAR(48);not null;uniqueIndex:idx_email" validate:"required,email"`
	Phone          string   `json:"phone" gorm:"type:VARCHAR(13)" validate:"omitempty,e164"`
	AccountEnabled *bool    `json:"account_enabled" gorm:"default:1; not null" validate:"omitempty"`
	RoleID         uint     `json:"role_id" validate:"omitempty,number"`
	Role           UserRole `json:"user_role" validate:"omitempty"`
}

/*
	Functionalities
*/

type LoginRequest struct {
	Username string `json:"username" form:"username" validate:"required,min=1,max=16"`
	Password string `json:"password" form:"password" validate:"required,min=1,max=32"`
}

func Login(c *fiber.Ctx) error {
	r := new(LoginRequest)
	err := c.BodyParser(r)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Input Fields"})
	}

	// Process Input
	r.Username = strings.ToLower(strings.TrimSpace(r.Username))
	r.Password = strings.TrimSpace(r.Password)

	if DEBUG {
		log.Printf("Login -> Username: %s\n", r.Username)
	}

	err = util.Validate(r)

	// Validation
	if err != nil {
		if DEBUG {
			log.Printf("Validation Error: %s\n", err.Error())
		}
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Fields"})
	}

	// Create User
	var user User
	user.Username = r.Username

	// Lookup User
	err = database.DB.Preload("Role").Where("username = ?", r.Username).First(&user).Error
	if err != nil {
		if DEBUG {
			log.Printf("Login Lookup Error: %s", err.Error())
		}
	}
	// Check if User Exists
	if user.ID == 0 {
		if DEBUG {
			log.Printf("Invalid Username: %s", user.Username)
		}
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Username or Password"})
	}
	// Check If Account Is Enabled
	if !*user.AccountEnabled {
		if DEBUG {
			log.Printf("Blocked Disabled Account: %s", user.Username)
		}
		return c.Status(403).JSON(fiber.Map{"error": "User Account Disabled"})
	}
	// Check Users Password
	pld := user.Username + r.Password + config.SALT
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pld)) != nil {
		if DEBUG {
			log.Printf("Failed Login Attempt: %s\n", user.Username)
		}
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Username or Password"})
	}

	// Create JWT Token For User
	token, err := auth.IssueJWT(user.ID, user.Role.Role)

	if err != nil {
		if DEBUG {
			log.Printf("Login JTW Error: %s\n", err.Error())
		}
		return c.SendStatus(500)
	}
	user.Password = ""
	return c.Status(200).JSON(fiber.Map{"token": token, "user": user})
}

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=1,max=16"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"omitempty,min=1,max=32"`
}

func CreateUser(c *fiber.Ctx) error {
	r := new(CreateUserRequest)
	err := c.BodyParser(r)

	if err != nil {
		if DEBUG {
			log.Printf("Create User Error: %s\n", err.Error())
		}
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Input Fields"})
	}

	// Validate Input
	err = util.Validate(r)
	// Check Input
	if err != nil {
		if DEBUG {
			log.Printf("Create User Error: %s\n", err.Error())
		}
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Input"})
	}

	var user User
	user.Username = strings.TrimSpace(strings.ToLower(r.Username))
	user.Password = strings.TrimSpace(r.Password)
	user.RoleID = 1 // Assign Default Role
	user.Email = strings.TrimSpace(r.Email)

	// Hash Password
	user.Password, err = HashPassword(user.Username, user.Password)
	if err != nil {
		if DEBUG {
			log.Printf("Create User Error: Failed To Hash Password: %s\n", err.Error())
		}
		return c.SendStatus(500)
	}
	// Saving User To DB
	err = database.DB.Create(&user).Error
	if err != nil || user.ID == 0 {
		if DEBUG {
			log.Printf("Create User Error: %s\n", err.Error())
		}
		return c.Status(409).JSON(fiber.Map{"error": "Username or Email Already Exists"})
	}
	// Pulling Out Data
	err = database.DB.Preload("Role").Omit("Password").First(&user).Error
	if err != nil {
		if DEBUG {
			log.Printf("Create User Failed To Retrieve Data: %s", err.Error())
		}
		return c.SendStatus(500)
	}
	// Create JWT Token For User
	token, err := auth.IssueJWT(user.ID, user.Role.Role)
	// Check The Token Didn't Explode
	if err != nil {
		if DEBUG {
			log.Printf("Create User: Failed To Generate Token: %s", err.Error())
		}
		return c.SendStatus(500)
	}
	user.Password = ""
	return c.Status(201).JSON(fiber.Map{"token": token, "user": user})
}

func GetUser(c *fiber.Ctx) error {
	// Get UserID From Locals
	user_id, err := strconv.ParseUint(fmt.Sprintf("%s", c.Locals("user_id")), 10, 32)
	if err != nil {
		if DEBUG {
			log.Printf("Get User Error: Failed Parsing Uint: %s\n", err.Error())
		}
		return c.SendStatus(500)
	}
	var user User
	user.ID = uint(user_id)
	err = database.DB.Preload("Role").First(&user).Error
	if err != nil {
		if DEBUG {
			log.Printf("Get User Error: User Doesn't Exist: %s\n", err.Error())
		}
		return c.Status(404).JSON(fiber.Map{"error": "User Doesn't Exist"})
	}
	user.Password = ""
	return c.Status(200).JSON(fiber.Map{"user": user})
}

type UserUpdateRequest struct {
	Email string `json:"email" validate:"omitempty,email"`
	Phone string `json:"phone" validate:"omitempty,e164"`
}

func UpdateUser(c *fiber.Ctx) error {
	// Get UserID From Locals
	user_id, err := strconv.ParseUint(fmt.Sprintf("%s", c.Locals("user_id")), 10, 32)
	if err != nil {
		if DEBUG {
			log.Printf(`Update User Error Parsing Uint: %s\n`, err.Error())
		}
		return c.SendStatus(500)
	}

	r := new(UserUpdateRequest)
	err = c.BodyParser(r)

	if err != nil {
		if DEBUG {
			log.Printf("Update User Error: %s", err.Error())
		}
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Input Fields"})
	}

	// Lookup Record
	var user User
	user.ID = uint(user_id)

	err = database.DB.First(&user).Error
	if user.ID == 0 || err != nil {
		if DEBUG {
			log.Printf("Update User: User Doesn't Exist: %s\n", err.Error())
		}
		return c.Status(404).JSON(fiber.Map{"error": "User Doesn't Exist"})
	}

	user.Email = r.Email
	user.Phone = r.Phone

	// Try to save the new fields
	err = database.DB.Save(&user).Error
	if err != nil {
		if DEBUG {
			log.Printf("Failed To Update User: %s\n", err.Error())
		}
		return c.SendStatus(500)
	}
	// Pull Out Updated User
	err = database.DB.Preload("Role").Omit("Password").First(&user).Error
	if err != nil {
		if DEBUG {
			log.Printf("Update User: Failed To Retrieve Data: %s\n", err.Error())
		}
		return c.SendStatus(500)
	}
	user.Password = ""
	return c.Status(200).JSON(fiber.Map{"user": user})
}

type UpdateUserPasswordRequest struct {
	Password string `json:"password" validate:"omitempty,min=1,max=32"`
}

// Must Rehash Password
func UpdatePassword(c *fiber.Ctx) error {
	// Get UserID From Locals
	user_id, err := strconv.ParseUint(fmt.Sprintf("%s", c.Locals("user_id")), 10, 32)
	if err != nil {
		if DEBUG {
			log.Printf("Get User Error: Failed Parsing Uint: %s\n", err.Error())
		}
		return c.SendStatus(500)
	}
	// Parse Request
	r := new(UpdateUserPasswordRequest)
	err = c.BodyParser(r)
	if err != nil {
		if DEBUG {
			log.Printf("Update Password Error: %s\n", err.Error())
		}
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Input Fields"})
	}

	err = util.Validate(r)

	if err != nil {
		if DEBUG {
			log.Printf("Update Password Error: %s\n", err.Error())
		}
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Input"})
	}

	// Lookup User
	var user User
	err = database.DB.Where("id = ?", user_id).First(&user).Error
	if err != nil {
		if DEBUG {
			log.Printf("Update Password Error: User Doesn't Exist: %s\n", err.Error())
		}
		return c.Status(404).JSON(fiber.Map{"error": "User Doesn't Exist"})
	}
	// Hash Password
	user.Password, err = HashPassword(user.Username, r.Password)
	if err != nil {
		if DEBUG {
			log.Printf("Update Password Error: %s\n", err.Error())
		}
		return c.SendStatus(500)
	}
	// Saving User To DB
	err = database.DB.Save(&user).Error
	if err != nil {
		if DEBUG {
			log.Printf("Update Password Error: Failed To Save: %s\n", err.Error())
		}
		return c.SendStatus(500)
	}
	return c.SendStatus(200)
}

// Soft Deletes A User
func DeleteUser(c *fiber.Ctx) error {
	// Get UserID From Locals
	user_id, err := strconv.ParseUint(fmt.Sprintf("%s", c.Locals("user_id")), 10, 32)
	if err != nil {
		if DEBUG {
			log.Printf("Get User Error: Failed Parsing Uint: %s\n", err.Error())
		}
		return c.SendStatus(500)
	}

	var user User
	user.ID = uint(user_id)

	err = database.DB.Delete(&user).Error
	if err != nil {
		if DEBUG {
			log.Printf("Delete User Error: User Doesn't Exist: %s\n", err.Error())
		}
		return c.Status(404).JSON(fiber.Map{"error": "User Doesn't Exist"})
	}

	return c.SendStatus(200)
}

// Not Implemented - May Limit To Administrator Only
func PermanentlyDeleteUser(c *fiber.Ctx) error {
	return c.SendStatus(501)
}

/*
	Admin Functions
*/

type AdminUserUpdateRequest struct {
	UserID          uint   `json:"user_id" validate:"required,number"`
	RoleID          uint   `json:"role_id" validate:"required,number"`
	Username        string `json:"username" validate:"omitempty,min=1,max=16"`
	Email           string `json:"email" validate:"omitempty,email"`
	Phone           string `json:"phone" validate:"omitempty,e164"`
	Account_enabled *bool  `json:"account_enabled" validate:"omitempty"`
}

// Allows a User To Update Their Settings
func AdminUpdateUser(c *fiber.Ctx) error {
	// Parse Data - Automatically Parses Sub Structures
	r := new(AdminUserUpdateRequest)
	err := c.BodyParser(r)

	if err != nil {
		if DEBUG {
			log.Println("Admin Update: ", err)
		}
		return c.Status(500).JSON(fiber.Map{"error": "Error Parsing Input"})
	}

	// Validate Input
	err = util.Validate(r)
	// Check Input
	if err != nil {
		log.Println(err)
		return c.Status(400).JSON(fiber.Map{"error": "Failed Validation"})
	}

	// Lookup Record
	var user User
	database.DB.First(&user, r.UserID)
	if user.ID == 0 {
		if DEBUG {
			log.Printf("Admin Update Error: User Doesn't Exist: %s", err.Error())
		}
		return c.Status(404).JSON(fiber.Map{"error": "User Doesn't Exist"})
	}

	user.ID = r.UserID
	user.AccountEnabled = r.Account_enabled
	user.Email = r.Email
	user.Phone = r.Phone
	user.RoleID = r.RoleID
	user.UpdatedAt = time.Now().Local()

	// Try to save the new fields
	err = database.DB.Save(&user).Error
	if err != nil {
		c.Status(500).JSON(fiber.Map{"error": "Failed To Update"})
	}
	// Pull Out Updated User
	database.DB.Preload("Role").First(&user)
	user.Password = ""
	return c.Status(200).JSON(fiber.Map{"user": user})
}

func GetAll(c *fiber.Ctx) error {
	var users []User
	err := database.DB.Preload("Role").Omit("Password").Find(&users).Error
	if err != nil {
		if DEBUG {
			log.Printf("Get All Users: %s", err.Error())
		}
		return c.SendStatus(500)
	}
	return c.Status(200).JSON(fiber.Map{"users": users})
}

func GetUserRoles(c *fiber.Ctx) error {
	var userRoles []UserRole
	err := database.DB.Find(&userRoles).Error
	if err != nil {
		if DEBUG {
			log.Printf("Get User Roles Error: %s\n", err.Error())
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed To Retrieve User Roles"})
	}
	return c.Status(200).JSON(fiber.Map{"user_roles": userRoles})
}
