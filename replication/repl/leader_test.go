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
