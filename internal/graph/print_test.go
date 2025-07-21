package graph

import (
	"strings"
	"testing"

	"github.com/averche/docker-compose-graph/internal/compose"
	"github.com/stretchr/testify/assert"
)

func TestPrintNode(t *testing.T) {
	var b1, b2, b3 strings.Builder

	printNode(&b1, "my-service", CategoryService, false, true)
	printNode(&b2, "cadence-service", CategoryCadence, false, false)
	printNode(&b3, "my-proxy", CategoryProxy, true, true)

	assert.Equal(
		t,
		`  my_service                 [shape = "box"        style = "rounded,bold,filled"    fillcolor = "/blues8/7"  color = "/blues8/8"  fontcolor = "white"      fontsize = "8pt"  label = "my-service"];`+"\n",
		b1.String(),
	)
	assert.Equal(
		t,
		`  cadence_service            [shape = "box"        style = "rounded,bold,filled"    fillcolor = "/brbg8/7"   color = "/brbg8/8"   fontcolor = "white"      label = "cadence-service"];`+"\n",
		b2.String(),
	)
	assert.Equal(
		t,
		`    my_proxy                 [shape = "diamond"    style = "rounded,bold,filled"    fillcolor = "/bupu8/7"   color = "/bupu8/8"   fontcolor = "white"      fontsize = "8pt"  label = "my-proxy"];`+"\n",
		b3.String(),
	)
}

func TestPrintDependencies(t *testing.T) {
	var b strings.Builder

	printDependencies(
		&b,
		"my-service",
		[]compose.ServiceDependency{{
			On:        "test-service-2",
			Condition: compose.ConditionServiceStarted,
		}, {
			On:        "test-service-1",
			Condition: compose.ConditionServiceHealthy,
		}, {
			On:        "test-service-3",
			Condition: compose.ConditionServiceHealthy,
		}, {
			On:        "test-service-0",
			Condition: compose.ConditionServiceCompletedSuccessfully,
		}},
	)
	t.Log(b.String())

	assert.Contains(
		t, `
  my_service                 -> test_service_0             [style="bold"];
  my_service                 -> test_service_1             [arrowhead="diamond" style="bold"];
  my_service                 -> test_service_2             [style="dashed"];
  my_service                 -> test_service_3             [arrowhead="diamond" style="bold"];
`,
		b.String(),
	)
}
