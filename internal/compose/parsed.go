package compose

import (
	"fmt"
	"strings"
)

type Parsed struct {
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
	default:
		return ConditionUnknown, fmt.Errorf("invalid condition: %s", s)
	}
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
	default:
		return VolumeTypeUnknown, fmt.Errorf("invalid volume type: %s", s)
	}
}

//
// helpers for parsing
//

// rawVolume supports both short & long formats
type rawVolume struct {
	Short string
	Long  *rawVolumeLong
}

type rawVolumeLong struct {
	Type     string `yaml:"type,omitempty"`
	Source   string `yaml:"source,omitempty"`
	Target   string `yaml:"target,omitempty"`
	ReadOnly bool   `yaml:"read_only,omitempty"`
}

func (r *rawVolume) UnmarshalYAML(unmarshal func(any) error) error {
	var short string
	if err := unmarshal(&short); err == nil {
		r.Short = short
		return nil
	}

	var long rawVolumeLong
	if err := unmarshal(&long); err == nil {
		r.Long = &long
		return nil
	}

	return fmt.Errorf("invalid volume format")
}

// rawDependsOn supports both list and map formats
type rawDependsOn struct {
	Map  map[string]rawDependsOnCondition
	List []string
}

type rawDependsOnCondition struct {
	Condition string `yaml:"condition,omitempty"`
}

func (d *rawDependsOn) UnmarshalYAML(unmarshal func(any) error) error {
	var fromMap map[string]rawDependsOnCondition
	if err := unmarshal(&fromMap); err == nil {
		d.Map = fromMap
		return nil
	}

	var fromList []string
	if err := unmarshal(&fromList); err == nil {
		d.List = fromList
		return nil
	}

	return fmt.Errorf("invalid depends_on format")
}

func (s *Service) UnmarshalYAML(unmarshal func(any) error) error {
	var raw struct {
		Volumes   []rawVolume  `yaml:"volumes,omitempty"`
		DependsOn rawDependsOn `yaml:"depends_on,omitempty"`
	}
	if err := unmarshal(&raw); err != nil {
		return err
	}

	// Normalize depends_on
	for _, dependency := range raw.DependsOn.List {
		s.Dependencies = append(s.Dependencies, Dependency{
			On:        dependency,
			Condition: ConditionServiceStarted, // default
		})
	}

	for name, condition := range raw.DependsOn.Map {
		c, err := parseCondition(condition.Condition)
		if err != nil {
			return err
		}
		s.Dependencies = append(s.Dependencies, Dependency{
			On:        name,
			Condition: c,
		})
	}

	// Normalize volumes
	for _, volume := range raw.Volumes {
		switch {
		case volume.Short != "":
			parts := strings.Split(volume.Short, ":")

			if len(parts) < 2 {
				return fmt.Errorf("invalid volume format")
			}

			t := VolumeTypeVolume
			if strings.HasPrefix(parts[0], "/") || strings.HasPrefix(parts[0], ".") {
				t = VolumeTypeBind
			}

			s.Volume = append(s.Volume, Volume{
				Type:     t,
				Source:   parts[0],
				Target:   parts[1],
				ReadOnly: len(parts) == 3 && parts[2] == "ro",
			})

		case volume.Long != nil:
			t, err := parseVolumeType(volume.Long.Type)
			if err != nil {
				return err
			}
			s.Volume = append(s.Volume, Volume{
				Type:     t,
				Source:   volume.Long.Source,
				Target:   volume.Long.Target,
				ReadOnly: volume.Long.ReadOnly,
			})
		}
	}

	return nil
}
