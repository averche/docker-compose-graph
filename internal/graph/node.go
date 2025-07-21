package graph

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/averche/docker-compose-graph/internal/compose"
)

type Node struct {
	name                string
	category            Category
	volumeMountsBind    []compose.VolumeMount
	volumeMountsTmpfs   []compose.VolumeMount
	volumeMountsVolume  []compose.VolumeMount
	serviceDependencies []compose.ServiceDependency
}

func NodesFromFiles(files []compose.File) []Node {
	var nodes []Node

	for _, file := range files {
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
					panic(fmt.Sprintf("unexpected volume mount type %s", v.Type))
				}
			}

			nodes = append(nodes, Node{
				name:                name,
				category:            DeterminteServiceCategory(name),
				volumeMountsBind:    volumeMountsBind,
				volumeMountsTmpfs:   volumeMountsTmpfs,
				volumeMountsVolume:  volumeMountsVolume,
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
