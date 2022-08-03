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

func Test_2SequentialCollapses(t *testing.T) {
	ctx := context.Background()
	ctx = StartDebug(ctx, true)
	_, ctx = New(ctx, "node1", "")

	{
		_, ctx2 := New(ctx, "node2.1", "")
		New(ctx2, "node2.2", "")
	}
	New(ctx, "node3", "")

	_, ctx = NewCollapse(ctx, "node4", "")

	New(ctx, "node5", "")
	{
		_, ctx2 := New(ctx, "node6.1", "")
		New(ctx2, "node6.2", "")
	}

	_, ctx = NewCollapse(ctx, "node7", "")

	assertTreesEqual(t, `
	root -> node1

	node1   -> node2.1
	node2.1 -> node2.2
	node1   -> node3

	node2.2 -> node4
	node3   -> node4

	node4   -> node5
	node4   -> node6.1
	node6.1 -> node6.2

	node5   -> node7
	node6.2 -> node7
	`, ctx)
}

// func Test_SubCollapse(t *testing.T) {
// 	ctx := context.Background()
// 	ctx = StartDebug(ctx, true)
// 	_, ctx = New(ctx, "node1", "")

// 	{
// 		_, ctx2 := New(ctx, "node2.1", "")

// 		New(ctx2, "node2.2.1", "")
// 		New(ctx2, "node2.2.2", "")

// 		NewCollapse(ctx2, "node2.3", "")
// 	}
// 	New(ctx, "node3", "")

// 	_, ctx = NewCollapse(ctx, "node4", "")

// 	assertTreesEqual(t, `
// 	root -> node1

// 	node1 -> node3
// 	node1 -> node2.1
// 	node2.1 -> node.2.2.1
// 	node2.1 -> node.2.2.2
// 	node.2.2.1 -> node2.3
// 	node.2.2.2 -> node2.3

// 	node3   -> node4
// 	node2.3 -> node4

// 	`, ctx)
// }
