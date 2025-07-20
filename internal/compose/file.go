package compose

import (
	"fmt"
)

type File struct {
	Version  string             `yaml:"version,omitempty"`
	Services map[string]Service `yaml:"services"`
}

type Service struct {
	Volume       []Volume
	Dependencies []Dependency
}

type Dependency struct {
	On        string
	Condition Condition
}

type Volume struct {
	Type     VolumeType
	Source   string
	Target   string
	ReadOnly bool
}

//
// enums
//

type Condition uint8

const (
	ConditionUnknown Condition = iota
	ConditionServiceStarted
	ConditionServiceHealthy
	ConditionServiceCompletedSuccessfully
)

func parseCondition(s string) (Condition, error) {
	switch s {
	case "", "service_started":
		return ConditionServiceStarted, nil

	case "service_healthy":
		return ConditionServiceHealthy, nil

	case "service_completed_successfully":
		return ConditionServiceCompletedSuccessfully, nil
	}

	return ConditionUnknown, fmt.Errorf("invalid condition: %s", s)
}

type VolumeType uint8

const (
	VolumeTypeUnknown VolumeType = iota
	VolumeTypeBind
	VolumeTypeVolume
	VolumeTypeTmpfs
)

func parseVolumeType(s string) (VolumeType, error) {
	switch s {
	case "", "volume":
		return VolumeTypeVolume, nil

	case "bind":
		return VolumeTypeBind, nil

	case "tmpfs":
		return VolumeTypeTmpfs, nil
	}

	return VolumeTypeUnknown, fmt.Errorf("invalid volume type: %s", s)
}
