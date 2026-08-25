package repl

import (
	"fmt"
	"sync"
	"time"
)

type Ack struct {
	FollowerID int
	AckedIndex int
}
type  Entries struct{
	CommitIndex int
	Logs []LogEntry
}

type Cluster struct {
	nodes    []*Node
	entryChs map[int]chan Entries
	ackCh    chan Ack
	mu       sync.Mutex
}

func NewCluster(n int) *Cluster {
	followerIDs := make([]int, n-1)
	for i := range followerIDs {
		followerIDs[i] = i + 2 // IDs: 2, 3, ..., n
	}

	leader := NewNode(1, followerIDs)
	nodes := []*Node{leader}

	entryChs := make(map[int]chan Entries)
	ackCh := make(chan Ack, n) // buffered to avoid goroutine blocking(inshort every follower can write idependently)

	c := &Cluster{
		nodes:    nodes,
		entryChs: entryChs,
		ackCh:    ackCh,
	}

	for _, fID := range followerIDs {
		follower := NewNode(fID, nil)
		nodes = append(nodes, follower)
		entryChs[fID] = make(chan Entries, 10) // buffered entry channel(so we cna put more log entry without worrying that follower read will block it)

		go c.runFollower(follower)
	}

	c.nodes = nodes
	return c
}

func (c *Cluster) runFollower(f *Node) {
	ch := c.entryChs[f.ID]
	for entries := range ch {
		f.ReceiveEntries(entries.Logs, entries.CommitIndex) // commitIndex propagation deferred for simplicity
		c.ackCh <- Ack{
			FollowerID: f.ID,
			AckedIndex: len(f.Log),
		}
	}
}

// Leader returns the leader node (node 0 in the cluster).
func (c *Cluster) Leader() *Node {
	return c.nodes[0]
}

// Followers returns all follower nodes.
func (c *Cluster) Followers() []*Node {
	return c.nodes[1:]
}

func (c *Cluster) Write(key, value string) error {
	leader := c.Leader()

	// Step 1: Append to leader's log
	c.mu.Lock()
	leader.Append(key, value)
	logSnapshot := make([]LogEntry, len(leader.Log))
	entries:=Entries{
		CommitIndex: leader.CommitIndex,
		Logs: logSnapshot,

	}
	target_index:=len(logSnapshot)
	copy(logSnapshot, leader.Log)
	c.mu.Unlock()

	// Step 2: Fan out to all followers
	for _, fID := range leader.Followers {
		c.entryChs[fID] <-entries
	}

	timeout := time.After(5 * time.Second)
	for target_index> leader.CommitIndex{
		select {
		case ack := <-c.ackCh:
			c.mu.Lock()
			leader.HandleAck(ack.FollowerID, ack.AckedIndex)
			c.mu.Unlock()
		case <-timeout:
			return fmt.Errorf("write timed out: majority not reached")
		}
	}

	return nil
}
