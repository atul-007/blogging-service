package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Blog struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserID string `json:"user_id"`
}

func main() {
	for {
		resp, err := http.Get("http://queue:3001/dequeue")
		if err != nil {
			log.Println("Error dequeuing:", err)
			time.Sleep(1 * time.Second)
			continue
		}

		if resp.StatusCode == http.StatusNoContent {
			time.Sleep(1 * time.Second)
			continue
		}

		var blog Blog
		if err := json.NewDecoder(resp.Body).Decode(&blog); err != nil {
			log.Println("Error decoding response:", err)
			time.Sleep(1 * time.Second)
			continue
		}

		blogJSON, err := json.Marshal(blog)
		if err != nil {
			log.Println("Error marshaling blog:", err)
			time.Sleep(1 * time.Second)
			continue
		}

		req, err := http.NewRequest("POST", "http://elasticsearch:9200/blogs/_doc", bytes.NewBuffer(blogJSON))
		if err != nil {
			log.Println("Error creating request:", err)
			time.Sleep(1 * time.Second)
			continue
		}
		req.Header.Set("Content-Type", "application/json")

		_, err = http.DefaultClient.Do(req)
		if err != nil {
			log.Println("Error sending request:", err)
		}

		time.Sleep(1 * time.Second)
	}
}
