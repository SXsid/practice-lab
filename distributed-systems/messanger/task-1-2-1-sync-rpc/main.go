package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

type Node struct {
	wg          sync.WaitGroup
	nodeId      string
	nodeIds     []string
	nextMsg     int
	pending_req map[int]chan map[string]any
	mu          sync.Mutex
	outMu       sync.Mutex
}

type Message struct {
	Src  string         `json:"src"`
	Dest string         `json:"dest"`
	Body map[string]any `json:"body"`
}

func (n *Node) NextMsgId() int {
	n.mu.Lock()
	defer n.mu.Unlock()
	id := n.nextMsg
	n.nextMsg++
	return id
}

func (n *Node) sendToNode(msg_id int, body map[string]any, dest string) {
	body["msg_id"] = msg_id
	b, _ := json.Marshal(Message{
		Dest: dest,
		Src:  n.nodeId,
		Body: body,
	})
	n.outMu.Lock()
	fmt.Println(string(b))
	n.outMu.Unlock()
}

func (n *Node) send(dest string, body map[string]any) {
	body["msg_id"] = n.NextMsgId()
	b, _ := json.Marshal(Message{
		Dest: dest,
		Src:  n.nodeId,
		Body: body,
	})
	n.outMu.Lock()
	fmt.Println(string(b))
	n.outMu.Unlock()
}

func (n *Node) syncRpc(ctx context.Context, msg Message) {
	target, _ := msg.Body["target"].(string)
	tartetBody, _ := msg.Body["inner"].(map[string]any)
	chanM := make(chan map[string]any, 1)
	n.mu.Lock()
	id := n.nextMsg
	n.nextMsg++
	n.pending_req[id] = chanM
	n.sendToNode(id, tartetBody, target)
	n.mu.Unlock()
	select {
	case body := <-chanM:
		n.send(msg.Src, body)
	case <-ctx.Done():
		fmt.Fprint(os.Stderr, "timeout error")
		return
	}
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
	msg_id, _ := msg.Body["msg_id"].(float64)
	body["in_reply_to"] = msg_id
	n.send(msg.Src, body)
}

func (node *Node) handLeMessage(msg Message) {
	Validate(msg)

	msgType, _ := msg.Body["type"].(string)

	switch msgType {
	case "init":

		node.nodeId, _ = msg.Body["node_id"].(string)
		node.nodeIds, _ = msg.Body["node_ids"].([]string)
		body := map[string]any{
			"type": "init_ok",
		}
		node.reply(msg, body)
	case "echo":
		node.reply(msg, map[string]any{
			"type": "echo_ok",
			"echo": msg.Body["echo"],
		})
	case "proxy":

		node.wg.Add(1)
		go func(msg Message) {
			ctx, cance := context.WithTimeout(context.Background(), time.Second*1)
			defer cance()
			defer node.wg.Done()
			node.syncRpc(ctx, msg)
		}(msg)
	case "proxy_ok":
		id, _ := msg.Body["in_reply_to"].(float64)
		node.mu.Lock()
		chanM := node.pending_req[int(id)]
		node.mu.Unlock()
		body, _ := msg.Body["result"].(map[string]any)
		chanM <- body

	default:
		fmt.Fprint(os.Stderr, "not a valid message type")
	}
}

func main() {
	// Your implementation here
	scanner := bufio.NewScanner(os.Stdin)

	node := &Node{
		pending_req: make(map[int]chan map[string]any, 0),
	}
	for scanner.Scan() {
		var msg Message
		if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
			fmt.Fprint(os.Stderr, err.Error())
		}
		node.handLeMessage(msg)

	}
	node.wg.Wait()
}
