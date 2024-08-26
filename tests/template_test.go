package template_test

import (
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"io"
	"net/http/httptest"
	"testing"
)

func TestGet(t *testing.T) {
	app := fiber.New()

	// Create a new endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!") // send text
	})

	req := httptest.NewRequest("GET", "/", nil)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Expected status code 200, but got: %d", resp.StatusCode)
	}

	bodyString, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if string(bodyString) != "Hello, World!" {
		t.Fatalf("Expected body to be %s, but got: %s", "Hello, World!", resp.Body)
	}
}

func TestPost(t *testing.T) {
	app := fiber.New()

	type payload struct {
		Name string `json:"name"`
	}

	// Create a new endpoint
	app.Post("/", func(c *fiber.Ctx) error {
		p := new(payload)
		if err := c.BodyParser(p); nil != err {
			return err
		}
		p.Name = p.Name + "!"
		return c.JSON(p)
	})

	reqBody, _ := json.Marshal(&payload{Name: "John"})

	req := httptest.NewRequest("POST", "/", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Expected status code 200, but got: %d", resp.StatusCode)
	}

	var userResp payload
	if err = json.NewDecoder(resp.Body).Decode(&userResp); nil != err {
		t.Fatalf("Error in parse response: %v", err)
	}

	if userResp.Name != "John!" {
		t.Fatalf("Expected body to be %s, but got: %s", "John!", resp.Body)
	}
}
