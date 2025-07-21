package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetermineCategory(t *testing.T) {
	assert.Equal(t, CategoryCadence, DeterminteCategory("my-cadence"))
	assert.Equal(t, CategoryFrontEnd, DeterminteCategory("my-ui"))
	assert.Equal(t, CategoryStorage, DeterminteCategory("my-storage"))
	assert.Equal(t, CategoryProxy, DeterminteCategory("my-proxy"))
	assert.Equal(t, CategoryDatabase, DeterminteCategory("my-postgres"))
	assert.Equal(t, CategoryDatabase, DeterminteCategory("my-database"))
	assert.Equal(t, CategoryService, DeterminteCategory("something"))
	assert.Equal(t, CategoryVault, DeterminteCategory("my-vault"))
}

func TestOrderedPresentCategories(t *testing.T) {
	entities := []entity{
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
	}, OrderedPresentCategories(entities))
}
