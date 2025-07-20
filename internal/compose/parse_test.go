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
    volumes:
      - type:   bind
        source: /dir/from
        target: /dir/to

  service2:
    image: service2:latest
    depends_on:
      - service1
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
	assert.Equal(t, "3.9", parsed.Version)

	service1, ok := parsed.Services["service1"]
	require.True(t, ok, "service1 missing from parsed result")
	assert.Equal(
		t,
		[]Volume{{
			Type:     VolumeTypeBind,
			Source:   "/dir/from",
			Target:   "/dir/to",
			ReadOnly: false,
		}},
		service1.Volumes,
	)

	service2, ok := parsed.Services["service2"]
	require.True(t, ok, "service2 missing from parsed result")
	assert.Equal(
		t,
		[]Volume{{
			Type:     VolumeTypeVolume,
			Source:   "my-volume",
			Target:   "/some/data",
			ReadOnly: true,
		}},
		service2.Volumes,
	)
	assert.Equal(
		t,
		[]Dependency{{
			On:        "service1",
			Condition: ConditionServiceStarted,
		}},
		service2.Dependencies,
	)

	service3, ok := parsed.Services["service3"]
	require.True(t, ok, "service3 missing from parsed result")
	assert.Equal(
		t,
		[]Volume{{
			Type:     VolumeTypeBind,
			Source:   "/dir/from",
			Target:   "/dir/to",
			ReadOnly: false,
		}},
		service3.Volumes,
	)
	assert.Equal(
		t,
		[]Dependency{{
			On:        "service2",
			Condition: ConditionServiceHealthy,
		}},
		service3.Dependencies,
	)
}
