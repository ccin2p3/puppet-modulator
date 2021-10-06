package module

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"
)

const (
	MetadataJSONModuleVersionKey = "version"
	MetadataJSONFilename         = "metadata.json"
)

const (
	MetadataJSONModuleRequirementsKey                        = "requirements"
	MetadataJSONModuleRequirementsPuppetVersionConstraintKey = "puppet"
)

var (
	MetadataHasherNewFunc = sha1.New
)

type MetadataJSON struct {
	m           map[string]interface{}
	OriginalSum []byte
}

type metadataJSONVersionConstraint struct {
	Name       string `json:"name"`
	Constraint string `json:"version_requirement,omitempty"`
}

func NewMetadataJSONFromFilename(filename string) (MetadataJSON, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return MetadataJSON{}, errors.Wrapf(err, "opening file %s for reading", filename)
	}
	defer fd.Close()

	return NewMetadataJSONFromReader(fd)
}

func NewMetadataJSONFromReader(r io.Reader) (MetadataJSON, error) {
	m := MetadataJSON{
		m: make(map[string]interface{}),
	}

	hasher := MetadataHasherNewFunc()
	teeReader := io.TeeReader(r, hasher)

	if err := json.NewDecoder(teeReader).Decode(&m.m); err != nil {
		return m, err
	}

	m.OriginalSum = hasher.Sum(nil)

	return m, nil
}

type MetadataJSONNoSuchKeyErr string

func (m MetadataJSONNoSuchKeyErr) Error() string {
	return fmt.Sprintf("no such key %q in metadata.json", string(m))
}

type MetadataJSONInvalidKeyDataTypeErr struct {
	Key              string
	GotDataType      string
	ExpectedDataType string
}

func (m MetadataJSONInvalidKeyDataTypeErr) Error() string {
	return fmt.Sprintf("key %q is of type %q while %q was expected", m.Key, m.GotDataType, m.ExpectedDataType)
}

func (m MetadataJSON) getVersionConstraintSliceValue(key string) ([]metadataJSONVersionConstraint, error) {
	v, ok := m.m[key]
	if !ok {
		return []metadataJSONVersionConstraint{}, MetadataJSONNoSuchKeyErr(key)
	}

	b, err := json.Marshal(v)
	if err != nil {
		return []metadataJSONVersionConstraint{}, errors.Wrapf(err, "JSON marshaling key %s", key)
	}

	var constraints []metadataJSONVersionConstraint
	if err := json.Unmarshal(b, &constraints); err != nil {
		return []metadataJSONVersionConstraint{}, errors.Wrapf(err, "JSON unmarshaling key %s", key)
	}

	return constraints, nil
}

func (m MetadataJSON) getStringValue(key string) (string, error) {
	v, ok := m.m[key]
	if !ok {
		return "", MetadataJSONNoSuchKeyErr(key)
	}

	vs, ok := v.(string)
	if !ok {
		return "", MetadataJSONInvalidKeyDataTypeErr{
			Key:              key,
			GotDataType:      reflect.ValueOf(v).Kind().String(),
			ExpectedDataType: "string",
		}
	}

	return vs, nil
}

func (m MetadataJSON) GetVersion() (*semver.Version, error) {
	v, err := m.getStringValue(MetadataJSONModuleVersionKey)
	if err != nil {
		return nil, err
	}

	sv, err := semver.NewVersion(v)
	if err != nil {
		return sv, errors.Wrapf(err, "parsing semver %s", err)
	}

	return sv, nil
}

func (m MetadataJSON) GetPuppetVersionRequirement() (*semver.Constraints, error) {
	requirements, err := m.getVersionConstraintSliceValue(MetadataJSONModuleRequirementsKey)
	if err != nil {
		return nil, errors.Wrapf(err, "parsing current version constraints at key %s", MetadataJSONModuleRequirementsKey)
	}

	if len(requirements) == 0 {
		return nil, fmt.Errorf("empty requirements")
	}

	for _, requirement := range requirements {
		if requirement.Name == MetadataJSONModuleRequirementsPuppetVersionConstraintKey {
			c, err := semver.NewConstraint(requirement.Constraint)
			if err != nil {
				return nil, errors.Wrap(err, "parsing puppet version constraint")
			}
			return c, nil
		}
	}

	return nil, fmt.Errorf("no puppet requirement specified")
}

func (m *MetadataJSON) SetPuppetVersionRequirement(v string) error {
	requirements, err := m.getVersionConstraintSliceValue(MetadataJSONModuleRequirementsKey)
	if err != nil {
		if _, ok := err.(MetadataJSONNoSuchKeyErr); !ok {
			return errors.Wrapf(err, "parsing current version constraints at key %s", MetadataJSONModuleRequirementsKey)
		}
		requirements = []metadataJSONVersionConstraint{}
	}

	var requirementFound bool
	for i := 0; i < len(requirements); i++ {
		if requirements[i].Name == MetadataJSONModuleRequirementsPuppetVersionConstraintKey {
			requirementFound = true
			requirements[i].Constraint = v
		}
	}

	if len(requirements) == 0 || !requirementFound {
		requirements = append(requirements, metadataJSONVersionConstraint{
			Name:       MetadataJSONModuleRequirementsPuppetVersionConstraintKey,
			Constraint: v,
		})
	}

	m.m[MetadataJSONModuleRequirementsKey] = requirements

	return nil
}

func (m *MetadataJSON) Set(key string, value interface{}) {
	m.m[key] = value
}

func (m MetadataJSON) WriteToWriter(w io.Writer, pretty bool) error {
	enc := json.NewEncoder(w)
	if pretty {
		enc.SetIndent("", "  ")
		enc.SetEscapeHTML(false)
	}

	return enc.Encode(m.m)
}

func (m MetadataJSON) Write(pretty bool) error {
	w, err := os.OpenFile(MetadataJSONFilename, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return errors.Wrapf(err, "opening file %q for writing", MetadataJSONFilename)
	}

	defer w.Close()
	return m.WriteToWriter(w, pretty)
}
