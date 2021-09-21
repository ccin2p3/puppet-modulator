package gfutils

import (
	"os"

	"github.com/sirupsen/logrus"
	"gitlab.in2p3.fr/rferrand/go-system-utils/command"
)

type HotfixStartOptions struct {
	Base            string
	GFOutputHandler func(string)
}

var (
	DefaultHotfixStartOptions = HotfixStartOptions{
		Base: "",
	}
)

func HotfixStart(version string, opts *HotfixStartOptions) error {
	if opts == nil {
		opts = &DefaultHotfixStartOptions
	}

	gfCmd := []string{"git", "flow", "hotfix", "start", version}
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

type HotfixFinishOptions struct {
	HotfixOrReleaseBaseOptions
	GFOutputHandler func(string)
}

var (
	DefaultHotfixFinishOptions = HotfixFinishOptions{
		HotfixOrReleaseBaseOptions: HotfixOrReleaseBaseOptions{
			Push:           false,
			NoEditorPrompt: false,
		},
	}
)

func HotfixFinish(opts *HotfixFinishOptions) error {
	if opts == nil {
		opts = &DefaultHotfixFinishOptions
	}

	gfCmd := []string{"git", "flow", "hotfix", "finish"}

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
