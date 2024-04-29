package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/kortbyhotel/reservation/api"
	"github.com/kortbyhotel/reservation/data"
	"github.com/kortbyhotel/reservation/data/fixtures"
	"github.com/kortbyhotel/reservation/types"
)

func insertTestUser(t *testing.T, userStore data.UserStore) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email: "test@test.com",
		FirstName: "john",
		LastName: "Doe",
		Password: "test1234",
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = userStore.CreateUser(context.TODO(), user)
	if err != nil {
		t.Fatal(err)
	}
	return user
}

func TestAuthenticateSuccess(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	insertedUser, _ := fixtures.AddUser(tdb.Store, "test", "test", false)

  
	app := fiber.New()
	authHandler := api.NewAuthHandler(tdb.User)
	app.Post("/auth", authHandler.HandleAuthenticate)

	params := api.AuthParams{
		Email: "test@test.com",
		Password: "test1234",
	}
	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected http status of 200, but got %d", resp.StatusCode)
	}

	var authResp api.AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Error(err)
	}

	if authResp.Token == ""{
		t.Fatalf("Expected the JWT to be present in the auth response")
	}

	insertedUser.EncryptedPassword = ""
	// set the encrypted password to and empty string 
	// we do not return the password
	if reflect.DeepEqual(insertedUser, authResp.User) {
		t.Fatalf("expected the user to be the inserted user")
	}
	// b, _ := io.ReadAll()
}

func TestAuthenticateNotSuccess(t *testing.T) {
    tdb := setup(t)
    defer tdb.teardown(t)

    // Insert a test user
    // insertTestUser(t, tdb.User)
	fixtures.AddUser(tdb.Store, "test", "test", false)

    app := fiber.New()
    authHandler := api.NewAuthHandler(tdb.User)
    app.Post("/auth", authHandler.HandleAuthenticate)

    // Provide incorrect password parameters
    params := api.AuthParams{
        Email: "test@test.com",
        Password: "wrongpassword", // Intentionally incorrect
    }
    b, _ := json.Marshal(params)

    req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
    req.Header.Add("Content-Type", "application/json")
    resp, err := app.Test(req)
    if err != nil {
        t.Fatal(err)
    }

    // Check for unauthorized status code
    if resp.StatusCode != fiber.StatusUnauthorized {
        t.Fatalf("expected http status of 401 Unauthorized, but got %d", resp.StatusCode)
    }

    var respBody map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
        t.Error(err)
    }

    // Verify the error message is about invalid credentials
    if respBody["error"] != "invalid credentials" {
        t.Errorf("expected error message 'invalid credentials', got '%v'", respBody["error"])
    }
}

