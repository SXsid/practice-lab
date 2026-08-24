package repl

// LogEntry represents a single entry in the replicated log.
type LogEntry struct {
	Index int    // Sequential position in the log (starts at 1)
	Term  int    // Leader term that created this entry
	Key   string // key being written
	Value string // value being written
}

// Node represents a single node in the replicated KV cluster.
type Node struct {
	ID          int        // unique node identifier
	Term        int        // current term (1 for now, no election)
	Log         []LogEntry // the replicated log
	CommitIndex int        // Highest log index known to be committed
}

// NewNode creates a node. The followers parameter is accepted but unused for now.
func NewNode(id int, followers []int) *Node {
	return &Node{
		ID:   id,
		Term: 1,
	}
}

// Append adds a new entry to this node's log with the next sequential index.
// The entry is uncommitted — commitIndex is not advanced.
func (n *Node) Append(key, value string) {
	nextIndex := len(n.Log) + 1
	n.Log = append(n.Log, LogEntry{
		Index: nextIndex,
		Term:  n.Term,
		Key:   key,
		Value: value,
	})
}
