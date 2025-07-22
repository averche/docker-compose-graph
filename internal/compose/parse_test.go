package compose

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFile(t *testing.T) {
	dockerComposeYaml := `
version: "3.9"

services:

  service1:
    image: service1:latest
    labels:
      - "graph.category=database"
    volumes:
      - type:   bind
        source: /dir/from
        target: /dir/to

  service2:
    image: service2:latest
    depends_on:
      - service1
    labels:
      graph.category: proxy
    volumes:
      - my-volume:/some/data:ro

  service3:
    image: service3:latest
    depends_on:
      service2:
        condition: service_healthy
    volumes:
      - /dir/from:/dir/to

volumes:
  my-volume:
`

	parsed, err := Parse(bytes.NewReader([]byte(dockerComposeYaml)))
	require.NoError(t, err)
	require.Len(t, parsed.Services, 3)
	require.Len(t, parsed.Volumes, 1)

	// service1
	service1, ok := parsed.Services["service1"]
	require.True(t, ok, "service1 missing from parsed result")
	assert.Equal(
		t,
		[]VolumeMount{{
			Type:     VolumeTypeBind,
			Source:   "/dir/from",
			Target:   "/dir/to",
			ReadOnly: false,
		}},
		service1.VolumeMounts,
	)
	assert.Equal(t, map[string]string{"graph.category": "database"}, service1.Labels)

	// service2
	service2, ok := parsed.Services["service2"]
	require.True(t, ok, "service2 missing from parsed result")
	assert.Equal(
		t,
		[]VolumeMount{{
			Type:     VolumeTypeVolume,
			Source:   "my-volume",
			Target:   "/some/data",
			ReadOnly: true,
		}},
		service2.VolumeMounts,
	)
	assert.Equal(
		t,
		[]ServiceDependency{{
			On:        "service1",
			Condition: ConditionServiceStarted,
		}},
		service2.ServiceDependencies,
	)
	assert.Equal(t, map[string]string{"graph.category": "proxy"}, service2.Labels)

	// service3
	service3, ok := parsed.Services["service3"]
	require.True(t, ok, "service3 missing from parsed result")
	assert.Equal(
		t,
		[]VolumeMount{{
			Type:     VolumeTypeBind,
			Source:   "/dir/from",
			Target:   "/dir/to",
			ReadOnly: false,
		}},
		service3.VolumeMounts,
	)
	assert.Equal(
		t,
		[]ServiceDependency{{
			On:        "service2",
			Condition: ConditionServiceHealthy,
		}},
		service3.ServiceDependencies,
	)

	// volumes
	assert.Equal(t, []string{"my-volume"}, parsed.Volumes)
}
