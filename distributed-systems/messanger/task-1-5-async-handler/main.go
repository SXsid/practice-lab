package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type Node struct {
	NodeID    string
	NodeIDs   []string
	NextMsgID int
	mu        sync.Mutex
	outMu     sync.Mutex
}

type Message struct {
	Src  string                 `json:"src"`
	Dest string                 `json:"dest"`
	Body map[string]interface{} `json:"body"`
}

func (n *Node) Send(dest string, body map[string]interface{}) {
	n.mu.Lock()
	body["msg_id"] = n.NextMsgID
	n.NextMsgID++
	n.mu.Unlock()

	msg := Message{Src: n.NodeID, Dest: dest, Body: body}
	output, _ := json.Marshal(msg)

	n.outMu.Lock()
	fmt.Println(string(output))
	n.outMu.Unlock()
}

func (n *Node) Reply(request Message, body map[string]interface{}) {
	if msgID, ok := request.Body["msg_id"].(float64); ok {
		body["in_reply_to"] = int(msgID)
	}
	n.Send(request.Src, body)
}

func (node *Node) HandleMessage(msg Message) {
	// TODO: Handle message
	// This function will be called from a goroutine
	msgType, _ := msg.Body["type"].(string)
	switch msgType {
	case "init":
		node.NodeID, _ = msg.Body["node_id"].(string)
		if ids, ok := msg.Body["node_ids"].([]interface{}); ok {
			for _, id := range ids {
				node.NodeIDs = append(node.NodeIDs, id.(string))
			}
		}
		node.Reply(msg, map[string]interface{}{"type": "init_ok"})
	case "echo":
		// TODO: Handle echo message
		// Reply with echo_ok containing the same echo value
		echo, _ := msg.Body["echo"].(string)
		node.Reply(msg, map[string]interface{}{"type": "echo_ok", "echo": echo})
	}
}

func main() {
	node := &Node{}
	scanner := bufio.NewScanner(os.Stdin)
	var wg sync.WaitGroup

	// TODO: Implement concurrent message handling
	// Use goroutines to handle messages concurrently

	for scanner.Scan() {
		var msg Message
		if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			continue
		}

		// TODO: Launch goroutine for concurrent handling

		go func(msg Message) {
			defer wg.Done()
			wg.Add(1)
			node.HandleMessage(msg)
		}(msg)

		// Currently synchronous
	}
	wg.Wait()
}
