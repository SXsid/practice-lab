package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type Node struct {
	nodeId  string
	nodeIds []string
	nextMsg int
	mu      sync.Mutex
}

type Message struct {
	Src  string         `json:"src"`
	Dest string         `json:"dest"`
	Body map[string]any `json:"body"`
}

func (n *Node) send(dest string, body map[string]any) {
	n.mu.Lock()
	defer n.mu.Unlock()
	body["msg_id"] = n.nextMsg
	n.nextMsg++

	b, _ := json.Marshal(&Message{
		Dest: dest,
		Src:  n.nodeId,
		Body: body,
	})
	fmt.Println(string(b))
}

func Validate(msg Message) {
	if msg.Src == "" || msg.Dest == "" || msg.Body == nil {
		fmt.Fprint(os.Stderr, "request is invalid")
		return
	}
	if v, ok := msg.Body["type"].(string); !ok || v == "" {
		fmt.Fprint(os.Stderr, "invalid body ")
		return
	}
}

func (n *Node) reply(msg Message, body map[string]any) {
	msg_id, _ := msg.Body["msg_id"].(string)
	body["in_reply_to"] = msg_id
	n.send(msg.Src, body)
}

func handLeMessage(msg Message) {
	Validate(msg)

	msgType, _ := msg.Body["type"].(string)
	switch msgType {
	case "init":
		node := &Node{}

		node.nodeId, _ = msg.Body["node_id"].(string)
		node.nodeIds, _ = msg.Body["node_ids"]
		body := map[string]any{
			"type": "init_ok",
		}
		node.reply(msg, body)
	case "echo":
	case "proxy":
	default:
		fmt.Fprint(os.Stderr, "not a valid message type")
	}
}

func main() {
	// Your implementation here
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var msg Message
		if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
			fmt.Fprint(os.Stderr, err.Error())
		}
		handLeMessage(msg)
	}
}
