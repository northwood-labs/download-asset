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
	"runtime"
	"text/tabwriter"

	"github.com/northwood-labs/golang-utils/exiterrorf"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	oafPattern string
	oafVerbose bool

	oafDarwin    string
	oafDragonfly string
	oafFreeBSD   string
	oafIllumos   string
	oafLinux     string
	oafNetBSD    string
	oafOpenBSD   string
	oafPlan9     string
	oafSolaris   string
	oafWindows   string

	oafArm32    string
	oafArm64    string
	oafIntel32  string
	oafIntel64  string
	oafLoong64  string
	oafMIPS     string
	oafMIPS64   string
	oafMIPS64LE string
	oafMIPSle   string
	oafPPC64    string
	oafPPC64LE  string
	oafRiscV64  string
	oafS390x    string

	// osArchCmd represents the osArch command
	osArchCmd = &cobra.Command{
		Use:   "os-arch",
		Short: "Normalizes the OS and CPU architecture values to a standard format",
		Long: `Stripped-down functionality of the 'get' command that only returns the OS and
CPU architecture values in a normalized format.

--------------------------------------------------------------------------------`,
		Run: func(cmd *cobra.Command, args []string) {
			switch runtime.GOOS {
			case "darwin":
				currentOS = oafDarwin
			case "dragonfly":
				currentOS = oafDragonfly
			case "freebsd":
				currentOS = oafFreeBSD
			case "illumos":
				currentOS = oafIllumos
			case "linux":
				currentOS = oafLinux
			case "netbsd":
				currentOS = oafNetBSD
			case "openbsd":
				currentOS = oafOpenBSD
			case "plan9":
				currentOS = oafPlan9
			case "solaris":
				currentOS = oafSolaris
			case "windows":
				currentOS = oafWindows
			default:
				exiterrorf.ExitErrorf(errors.New("unknown operating system"))
			}

			switch runtime.GOARCH {
			case "arm":
				currentCPU = oafArm32
			case "arm64":
				currentCPU = oafArm64
			case "386":
				currentCPU = oafIntel32
			case "amd64":
				currentCPU = oafIntel64
			case "loong64":
				currentCPU = oafLoong64
			case "mips":
				currentCPU = oafMIPS
			case "mips64":
				currentCPU = oafMIPS64
			case "mips64le":
				currentCPU = oafMIPS64LE
			case "mipsle":
				currentCPU = oafMIPSle
			case "ppc64":
				currentCPU = oafPPC64
			case "ppc64le":
				currentCPU = oafPPC64LE
			case "riscv64":
				currentCPU = oafRiscV64
			case "s390x":
				currentCPU = oafS390x
			default:
				exiterrorf.ExitErrorf(errors.New("unknown CPU architecture"))
			}

			patternVars := PatternMatches{
				OS:   currentOS,
				Arch: currentCPU,
			}

			resolvedAssetPattern, err := replacePatternVariables(oafPattern, patternVars)
			if err != nil {
				exiterrorf.ExitErrorf(err)
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

			if oafVerbose {
				fmt.Fprintf(w, " Current OS ident:\t%s\t\n", currentOS)
				fmt.Fprintf(w, " Current CPU ident:\t%s\t\n", currentCPU)
				fmt.Fprintf(w, " Asset pattern:\t%s\t\n", oafPattern)
				fmt.Fprintf(w, " Resolved pattern:\t%s\t\n", resolvedAssetPattern)
				fmt.Fprintln(w, "")
			}

			err = w.Flush()
			if err != nil {
				exiterrorf.ExitErrorf(err)
			}

			fmt.Println(resolvedAssetPattern)
		},
	}
)

func init() {
	rootCmd.AddCommand(osArchCmd)

	osArchCmd.Flags().StringVarP(
		&oafPattern,
		"pattern",
		"p",
		"{{.OS}}/{{.Arch}}",
		"The naming pattern of the asset name to match. Supports a substring or regexp. "+
			"Supported variables are .OS and .Arch.",
	)
	osArchCmd.Flags().BoolVarP(
		&oafVerbose,
		"verbose",
		"v",
		false,
		"Display verbose output.",
	)

	// OS-specific options.
	osArchCmd.Flags().StringVarP(
		&oafDarwin,
		"darwin",
		"",
		"darwin",
		"When the current OS is Darwin, use this value when looking up the correct asset.",
	)
	osArchCmd.Flags().StringVarP(
		&oafDragonfly,
		"dragonfly",
		"",
		"dragonfly",
		"When the current OS is Dragonfly, use this value when looking up the correct asset.",
	)
	osArchCmd.Flags().StringVarP(
		&oafFreeBSD,
		"freebsd",
		"",
		"freebsd",
		"When the current OS is FreeBSD, use this value when looking up the correct asset.",
	)
	osArchCmd.Flags().StringVarP(
		&oafIllumos,
		"illumos",
		"",
		"illumos",
		"When the current OS is Illumos, use this value when looking up the correct asset.",
	)
	osArchCmd.Flags().StringVarP(
		&oafLinux,
		"linux",
		"",
		"linux",
		"When the current OS is Linux, use this value when looking up the correct asset.",
	)
	osArchCmd.Flags().StringVarP(
		&oafNetBSD,
		"netbsd",
		"",
		"netbsd",
		"When the current OS is NetBSD, use this value when looking up the correct asset.",
	)
	osArchCmd.Flags().StringVarP(
		&oafOpenBSD,
		"openbsd",
		"",
		"openbsd",
		"When the current OS is OpenBSD, use this value when looking up the correct asset.",
	)
	osArchCmd.Flags().StringVarP(
		&oafPlan9,
		"plan9",
		"",
		"plan9",
		"When the current OS is Plan9, use this value when looking up the correct asset.",
	)
	osArchCmd.Flags().StringVarP(
		&oafSolaris,
		"solaris",
		"",
		"solaris",
		"When the current OS is Solaris, use this value when looking up the correct asset.",
	)
	osArchCmd.Flags().StringVarP(
		&oafWindows,
		"windows",
		"",
		"windows",
		"When the current OS is Windows, use this value when looking up the correct asset.",
	)

	// CPU Arch-specific options.
	osArchCmd.Flags().StringVarP(
		&oafArm32,
		"arm32",
		"",
		"arm",
		"When the current CPU architecture is 32-bit ARM, use this value when looking up the correct asset.",
	)
	osArchCmd.Flags().StringVarP(
		&oafArm64,
		"arm64",
		"",
		"arm64",
		"When the current CPU architecture is 64-bit ARM, use this value when looking up the correct asset.",
	)
	osArchCmd.Flags().StringVarP(
		&oafIntel32,
		"intel32",
		"",
		"386",
		"When the current CPU architecture is 32-bit Intel-compatible, use this value when looking up the correct asset.",
	)
	osArchCmd.Flags().StringVarP(
		&oafIntel64,
		"intel64",
		"",
		"amd64",
		"When the current CPU architecture is 64-bit Intel-compatible, use this value when looking up the correct asset.",
	)
	osArchCmd.Flags().StringVarP(
		&oafLoong64,
		"loong64",
		"",
		"loong64",
		"When the current CPU architecture is 64-bit Loongson, use this value when looking up the correct asset.",
	)
	osArchCmd.Flags().StringVarP(
		&oafMIPS,
		"mips32",
		"",
		"mips",
		"When the current CPU architecture is 32-bit MIPS, use this value when looking up the correct asset.",
	)
	osArchCmd.Flags().StringVarP(
		&oafMIPSle,
		"mips32-le",
		"",
		"mipsle",
		"When the current CPU architecture is 32-bit MIPS (LE), use this value when looking up the correct asset.",
	)
	osArchCmd.Flags().StringVarP(
		&oafMIPS64,
		"mips64",
		"",
		"mips64",
		"When the current CPU architecture is 64-bit MIPS, use this value when looking up the correct asset.",
	)
	osArchCmd.Flags().StringVarP(
		&oafMIPS64LE,
		"mips64-le",
		"",
		"mips64le",
		"When the current CPU architecture is 64-bit MIPS (LE), use this value when looking up the correct asset.",
	)
	osArchCmd.Flags().StringVarP(
		&oafPPC64,
		"ppc64",
		"",
		"ppc64",
		"When the current CPU architecture is 64-bit PowerPC, use this value when looking up the correct asset.",
	)
	osArchCmd.Flags().StringVarP(
		&oafPPC64LE,
		"ppc64le",
		"",
		"ppc64le",
		"When the current CPU architecture is 64-bit PowerPC (LE), use this value when looking up the correct asset.",
	)
	osArchCmd.Flags().StringVarP(
		&oafRiscV64,
		"riscv64",
		"",
		"riscv64",
		"When the current CPU architecture is 64-bit RISC-V, use this value when looking up the correct asset.",
	)
	osArchCmd.Flags().StringVarP(
		&oafS390x,
		"s390x",
		"",
		"s390x",
		"When the current CPU architecture is 64-bit s390x, use this value when looking up the correct asset.",
	)
}
