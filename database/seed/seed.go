package seed

import (
	"fmt"
)

func Seed() {
	// AutoMigrate And Seed Tables
	SeedUserRoleTable()
	SeedUserTable()
	fmt.Println("Successfully Seeded Database")
}
