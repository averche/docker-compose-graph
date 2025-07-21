package compose

import (
	"cmp"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/goccy/go-yaml"
)

func ParseFile(path string) (_ File, errs error) {
	f, err := os.Open(path)
	if err != nil {
		return File{}, fmt.Errorf("could not open: %w", err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			errs = errors.Join(errs, fmt.Errorf("could not close: %w", err))
		}
	}()

	return Parse(f)
}

func Parse(r io.Reader) (File, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return File{}, fmt.Errorf("could not read contents: %w", err)
	}

	var parsed File

	if err := yaml.Unmarshal(b, &parsed); err != nil {
		return File{}, fmt.Errorf("could not unmarshal yaml contents: %w", err)
	}

	return parsed, nil
}

//
// helpers for parsing
//

// rawVolumeMount supports both short & long 'volume' formats
type rawVolumeMount struct {
	Short string
	Long  *rawVolumeMountLong
}

type rawVolumeMountLong struct {
	Type     string `yaml:"type,omitempty"`
	Source   string `yaml:"source,omitempty"`
	Target   string `yaml:"target,omitempty"`
	ReadOnly bool   `yaml:"read_only,omitempty"`
}

func (r *rawVolumeMount) UnmarshalYAML(unmarshal func(any) error) error {
	var short string
	if err := unmarshal(&short); err == nil {
		r.Short = short
		return nil
	}

	var long rawVolumeMountLong
	if err := unmarshal(&long); err == nil {
		r.Long = &long
		return nil
	}

	return fmt.Errorf("invalid volume format")
}

// rawDependsOn supports both list and map formats
type rawDependsOn struct {
	List []string
	Map  map[string]rawDependsOnCondition
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
		VolumeMounts []rawVolumeMount `yaml:"volumes,omitempty"`
		DependsOn    rawDependsOn     `yaml:"depends_on,omitempty"`
	}
	if err := unmarshal(&raw); err != nil {
		return err
	}

	// normalize depends_on
	for _, dependency := range raw.DependsOn.List {
		s.ServiceDependencies = append(s.ServiceDependencies, ServiceDependency{
			On:        dependency,
			Condition: ConditionServiceStarted, // default
		})
	}

	for name, condition := range raw.DependsOn.Map {
		c, err := parseCondition(condition.Condition)
		if err != nil {
			return err
		}
		s.ServiceDependencies = append(s.ServiceDependencies, ServiceDependency{
			On:        name,
			Condition: c,
		})
	}

	slices.SortFunc(s.ServiceDependencies, func(a, b ServiceDependency) int {
		return cmp.Compare(a.On, b.On)
	})

	// normalize volume mounts
	for _, volume := range raw.VolumeMounts {
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

			s.VolumeMounts = append(s.VolumeMounts, VolumeMount{
				Type:     t,
				Source:   parts[0],
				Target:   parts[1],
				ReadOnly: len(parts) == 3 && parts[2] == "ro",
			})

		case volume.Long != nil:
			t, err := parseVolumeMountType(volume.Long.Type)
			if err != nil {
				return err
			}
			s.VolumeMounts = append(s.VolumeMounts, VolumeMount{
				Type:     t,
				Source:   volume.Long.Source,
				Target:   volume.Long.Target,
				ReadOnly: volume.Long.ReadOnly,
			})
		}
	}

	return nil
}

func (s *File) UnmarshalYAML(unmarshal func(any) error) error {
	var raw struct {
		Services map[string]Service `yaml:"services"`
		Volumes  map[string]any     `yaml:"volumes,omitempty"`
	}
	if err := unmarshal(&raw); err != nil {
		return err
	}

	volumes := make([]string, 0, len(raw.Volumes))

	for volume := range raw.Volumes {
		volumes = append(volumes, volume)
	}

	slices.Sort(volumes)

	s.Services = raw.Services
	s.Volumes = volumes

	return nil
}
