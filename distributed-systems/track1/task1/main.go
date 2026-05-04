package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

// Message represents a Maelstrom message
type Message struct {
	Src  string                 `json:"src"`
	Dest string                 `json:"dest"`
	Body map[string]interface{} `json:"body"`
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		var msg Message
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			fmt.Fprintln(os.Stderr, "Error parsing JSON:", err)
			continue
		}

		value, ok := msg.Body["type"]
		if !ok {
			fmt.Fprintln(os.Stderr, "No type field in body")
			continue
		}
		fmt.Printf("PARSED: %s|%s|%s", msg.Src, msg.Dest, value)
	}
}
