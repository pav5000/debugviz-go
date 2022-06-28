package node

import (
	"context"
	"sync"
)

type DebugViz struct {
	lock  sync.Mutex
	nodes []*Node
}

// StartDebug стартует отладку, если передано true
// сделано так, чтобы сюда можно было просто передать поле debug из запроса
// после вызова этой функции с true, начинет работать создание нод ниже по контексту
func StartDebug(ctx context.Context, debug bool) context.Context {
	if !debug {
		return ctx
	}
	debugviz := &DebugViz{}
	ctx = injectDebugIntoContext(ctx, debugviz)
	ctx, root := initRootNode(ctx)
	debugviz.registerNode(root)
	return ctx
}

// FinishDebug строит финальный JSON со всей отладочной информацией
func FinishDebug(ctx context.Context) []byte {
	debugviz := getDebugFromContext(ctx)
	if debugviz == nil {
		return nil
	}
	return debugviz.JSON()
}

func (d *DebugViz) registerNode(node *Node) {
	d.lock.Lock()
	d.nodes = append(d.nodes, node)
	d.lock.Unlock()
}
