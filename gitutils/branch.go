package gitutils

import (
	"strings"

	"github.com/bitfield/script"
)

func GetCurrentBranch() (string, error) {
	p := script.Exec("git branch --show-current")
	pOut, err := p.String()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(pOut), nil
}
