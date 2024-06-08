package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/segmentio/kafka-go"
)

type Blog struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserID string `json:"user_id"`
}

func main() {
	app := fiber.New()

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"kafka:9092"},
		Topic:   "blog-submissions",
	})

	app.Post("/submit", func(c *fiber.Ctx) error {
		var blog Blog
		if err := c.BodyParser(&blog); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		blogJSON, err := json.Marshal(blog)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		err = writer.WriteMessages(context.Background(), kafka.Message{
			Value: blogJSON,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"status": "Blog submission queued"})
	})

	log.Fatal(app.Listen(":3000"))
}
