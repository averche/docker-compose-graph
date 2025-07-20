package compose

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
)

// ParseConfigurations parses the given slice of files
func ParseMultiple(paths []string) ([]File, error) {
	var parsed []File

	for _, path := range paths {
		p, err := ParseFile(path)
		if err != nil {
			return nil, fmt.Errorf("could not parse %q: %w", path, err)
		}
		parsed = append(parsed, p)
	}

	return parsed, nil
}

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

// rawVolume supports both short & long 'volume' formats
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
		Volumes   []rawVolume  `yaml:"volumes,omitempty"`
		DependsOn rawDependsOn `yaml:"depends_on,omitempty"`
	}
	if err := unmarshal(&raw); err != nil {
		return err
	}

	// normalize depends_on
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

	// normalize volumes
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

			s.Volumes = append(s.Volumes, Volume{
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
			s.Volumes = append(s.Volumes, Volume{
				Type:     t,
				Source:   volume.Long.Source,
				Target:   volume.Long.Target,
				ReadOnly: volume.Long.ReadOnly,
			})
		}
	}

	return nil
}
