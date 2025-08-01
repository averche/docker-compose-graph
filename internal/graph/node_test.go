package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderedPresentCategories(t *testing.T) {
	groups := []NodeGroup{{
		Label: "docker-compose-1.yaml",
		Nodes: []Node{
			{Name: "my-service1", Category: CategoryService1},
			{Name: "my-database", Category: CategoryDatabase},
		},
	}, {
		Label: "docker-compose-2.yaml",
		Nodes: []Node{
			{Name: "my-proxy", Category: CategoryTool},
			{Name: "my-storage", Category: CategoryStorage},
			{Name: "my-service2", Category: CategoryService1},
		},
	}}

	assert.Equal(t, []Category{
		CategoryService1,
		CategoryTool,
		CategoryDatabase,
		CategoryStorage,
	}, orderedPresentCategories(groups))
}
