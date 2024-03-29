// Copyright 2023–2024, Northwood Labs
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/hashicorp/go-version"
	"github.com/northwood-labs/download-asset/github"
	"github.com/northwood-labs/golang-utils/exiterrorf"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	fStrip      bool
	fSkipToTags bool

	ref *version.Version
	tag string

	// latestTagCmd represents the latestTag command
	latestTagCmd = &cobra.Command{
		Use:   "latest-tag",
		Short: "Checks GitHub for the latest release or tag for a package",
		Long: `Checks the GitHub API for the latest release for a repository. If the repository
does not use the 'releases' feature, the latest Git tag is returned instead.

If you would prefer the latest tag over the latest release, use the
--skip-to-tags flag.

--------------------------------------------------------------------------------`,
		Run: func(cmd *cobra.Command, args []string) {
			if apiToken == "" {
				exiterrorf.ExitErrorf(errors.New("GitHub token not found; set GITHUB_TOKEN environment variable"))
			}

			if fVerbose {
				colorHeader.Println(" VERBOSE ")
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

			apiEndpoint, _, _ = github.ParseDomain(fEndpoint)

			if fVerbose {
				fmt.Fprintf(w, " GitHub endpoint:\t%s\t\n", apiEndpoint)
				fmt.Fprintf(w, " GitHub token:\t%s\t\n", apiToken[0:8]+".................................")
			}

			client, err := github.NewClient(&github.NewClientInput{
				Token:    apiToken,
				Endpoint: fEndpoint,
			})
			if err != nil {
				exiterrorf.ExitErrorf(errors.Wrap(err, "failed to create GitHub client"))
			}

			ownerRepo := strings.Split(fOwnerRepo, "/")
			if len(ownerRepo) != 2 { // lint:allow_raw_number
				exiterrorf.ExitErrorf(errors.New("invalid owner/repo"))
			}

			if fVerbose {
				fmt.Fprintf(w, " Owner:\t%s\t\n", ownerRepo[0])
				fmt.Fprintf(w, " Repository:\t%s\t\n", ownerRepo[1])
			}

			if fSkipToTags {
				ref, err = github.GetLatestTag(client, ownerRepo[0], ownerRepo[1])
				if err != nil {
					exiterrorf.ExitErrorf(errors.Wrap(err, "failed to discover the release"))
				}

				tag = github.RemoveVFromTag(ref.String())
			} else {
				release, err = github.GetLatestRelease(client, ownerRepo[0], ownerRepo[1])
				if err != nil {
					ref, err = github.GetLatestTag(client, ownerRepo[0], ownerRepo[1])
					if err != nil {
						exiterrorf.ExitErrorf(errors.Wrap(err, "failed to discover the release"))
					}
					tag = github.RemoveVFromTag(ref.String())
				} else {
					tag = github.RemoveVFromTag(*release.TagName)
				}
			}

			if fVerbose {
				fmt.Fprintf(w, " Latest release:\t%s\t\n", tag)
				fmt.Fprintln(w, "")
			}

			err = w.Flush()
			if err != nil {
				exiterrorf.ExitErrorf(err)
			}

			if fStrip {
				fmt.Print(tag)
			} else {
				fmt.Println(tag)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(latestTagCmd)

	// Standard GitHub options.
	latestTagCmd.Flags().StringVarP(
		&fOwnerRepo,
		"owner-repo",
		"r",
		"",
		"The owner and repository name in the format of 'owner/repo'.",
	)
	latestTagCmd.Flags().StringVarP(
		&fEndpoint,
		"endpoint",
		"e",
		"https://api.github.com",
		"The GitHub API domain to use. See https://bit.ly/3P1O9Rt for more information.",
	)
	latestTagCmd.Flags().BoolVarP(
		&fVerbose,
		"verbose",
		"v",
		false,
		"Display verbose output.",
	)
	latestTagCmd.Flags().BoolVarP(
		&fStrip,
		"strip",
		"s",
		false,
		"Strip the trailing line ending.",
	)
	latestTagCmd.Flags().BoolVarP(
		&fSkipToTags,
		"skip-to-tags",
		"t",
		false,
		"Skip looking up releases, and just look at tags.",
	)
}
