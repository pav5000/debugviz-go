package node

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
)

type DataBlock interface {
	Type() string
	Data() interface{}
}

type NodeID int64

func (n NodeID) String() string {
	return strconv.FormatInt(int64(n), 10)
}

type Node struct {
	lock     sync.Mutex
	children []*Node // нужно для механизма collapse

	id         NodeID
	name       string
	descr      string
	dependsOn  []NodeID
	text       []string
	dataBlocks []DataBlock

	encapsulatedDebugViz json.RawMessage // для будущего механизма межсервисного дебага
}

func new(ctx context.Context, name, descr string) (thisNode *Node, parentNode *Node, newCtx context.Context) {
	debugviz := getDebugFromContext(ctx)
	if debugviz == nil {
		return nil, nil, ctx
	}

	parentNode = getNodeFromContext(ctx)
	if parentNode == nil {
		return nil, nil, ctx
	}

	thisNode = &Node{
		id:    NodeID(rand.Int63()),
		name:  name,
		descr: descr,
	}

	debugviz.registerNode(thisNode)

	newCtx = injectNodeIntoContext(ctx, thisNode)
	return thisNode, parentNode, newCtx
}

func initRootNode(ctx context.Context) (context.Context, *Node) {
	node := &Node{
		id:    1,
		name:  "root",
		descr: "Начало отладки",
	}
	return injectNodeIntoContext(ctx, node), node
}

// New создает новую отладочную ноду
// возвращенный контекст следует обязательно использовать в последующей работе приложения
// ибо через него будет определяться, какая нода от какой зависит
func New(ctx context.Context, name, descr string) (*Node, context.Context) {
	node, parent, ctx := new(ctx, name, descr)
	if node == nil || parent == nil {
		return nil, ctx
	}

	parent.addChild(node)
	node.dependsOn = append(node.dependsOn, parent.id)

	return node, ctx
}

// NewCollapse создает новую ноду, которая идет после всех детей текущей ноды контекста
// используется, когда у вас в функции было запущено несколько горутин с нодами,
// они отработали и нужно, чтобы создаваемая нода шла после них в схеме
// то есть, когда выполнение программы разветвилось, а потом результат выполнения соединился обратно
//                ┌───────┐
//                │ node1 │ <- контекст этой ноды передаём в NewCollapse(...)
//                └┬──┬──┬┘
//     ┌───────────┘  │  └───────────┐
// ┌───▼───┐      ┌───▼───┐      ┌───▼───┐
// │ node2 │      │ node3 │      │ node4 │ <- эти ребята отработали независимо в горутинах
// └───┬───┘      └───┬───┘      └───┬───┘
//     └────────────┐ │ ┌────────────┘
//           -------│ │ │-------- collapse
//                ┌─▼─▼─▼─┐
//                │ node5 │ <- нода, созданная функцией NewCollapse(...)
//                └───────┘
func NewCollapse(ctx context.Context, name, descr string) (*Node, context.Context) {
	fmt.Println("### NewCollapse:", name)
	node, parent, ctx := new(ctx, name, descr)
	if node == nil || parent == nil {
		return nil, ctx
	}

	childlessNodes := parent.findAllChildlessChildren()
	dedupMap := make(map[NodeID]struct{}, len(childlessNodes))

	for _, childlessNode := range childlessNodes {
		if _, ok := dedupMap[childlessNode.id]; ok {
			continue
		}
		dedupMap[childlessNode.id] = struct{}{}

		childlessNode.addChild(node)
		node.dependsOn = append(node.dependsOn, childlessNode.id)
	}

	return node, ctx
}

func (n *Node) addChild(node *Node) {
	if n == nil {
		return
	}
	n.lock.Lock()
	n.children = append(n.children, node)
	n.lock.Unlock()
}

func (n *Node) findAllChildlessChildren() []*Node {
	if n == nil {
		return nil
	}
	children := n.getChildren()
	if len(children) == 0 {
		return []*Node{n}
	}

	var nodes []*Node
	for _, child := range children {
		nodes = append(nodes, child.findAllChildlessChildren()...)
	}
	return nodes
}

func (n *Node) getChildren() []*Node {
	if n == nil {
		return nil
	}
	n.lock.Lock()
	children := n.children
	n.lock.Unlock()
	return children
}
