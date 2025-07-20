package compose

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCondition(t *testing.T) {
	tests := []struct {
		input       string
		expected    Condition
		expectedErr bool
	}{{
		input:       "",
		expected:    ConditionServiceStarted,
		expectedErr: false,
	}, {
		input:       "service_started",
		expected:    ConditionServiceStarted,
		expectedErr: false,
	}, {
		input:       "service_healthy",
		expected:    ConditionServiceHealthy,
		expectedErr: false,
	}, {
		input:       "service_completed_successfully",
		expected:    ConditionServiceCompletedSuccessfully,
		expectedErr: false,
	}, {
		input:       "invalid_condition",
		expected:    ConditionUnknown,
		expectedErr: true,
	}}

	for _, tt := range tests {
		actual, err := parseCondition(tt.input)
		if tt.expectedErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, actual)
		}
	}
}

func TestParseVolumeType(t *testing.T) {
	tests := []struct {
		input       string
		expected    VolumeType
		expectedErr bool
	}{{
		input:       "",
		expected:    VolumeTypeVolume,
		expectedErr: false,
	}, {
		input:       "volume",
		expected:    VolumeTypeVolume,
		expectedErr: false,
	}, {
		input:       "bind",
		expected:    VolumeTypeBind,
		expectedErr: false,
	}, {
		input:       "tmpfs",
		expected:    VolumeTypeTmpfs,
		expectedErr: false,
	}, {
		input:       "unknown",
		expected:    VolumeTypeUnknown,
		expectedErr: true,
	}}

	for _, tt := range tests {
		got, err := parseVolumeType(tt.input)
		if tt.expectedErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, got)
		}
	}
}
