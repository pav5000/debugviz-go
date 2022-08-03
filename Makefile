
test:
	go test ./...

convert:
	go run github.com/pav5000/debugviz-go/cmd/graphviz
	cat out.digraph | dot -Tsvg > output.svg
