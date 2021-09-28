package git

import (
	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"
	"gitlab.in2p3.fr/cc-in2p3-puppet-master-tools/puppet-modulator/gitutils"
	"gitlab.in2p3.fr/cc-in2p3-puppet-master-tools/puppet-modulator/puppet/module"
)

func GetModuleVersionAtRef(ref string) (*semver.Version, error) {
	metadataAtRef, err := gitutils.GitShowFileAtRef(ref, module.MetadataJSONFilename)
	if err != nil {
		return nil, errors.Wrapf(err, "getting %s file at reference %s", module.MetadataJSONFilename, ref)
	}

	metadata, err := module.NewMetadataJSONFromReader(metadataAtRef)
	if err != nil {
		return nil, errors.Wrapf(err, "parsing %s file", module.MetadataJSONFilename)
	}

	sVer, err := metadata.GetVersion()
	if err != nil {
		return nil, errors.Wrap(err, "extracting module version from metadata")
	}

	return sVer, nil
}
