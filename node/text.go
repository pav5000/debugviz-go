package node

import "fmt"

func (n *Node) addTextLine(line string) *Node {
	if n == nil {
		return nil
	}
	n.lock.Lock()
	n.text = append(n.text, line)
	n.lock.Unlock()
	return n
}

func (n *Node) Print(values ...interface{}) *Node {
	if n == nil {
		return nil
	}
	text := fmt.Sprint(values...)
	n.addTextLine(text)
	return n
}

func (n *Node) Printf(template string, params ...interface{}) *Node {
	if n == nil {
		return nil
	}
	n.addTextLine(fmt.Sprintf(template, params...))
	return n
}
