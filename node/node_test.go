package node

import (
	"context"
	"testing"
)

func Test_SequentialNodes(t *testing.T) {
	ctx := context.Background()
	ctx = StartDebug(ctx, true)
	_, ctx = New(ctx, "node1", "")
	_, ctx = New(ctx, "node2", "")
	_, ctx = New(ctx, "node3", "")

	assertTreesEqual(t, `
	root -> node1
	node1 -> node2
	node2 -> node3
	`, ctx)
}

func Test_SimpleCollapse(t *testing.T) {
	ctx := context.Background()
	ctx = StartDebug(ctx, true)
	_, ctx = New(ctx, "node1", "")

	_, _ = New(ctx, "node2", "")
	_, _ = New(ctx, "node3", "")

	_, ctx = NewCollapse(ctx, "node4", "")

	assertTreesEqual(t, `
	root -> node1

	node1 -> node2
	node1 -> node3

	node2 -> node4
	node3 -> node4
	`, ctx)
}

func Test_AsymmetricCollapse(t *testing.T) {
	ctx := context.Background()
	ctx = StartDebug(ctx, true)
	_, ctx = New(ctx, "node1", "")

	{
		_, ctx2 := New(ctx, "node2.1", "")
		New(ctx2, "node2.2", "")
	}
	New(ctx, "node3", "")

	_, ctx = NewCollapse(ctx, "node4", "")

	assertTreesEqual(t, `
	root -> node1

	node1   -> node2.1
	node2.1 -> node2.2
	node1   -> node3

	node2.2 -> node4
	node3   -> node4
	`, ctx)
}
