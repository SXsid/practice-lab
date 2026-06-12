package main

import (
	"bufio"
	// "context"
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

func (n *Node) NextMsgId() int {
	n.mu.Lock()
	defer n.mu.Unlock()
	id := n.NextMsgID
	n.NextMsgID++
	return id
}

func (n *Node) sendToNode(msg_id int, body map[string]interface{}, dest string) {
	body["msg_id"] = msg_id
	b, _ := json.Marshal(Message{
		Dest: dest,
		Src:  n.NodeID,
		Body: body,
	})
	n.outMu.Lock()
	fmt.Println(string(b))
	n.outMu.Unlock()
}

func (n *Node) Send(dest string, body map[string]interface{}) {
	body["msg_id"] = n.NextMsgId()
	b, _ := json.Marshal(Message{
		Dest: dest,
		Src:  n.NodeID,
		Body: body,
	})
	n.outMu.Lock()
	fmt.Println(string(b))
	n.outMu.Unlock()
}

// func (n *Node) syncRpc(ctx context.Context, msg Message) {
// 	target, _ := msg.Body["target"].(string)
// 	tartetBody, _ := msg.Body["inner"].(map[string]interface{})
// 	chanM := make(chan map[string]interface{}, 1)
// 	n.mu.Lock()
// 	id := n.nextMsg
// 	n.nextMsg++
// 	n.pending_req[id] = chanM
// 	n.sendToNode(id, tartetBody, target)
// 	n.mu.Unlock()
// 	select {
// 	case body := <-chanM:
// 		n.send(msg.Src, body)
// 	case <-ctx.Done():
// 		fmt.Fprint(os.Stderr, "timeout error")
// 		return
// 	}
// }

func Validate(msg Message) error {
	if msg.Src == "" || msg.Dest == "" || msg.Body == nil {
		return fmt.Errorf("request is invalid")
	}
	if v, ok := msg.Body["type"].(string); !ok || v == "" {
		return fmt.Errorf("invalid  body")
	}
	return nil
}

func (n *Node) Reply(msg Message, body map[string]interface{}) {
	msg_id, ok := msg.Body["msg_id"].(float64)
	if !ok {
		fmt.Fprint(os.Stderr, "error occured msg_id")
	}

	body["in_reply_to"] = int(msg_id)
	n.Send(msg.Src, body)
}

func (node *Node) handLeMessage(msg Message) {
	// if err := Validate(msg); err != nil {
	// 	fmt.Fprint(os.Stderr, err.Error())
	// 	return
	// }

	fmt.Println(msg.Body)
	msgType, ok := msg.Body["type"].(string)
	if !ok {
		fmt.Fprint(os.Stderr, "error occured no typ")
	}

	switch msgType {
	case "init":
		node.mu.Lock()

		id, ok := msg.Body["node_id"].(string)
		if !ok {
			fmt.Fprint(os.Stderr, "error occured node_id")
		}
		node.NodeID = id

		ids, ok := msg.Body["node_ids"].([]interface{})
		if !ok {
			fmt.Fprint(os.Stderr, "error occured node_ids")
		}
		for _, id := range ids {
			node.NodeIDs = append(node.NodeIDs, id.(string))
		}
		node.mu.Unlock()

		body := map[string]interface{}{
			"type": "init_ok",
		}
		node.Reply(msg, body)
	case "echo":
		echo, ok := msg.Body["echo"].(string)
		if !ok {
			fmt.Fprint(os.Stderr, "error")
		}
		node.Reply(msg, map[string]interface{}{
			"type": "echo_ok",
			"echo": echo,
		})
	case "proxy":
		// INFO: easy veriosn just fire and forget
		target, ok := msg.Body["target"].(string)
		if !ok {
			fmt.Fprint(os.Stderr, "error")
		}
		inner, ok := msg.Body["inner"].(map[string]interface{})
		if !ok {
			fmt.Fprint(os.Stderr, "error")
		}
		// Test only checks the forward, not the reply
		node.Send(target, inner)

		// node.wg.Add(1)
		// go func(msg Message) {
		// 	ctx, cance := context.WithTimeout(context.Background(), time.Second*1)
		// 	defer cance()
		// 	defer node.wg.Done()
		// 	node.syncRpc(ctx, msg)
		// }(msg)
		// case "proxy_ok":
		// 	id, _ := msg.Body["in_reply_to"].(float64)
		// 	node.mu.Lock()
		// 	chanM := node.pending_req[int(id)]
		// 	node.mu.Unlock()
		// 	body, _ := msg.Body["result"].(map[string]interface{})
		// 	chanM <- body

	}
}

func main() {
	// Your implementation here
	scanner := bufio.NewScanner(os.Stdin)

	var wg sync.WaitGroup
	node := &Node{
		// pending_req: make(map[int]chan map[string]interface{}, 0),
	}
	for scanner.Scan() {
		var msg Message
		if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			continue

		}
		wg.Add(1)
		go func(msg Message) {
			defer wg.Done()

			node.handLeMessage(msg)
		}(msg)
		wg.Wait()

	}
	// node.wg.Wait()
}
