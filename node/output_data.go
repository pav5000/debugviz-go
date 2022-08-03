package node

import (
	"encoding/json"
	"strconv"
)

const version = "0.0.1"

type DebugJSON struct {
	ThisIsDebugger  bool
	DebuggerVersion string
	Nodes           map[string]NodeJSON
}

type NodeJSON struct {
	ID         string
	Name       string
	Descr      string
	DependsOn  []string
	Text       []string
	DataBlocks []DataBlockJSON
	DebugViz   json.RawMessage
}

type DataBlockJSON struct {
	Type string
	Data json.RawMessage
}

type Dummy struct{}

func (n *Node) Data() NodeJSON {
	if n == nil {
		return NodeJSON{}
	}
	dataBlocksJSON := make([]DataBlockJSON, 0, len(n.dataBlocks))

	for _, block := range n.dataBlocks {
		rawJSON, err := json.Marshal(block.Data())
		if err != nil {
			panic(err)
		}

		dataBlocksJSON = append(dataBlocksJSON, DataBlockJSON{
			Type: block.Type(),
			Data: rawJSON,
		})
	}

	deps := make([]string, 0, len(n.dependsOn))
	for _, id := range n.dependsOn {
		deps = append(deps, id.String())
	}

	return NodeJSON{
		ID:         strconv.FormatInt(int64(n.id), 10),
		Name:       n.name,
		Descr:      n.descr,
		DependsOn:  deps,
		Text:       n.text,
		DataBlocks: dataBlocksJSON,
		DebugViz:   n.encapsulatedDebugViz,
	}
}

func (d *DebugViz) JSON() []byte {
	if d == nil {
		return []byte("null")
	}
	d.lock.Lock()
	defer d.lock.Unlock()

	out := DebugJSON{
		ThisIsDebugger:  true,
		DebuggerVersion: version,
		Nodes:           make(map[string]NodeJSON, len(d.nodes)),
	}

	for _, node := range d.nodes {
		if node == nil {
			continue
		}
		data := node.Data()
		if data.ID != "" {
			out.Nodes[node.id.String()] = data
		}
	}

	rawJSON, _ := json.Marshal(out)
	return rawJSON
}
