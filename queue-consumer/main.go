package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/segmentio/kafka-go"
)

type Blog struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserID string `json:"user_id"`
}

func main() {

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"kafka:9092"},
		Topic:    "blog-submissions",
		GroupID:  "blog-consumers",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("Error reading message:", err)
			time.Sleep(1 * time.Second)
			continue
		}

		var blog Blog
		if err := json.Unmarshal(msg.Value, &blog); err != nil {
			log.Println("Error unmarshaling message:", err)
			continue
		}

		blogJSON, err := json.Marshal(blog)
		if err != nil {
			log.Println("Error marshaling blog:", err)
			continue
		}

		req, err := http.NewRequest("POST", "http://elasticsearch:9200/blogs/_doc", bytes.NewBuffer(blogJSON))
		if err != nil {
			log.Println("Error creating request:", err)
			continue
		}
		req.Header.Set("Content-Type", "application/json")

		_, err = http.DefaultClient.Do(req)
		if err != nil {
			log.Println("Error sending request:", err)
		}
	}
}
