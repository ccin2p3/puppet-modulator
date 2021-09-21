package gitutils

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gitlab.in2p3.fr/rferrand/go-system-utils/command"
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

	gfCmd := []string{"git", "commit", "-F", cFile.Name()}
	for _, fileToCommit := range files {
		gfCmd = append(gfCmd, fileToCommit)
	}

	log.Debugf("executing %s", gfCmd)

	argv := &command.ExecArgv{
		Command: gfCmd,
	}

	result, err := command.Execute(argv)
	if err != nil {
		return errors.Wrap(err, "commiting modifications")
	}

	log.Debugf("git commit output: %s", result.Stdout)

	return nil
}
