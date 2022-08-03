package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/pav5000/debugviz-go/node"
	"github.com/pkg/errors"
)

func main() {
	infile, err := os.Open("input.json")
	if err != nil {
		panic(err)
	}
	defer infile.Close()

	outfile, err := os.Create("out.digraph")
	if err != nil {
		panic(err)
	}
	defer outfile.Close()

	err = jsonToGraphviz(infile, outfile)
	if err != nil {
		log.Fatal(err)
	}
}

func jsonToGraphviz(in io.Reader, out io.Writer) error {
	decoder := json.NewDecoder(in)
	var data node.DebugJSON
	err := decoder.Decode(&data)
	if err != nil {
		return errors.Wrap(err, "decode JSON")
	}

	w := bufio.NewWriter(out)
	defer w.Flush()

	w.WriteString("digraph {\n")
	for id, node := range data.Nodes {
		w.WriteString("  ")
		w.WriteString(id)
		w.WriteString(` [label="`)
		w.WriteString(node.Name)
		w.WriteString("\" shape=rect]\n")
	}

	w.WriteByte('\n')

	for idDst, node := range data.Nodes {
		for _, idSrc := range node.DependsOn {
			w.WriteString(`  `)
			w.WriteString(idSrc)
			w.WriteString(` -> `)
			w.WriteString(idDst)
			w.WriteString("\n")
		}
	}
	w.WriteByte('}')

	return nil
}
