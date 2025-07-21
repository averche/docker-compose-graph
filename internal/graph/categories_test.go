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

func TestConsitency(t *testing.T) {
	assert.Equal(t, int(categoryCount), len(categoryStrings), "inconsitent number of category strings")
	assert.Equal(t, int(categoryCount), len(categoryDecorations), "inconsitent number of category decorations")
}
