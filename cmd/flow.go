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
	"os"
	"strings"

	"github.com/Masterminds/semver/v3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.in2p3.fr/cc-in2p3-puppet-master-tools/puppet-modulator/cmd/summaries"
	"gitlab.in2p3.fr/cc-in2p3-puppet-master-tools/puppet-modulator/flagutils"
	"gitlab.in2p3.fr/cc-in2p3-puppet-master-tools/puppet-modulator/gfutils"
	"gitlab.in2p3.fr/cc-in2p3-puppet-master-tools/puppet-modulator/gitutils"
	"gitlab.in2p3.fr/cc-in2p3-puppet-master-tools/puppet-modulator/puppet/module"
	mgit "gitlab.in2p3.fr/cc-in2p3-puppet-master-tools/puppet-modulator/puppet/module/git"
	mmodifier "gitlab.in2p3.fr/cc-in2p3-puppet-master-tools/puppet-modulator/puppet/module/modifier"
)

type cobraGitFlowFuncAdapterOptions struct {
}

// gflowCmd represents the flow command
var (
	gflowCmd = &cobra.Command{
		Use:   "flow hotfix|release start|finish [version] [base-ref]",
		Short: "A git-flow high-level wrapper for hotfixes and releases",
	}

	gflowHotfixCmd = &cobra.Command{
		Use:   "hotfix start|finish [version] [base-ref]",
		Short: "A git-flow high-level wrapper for hotfixes",
	}

	gflowReleaseCmd = &cobra.Command{
		Use:   "release start|finish [version] [base-ref]",
		Short: "A git-flow high-level wrapper for releases",
	}

	gflowHotfixStartCmd = &cobra.Command{
		Use:   "start [version] [base-ref]",
		Short: "A git-flow high-level wrapper to start hotfixes",
		Args:  cobra.MaximumNArgs(2),
		Run:   newCobraGFlowVersionBaseRefHandlerAdapter(gflowCLIStartHotfix),
	}

	gflowHotfixFinishCmd = &cobra.Command{
		Use:   "finish",
		Short: "A git-flow high-level wrapper to finish hotfixes",
		Args:  cobra.NoArgs,
		Run:   cobraGFlowFinishHotfixOrReleaseAdapter,
	}

	gflowReleaseStartCmd = &cobra.Command{
		Use:   "start [version] [base-ref]",
		Short: "A git-flow high-level wrapper to start releases",
		Args:  cobra.MaximumNArgs(2),
		Run:   newCobraGFlowVersionBaseRefHandlerAdapter(gflowCLIStartRelease),
	}

	gflowReleaseFinishCmd = &cobra.Command{
		Use:   "finish",
		Short: "A git-flow high-level wrapper to finish releases",
		Args:  cobra.NoArgs,
		Run:   cobraGFlowFinishHotfixOrReleaseAdapter,
	}
)

func newCobraGFlowVersionBaseRefHandlerAdapter(next func(string, string)) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		var version string
		var baseRef string

		if len(args) >= 1 {
			version = args[0]
		}
		if len(args) == 2 {
			baseRef = args[1]
		}

		log.WithFields(log.Fields{
			"version":  version,
			"baseRef":  baseRef,
			"cmd-name": cmd.Parent() == gflowHotfixCmd,
		}).Debug("CLI flags")

		if version == "" || version == "?" {
			var versModifierFunc func(v *semver.Version) semver.Version
			if cmd.Parent() == gflowHotfixCmd {
				if baseRef == "" {
					baseRef = "master"
				}
				versModifierFunc = func(v *semver.Version) semver.Version {
					return v.IncPatch()
				}
			} else {
				// this is a release
				if baseRef == "" {
					baseRef = "develop"
				}
				versModifierFunc = func(v *semver.Version) semver.Version {
					return v.IncMinor()
				}
			}

			sVersion, err := mgit.GetModuleVersionAtRef(baseRef)
			if err != nil {
				log.WithFields(log.Fields{
					"baseRef": baseRef,
					"error":   err,
				}).Fatal("fail to get module version in reference branch")
			}

			version = versModifierFunc(sVersion).String()

			log.WithFields(log.Fields{
				"version": version,
				"baseRef": baseRef,
			}).Debug("version discovered")

		} else {
			if err := module.ValidateSemverString(version); err != nil {
				log.WithFields(log.Fields{
					"version": version,
					"error":   err,
				}).Fatal("invalid version specified")
			}
		}

		next(version, baseRef)
	}
}

func gflowOutputLogHandler(gfOutput string) {
	log.Debugf("git-flow summary: %s", gfOutput)
}

func gflowCLIStartHotfix(version string, baseRef string) {
	opts := &gfutils.HotfixStartOptions{
		GFOutputHandler: gflowOutputLogHandler,
		Base:            baseRef,
	}
	if err := gfutils.HotfixStart(version, opts); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("fail to start hotfix")
	}

	// immediatly bump the version
	metadataCLIRun(metadataCLIOptions{
		commitModifications: true,
		commitPolicy:        metadataKeysPreSortAndCommitPolicyName,
		commitMessage:       metadataCLIBumpVersionDefaultCommitMsg,
		modifierFunc:        mmodifier.SetVersionModifierFunc(version),
	})

	rContext := summaries.RenderHotfixStartRendererContext{
		Version: version,
		BaseRef: baseRef,
	}
	if err := summaries.RenderHotfixStartSummary(os.Stdout, rContext); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("fail to print out summary")
	}
}

// TODO: factorize with gflowCLIStartHotfix
func gflowCLIStartRelease(version string, baseRef string) {
	opts := &gfutils.ReleaseStartOptions{
		GFOutputHandler: gflowOutputLogHandler,
		Base:            baseRef,
	}
	if err := gfutils.ReleaseStart(version, opts); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("fail to start release")
	}

	// immediatly bump the version
	metadataCLIRun(metadataCLIOptions{
		commitModifications: true,
		commitPolicy:        metadataKeysPreSortAndCommitPolicyName,
		commitMessage:       metadataCLIBumpVersionDefaultCommitMsg,
		modifierFunc:        mmodifier.SetVersionModifierFunc(version),
	})

	rContext := summaries.RenderReleaseStartRendererContext{
		Version: version,
		BaseRef: baseRef,
	}
	if err := summaries.RenderReleaseStartSummary(os.Stdout, rContext); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("fail to print out summary")
	}
}

func cobraGFlowFinishHotfixOrReleaseAdapter(cmd *cobra.Command, args []string) {

	baseOptions := gfutils.HotfixOrReleaseBaseOptions{}
	baseOptions.Push, _ = cmd.Flags().GetBool("push")
	baseOptions.NoEditorPrompt, _ = cmd.Flags().GetBool("no-prompt")

	if flagutils.HasFlag(cmd.Flags(), "message") {
		bMessage, _ := cmd.Flags().GetString("message")
		baseOptions.Message = &bMessage
	}
	if flagutils.HasFlag(cmd.Flags(), "messagefile") {
		bMessageFile, _ := cmd.Flags().GetString("messagefile")
		baseOptions.MessageFile = &bMessageFile
	}

	brType := gflowGuessBranchType()
	switch brType {
	case "hotfix":
		err := gfutils.HotfixFinish(&gfutils.HotfixFinishOptions{
			GFOutputHandler:            gflowOutputLogHandler,
			HotfixOrReleaseBaseOptions: baseOptions,
		})
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Fatal("fail to finish hotfix")
		}

	case "release":
		err := gfutils.ReleaseFinish(&gfutils.ReleaseFinishOptions{
			GFOutputHandler:            gflowOutputLogHandler,
			HotfixOrReleaseBaseOptions: baseOptions,
		})
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Fatal("fail to finish release")
		}

	default:
		log.WithFields(log.Fields{
			"branch-type": brType,
		}).Fatal("invalid branch type")
	}
}

func gflowGuessBranchType() string {
	bType, err := gitutils.GetCurrentBranch()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("fail to get current branch name")
	}

	if strings.HasPrefix(bType, "hotfix/") {
		return "hotfix"
	}
	if strings.HasPrefix(bType, "release/") {
		return "release"
	}

	log.WithFields(log.Fields{
		"current-branch": bType,
	}).Fatal("unable to guess operation type based on current branch name")

	return "" // never reached
}

func init() {
	rootCmd.AddCommand(gflowCmd)
	gflowCmd.AddCommand(gflowHotfixCmd)
	gflowCmd.AddCommand(gflowReleaseCmd)
	gflowHotfixCmd.AddCommand(gflowHotfixStartCmd)
	gflowHotfixCmd.AddCommand(gflowHotfixFinishCmd)
	gflowReleaseCmd.AddCommand(gflowReleaseStartCmd)
	gflowReleaseCmd.AddCommand(gflowReleaseFinishCmd)

	for _, cmd := range []*cobra.Command{gflowHotfixFinishCmd, gflowReleaseFinishCmd} {
		// git-flow wrapped flags
		cmd.Flags().BoolP("push", "p", false, "Push to origin after performing finish")
		cmd.Flags().StringP("message", "m", "", "Use the given tag message")
		cmd.Flags().StringP("messagefile", "f", "", "Use the contents of the given file as tag message")
		cmd.Flags().BoolP("no-prompt", "q", false, "No prompt for editor (set GIT_MERGE_AUTOEDIT=no)")
	}
}
