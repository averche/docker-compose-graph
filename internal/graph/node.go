package graph

import (
	"cmp"
	"slices"

	"github.com/averche/docker-compose-graph/internal/compose"
)

const (
	graphNodeCategory = "graph.node.category"
	graphNodeLabel    = "graph.node.label"
)

type NodeGroup struct {
	Label string
	Nodes []Node
}

type Node struct {
	Name                string
	Label               string
	Category            Category
	VolumeMounts        []compose.VolumeMount
	ServiceDependencies []compose.ServiceDependency
}

func NodesFromFile(file compose.File) []Node {
	var nodes []Node

	for name, service := range file.Services {
		var volumeMounts []compose.VolumeMount

		for _, v := range service.VolumeMounts {
			// we only care about type "volume" for now
			if v.Type == compose.VolumeTypeVolume {
				volumeMounts = append(volumeMounts, v)
			}
		}

		label, ok := service.Labels[graphNodeLabel]
		if !ok {
			label = name
		}

		nodes = append(nodes, Node{
			Name:                name,
			Label:               label,
			Category:            DeterminteServiceCategory(name, service.Labels[graphNodeCategory]),
			VolumeMounts:        volumeMounts,
			ServiceDependencies: service.ServiceDependencies,
		})
	}

	for _, name := range file.Volumes {
		nodes = append(nodes, Node{
			Name:     name,
			Category: CategoryVolume,
		})
	}

	// sort the nodes to achieve a reproducible output
	slices.SortFunc(nodes, func(a, b Node) int {
		return cmp.Compare(a.Name, b.Name)
	})

	return nodes
}

// orderedPresentCategories returns an ordered list of categories that are present in the given slice
func orderedPresentCategories(groups []NodeGroup) []Category {
	// bitmap intexed by category
	var exists [categoryCount]bool

	for _, group := range groups {
		for _, node := range group.Nodes {
			exists[int(node.Category)] = true
		}
	}

	var present []Category

	for category := CategoryNone; category < categoryCount; category++ {
		if exists[category] {
			present = append(present, category)
		}
	}

	return present
}
