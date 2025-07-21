package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderedPresentCategories(t *testing.T) {
	nodes := []Node{
		{name: "my-service1", category: CategoryService},
		{name: "my-database", category: CategoryDatabase},
		{name: "my-proxy", category: CategoryProxy},
		{name: "my-storage", category: CategoryStorage},
		{name: "my-service2", category: CategoryService},
	}

	assert.Equal(t, []Category{
		CategoryService,
		CategoryProxy,
		CategoryDatabase,
		CategoryStorage,
	}, orderedPresentCategories(nodes))
}
