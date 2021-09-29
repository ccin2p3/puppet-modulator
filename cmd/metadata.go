/*
Copyright © 2021 IN2P3 Computing Centre, IN2P3, CNRS

Contributor(s): Remi Ferrand <remi.ferrand_at_cc.in2p3.fr>, 2021

This software is governed by the CeCILL-B license under French law and
abiding by the rules of distribution of free software.  You can  use,
modify and/ or redistribute the software under the terms of the CeCILL-B
license as circulated by CEA, CNRS and INRIA at the following URL
"http://www.cecill.info".

As a counterpart to the access to the source code and  rights to copy,
modify and redistribute granted by the license, users are provided only
with a limited warranty  and the software's author,  the holder of the
economic rights,  and the successive licensors  have only  limited
liability.

In this respect, the user's attention is drawn to the risks associated
with loading,  using,  modifying and/or developing or reproducing the
software by the user in light of its specific status of free software,
that may mean  that it is complicated to manipulate,  and  that  also
therefore means  that it is reserved for developers  and  experienced
professionals having in-depth computer knowledge. Users are therefore
encouraged to load and test the software's suitability as regards their
requirements in conditions enabling the security of their systems and/or
data to be ensured and,  more generally, to use and operate it in the
same conditions as regards security.

The fact that you are presently reading this means that you have had
knowledge of the CeCILL-B license and that you accept its terms.

*/
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/Masterminds/semver/v3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.in2p3.fr/cc-in2p3-puppet-master-tools/puppet-modulator/gitutils"
	"gitlab.in2p3.fr/cc-in2p3-puppet-master-tools/puppet-modulator/ioutils"
	"gitlab.in2p3.fr/cc-in2p3-puppet-master-tools/puppet-modulator/puppet/module"
	mgit "gitlab.in2p3.fr/cc-in2p3-puppet-master-tools/puppet-modulator/puppet/module/git"
	mmodifier "gitlab.in2p3.fr/cc-in2p3-puppet-master-tools/puppet-modulator/puppet/module/modifier"
)

const (
	metadataKeysPreSortAndCommitPolicyName   = "pre-commit"
	metadataKeysNoPreSortAndCommitPolicyName = "no-pre-commit"
)

var (
	metadataCmd = &cobra.Command{
		Use:   "metadata",
		Short: "Manipulate module metadata.json file",
	}

	metadataSetVersionCmd = &cobra.Command{
		Use:   "set-version VERSION",
		Short: "Set exact module version",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			newCobraModuleMetadataModifierFuncAdapter(mmodifier.SetVersionModifierFunc(args[0]))(cmd, args)
		},
	}

	metadataBumpCmd = &cobra.Command{
		Use: "bump",
	}

	metadataBumpPatchCmd = &cobra.Command{
		Use:   "patch",
		Short: "Bump module version to the next patch",
		Run:   newCobraModuleMetadataModifierFuncAdapter(mmodifier.IncPatchVersionModifierFunc),
	}

	metadataBumpMinorCmd = &cobra.Command{
		Use:   "minor",
		Short: "Bump module version to the next minor",
		Run:   newCobraModuleMetadataModifierFuncAdapter(mmodifier.IncMinorVersionModifierFunc),
	}

	metadataBumpMajorCmd = &cobra.Command{
		Use:   "major",
		Short: "Bump module version to the next major",
		Run:   newCobraModuleMetadataModifierFuncAdapter(mmodifier.IncMajorVersionModifierFunc),
	}

	metadataVersionCmd = &cobra.Command{
		Use:   "version",
		Short: "Get module version (current or next)",
	}

	metadataVersionGetCurrentCmd = &cobra.Command{
		Use:   "get",
		Short: "Get current module version",
		Run:   metadataGetVersionCLIRun(nil),
	}

	metadataVersionGetNextCmd = &cobra.Command{
		Use:   "get-next",
		Short: "Get next module version",
		Run:   metadataGetVersionCLIRun(nil),
	}

	metadataVersionGetNextPatchCmd = &cobra.Command{
		Use:   "patch",
		Short: "Get next patch module version",
		Run: metadataGetVersionCLIRun(func(v *semver.Version) semver.Version {
			return v.IncPatch()
		}),
	}

	metadataVersionGetNextMinorCmd = &cobra.Command{
		Use:   "minor",
		Short: "Get next minor module version",
		Run: metadataGetVersionCLIRun(func(v *semver.Version) semver.Version {
			return v.IncMinor()
		}),
	}

	metadataVersionGetNextMajorCmd = &cobra.Command{
		Use:   "major",
		Short: "Get next major module version",
		Run: metadataGetVersionCLIRun(func(v *semver.Version) semver.Version {
			return v.IncMajor()
		}),
	}
)

func metadataGetVersionCLIRun(vModifier func(*semver.Version) semver.Version) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, s []string) {
		v, err := mgit.GetModuleVersionAtRef("HEAD")
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
				"ref":   "HEAD",
			}).Fatal("fail to get module version")
		}

		if vModifier != nil {
			nv := vModifier(v)
			v = &nv
		}

		fmt.Printf("%s\n", v.String())
	}
}

type metadataBumpCLIOptions struct {
	destFile            string
	commitModifications bool
	commitPolicy        string
	commitMessage       string
	modifierFunc        mmodifier.MetadataModifierFunc
}

const (
	metadataCLIPreRealCommitCommitMsg = "[meta] metadata.json automated modifications (pre-real modifications)"
)

func metadataBumpCLIRun(opts metadataBumpCLIOptions) {
	mdModifier := mmodifier.NewMetadataModifier(opts.modifierFunc)

	destFile := opts.destFile
	if destFile != "" {
		var writer io.WriteCloser
		if destFile == "-" {
			// write to STDOUT
			writer = ioutils.NopWriteCloser(os.Stdout)
		} else {
			mdFileFd, err := os.OpenFile(destFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Fatalf("fail to open metadata file %q for writing", destFile)
			}

			writer = mdFileFd
		}

		mdModifier.SetWriter(writer)
	}

	commit := opts.commitModifications
	if commit {
		if ksp := opts.commitPolicy; ksp == metadataKeysPreSortAndCommitPolicyName {
			mdModifier.SetPostRewriteFunc(func() error {
				commitMsg := fmt.Sprintf(metadataCLIPreRealCommitCommitMsg)
				if err := gitutils.GitCommitFile(commitMsg, module.MetadataJSONFilename); err != nil {
					log.WithFields(log.Fields{
						"error": err,
					}).Fatal("fail to commit modifications")
				}

				log.Debug("pre-modifications modifications only commited")

				return nil
			})
		}
	}

	if err := mdModifier.Modify(); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("fail to modify module metadata")
	}

	if commit {
		commitMsg := opts.commitMessage

		if err := gitutils.GitCommitFile(commitMsg, module.MetadataJSONFilename); err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Fatal("fail to commit modifications")
		}
	}
}

func newCobraModuleMetadataModifierFuncAdapter(mdModifierFunc mmodifier.MetadataModifierFunc) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		destFile, _ := cmd.Flags().GetString("output")
		commitPolicy, _ := cmd.Flags().GetString("keys-sort-commit-policy")
		commitMsg, _ := cmd.Flags().GetString("git-commit-msg")
		commit, _ := cmd.Flags().GetBool("git-commit")

		cliOpts := metadataBumpCLIOptions{
			destFile:            destFile,
			modifierFunc:        mdModifierFunc,
			commitModifications: commit,
			commitPolicy:        commitPolicy,
			commitMessage:       commitMsg,
		}

		metadataBumpCLIRun(cliOpts)
	}
}

func init() {
	rootCmd.AddCommand(metadataCmd)
	metadataCmd.AddCommand(metadataBumpCmd)
	metadataCmd.AddCommand(metadataSetVersionCmd)
	metadataCmd.AddCommand(metadataVersionCmd)

	metadataVersionCmd.AddCommand(metadataVersionGetCurrentCmd)
	metadataVersionCmd.AddCommand(metadataVersionGetNextCmd)

	metadataVersionGetNextCmd.AddCommand(metadataVersionGetNextPatchCmd)
	metadataVersionGetNextCmd.AddCommand(metadataVersionGetNextMinorCmd)
	metadataVersionGetNextCmd.AddCommand(metadataVersionGetNextMajorCmd)

	metadataBumpCmd.AddCommand(metadataBumpPatchCmd)
	metadataBumpCmd.AddCommand(metadataBumpMinorCmd)
	metadataBumpCmd.AddCommand(metadataBumpMajorCmd)

	metadataCmd.PersistentFlags().StringP("keys-sort-commit-policy", "p", metadataKeysPreSortAndCommitPolicyName,
		fmt.Sprintf("policy related to metadata keys sort commit. If %s is used, then a dedicated commit will be created dedicated to metadata keys sorting. If %s is used, metadata keys sorting will still occurs, but no dedicated commit will be created", metadataKeysPreSortAndCommitPolicyName, metadataKeysNoPreSortAndCommitPolicyName))

	for _, cmd := range []*cobra.Command{metadataCmd} {
		cmd.PersistentFlags().StringP("output", "o", "", "Where to write metadata to. Defaults to modify metadata in-place")
	}

	metadataCmd.PersistentFlags().BoolP("git-commit", "g", false, "Commit changes to git")
	metadataCmd.PersistentFlags().StringP("git-commit-msg", "m", metadataCLIBumpVersionDefaultCommitMsg, "Git commit message")
}

const (
	metadataCLIBumpVersionDefaultCommitMsg = "[meta] Bump version"
)
