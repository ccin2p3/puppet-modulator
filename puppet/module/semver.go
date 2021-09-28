package module

import "github.com/Masterminds/semver/v3"

func ValidateSemverString(s string) error {
	_, err := semver.NewVersion(s)
	if err != nil {
		return err
	}
	return nil
}
