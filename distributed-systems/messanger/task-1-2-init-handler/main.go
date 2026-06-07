package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// Node represents a Maelstrom node
type Node struct {
	NodeID    string
	NodeIDs   []string
	NextMsgID int
	mu        sync.Mutex
}

// Message represents a Maelstrom message
type Message struct {
	Src  string                 `json:"src"`
	Dest string                 `json:"dest"`
	Body map[string]interface{} `json:"body"`
}

// Send sends a message to a destination node
func (n *Node) Send(dest string, body map[string]interface{}) {
	n.mu.Lock()
	defer n.mu.Unlock()
	body["msg_id"] = n.NextMsgID
	msg := Message{
		Src:  n.NodeID,
		Dest: dest,
		Body: body,
	}
	byte, err := json.Marshal(msg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	fmt.Println(string(byte))
	n.NextMsgID += 1
}

// Reply sends a response to an incoming request
func (n *Node) Reply(request Message, body map[string]interface{}) {
	body["in_reply_to"] = request.Body["msg_id"]
	n.Send(request.Src, body)
}

func main() {
	node := &Node{}
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		var msg Message
		if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			continue
		}

		msgType, _ := msg.Body["type"].(string)
		if msgType == "init" {
			node.mu.Lock()
			nodeIds, ok := msg.Body["node_ids"].([]string)
			if !ok {
				fmt.Fprintf(os.Stderr, "can't nodeids")
			}
			node.NodeIDs = nodeIds
			nodeId, ok := msg.Body["node_id"].(string)
			if !ok {
				fmt.Fprintln(os.Stderr, "can't nodeid")
			}
			node.NodeID = nodeId
			body := map[string]interface{}{
				"type": "init_ok",
			}
			node.mu.Unlock()
			node.Reply(msg, body)

			// TODO: Handle init message
			// 1. Store node_id and node_ids
			// 2. Reply with init_ok
		}
	}
}
