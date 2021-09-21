package mmodifier

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"gitlab.in2p3.fr/cc-in2p3-puppet-master-tools/puppet-modulator/ioutils"
	"gitlab.in2p3.fr/cc-in2p3-puppet-master-tools/puppet-modulator/puppet/module"
)

var (
	simpleMetadataJSON1 string = `
{
	"name": "ccin2p3-site_nrpe",
	"version": "3.0.1",
	"author": "ccin2p3",
	"license": "CeCILL B",
	"dependencies": [
		{
			"name": "puppetlabs/stdlib",
			"version_requirement": ">= 4.11.0 < 5.0.0"
		}
	],
	"requirements": [
		{
			"name": "puppet",
			"version_requirement": ">= 7.0.0 < 8.0.0"
		}
	]
}`
)

func testHelperMdInspectEnsureMetadataVersion(expected string) func(*testing.T, module.MetadataJSON) {
	return func(t *testing.T, md module.MetadataJSON) {
		v, err := md.GetVersion()
		if err != nil {
			t.Fatalf("getting metadata version: %v", err)
		}

		if expected != v.String() {
			t.Errorf("got %q. %q was expected", v.String(), expected)
		}
	}
}

func TestMetadataModifier(t *testing.T) {
	type testCase struct {
		name                string
		modifyFunc          MetadataModifierFunc
		metadataString      string
		metadataInspectFunc func(*testing.T, module.MetadataJSON)
		wantErr             bool
		errInspectFunc      func(*testing.T, error)
	}

	for _, tc := range []testCase{
		{
			name:                "set version",
			modifyFunc:          SetVersionModifierFunc("1.2.3"),
			metadataString:      simpleMetadataJSON1,
			metadataInspectFunc: testHelperMdInspectEnsureMetadataVersion("1.2.3"),
		},
		{
			name:                "increment patch version",
			modifyFunc:          IncPatchVersionModifierFunc,
			metadataString:      simpleMetadataJSON1,
			metadataInspectFunc: testHelperMdInspectEnsureMetadataVersion("3.0.2"),
		},
		{
			name:                "increment minor version",
			modifyFunc:          IncMinorVersionModifierFunc,
			metadataString:      simpleMetadataJSON1,
			metadataInspectFunc: testHelperMdInspectEnsureMetadataVersion("3.1.0"),
		},
		{
			name:                "increment major version",
			modifyFunc:          IncMajorVersionModifierFunc,
			metadataString:      simpleMetadataJSON1,
			metadataInspectFunc: testHelperMdInspectEnsureMetadataVersion("4.0.0"),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			modifier := NewMetadataModifier(tc.modifyFunc)

			modifier.getReaderFunc = func() (io.ReadCloser, error) {
				return ioutil.NopCloser(strings.NewReader(tc.metadataString)), nil
			}

			modifiedMd := &bytes.Buffer{}
			modifier.getWriterFunc = func() (io.WriteCloser, error) {
				return ioutils.NopWriteCloser(modifiedMd), nil
			}

			if err := modifier.Modify(); err != nil {
				if tc.wantErr {
					if tc.errInspectFunc != nil {
						tc.errInspectFunc(t, err)
						return
					}
					return
				}

				t.Fatalf("Modify() failed: %v", err)
			} else {
				if tc.wantErr {
					t.Fatal("no error raised")
				}
			}

			md, err := module.NewMetadataJSONFromReader(modifiedMd)
			if err != nil {
				t.Fatalf("NewMetadataJSONFromReader() failed: %v", err)
			}

			tc.metadataInspectFunc(t, md)
		})

	}

}
