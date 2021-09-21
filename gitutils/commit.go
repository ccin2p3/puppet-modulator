package gitutils

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/bitfield/script"
)

func GitCommitFile(msg string, files ...string) error {

	cFile, err := ioutil.TempFile("", "puppet-modulator-bump")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("fail to create temporary file")
	}
	defer os.Remove(cFile.Name())

	if _, err := fmt.Fprintf(cFile, msg); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("fail to write commit message to temporary file")
	}
	if err := cFile.Close(); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("fail to close commit message temporary file")
	}

	gCmd := fmt.Sprintf("git commit -F %s %s", cFile.Name(), strings.Join(files, " "))
	log.Debugf("executing %s", gCmd)

	p := script.Exec(gCmd)
	pOut, err := p.String()
	if err != nil {
		return errors.Wrap(err, "commiting modifications")
	}

	log.Debugf("git commit output: %s", pOut)

	return nil
}
