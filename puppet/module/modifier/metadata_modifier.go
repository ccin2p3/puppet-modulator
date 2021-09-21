package mmodifier

import (
	"bytes"
	"encoding/hex"
	"io"
	"os"

	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gitlab.in2p3.fr/cc-in2p3-puppet-master-tools/puppet-modulator/puppet/module"
)

type metadataModifier struct {
	modifyFunc      MetadataModifierFunc
	getReaderFunc   func() (io.ReadCloser, error)
	getWriterFunc   func() (io.WriteCloser, error)
	postRewriteFunc func() error
	opts            MetadataModifierOptions
}

type MetadataModifierOptions struct {
	JSONPretty bool
}

type MetadataModifier interface {
	Modify() error
	SetWriter(io.WriteCloser)
	SetPostRewriteFunc(func() error)
}

type MetadataModifierFunc func(md module.MetadataJSON) error

var (
	defaultGetReaderFunc = func() (io.ReadCloser, error) {
		r, err := os.Open(module.MetadataJSONFilename)
		if err != nil {
			return nil, errors.Wrapf(err, "opening file %q for reading", module.MetadataJSONFilename)
		}

		return r, nil
	}

	defaultGetWriterFunc = func() (io.WriteCloser, error) {
		w, err := os.OpenFile(module.MetadataJSONFilename, os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return nil, errors.Wrapf(err, "opening file %q for writing", module.MetadataJSONFilename)
		}

		return w, nil
	}

	DefaultMetadataModifierOptions = MetadataModifierOptions{
		JSONPretty: true,
	}
)

func NewMetadataModifierWithOptions(mfunc MetadataModifierFunc, opts MetadataModifierOptions) metadataModifier {
	return metadataModifier{
		modifyFunc:    mfunc,
		getWriterFunc: defaultGetWriterFunc,
		getReaderFunc: defaultGetReaderFunc,
		opts:          opts,
	}
}

func NewMetadataModifier(mfunc MetadataModifierFunc) metadataModifier {
	return NewMetadataModifierWithOptions(mfunc, DefaultMetadataModifierOptions)
}

func (m metadataModifier) readMetadata() (module.MetadataJSON, error) {
	mdReader, err := m.getReaderFunc()
	if err != nil {
		return module.MetadataJSON{}, errors.Wrap(err, "reading metadata")
	}
	defer mdReader.Close()

	return module.NewMetadataJSONFromReader(mdReader)
}

func (m *metadataModifier) SetWriter(w io.WriteCloser) {
	m.getWriterFunc = func() (io.WriteCloser, error) {
		return w, nil
	}
}

func (m *metadataModifier) SetPostRewriteFunc(cb func() error) {
	m.postRewriteFunc = cb
}

func (m metadataModifier) wouldRewriteModifyFile(md module.MetadataJSON) bool {
	hasher := module.MetadataHasherNewFunc()
	m.writeMetadataToWriter(hasher, md)
	newSum := hasher.Sum(nil)

	logrus.WithFields(logrus.Fields{
		"old-checksum": hex.EncodeToString(md.OriginalSum),
		"new-checksum": hex.EncodeToString(newSum),
	}).Debug("wouldRewriteModifyFile()")

	return bytes.Compare(md.OriginalSum, newSum) != 0
}

func (m metadataModifier) Modify() error {
	md, err := m.readMetadata()
	if err != nil {
		return errors.Wrap(err, "reading metadata")
	}

	// we're checking if a simple rewrite of the metadata file "as-this"
	// would still modify the content of the metadata file
	// This may due to Go map[] keys sort or other "cosmetic" human change
	if m.wouldRewriteModifyFile(md) && m.postRewriteFunc != nil {
		logrus.Debug("metadata file would be rewritten even with no modifications")

		// The callback is quite useful to allow caller to commit
		// "cosmetic" modifications without the real change
		// This is really useful to differenciate a cosmetic change
		// due to the fact that the metadata has been updated by a program
		// V.S the real change (like a module version change) that is
		// really meaningful
		if err := m.writeMetadata(md); err != nil {
			return errors.Wrap(err, "writing metadata file before modification")
		}

		if err := m.postRewriteFunc(); err != nil {
			return errors.Wrap(err, "executing post rewrite callback")
		}
	}

	m.modifyFunc(md)

	return m.writeMetadata(md)
}

func (m metadataModifier) writeMetadata(md module.MetadataJSON) error {
	w, err := m.getWriterFunc()
	if err != nil {
		return errors.Wrap(err, "getting metadata writer")
	}
	defer w.Close()

	if err := m.writeMetadataToWriter(w, md); err != nil {
		return errors.Wrap(err, "writing metadata")
	}

	return nil
}

func (m metadataModifier) writeMetadataToWriter(w io.Writer, md module.MetadataJSON) error {
	return md.WriteToWriter(w, m.opts.JSONPretty)
}

func metadataVersionModify(md module.MetadataJSON, vModifier func(*semver.Version) semver.Version) error {
	v, err := md.GetVersion()
	if err != nil {
		return err
	}

	newV := vModifier(v)
	md.Set(module.MetadataJSONModuleVersionKey, newV.String())
	return nil
}

func IncPatchVersionModifierFunc(md module.MetadataJSON) error {
	return metadataVersionModify(md, func(v *semver.Version) semver.Version {
		return v.IncPatch()
	})
}

func IncMinorVersionModifierFunc(md module.MetadataJSON) error {
	return metadataVersionModify(md, func(v *semver.Version) semver.Version {
		return v.IncMinor()
	})
}

func IncMajorVersionModifierFunc(md module.MetadataJSON) error {
	return metadataVersionModify(md, func(v *semver.Version) semver.Version {
		return v.IncMajor()
	})
}

func SetVersionModifierFunc(v string) MetadataModifierFunc {
	return func(md module.MetadataJSON) error {
		md.Set(module.MetadataJSONModuleVersionKey, v)
		return nil
	}
}
