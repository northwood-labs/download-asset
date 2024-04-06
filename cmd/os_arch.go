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
	"text/tabwriter"

	"github.com/northwood-labs/golang-utils/exiterrorf"
	"github.com/spf13/cobra"
)

// osArchCmd represents the osArch command
var osArchCmd = &cobra.Command{
	Use:   "os-arch",
	Short: "Normalizes the OS and CPU architecture values to a standard format",
	Long: LongHelpText(`
	Stripped-down functionality of the 'get' command that only returns the OS and
	CPU architecture values in a normalized format. This is useful for scripting.

	--------------------------------------------------------------------------------

	Supported variables are {{.Ver}}, {{.OS}}, {{.Arch}}, and {{.Ext}}. These can
	be used with:
		--pattern.

	--------------------------------------------------------------------------------

	Less common operating system flags not listed below are:
		--dragonfly, --freebsd, --illumos, --netbsd, --openbsd, --plan9, --solaris

	Less common CPU architecture flags not listed below are:
		--loong64, --mips32, --mips32le, --mips64, --mips64le, --ppc64, --ppc64le,
		--riscv64`),
	Run: func(cmd *cobra.Command, args []string) {
		err := handleCurrentOSArch()
		if err != nil {
			exiterrorf.ExitErrorf(err)
		}

		patternVars := PatternMatches{
			OS:   currentOS,
			Arch: currentCPU,
		}

		resolvedAssetPattern, err := replacePatternVariables(fPattern, patternVars)
		if err != nil {
			exiterrorf.ExitErrorf(err)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

		if fVerbose {
			fmt.Fprintf(w, " Current OS ident:\t%s\t\n", currentOS)
			fmt.Fprintf(w, " Current CPU ident:\t%s\t\n", currentCPU)
			fmt.Fprintf(w, " Asset pattern:\t%s\t\n", fPattern)
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

func init() {
	rootCmd.AddCommand(osArchCmd)

	osArchCmd.Flags().StringVarP(
		&fPattern,
		"pattern",
		"p",
		"{{.OS}}/{{.Arch}}",
		"The naming pattern of the asset name to match.",
	)
	osArchCmd.Flags().BoolVarP(
		&fVerbose,
		"verbose",
		"v",
		false,
		"Display verbose output.",
	)

	handleFlags(osArchCmd)
}
