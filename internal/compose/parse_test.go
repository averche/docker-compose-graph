package compose

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFile(t *testing.T) {
	yamlData := `
version: "3.9"

services:

  nginx:
    image: nginx:latest
    depends_on:
      - postgres
    volumes:
      - nginx_data:/app:ro
      - type: bind
        source: /logs
        target: /var/log/nginx

  postgres:
    image: postgres:17.5
    depends_on:
      kafka:
        condition: service_started
    volumes:
      - db_data:/var/lib/postgresql/data

  kafka:
    image: kafka:latest
    depends_on:
      nginx:
        condition: service_healthy
    volumes:
      - type: volume
        source: kafka_data
        target: /data
`

	parsed, err := Parse(bytes.NewReader([]byte(yamlData)))
	require.NoError(t, err)
	require.Len(t, parsed.Services, 3)
	assert.Equal(t, "3.9", parsed.Version)

	nginx := parsed.Services["nginx"]
	require.Len(t, nginx.Volume, 2)
	assert.Equal(t, Volume{Type: VolumeTypeVolume, Source: "nginx_data", Target: "/app", ReadOnly: true}, nginx.Volume[0])
	assert.Equal(t, Volume{Type: VolumeTypeBind, Source: "/logs", Target: "/var/log/nginx", ReadOnly: false}, nginx.Volume[1])
	assert.Equal(t, []Dependency{{On: "postgres", Condition: ConditionServiceStarted}}, nginx.Dependencies)

	postgres := parsed.Services["postgres"]
	require.Len(t, postgres.Volume, 1)
	assert.Equal(t, Volume{Type: VolumeTypeVolume, Source: "db_data", Target: "/var/lib/postgresql/data", ReadOnly: false}, postgres.Volume[0])
	assert.Equal(t, []Dependency{{On: "kafka", Condition: ConditionServiceStarted}}, postgres.Dependencies)

	kafka := parsed.Services["kafka"]
	require.Len(t, kafka.Volume, 1)
	assert.Equal(t, Volume{Type: VolumeTypeVolume, Source: "kafka_data", Target: "/data", ReadOnly: false}, kafka.Volume[0])
	assert.Equal(t, []Dependency{{On: "nginx", Condition: ConditionServiceHealthy}}, kafka.Dependencies)
}
