package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Blog struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserID string `json:"user_id"`
}

func main() {
	app := fiber.New()

	app.Post("/submit", func(c *fiber.Ctx) error {
		var blog Blog
		if err := c.BodyParser(&blog); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		blogJSON, err := json.Marshal(blog)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		_, err = http.Post("http://queue:3001/enqueue", "application/json", bytes.NewBuffer(blogJSON))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"status": "Blog submission queued"})
	})

	log.Fatal(app.Listen(":3000"))
}
