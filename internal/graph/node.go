package graph

import (
	"sort"

	"github.com/averche/docker-compose-graph/internal/compose"
)

type Node struct {
	name         string
	category     Category
	volumes      []compose.Volume
	dependencies []compose.Dependency
}

func NodesFromFiles(files []compose.File) []Node {
	var nodes []Node

	for _, file := range files {
		for name, service := range file.Services {
			nodes = append(nodes, Node{
				name:         name,
				category:     DeterminteCategory(name),
				volumes:      service.Volumes,
				dependencies: service.Dependencies,
			})
		}
	}

	// sort the nodes to achieve a reproducible output
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].name < nodes[j].name
	})

	return nodes
}

// orderedPresentCategories returns an ordered list of categories that are present in the given slice
func orderedPresentCategories(nodes []Node) []Category {
	// bitmap intexed by category
	var exists [categoryCount]bool

	for _, e := range nodes {
		exists[int(e.category)] = true
	}

	var present []Category

	for category := CategoryNone; category < categoryCount; category++ {
		if exists[category] {
			present = append(present, category)
		}
	}

	return present
}
