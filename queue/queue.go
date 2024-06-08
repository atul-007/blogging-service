package main

import (
	"log"
	"sync"

	"github.com/gofiber/fiber/v2"
)

type Blog struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserID string `json:"user_id"`
}

var (
	queue []Blog
	mu    sync.Mutex
)

func enqueue(c *fiber.Ctx) error {
	var blog Blog
	if err := c.BodyParser(&blog); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	mu.Lock()
	queue = append(queue, blog)
	mu.Unlock()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "Blog enqueued"})
}

func dequeue() *Blog {
	mu.Lock()
	defer mu.Unlock()
	if len(queue) == 0 {
		return nil
	}
	blog := queue[0]
	queue = queue[1:]
	return &blog
}

func main() {
	app := fiber.New()

	app.Post("/enqueue", enqueue)

	app.Get("/dequeue", func(c *fiber.Ctx) error {
		blog := dequeue()
		if blog == nil {
			return c.Status(fiber.StatusNoContent).SendString("Queue is empty")
		}
		return c.JSON(blog)
	})

	log.Fatal(app.Listen(":3001"))
}
