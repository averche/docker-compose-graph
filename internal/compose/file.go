package compose

import (
	"fmt"
)

type File struct {
	Services map[string]Service
	Volumes  []string
}

type Service struct {
	VolumeMounts        []VolumeMount
	ServiceDependencies []ServiceDependency
	Labels              map[string]string
}

type ServiceDependency struct {
	On        string
	Condition Condition
}

type VolumeMount struct {
	Type     VolumeMountType
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

type VolumeMountType uint8

const (
	VolumeTypeUnknown VolumeMountType = iota
	VolumeTypeBind
	VolumeTypeVolume
	VolumeTypeTmpfs
)

func parseVolumeMountType(s string) (VolumeMountType, error) {
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
