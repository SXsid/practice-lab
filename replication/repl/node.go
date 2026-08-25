package repl

// LogEntry represents a single entry in the replicated log.
type LogEntry struct {
	Index int    // sequential position in the log (starts at 1)
	Term  int    // leader term that created this entry
	Key   string // key being written
	Value string // value being written
}

// Node represents a single node in the replicated KV cluster.
type Node struct {
	ID          int        // unique node identifier
	Term        int        // current term (1 for now, no election)
	Log         []LogEntry // the replicated log
	CommitIndex int        // highest log index known to be committed

	// Leader-only state
	Followers  []int      // IDs of follower nodes in the cluster
	MatchIndex map[int]int // followerID → highest log index known to be replicated on that follower
}

// NewNode creates a node. The followers parameter defines the cluster membership
// (nil for followers themselves — they don't track other followers).
func NewNode(id int, followers []int) *Node {
	n := &Node{
		ID:        id,
		Term:      1,
		Followers: followers,
	}
	if len(followers) > 0 {
		n.MatchIndex = make(map[int]int)
		for _, f := range followers {
			n.MatchIndex[f] = 0 // followers start with nothing replicated
		}
	}
	return n
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

// ReceiveEntries is the follower-side handler. The leader sends its log entries
// and its current commitIndex. The follower appends any new entries it doesn't
// have yet and advances its own commitIndex.
func (n *Node) ReceiveEntries(entries []LogEntry, leaderCommitIndex int) {
	for _, e := range entries {
		// Only append entries we don't already have
		if e.Index > len(n.Log) {
			n.Log = append(n.Log, e)
		}
	}
	// Advance commitIndex to the leader's, but no further than our log length
	if leaderCommitIndex > n.CommitIndex {
		if leaderCommitIndex <= len(n.Log) {
			n.CommitIndex = leaderCommitIndex
		} else {
			n.CommitIndex = len(n.Log)
		}
	}
}

func (n *Node) HandleAck(followerID int, ackedIndex int) {
	if ackedIndex > n.MatchIndex[followerID] {
		n.MatchIndex[followerID] = ackedIndex
	}
	for idx := n.CommitIndex + 1; idx <= len(n.Log); idx++ {
		ackCount := 1 // leader always has it
		for _, fID := range n.Followers {
			if n.MatchIndex[fID] >= idx {
				ackCount++
			}
		}
		majority := (len(n.Followers)+1)/2 + 1
		if ackCount >= majority {
			n.CommitIndex = idx
		} else {
			break 
		}
	}
}
