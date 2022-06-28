package node

import "context"

type debugvizKeyType struct{}

var debugvizKey = debugvizKeyType{}

func getDebugFromContext(ctx context.Context) *DebugViz {
	debugRaw := ctx.Value(debugvizKey)
	if debugRaw == nil {
		return nil
	}
	debugCasted, ok := debugRaw.(*DebugViz)
	if !ok {
		return nil
	}
	return debugCasted
}

func injectDebugIntoContext(ctx context.Context, node *DebugViz) context.Context {
	return context.WithValue(ctx, debugvizKey, node)
}

type nodeKeyType struct{}

var nodeKey = nodeKeyType{}

func getNodeFromContext(ctx context.Context) *Node {
	nodeRaw := ctx.Value(nodeKey)
	if nodeRaw == nil {
		return nil
	}
	castedNode, ok := nodeRaw.(*Node)
	if !ok {
		return nil
	}
	return castedNode
}

func injectNodeIntoContext(ctx context.Context, node *Node) context.Context {
	return context.WithValue(ctx, nodeKey, node)
}
