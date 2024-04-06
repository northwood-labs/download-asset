package cmd

import (
	"runtime"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func handleFlags(cmd *cobra.Command) {
	// OS-specific options.
	cmd.Flags().StringVarP(
		&fDarwin,
		"darwin",
		"",
		"darwin",
		"When Darwin, set .OS to this value.",
	)
	cmd.Flags().StringVarP(
		&fDragonfly,
		"dragonfly",
		"",
		"dragonfly",
		"When Dragonfly, set .OS to this value.",
	)
	cmd.Flags().StringVarP(
		&fFreeBSD,
		"freebsd",
		"",
		"freebsd",
		"When FreeBSD, set .OS to this value.",
	)
	cmd.Flags().StringVarP(
		&fIllumos,
		"illumos",
		"",
		"illumos",
		"When Illumos, set .OS to this value.",
	)
	cmd.Flags().StringVarP(
		&fLinux,
		"linux",
		"",
		"linux",
		"When Linux, set .OS to this value.",
	)
	cmd.Flags().StringVarP(
		&fNetBSD,
		"netbsd",
		"",
		"netbsd",
		"When NetBSD, set .OS to this value.",
	)
	cmd.Flags().StringVarP(
		&fOpenBSD,
		"openbsd",
		"",
		"openbsd",
		"When OpenBSD, set .OS to this value.",
	)
	cmd.Flags().StringVarP(
		&fPlan9,
		"plan9",
		"",
		"plan9",
		"When Plan9, set .OS to this value.",
	)
	cmd.Flags().StringVarP(
		&fSolaris,
		"solaris",
		"",
		"solaris",
		"When Solaris, set .OS to this value.",
	)
	cmd.Flags().StringVarP(
		&fWindows,
		"windows",
		"",
		"windows",
		"When Windows, set .OS to this value.",
	)

	_ = cmd.Flags().MarkHidden("dragonfly") // lint:allow_unhandled
	_ = cmd.Flags().MarkHidden("freebsd")   // lint:allow_unhandled
	_ = cmd.Flags().MarkHidden("illumos")   // lint:allow_unhandled
	_ = cmd.Flags().MarkHidden("netbsd")    // lint:allow_unhandled
	_ = cmd.Flags().MarkHidden("openbsd")   // lint:allow_unhandled
	_ = cmd.Flags().MarkHidden("plan9")     // lint:allow_unhandled
	_ = cmd.Flags().MarkHidden("solaris")   // lint:allow_unhandled

	// CPU Arch-specific options.
	cmd.Flags().StringVarP(
		&fArm32,
		"arm32",
		"",
		"arm",
		"When 32-bit ARM, set .Arch to this value.",
	)
	cmd.Flags().StringVarP(
		&fArm64,
		"arm64",
		"",
		"arm64",
		"When 64-bit ARM, set .Arch to this value.",
	)
	cmd.Flags().StringVarP(
		&fIntel32,
		"intel32",
		"",
		"386",
		"When 32-bit Intel-compat, set .Arch to this value.",
	)
	cmd.Flags().StringVarP(
		&fIntel64,
		"intel64",
		"",
		"amd64",
		"When 64-bit Intel-compat, set .Arch to this value.",
	)
	cmd.Flags().StringVarP(
		&fLoong64,
		"loong64",
		"",
		"loong64",
		"When 64-bit Loongson, set .Arch to this value.",
	)
	cmd.Flags().StringVarP(
		&fMIPS32,
		"mips32",
		"",
		"mips",
		"When 32-bit MIPS, set .Arch to this value.",
	)
	cmd.Flags().StringVarP(
		&fMIPS32LE,
		"mips32le",
		"",
		"mipsle",
		"When 32-bit MIPS (LE), set .Arch to this value.",
	)
	cmd.Flags().StringVarP(
		&fMIPS64,
		"mips64",
		"",
		"mips64",
		"When 64-bit MIPS, set .Arch to this value.",
	)
	cmd.Flags().StringVarP(
		&fMIPS64LE,
		"mips64le",
		"",
		"mips64le",
		"When 64-bit MIPS (LE), set .Arch to this value.",
	)
	cmd.Flags().StringVarP(
		&fPPC64,
		"ppc64",
		"",
		"ppc64",
		"When 64-bit PowerPC, set .Arch to this value.",
	)
	cmd.Flags().StringVarP(
		&fPPC64LE,
		"ppc64le",
		"",
		"ppc64le",
		"When 64-bit PowerPC (LE), set .Arch to this value.",
	)
	cmd.Flags().StringVarP(
		&fRiscV64,
		"riscv64",
		"",
		"riscv64",
		"When 64-bit RISC-V, set .Arch to this value.",
	)
	cmd.Flags().StringVarP(
		&fS390x,
		"s390x",
		"",
		"s390x",
		"When 64-bit s390x, set .Arch to this value.",
	)

	_ = cmd.Flags().MarkHidden("loong64")  // lint:allow_unhandled
	_ = cmd.Flags().MarkHidden("mips32")   // lint:allow_unhandled
	_ = cmd.Flags().MarkHidden("mips32le") // lint:allow_unhandled
	_ = cmd.Flags().MarkHidden("mips64")   // lint:allow_unhandled
	_ = cmd.Flags().MarkHidden("mips64le") // lint:allow_unhandled
	_ = cmd.Flags().MarkHidden("ppc64")    // lint:allow_unhandled
	_ = cmd.Flags().MarkHidden("ppc64le")  // lint:allow_unhandled
	_ = cmd.Flags().MarkHidden("riscv64")  // lint:allow_unhandled
	_ = cmd.Flags().MarkHidden("s390x")    // lint:allow_unhandled
}

func handleCurrentOSArch() error {
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
		return errors.New("unknown operating system")
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
		currentCPU = fMIPS32
	case "mips64":
		currentCPU = fMIPS64
	case "mips64le":
		currentCPU = fMIPS64LE
	case "mipsle":
		currentCPU = fMIPS32LE
	case "ppc64":
		currentCPU = fPPC64
	case "ppc64le":
		currentCPU = fPPC64LE
	case "riscv64":
		currentCPU = fRiscV64
	case "s390x":
		currentCPU = fS390x
	default:
		return errors.New("unknown CPU architecture")
	}

	return nil
}
