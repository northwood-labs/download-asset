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
	"runtime"
	"strings"
	"text/tabwriter"

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
	fMIPS     string
	fMIPS64   string
	fMIPS64LE string
	fMIPSle   string
	fPPC64    string
	fPPC64LE  string
	fRiscV64  string
	fS390x    string

	apiToken    = os.Getenv("GITHUB_TOKEN")
	apiEndpoint = ""
	release     *gh.RepositoryRelease

	currentOS  string
	currentCPU string

	// getCmd represents the get command
	getCmd = &cobra.Command{
		Use:   "get",
		Short: "Download an asset from a GitHub release",
		Long: `Download an asset from a GitHub release.

--------------------------------------------------------------------------------`,
		Run: func(cmd *cobra.Command, args []string) {
			if apiToken == "" {
				exiterrorf.ExitErrorf(errors.New("GitHub token not found; set GITHUB_TOKEN environment variable"))
			}

			err := readConfig()
			if err != nil {
				exiterrorf.ExitErrorf(err)
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

				if viper.ConfigFileUsed() != "" {
					fmt.Fprintf(w, " Config file:\t%s\t\n", viper.ConfigFileUsed())
				}
			}

			// Apply values from configuration file.
			applyConfigValues(ownerRepo)

			if fTag == "latest" {
				release, err = github.GetLatestRelease(client, ownerRepo[0], ownerRepo[1])
				if err != nil {
					exiterrorf.ExitErrorf(errors.Wrap(err, "failed to discover the release"))
				}

				if fVerbose {
					fmt.Fprintf(w, " Latest release:\t%s\t\n", *release.TagName)
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

				if fVerbose {
					fmt.Fprintf(w, " Selected release:\t%s\t\n", *release.TagName)
				}
			}

			switch runtime.GOOS {
			case "darwin":
				currentOS = fDarwin
			case "dragonfly":
				currentOS = fDragonfly
			case "freebsd":
				currentOS = fFreeBSD
			case "illumos":
				currentOS = fIllumos
			case "linux":
				currentOS = fLinux
			case "netbsd":
				currentOS = fNetBSD
			case "openbsd":
				currentOS = fOpenBSD
			case "plan9":
				currentOS = fPlan9
			case "solaris":
				currentOS = fSolaris
			case "windows":
				currentOS = fWindows
			default:
				exiterrorf.ExitErrorf(errors.New("unknown operating system"))
			}

			switch runtime.GOARCH {
			case "arm":
				currentCPU = fArm32
			case "arm64":
				currentCPU = fArm64
			case "386":
				currentCPU = fIntel32
			case "amd64":
				currentCPU = fIntel64
			case "loong64":
				currentCPU = fLoong64
			case "mips":
				currentCPU = fMIPS
			case "mips64":
				currentCPU = fMIPS64
			case "mips64le":
				currentCPU = fMIPS64LE
			case "mipsle":
				currentCPU = fMIPSle
			case "ppc64":
				currentCPU = fPPC64
			case "ppc64le":
				currentCPU = fPPC64LE
			case "riscv64":
				currentCPU = fRiscV64
			case "s390x":
				currentCPU = fS390x
			default:
				exiterrorf.ExitErrorf(errors.New("unknown CPU architecture"))
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
				fmt.Fprintf(w, " Current OS ident:\t%s\t\n", currentOS)
				fmt.Fprintf(w, " Current CPU ident:\t%s\t\n", currentCPU)
				fmt.Fprintf(w, " Asset pattern:\t%s\t\n", fPattern)
				fmt.Fprintf(w, " Resolved pattern:\t%s\t\n", resolvedAssetPattern)
			}

			err = w.Flush()
			if err != nil {
				exiterrorf.ExitErrorf(err)
			}

			// Check that we have everything before we trigger downloads
			if fPattern == "" || fWriteToBin == "" {
				exiterrorf.ExitErrorf(errors.New("missing one of pattern or write-to-bin"))
			}

			if fVerbose {
				fmt.Fprintf(w, " File inside archive:\t%s\t\n", resolvedArchivePath)
				fmt.Fprintf(w, " Binary added to PATH:\t%s\t\n", fWriteToBin)
				fmt.Fprintln(w, "")
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

			binPath, err := github.DownloadStream(archiveStream, name, resolvedArchivePath, fWriteToBin)
			if err != nil {
				exiterrorf.ExitErrorf(err)
			}

			fmt.Printf(
				"Downloaded %s; copied %s → %s\n",
				colorUnderlined.Sprintf(name),
				colorUnderlined.Sprintf(resolvedArchivePath),
				colorUnderlined.Sprintf(binPath),
			)

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
		"The GitHub API domain to use. See https://bit.ly/3P1O9Rt for more information.",
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
		"The naming pattern of the asset name to match. Supports a substring or regexp. "+
			"Supported variables are .Ver, .OS, .Arch, and .Ext.",
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
		"The final name of the binary. Will attempt to save to /usr/local/bin/NAME, but will fall back "+
			"to $HOME/bin/NAME if /usr/local/bin is not writable.",
	)

	// OS-specific options.
	getCmd.Flags().StringVarP(
		&fDarwin,
		"darwin",
		"",
		"darwin",
		"When the current OS is Darwin, use this value when looking up the correct asset.",
	)
	getCmd.Flags().StringVarP(
		&fDragonfly,
		"dragonfly",
		"",
		"dragonfly",
		"When the current OS is Dragonfly, use this value when looking up the correct asset.",
	)
	getCmd.Flags().StringVarP(
		&fFreeBSD,
		"freebsd",
		"",
		"freebsd",
		"When the current OS is FreeBSD, use this value when looking up the correct asset.",
	)
	getCmd.Flags().StringVarP(
		&fIllumos,
		"illumos",
		"",
		"illumos",
		"When the current OS is Illumos, use this value when looking up the correct asset.",
	)
	getCmd.Flags().StringVarP(
		&fLinux,
		"linux",
		"",
		"linux",
		"When the current OS is Linux, use this value when looking up the correct asset.",
	)
	getCmd.Flags().StringVarP(
		&fNetBSD,
		"netbsd",
		"",
		"netbsd",
		"When the current OS is NetBSD, use this value when looking up the correct asset.",
	)
	getCmd.Flags().StringVarP(
		&fOpenBSD,
		"openbsd",
		"",
		"openbsd",
		"When the current OS is OpenBSD, use this value when looking up the correct asset.",
	)
	getCmd.Flags().StringVarP(
		&fPlan9,
		"plan9",
		"",
		"plan9",
		"When the current OS is Plan9, use this value when looking up the correct asset.",
	)
	getCmd.Flags().StringVarP(
		&fSolaris,
		"solaris",
		"",
		"solaris",
		"When the current OS is Solaris, use this value when looking up the correct asset.",
	)
	getCmd.Flags().StringVarP(
		&fWindows,
		"windows",
		"",
		"windows",
		"When the current OS is Windows, use this value when looking up the correct asset.",
	)

	// CPU Arch-specific options.
	getCmd.Flags().StringVarP(
		&fArm32,
		"arm32",
		"",
		"arm",
		"When the current CPU architecture is 32-bit ARM, use this value when looking up the correct asset.",
	)
	getCmd.Flags().StringVarP(
		&fArm64,
		"arm64",
		"",
		"arm64",
		"When the current CPU architecture is 64-bit ARM, use this value when looking up the correct asset.",
	)
	getCmd.Flags().StringVarP(
		&fIntel32,
		"intel32",
		"",
		"386",
		"When the current CPU architecture is 32-bit Intel-compatible, use this value when looking up the correct asset.",
	)
	getCmd.Flags().StringVarP(
		&fIntel64,
		"intel64",
		"",
		"amd64",
		"When the current CPU architecture is 64-bit Intel-compatible, use this value when looking up the correct asset.",
	)
	getCmd.Flags().StringVarP(
		&fLoong64,
		"loong64",
		"",
		"loong64",
		"When the current CPU architecture is 64-bit Loongson, use this value when looking up the correct asset.",
	)
	getCmd.Flags().StringVarP(
		&fMIPS,
		"mips32",
		"",
		"mips",
		"When the current CPU architecture is 32-bit MIPS, use this value when looking up the correct asset.",
	)
	getCmd.Flags().StringVarP(
		&fMIPSle,
		"mips32-le",
		"",
		"mipsle",
		"When the current CPU architecture is 32-bit MIPS (LE), use this value when looking up the correct asset.",
	)
	getCmd.Flags().StringVarP(
		&fMIPS64,
		"mips64",
		"",
		"mips64",
		"When the current CPU architecture is 64-bit MIPS, use this value when looking up the correct asset.",
	)
	getCmd.Flags().StringVarP(
		&fMIPS64LE,
		"mips64-le",
		"",
		"mips64le",
		"When the current CPU architecture is 64-bit MIPS (LE), use this value when looking up the correct asset.",
	)
	getCmd.Flags().StringVarP(
		&fPPC64,
		"ppc64",
		"",
		"ppc64",
		"When the current CPU architecture is 64-bit PowerPC, use this value when looking up the correct asset.",
	)
	getCmd.Flags().StringVarP(
		&fPPC64LE,
		"ppc64le",
		"",
		"ppc64le",
		"When the current CPU architecture is 64-bit PowerPC (LE), use this value when looking up the correct asset.",
	)
	getCmd.Flags().StringVarP(
		&fRiscV64,
		"riscv64",
		"",
		"riscv64",
		"When the current CPU architecture is 64-bit RISC-V, use this value when looking up the correct asset.",
	)
	getCmd.Flags().StringVarP(
		&fS390x,
		"s390x",
		"",
		"s390x",
		"When the current CPU architecture is 64-bit s390x, use this value when looking up the correct asset.",
	)
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
			"mips32":    &fMIPS,
			"mips32-le": &fMIPSle,
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
