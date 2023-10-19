package seed

import (
	"app/database"
	"app/models/user"
	"log"
)

func SeedUserRoleTable() {
	err := database.DB.AutoMigrate(&user.UserRole{})
	if err != nil {
		log.Fatalf(`Unable To Migrate UserRole: %v`, err.Error())
	}
	// Check To See If Already Seeded
	var userRole user.UserRole
	err = database.DB.Take(&userRole).Error

	// Table Is Already Seeded
	if userRole.ID != 0 {
		return
	}

	userRoles := []user.UserRole{
		{Role: "default", Description: "The default role for all new accounts"},
		{Role: "admin", Description: "Adminstrator"},
	}

	err = database.DB.Create(userRoles).Error

	if err != nil {
		log.Fatalf(`Error Seeding UserRole: %v`, err.Error())
	}

}

func SeedUserTable() {
	err := database.DB.AutoMigrate(&user.User{})

	if err != nil {
		log.Fatalf(`Error Migrating User: %v`, err.Error())
	}
}
