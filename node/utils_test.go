package node

import (
	"context"
	"encoding/json"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var treeLineRe = regexp.MustCompile(`^(\S+)\s*->\s*(\S+)`)

type Links []string

func appendLink(links Links, from, to string) Links {
	return append(links, from+" -> "+to)
}

func assertTreesEqual(t *testing.T, expected string, actual context.Context) {
	rawJSON := FinishDebug(actual)
	var debugJSON DebugJSON
	err := json.Unmarshal(rawJSON, &debugJSON)
	require.NoError(t, err)

	lines := strings.Split(expected, "\n")
	expectedLinks := make(Links, 0, len(lines))
	expectedNamesMap := make(map[string]struct{}, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		matches := treeLineRe.FindStringSubmatch(line)
		if len(matches) != 3 {
			continue
		}
		srcNode := matches[1]
		dstNode := matches[2]
		expectedNamesMap[srcNode] = struct{}{}
		expectedNamesMap[dstNode] = struct{}{}

		expectedLinks = appendLink(expectedLinks, srcNode, dstNode)
	}
	expectedNames := make([]string, 0, len(expectedNamesMap))
	for name := range expectedNamesMap {
		expectedNames = append(expectedNames, name)
	}

	actualLinks := make(Links, 0, len(expectedLinks))
	actualNames := make([]string, 0, len(debugJSON.Nodes))
	for _, node := range debugJSON.Nodes {
		actualNames = append(actualNames, node.Name)
		dstNode := node.Name
		for _, parent := range node.DependsOn {
			srcNode := debugJSON.Nodes[parent].Name
			actualLinks = appendLink(actualLinks, srcNode, dstNode)
		}
	}

	assert.ElementsMatch(t, expectedNames, actualNames)
	assert.ElementsMatch(t, expectedLinks, actualLinks)
}
