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
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	gh "github.com/google/go-github/v60/github"
	"github.com/northwood-labs/download-asset/github"
	"github.com/northwood-labs/golang-utils/exiterrorf"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	fArchivePath string
	fEndpoint    string
	fOwnerRepo   string
	fPattern     string
	fTag         string
	fVerbose     bool
	fWriteToBin  string
	fConstraint  string

	fDarwin    string
	fDragonfly string
	fFreeBSD   string
	fIllumos   string
	fLinux     string
	fNetBSD    string
	fOpenBSD   string
	fPlan9     string
	fSolaris   string
	fWindows   string

	fArm32    string
	fArm64    string
	fIntel32  string
	fIntel64  string
	fLoong64  string
	fMIPS32   string
	fMIPS64   string
	fMIPS64LE string
	fMIPS32LE string
	fPPC64    string
	fPPC64LE  string
	fRiscV64  string
	fS390x    string

	apiToken    = os.Getenv("GITHUB_TOKEN")
	apiEndpoint = ""
	release     *gh.RepositoryRelease

	currentOS  string
	currentCPU string

	textUnderline = lipgloss.NewStyle().
			Underline(true)

	// getCmd represents the get command
	getCmd = &cobra.Command{
		Use:   "get",
		Short: "Download an asset from a GitHub release",
		Long: LongHelpText(`
		Download an asset from a GitHub release.

		Can identify the current OS and current CPU architecture, then allows you to
		replace those matches with a more appropriate pattern to download the correct
		asset.

		--------------------------------------------------------------------------------

		Supported variables are {{.Ver}}, {{.OS}}, {{.Arch}}, and {{.Ext}}. These can
		be used with:
		    --pattern, --archive-path, --write-to-bin.

		Set --archive-path to the path of the binary inside of a compressed archive.
		Leave blank if the release asset is a binary itself.

		Set --write-to-bin to the name of the final binary. Will attempt to save to
		/usr/local/bin/NAME, but will fall back to $HOME/bin/NAME if /usr/local/bin is
		not writable.

		See https://bit.ly/3P1O9Rt for more information about setting GitHub API endpoints
		for GitHub Enterprise Server.

		--------------------------------------------------------------------------------

		Less common operating system flags not listed below are:
		    --dragonfly, --freebsd, --illumos, --netbsd, --openbsd, --plan9, --solaris

		Less common CPU architecture flags not listed below are:
		    --loong64, --mips32, --mips32le, --mips64, --mips64le, --ppc64, --ppc64le,
		    --riscv64`),
		Run: func(cmd *cobra.Command, args []string) {
			if apiToken == "" {
				exiterrorf.ExitErrorf(errors.New("GitHub token not found; set GITHUB_TOKEN environment variable"))
			}

			err := readConfig()
			if err != nil {
				exiterrorf.ExitErrorf(err)
			}

			t := table.New().
				Border(lipgloss.RoundedBorder()).
				BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
				BorderColumn(true).
				StyleFunc(func(row, col int) lipgloss.Style {
					return lipgloss.NewStyle().Padding(0, 1)
				}).
				Headers("FIELD", "VALUE")

			apiEndpoint, _, _ = github.ParseDomain(fEndpoint)

			if fVerbose {
				t.Row("GitHub endpoint", apiEndpoint)
				t.Row("GitHub token", apiToken[0:8]+".................................")
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
				t.Row("Owner", ownerRepo[0])
				t.Row("Repository", ownerRepo[1])
				if viper.ConfigFileUsed() != "" {
					t.Row("Config file", viper.ConfigFileUsed())
				}
			}

			// Apply values from configuration file.
			applyConfigValues(ownerRepo)

			// If we have a constraint, we need to find the latest tag that satisfies it.
			if fConstraint != "" {
				ref, err = github.GetLatestTag(client, ownerRepo[0], ownerRepo[1], fConstraint)
				if err != nil {
					exiterrorf.ExitErrorf(errors.Wrap(err, "failed to discover the release"))
				}

				fTag = github.RemoveVFromTag(ref.String())
			}

			if fTag == "latest" {
				release, err = github.GetLatestRelease(client, ownerRepo[0], ownerRepo[1])
				if err != nil {
					exiterrorf.ExitErrorf(errors.Wrap(err, "failed to discover the release"))
				}

				if fVerbose {
					t.Row("Latest release", *release.TagName)
				}
			} else {
				toTry := []string{
					fTag,
					github.InvertTag(fTag),
				}

				for i := range toTry {
					tag := toTry[i]

					release, err = github.GetReleaseVersion(client, ownerRepo[0], ownerRepo[1], tag)
					if err != nil {
						continue // skip to the next loop
					} else {
						break // we found a match; break out of loop
					}
				}

				if err != nil {
					exiterrorf.ExitErrorf(errors.Wrap(err, "failed to discover the release"))
				}
			}

			err = handleCurrentOSArch()
			if err != nil {
				exiterrorf.ExitErrorf(err)
			}

			patternVars := PatternMatches{
				Ver:  github.RemoveVFromTag(*release.TagName),
				OS:   currentOS,
				Arch: currentCPU,
				Ext: fmt.Sprintf("(%s)", strings.Join(
					[]string{
						// "7z",
						// "bz2",
						"exe",
						"gz",
						"tar.bz2",
						"tar.gz",
						// "tar.lz",
						"tar.xz",
						// "tar.Z",
						// "tar",
						"tbz2",
						"tgz",
						// "tlz",
						"txz",
						// "xz",
						"zip",
					}, "|",
				)),
			}

			resolvedArchivePath, err := replacePatternVariables(fArchivePath, patternVars)
			if err != nil {
				exiterrorf.ExitErrorf(err)
			}

			resolvedAssetPattern, err := replacePatternVariables(fPattern, patternVars)
			if err != nil {
				exiterrorf.ExitErrorf(err)
			}

			if fVerbose {
				t.Row("Current OS ident", currentOS)
				t.Row("Current CPU ident", currentCPU)
				t.Row("Asset pattern", fPattern)
				t.Row("Resolved pattern", resolvedAssetPattern)
			}

			// Check that we have everything before we trigger downloads
			if fPattern == "" || fWriteToBin == "" {
				exiterrorf.ExitErrorf(errors.New("missing one of pattern or write-to-bin"))
			}

			// Ready to download the asset
			archiveStream, name, err := github.GetAssetStream(
				client,
				ownerRepo,
				release,
				resolvedAssetPattern,
			)
			if err != nil {
				exiterrorf.ExitErrorf(err)
			}

			if fVerbose {
				t.Row("Matched asset name", name)
				t.Row("File inside archive", resolvedArchivePath)
				t.Row("Binary added to PATH", fWriteToBin)

				fmt.Println(t.Render())
			}

			binPath, err := github.DownloadStream(archiveStream, name, resolvedArchivePath, fWriteToBin)
			if err != nil {
				exiterrorf.ExitErrorf(err)
			}

			if resolvedArchivePath == "" {
				fmt.Printf(
					"Downloaded %s; renamed name → %s\n",
					textUnderline.Render(name),
					textUnderline.Render(binPath),
				)
			} else {
				fmt.Printf(
					"Downloaded %s; copied %s → %s\n",
					textUnderline.Render(name),
					textUnderline.Render(resolvedArchivePath),
					textUnderline.Render(binPath),
				)
			}

			err = archiveStream.Close()
			if err != nil {
				exiterrorf.ExitErrorf(err)
			}
		},
	}
)

type PatternMatches struct {
	Ver  string
	OS   string
	Arch string
	Ext  string
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Standard GitHub options.
	getCmd.Flags().StringVarP(
		&fOwnerRepo,
		"owner-repo",
		"r",
		"",
		"The owner and repository name in the format of 'owner/repo'.",
	)
	getCmd.Flags().StringVarP(
		&fEndpoint,
		"endpoint",
		"e",
		"https://api.github.com",
		"The GitHub API domain to use.",
	)
	getCmd.Flags().StringVarP(
		&fTag,
		"tag",
		"t",
		"latest",
		"The Git tag for which to check releases.",
	)
	getCmd.Flags().StringVarP(
		&fPattern,
		"pattern",
		"p",
		"",
		"The naming pattern of the asset name to match.",
	)
	getCmd.Flags().BoolVarP(
		&fVerbose,
		"verbose",
		"v",
		false,
		"Display verbose output.",
	)
	getCmd.Flags().StringVarP(
		&fArchivePath,
		"archive-path",
		"a",
		"",
		"The path to the file inside the archive.",
	)
	getCmd.Flags().StringVarP(
		&fWriteToBin,
		"write-to-bin",
		"w",
		"",
		"The final name of the binary.",
	)
	getCmd.Flags().StringVarP(
		&fConstraint,
		"constraint",
		"c",
		"",
		"Constrain the version to a particular range.",
	)

	handleFlags(getCmd)
}

func replacePatternVariables(pattern string, patternVars PatternMatches) (string, error) {
	tmpl, err := template.New("test").Parse(pattern)
	if err != nil {
		return "", errors.Wrapf(err, "failed to parse pattern: %s", pattern)
	}

	buf := new(bytes.Buffer)

	err = tmpl.Execute(buf, patternVars)
	if err != nil {
		return "", errors.Wrapf(err, "failed to apply values to pattern: %v", patternVars)
	}

	return buf.String(), nil
}

func readConfig() error {
	viper.SetConfigName("download-asset")
	viper.SetConfigType("toml")

	viper.AddConfigPath(".")                      // Current directory first.
	viper.AddConfigPath("$HOME/.download-asset/") // Then the HOME directory.
	viper.AddConfigPath("/etc/download-asset/")   // Then the system directory.

	errConfigFileNotFoundError := viper.ConfigFileNotFoundError{}

	if err := viper.ReadInConfig(); err != nil {
		if ok := errors.As(err, &errConfigFileNotFoundError); ok {
			// Config file not found; ignore error.
		} else {
			return errors.Wrap(err, "failed to read config file")
		}
	}

	return nil
}

func applyConfigValues(ownerRepo []string) {
	if viper.IsSet(strings.Join(ownerRepo, ".")) {
		flagMap := map[string]*string{
			// Config
			"endpoint":     &fEndpoint,
			"pattern":      &fPattern,
			"archive-path": &fArchivePath,
			"write-to-bin": &fWriteToBin,

			// OS
			"darwin":    &fDarwin,
			"dragonfly": &fDragonfly,
			"freebsd":   &fFreeBSD,
			"illumos":   &fIllumos,
			"linux":     &fLinux,
			"netbsd":    &fNetBSD,
			"openbsd":   &fOpenBSD,
			"plan9":     &fPlan9,
			"solaris":   &fSolaris,
			"windows":   &fWindows,

			// CPU Architectures
			"arm32":     &fArm32,
			"arm64":     &fArm64,
			"intel32":   &fIntel32,
			"intel64":   &fIntel64,
			"loong64":   &fLoong64,
			"mips32":    &fMIPS32,
			"mips32-le": &fMIPS32LE,
			"mips64":    &fMIPS64,
			"mips64-le": &fMIPS64LE,
			"ppc64":     &fPPC64,
			"ppc64le":   &fPPC64LE,
			"riscv64":   &fRiscV64,
			"s390x":     &fS390x,
		}

		for k := range flagMap {
			v := flagMap[k]

			if viper.IsSet(strings.Join(ownerRepo, ".") + "." + k) {
				*v = viper.GetString(strings.Join(ownerRepo, ".") + "." + k)
			}
		}
	}
}
