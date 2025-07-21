package graph

import (
	"cmp"
	"slices"

	"github.com/averche/docker-compose-graph/internal/compose"
)

type Node struct {
	name                string
	category            Category
	volumeMounts        []compose.VolumeMount
	serviceDependencies []compose.ServiceDependency
}

func NodesFromFiles(files []compose.File) []Node {
	var nodes []Node

	for _, file := range files {
		for name, service := range file.Services {
			nodes = append(nodes, Node{
				name:                name,
				category:            DeterminteCategory(name),
				volumeMounts:        service.VolumeMounts,
				serviceDependencies: service.ServiceDependencies,
			})
		}

		for _, name := range file.Volumes {
			nodes = append(nodes, Node{
				name:     name,
				category: CategoryVolume,
			})
		}
	}

	// sort the nodes to achieve a reproducible output
	slices.SortFunc(nodes, func(a, b Node) int {
		return cmp.Compare(a.name, b.name)
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
