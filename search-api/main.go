package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type SearchResponse struct {
	Hits struct {
		Hits []struct {
			Source struct {
				Title string `json:"title"`
				Text  string `json:"text"`
			} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

func main() {
	app := fiber.New()

	app.Get("/search", func(c *fiber.Ctx) error {
		query := c.Query("q")
		if query == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Query parameter is required"})
		}

		esQuery := map[string]interface{}{
			"query": map[string]interface{}{
				"multi_match": map[string]interface{}{
					"query":  query,
					"fields": []string{"title", "text"},
				},
			},
		}

		var buf strings.Builder
		if err := json.NewEncoder(&buf).Encode(esQuery); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		resp, err := http.Post("http://elasticsearch:9200/blogs/_search", "application/json", strings.NewReader(buf.String()))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error querying Elasticsearch"})
		}

		var searchResponse SearchResponse
		if err := json.NewDecoder(resp.Body).Decode(&searchResponse); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		results := []fiber.Map{}
		for _, hit := range searchResponse.Hits.Hits {
			results = append(results, fiber.Map{
				"title": hit.Source.Title,
				"text":  hit.Source.Text,
			})
		}

		return c.JSON(results)
	})

	log.Fatal(app.Listen(":3002"))
}
