package main

import (
	"fmt"
	"os"

	"github.com/averche/docker-compose-graph/internal/compose"
	"github.com/averche/docker-compose-graph/internal/graph"
)

func main() {
	files, err := compose.ParseMultiple(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error :: could not parse docker-compose yaml configuration(s): %v\n", err)
		os.Exit(1)
	}

	graph.Print(os.Stdout, graph.NodesFromFiles(files))
}
