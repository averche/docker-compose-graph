package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJoin(t *testing.T) {
	// edge cases
	assert.Equal(t, "", JoinStyles(nil, "something"))
	assert.Equal(t, "", JoinStyles([]Style{}, ","))

	// single
	assert.Equal(t, "bold", JoinStyles([]Style{Bold}, ","))

	// multiple
	assert.Equal(t, "filled,dotted", JoinStyles([]Style{Filled, Dotted}, ","))
	assert.Equal(t, "filled,dotted,rounded", JoinStyles([]Style{Filled, Dotted, Rounded}, ","))
}
