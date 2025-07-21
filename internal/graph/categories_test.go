package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetermineCategory(t *testing.T) {
	assert.Equal(t, CategoryCadence, DeterminteServiceCategory("my-cadence", ""))
	assert.Equal(t, CategoryFrontEnd, DeterminteServiceCategory("my-ui", ""))
	assert.Equal(t, CategoryStorage, DeterminteServiceCategory("my-storage", ""))
	assert.Equal(t, CategoryProxy, DeterminteServiceCategory("my-proxy", ""))
	assert.Equal(t, CategoryDatabase, DeterminteServiceCategory("my-postgres", ""))
	assert.Equal(t, CategoryDatabase, DeterminteServiceCategory("my-database", ""))
	assert.Equal(t, CategoryService, DeterminteServiceCategory("something", ""))
	assert.Equal(t, CategoryVault, DeterminteServiceCategory("my-vault", ""))

	assert.Equal(t, CategoryService, DeterminteServiceCategory("my-cadence", "service"))
	assert.Equal(t, CategoryDatabase, DeterminteServiceCategory("my-storage", "database"))
}

func TestConsitency(t *testing.T) {
	assert.Equal(t, int(categoryCount), len(categoryStrings), "inconsitent number of category strings")
	assert.Equal(t, int(categoryCount), len(categoryDecorations), "inconsitent number of category decorations")
}
