package gitutils

import (
	"io"
	"strings"

	"gitlab.in2p3.fr/rferrand/go-system-utils/command"
)

func GitShowFileAtRef(ref, filename string) (io.Reader, error) {
	cmd := []string{"git", "show", ref + ":" + filename}
	argv := &command.ExecArgv{
		Command: cmd,
	}

	result, err := command.Execute(argv)
	if err != nil {
		return nil, err
	}

	return strings.NewReader(result.Stdout), nil
}
