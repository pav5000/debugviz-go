package node

import (
	"github.com/pav5000/debugviz-go/datablocks/request"
	"github.com/pav5000/debugviz-go/datablocks/table"
)

func (n *Node) addDataBlock(block DataBlock) {
	if n == nil {
		return
	}
	n.lock.Lock()
	n.dataBlocks = append(n.dataBlocks, block)
	n.lock.Unlock()
}

func (n *Node) NewRequestBlock(serviceName, handlerName string) *request.Block {
	if n == nil {
		return nil
	}
	block := request.New(serviceName, handlerName)
	n.addDataBlock(block)
	return block
}

func (n *Node) NewTableBlock(caption string) *table.Block {
	if n == nil {
		return nil
	}
	block := table.New(caption)
	n.addDataBlock(block)
	return block
}
