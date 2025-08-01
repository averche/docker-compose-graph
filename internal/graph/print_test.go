package graph

import (
	"strings"
	"testing"

	"github.com/averche/docker-compose-graph/internal/compose"
	"github.com/stretchr/testify/assert"
)

func TestPrintNode(t *testing.T) {
	var b1, b2, b3 strings.Builder

	printNode(&b1, "my-service", "my-service", CategoryService1, true)
	printNode(&b2, "cadence-service", "cadence", CategoryCadence, false)
	printNode(&b3, "my-tool", "tool1", CategoryTool, true)

	assert.Equal(
		t,
		`    my_service                           [shape = "box"        style = "rounded,bold,filled"    fillcolor = "/blues8/7"  color = "/blues8/8"  fontcolor = "white"      fontsize = "8pt"  label = "my-service"];`+"\n",
		b1.String(),
	)
	assert.Equal(
		t,
		`    cadence_service                      [shape = "box"        style = "rounded,bold,filled"    fillcolor = "/orrd8/7"   color = "/orrd8/8"   fontcolor = "white"      label = "cadence"];`+"\n",
		b2.String(),
	)
	assert.Equal(
		t,
		`    my_tool                              [shape = "octagon"    style = "rounded,bold,filled"    fillcolor = "/blues8/7"  color = "/blues8/8"  fontcolor = "white"      fontsize = "8pt"  label = "tool1"];`+"\n",
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
		[]compose.VolumeMount{{
			Type:     compose.VolumeTypeVolume,
			Source:   "my-volume",
			Target:   "/var/log",
			ReadOnly: true,
		}},
	)

	assert.Contains(
		t, `
  my_service                             -> test_service_2                         [style="dashed"];
  my_service                             -> test_service_1                         [style="bold" arrowhead="diamond"];
  my_service                             -> test_service_3                         [style="bold" arrowhead="diamond"];
  my_service                             -> test_service_0                         [style="bold"];
  my_service                             -> my_volume                              [style="dashed"];
`,
		b.String(),
	)
}
