package graph

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/averche/docker-compose-graph/internal/compose"
)

type NodeGroup struct {
	Name  string
	Nodes []Node
}

type Node struct {
	Name                string
	Category            Category
	VolumeMountsBind    []compose.VolumeMount
	VolumeMountsTmpfs   []compose.VolumeMount
	VolumeMountsVolume  []compose.VolumeMount
	ServiceDependencies []compose.ServiceDependency
}

func NodesFromFile(file compose.File) []Node {
	var nodes []Node

	for name, service := range file.Services {
		var (
			volumeMountsBind   []compose.VolumeMount
			volumeMountsTmpfs  []compose.VolumeMount
			volumeMountsVolume []compose.VolumeMount
		)

		for _, v := range service.VolumeMounts {
			switch v.Type {
			case compose.VolumeTypeBind:
				volumeMountsBind = append(volumeMountsBind, v)
			case compose.VolumeTypeTmpfs:
				volumeMountsTmpfs = append(volumeMountsTmpfs, v)
			case compose.VolumeTypeVolume:
				volumeMountsVolume = append(volumeMountsVolume, v)
			default:
				panic(fmt.Sprintf("unexpected volume mount type %q", v.Type))
			}
		}

		nodes = append(nodes, Node{
			Name:                name,
			Category:            DeterminteServiceCategory(name),
			VolumeMountsBind:    volumeMountsBind,
			VolumeMountsTmpfs:   volumeMountsTmpfs,
			VolumeMountsVolume:  volumeMountsVolume,
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
