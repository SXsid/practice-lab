package repl

import "testing"

// Test 1: A leader appends a client write to its log with a sequential index.
//
// Invariant guarded:
//   - Log entries have monotonically increasing indices starting at 1.
//   - Entries are uncommitted when first appended (commitIndex stays behind).
//
// What a naive/wrong impl would do that this catches:
//   - Directly writing to a map (no log, no ordering)
//   - Committing immediately without waiting for replication
//   - Assigning non-sequential indices
func TestLeader_AppendCreatesSequentialUncommittedEntries(t *testing.T) {
	leader := NewNode(1, nil) // node ID=1, no followers yet

	// First write
	leader.Append("x", "42")

	if len(leader.Log) != 1 {
		t.Fatalf("expected 1 log entry, got %d", len(leader.Log))
	}
	if leader.Log[0].Index != 1 {
		t.Errorf("first entry index: want 1, got %d", leader.Log[0].Index)
	}
	if leader.Log[0].Key != "x" || leader.Log[0].Value != "42" {
		t.Errorf("entry content: want x=42, got %s=%s", leader.Log[0].Key, leader.Log[0].Value)
	}
	if leader.Log[0].Term != 1 {
		t.Errorf("term: want 1, got %d", leader.Log[0].Term)
	}

	// Second write — index must be 2
	leader.Append("y", "99")

	if len(leader.Log) != 2 {
		t.Fatalf("expected 2 log entries, got %d", len(leader.Log))
	}
	if leader.Log[1].Index != 2 {
		t.Errorf("second entry index: want 2, got %d", leader.Log[1].Index)
	}

	// commitIndex must still be 0 — nothing replicated, nothing committed
	if leader.CommitIndex != 0 {
		t.Errorf("commitIndex should be 0 before replication, got %d", leader.CommitIndex)
	}
}

// Test 2: CommitIndex advances only when a majority of nodes have the entry.
//
// Invariant guarded:
//   - A write is committed if and only if it exists on a majority of nodes.
//   - commitIndex must NOT advance with fewer than majority acks.
//   - commitIndex MUST advance once majority is reached.
//
// What a naive/wrong impl would do that this catches:
//   - Committing on leader-only ack (commitIndex advances on Append)
//   - Committing on ALL-node ack (commitIndex never advances until every node acks)
//   - Not counting the leader as part of the quorum
func TestLeader_CommitsOnlyOnMajorityAck(t *testing.T) {
	// 3-node cluster: leader (ID=1) + 2 followers (IDs 2,3)
	leader := NewNode(1, []int{2, 3})
	follower1 := NewNode(2, nil)
	follower2 := NewNode(3, nil)
	// Leader accepts a write
	leader.Append("x", "42")

	// Before any replication — commitIndex must be 0
	if leader.CommitIndex != 0 {
		t.Fatalf("commitIndex should be 0 before replication, got %d", leader.CommitIndex)
	}

	// Simulate: leader replicates entry to follower1, follower1 acks
	// In a real system this is an RPC. Here we simulate the mechanics.
	follower1.ReceiveEntries(leader.Log, leader.CommitIndex)
	leader.HandleAck(follower1.ID, 1) // follower1 acked up to index 1

	// Leader (has it) + follower1 (acked) = 2/3 = majority → committed
	if leader.CommitIndex != 1 {
		t.Errorf("commitIndex should be 1 after majority ack, got %d", leader.CommitIndex)
	}

	// follower2 hasn't received anything yet — that's fine, majority already reached
	_ = follower2
}

// Test 3: CommitIndex does NOT advance without majority.
//
// Invariant guarded:
//   - In a 5-node cluster, 2 follower acks are needed (leader + 2 = 3/5).
//   - 1 follower ack alone must NOT advance commitIndex.
func TestLeader_DoesNotCommitWithoutMajority(t *testing.T) {
	// 5-node cluster: leader + 4 followers
	leader := NewNode(1, []int{2, 3, 4, 5})
	follower1 := NewNode(2, nil)

	leader.Append("x", "42")

	// Only 1 follower acks — leader + 1 = 2/5, not majority
	follower1.ReceiveEntries(leader.Log, leader.CommitIndex)
	leader.HandleAck(follower1.ID, 1)

	if leader.CommitIndex != 0 {
		t.Errorf("commitIndex should still be 0 with only 1 follower ack in 5-node cluster, got %d", leader.CommitIndex)
	}

	// Second follower acks — leader + 2 = 3/5 = majority → commit
	follower2 := NewNode(3, nil)
	follower2.ReceiveEntries(leader.Log, leader.CommitIndex)
	leader.HandleAck(follower2.ID, 1)

	if leader.CommitIndex != 1 {
		t.Errorf("commitIndex should be 1 after 2 follower acks in 5-node cluster, got %d", leader.CommitIndex)
	}
}

// Test 4: Full concurrent write — leader.Write() fans out to follower goroutines,
// blocks until majority ack, then returns success.
//
// Invariant guarded:
//   - Write() is a blocking call that only returns once the entry is committed.
//   - Followers run as independent goroutines processing entries concurrently.
//   - The leader does not return to the client before majority is confirmed.
//
// What a naive/wrong impl would do that this catches:
//   - Returning success before any replication (fire-and-forget)
//   - Deadlocking because ack channel is unbuffered / goroutine leaks
//   - Not wiring up the channel-based transport correctly
func TestLeader_WriteConcurrentMajorityAck(t *testing.T) {
	// Build a 3-node cluster with channels connecting them
	cluster := NewCluster(3)

	// Write should block until majority ack, then return nil error
	err := cluster.Write("x", "42")
	if err != nil {
		t.Fatalf("Write returned error: %v", err)
	}

	// After Write returns, commitIndex must be 1
	if cluster.Leader().CommitIndex != 1 {
		t.Errorf("leader commitIndex: want 1, got %d", cluster.Leader().CommitIndex)
	}

	// At least one follower must have the entry (majority = leader + 1 follower)
	acked := 0
	for _, f := range cluster.Followers() {
		if len(f.Log) >= 1 && f.Log[0].Key == "x" {
			acked++
		}
	}
	if acked < 1 {
		t.Errorf("expected at least 1 follower to have the entry, got %d", acked)
	}
}
