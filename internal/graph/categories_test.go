package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetermineCategory(t *testing.T) {
	assert.Equal(t, CategoryService3, DeterminteServiceCategory("my-cadence", ""))
	assert.Equal(t, CategoryUserInterface, DeterminteServiceCategory("my-ui", ""))
	assert.Equal(t, CategoryStorage, DeterminteServiceCategory("my-storage", ""))
	assert.Equal(t, CategoryTool, DeterminteServiceCategory("my-tool", ""))
	assert.Equal(t, CategoryDatabase, DeterminteServiceCategory("my-postgres", ""))
	assert.Equal(t, CategoryDatabase, DeterminteServiceCategory("my-database", ""))

	assert.Equal(t, CategoryService1, DeterminteServiceCategory("my-cool-service", "service1"))
	assert.Equal(t, CategoryService2, DeterminteServiceCategory("my-cool-service", "service2"))
	assert.Equal(t, CategoryDatabase, DeterminteServiceCategory("my-storage", "database"))
}

func TestConsitency(t *testing.T) {
	assert.Equal(t, int(categoryCount), len(categoryStrings), "inconsitent number of category strings")
	assert.Equal(t, int(categoryCount), len(categoryDecorations), "inconsitent number of category decorations")
}
