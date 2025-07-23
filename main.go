package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/averche/docker-compose-graph/internal/compose"
	"github.com/averche/docker-compose-graph/internal/graph"
)

func main() {
	var groups []graph.NodeGroup

	for _, path := range os.Args[1:] {
		// parse each file & construct graph nodes
		f, err := compose.ParseFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error :: could not parse '%s': %v\n", path, err)
			os.Exit(1)
		}

		groups = append(groups, graph.NodeGroup{
			Label: filepath.Base(path),
			Nodes: graph.NodesFromFile(f),
		})
	}

	graph.Print(os.Stdout, groups)
}
