package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/kortbyhotel/reservation/api"
	"github.com/kortbyhotel/reservation/types"
	"github.com/stretchr/testify/assert"
)

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
  
	app := fiber.New()
	userHandler := api.NewUserHandler(tdb.User)
	app.Post("/", userHandler.HandlePostUser)
  
	params := types.CreateUserParams{
	  Email:     "some@test.com",
	  FirstName: "some",
	  LastName:  "test",
	  Password:  "1234test",
	}
  
	b, err := json.Marshal(params)
	if err != nil {
	  t.Fatal(err)
	}
  
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
  
	resp, err := app.Test(req)
	if err != nil {
	  t.Error(err)
	}

	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)
	if user.FirstName != params.FirstName {
		t.Errorf("expected firstname %s but got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("expected firstname %s but got %s", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		t.Errorf("expected firstname %s but got %s", params.Email, user.Email)
	}
  
	// Assertions
	assert := assert.New(t)
	assert.Equal(http.StatusOK, resp.StatusCode) // Assert successful creation (status code 200)
	// You can add further assertions based on your logic (e.g., response body content)
  }
