package user

import (
	"log"
	"testing"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

const URL = "http://localhost:5000/user"

var TOKEN = "Bearer: "

// Create User And Login Response
type CreateUserResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type UserResponse struct {
	User User `json:"user"`
}

func TestCreateUser(t *testing.T) {
	agent := fiber.Post(URL + "/create")
	if DEBUG {
		log.Println("Running Create User Test")
		agent.Debug()
	}

	// Create New Request
	req := new(CreateUserRequest)
	req.Username = "tester"
	req.Password = "123"
	req.Email = "test@tester.com"

	// Send JSON Request
	agent.JSON(req)

	// Destructure Response
	statusCode, body, errs := agent.Bytes()

	expectedStatus := 201
	if statusCode != expectedStatus {
		t.Fatalf("\nInvalid Status Code: %d Expected: %d\n", statusCode, expectedStatus)
	}

	// Parse JSON
	res := new(CreateUserResponse)
	err := json.Unmarshal(body, res)

	if err != nil {
		t.Fatalf("\nFailed To Unmarshal JSON\n")
	}

	if len(errs) > 0 {
		t.Fatalf("\nFailed: %v+\n", errs)
	}

	// Expecting Username = tester
	if req.Username != res.User.Username {
		t.Fatalf("\nFailed To Create User: %s\n", req.Username)
	}

}

func TestLogin(t *testing.T) {
	agent := fiber.Post(URL + "/login")
	if DEBUG {
		log.Println("Running Login Test")
		agent.Debug()
	}

	req := new(LoginRequest)
	req.Username = "tester"
	req.Password = "123"

	agent.JSON(req)

	// Destructure Response
	statusCode, body, errs := agent.Bytes()

	expectedStatus := 200
	if statusCode != expectedStatus {
		t.Fatalf("\nInvalid Status Code: %d Expected: %d\n", statusCode, expectedStatus)
	}

	// Parse JSON
	res := new(CreateUserResponse)
	err := json.Unmarshal(body, res)

	// Set Token For Other Tests
	TOKEN += res.Token

	if err != nil {
		t.Fatalf("\nFailed To Unmarshal JSON\n")
	}

	if len(errs) > 0 {
		t.Fatalf("\nFailed: %v+\n", errs)
	}

	// Expecting Username = tester
	if req.Username != res.User.Username {
		t.Fatalf("\nFailed To Login: %s\n", req.Username)
	}
}

func TestGet(t *testing.T) {
	agent := fiber.Get(URL)
	if DEBUG {
		log.Println("Running Get Test")
		// log.Printf("TOKEN = %s\n", TOKEN)
		agent.Debug()
	}

	agent.Set("Authorization", TOKEN)

	// Destructure Response
	statusCode, body, errs := agent.Bytes()

	expectedStatus := 200
	if statusCode != expectedStatus {
		t.Fatalf("\nInvalid Status Code: %d Expected: %d\n", statusCode, expectedStatus)
	}

	res := new(UserResponse)
	err := json.Unmarshal(body, res)

	if err != nil {
		t.Fatalf("\nFailed To Unmarshal JSON\n")
	}

	if len(errs) > 0 {
		t.Fatalf("\nFailed: %v+\n", errs)
	}

	// Expecting Username = tester
	expectedUsername := "tester"
	if expectedUsername != res.User.Username {
		t.Fatalf("\nFailed To Get User: %s\n", expectedUsername)
	}
}

func TestUpdateUser(t *testing.T) {
	agent := fiber.Put(URL + "/update-user")
	if DEBUG {
		log.Println("Running Update User Test")
		agent.Debug()
	}

	agent.Request().Header.Add("Authorization", TOKEN)

	req := new(User)
	req.Email = "testerupdated@tester.com"
	req.Phone = "+19998675309"

	agent.JSON(req)

	// Destructure Response
	statusCode, body, errs := agent.Bytes()

	expectedStatus := 200
	if statusCode != expectedStatus {
		t.Fatalf("\nInvalid Status Code: %d Expected: %d\n", statusCode, expectedStatus)
	}

	// Parse JSON
	res := new(UserResponse)
	err := json.Unmarshal(body, res)

	if err != nil {
		t.Fatalf("\nFailed To Unmarshal JSON\n")
	}

	if len(errs) > 0 {
		t.Fatalf("\nFailed: %v+\n", errs)
	}

	// Expecting Username = tester
	if req.Email != res.User.Email {
		t.Fatalf("\nFailed To Update: %s\n", req.Username)
	}
}

func TestUpdatePassword(t *testing.T) {
	agent := fiber.Put(URL + "/update-password")
	if DEBUG {
		log.Println("Running Update Password Test")
		agent.Debug()
	}

	agent.Request().Header.Add("Authorization", TOKEN)

	req := new(UpdateUserPasswordRequest)
	req.Password = "1234"

	agent.JSON(req)

	statusCode, _, errs := agent.Bytes()

	expectedStatus := 200
	if statusCode != expectedStatus {
		t.Fatalf("\nInvalid Status Code: %d Expected: %d\n", statusCode, expectedStatus)
	}

	if len(errs) > 0 {
		t.Fatalf("\nFailed: %v+\n", errs)
	}
}

// Expected To Fail
func TestUpdatedPassword_LoginFail(t *testing.T) {
	agent := fiber.Post(URL + "/login")
	if DEBUG {
		log.Println("Running Updated Password Login Fail Test")
		agent.Debug()
	}

	req := new(LoginRequest)
	req.Username = "tester"
	req.Password = "123"

	agent.JSON(req)

	// Destructure Response
	statusCode, _, errs := agent.Bytes()

	expectedStatus := 400
	if statusCode != expectedStatus {
		t.Fatalf("\nInvalid Status Code: %d Expected: %d\n", statusCode, expectedStatus)
	}

	if len(errs) > 0 {
		t.Fatalf("\nFailed: %v+\n", errs)
	}
}

// Expected To Pass
func TestUpdatedPassword_LoginPass(t *testing.T) {
	agent := fiber.Post(URL + "/login")
	if DEBUG {
		log.Println("Running Updated Password Login Pass Test")
		agent.Debug()
	}

	req := new(LoginRequest)
	req.Username = "tester"
	req.Password = "1234"

	agent.JSON(req)

	// Destructure Response
	statusCode, body, errs := agent.Bytes()

	expectedStatus := 200
	if statusCode != expectedStatus {
		t.Fatalf("\nInvalid Status Code: %d Expected: %d\n", statusCode, expectedStatus)
	}

	// Parse JSON
	res := new(CreateUserResponse)
	err := json.Unmarshal(body, res)

	if err != nil {
		t.Fatalf("\nFailed To Unmarshal JSON\n")
	}

	if len(errs) > 0 {
		t.Fatalf("\nFailed: %v+\n", errs)
	}

	// Expecting Username = tester
	if req.Username != res.User.Username {
		t.Fatalf("\nFailed To Login: %s\n", req.Username)
	}
}

// Tests Soft Delete Functionality
func TestDeleteUser(t *testing.T) {
	agent := fiber.Delete(URL + "/")
	if DEBUG {
		log.Println("Running Delete User Test")
		agent.Debug()
	}

	agent.Request().Header.Add("Authorization", TOKEN)

	statusCode, _, errs := agent.Bytes()

	expectedStatus := 200
	if statusCode != expectedStatus {
		t.Fatalf("\nInvalid Status Code: %d Expected: %d\n", statusCode, expectedStatus)
	}

	if len(errs) > 0 {
		t.Fatalf("\nFailed: %v+\n", errs)
	}
}
