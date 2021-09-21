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

var (
	MetadataHasherNewFunc = sha1.New
)

type MetadataJSON struct {
	m           map[string]interface{}
	OriginalSum []byte
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
