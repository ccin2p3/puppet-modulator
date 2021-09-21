package gfutils

import (
	"os"

	"github.com/sirupsen/logrus"
	"gitlab.in2p3.fr/rferrand/go-system-utils/command"
)

type ReleaseStartOptions struct {
	Base            string
	GFOutputHandler func(string)
}

var (
	DefaultReleaseStartOptions = ReleaseStartOptions{
		Base: "",
	}
)

func ReleaseStart(version string, opts *ReleaseStartOptions) error {
	if opts == nil {
		opts = &DefaultReleaseStartOptions
	}

	gfCmd := []string{"git", "flow", "release", "start", version}
	if opts.Base != "" {
		gfCmd = append(gfCmd, opts.Base)
	}

	argv := &command.ExecArgv{
		Command: gfCmd,
	}

	result, err := command.Execute(argv)
	if err != nil {
		return err
	}

	if opts.GFOutputHandler != nil {
		opts.GFOutputHandler(result.Stdout)
	}

	return nil
}

type ReleaseFinishOptions struct {
	HotfixOrReleaseBaseOptions
	GFOutputHandler func(string)
}

var (
	DefaultReleaseFinishOptions = ReleaseFinishOptions{
		HotfixOrReleaseBaseOptions: HotfixOrReleaseBaseOptions{
			Push:           false,
			NoEditorPrompt: false,
		},
	}
)

func ReleaseFinish(opts *ReleaseFinishOptions) error {
	if opts == nil {
		opts = &DefaultReleaseFinishOptions
	}

	gfCmd := []string{"git", "flow", "release", "finish"}

	if opts.Push {
		gfCmd = append(gfCmd, "-p")
	}

	if opts.Message != nil {
		gfCmd = append(gfCmd, "-m", *opts.Message)
	}

	if opts.MessageFile != nil {
		gfCmd = append(gfCmd, "-f", *opts.MessageFile)
	}

	var cmdEnv []string
	if opts.NoEditorPrompt {
		cmdEnv = append(cmdEnv, "GIT_MERGE_AUTOEDIT=no")
	}

	logrus.WithFields(logrus.Fields{
		"command":     gfCmd,
		"environment": cmdEnv,
	}).Debug("git-flow raw command")

	argv := &command.ExecArgv{
		Command: gfCmd,
		Stdin:   os.Stdin,
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
		Environ: cmdEnv,
	}

	_, err := command.Execute(argv)
	if err != nil {
		return err
	}

	return nil
}
