# Download Assets

Simplifies the process of downloading release assets from GitHub for the current operating system and current CPU architecture.

```gherkin
Given the GitHub repository identifier
  And a mapping of the current OS to the asset's naming pattern
  And a mapping of the current CPU to the asset's naming pattern
  And a pattern to match for the asset names
  And the location of the binary inside the requested archive
 Then download the archive
  And read the files inside of the archive
  And extract the requested file from the archive
  And save it to your $PATH.
```

## The problem this helps with

I work with Linux in Docker on the daily. And the people and machines I support have a blend of Intel/AMD and ARM/Graviton/Apple Silicon chips.

<details>
<summary>Read more…</summary><br>

* We have users on macOS, Windows, and Linux.
* We have a blend of worker laptops using both Intel/AMD and Apple Silicon CPU architectures.
* We have cloud servers in AWS, GCP, Azure, and Oracle Cloud.
* We have a blend of user machines using both Intel/AMD and Apple Silicon CPU architectures.
* We rely heavily on Docker/Terraform/OpenTofu for consistency/repeatability, and to better scale the perpetually-limited resources of our DevOps/SRE/Cloud/Platform engineering teams.
* Docker runs natively on Linux.
* Docker runs virtualized in macOS and Windows.
* Software running inside the Linux-based Docker containers is most efficient when compiled for the current CPU architecture.
* Out on the internet, people build packages that can be installed. Many are not inside the Linux system’s package manager, and must be installed from the web. The people who publish these packages use a variety of identifiers for Intel-compatible vs ARM-compatible CPU architectures. There is no consistency.

When building tooling/solutions for a heterogenous set of machines across an enterprise, you need to solve for (at least) the following matrix.

* Current OS
* Current CPU architecture
* Package filenames on the internet

Deploying software as Docker containers (running Linux) helps normalize things like:

* Relying on GNU vs BSD-flavored CLI tools
* Download packages into the Docker container, worrying only about Linux
* Deploying software across worker laptops running different host operating systems
* Deploying software to Linux servers in the cloud

But these solutions don't solve the (relatively new) problem of an uptick of 64-bit ARM software/CPUs being added to the matrix — _and the fact that these are not referred-to in a unified, consistent way_.

### Common values for `uname -m`

| OS              | 64-bit Intel-compat | 64-bit ARM |
|-----------------|---------------------|------------|
| macOS           | `x86_64`            | `arm64`    |
| Red Hat family¹ | `x86_64`            | `aarch64`  |
| Debian family²  | `amd64`             | `arm64`    |
| Busybox family³ | `x86_64`            | `aarch64`  |
| Windows WSL2⁴   | _Varies_            | _Varies_   |

<footnote>

* ¹ Red Hat family includes Red Hat Enterprise Linux, CentOS, Fedora, Amazon Linux, and others.
* ² Debian family includes Debian, Ubuntu, Linux Mint, and others.
* ³ Busybox family includes Busybox, Alpine Linux, and others.
* ⁴ Windows WSL2 returns whatever the underlying Linux installation says.

</footnote>

</details>

## Extremely brief overview of CPU architectures

There are different names for (essentially) the same CPU architectures. Different vendors use different names for the same thing.

<details>
<summary>Read more…</summary><br>

Here's an (extremely) brief overview of modern CPU architectures that you most commonly find in cloud service providers and modern desktops/laptops.

This is meant to be _illustrative_, not _comprehensive_. As of today, these are the top 2 by a large margin.

| Family | Arch IDs                                 | Description                                                                                                                                                                                                                                                          |
|--------|------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `x86`  | `x86_64`, `amd64`, `x64`                 | Intel’s 80x86 line of CPUs, and AMD clones. Shortened to `x86` (or sometimes `x64`), these are the newer 64-bit models. Includes Amazon EC2 instances powered by Intel Xeon™ or AMD EPYC™ CPUs, and Intel i-Series Macs.                                             |
| `arm`  | `arm64`, `arm64v8`, `arm64v9`, `aarch64` | ARM v8/v9, 64-bit. AWS Graviton, Apple A7 and newer (including M-series). All 64-bit ARM chips are ARM v8/v9, but the inverse is not true. `arm64 == ( arm64v8 \|\| arm64v9 )`. Includes Amazon EC2 instances powered by AWS Graviton™ CPUs and Apple M-Series Macs. |

</details>

## Installation

With [Go](https://go.dev) installed:

```bash
go install github.com/northwood-labs/download-asset@latest
```

This will download and compile `download-asset` on-demand for your current OS and CPU architecture.

> [!IMPORTANT]
> We are _very intentionally_ NOT attaching pre-built assets to releases because it creates a chicken-and-egg problem. You'd have to select your OS and CPU architecture in order to install the code, which is designed to dynamically _figure-out_ your OS and CPU architecture in order to download other assets. That seems silly to us.

## Usage

> [!NOTE]
> Using [`aquasecurity/trivy`](https://github.com/aquasecurity/trivy/releases) as an example repository since they build for several systems, and use many non-standard names. But this will work for any GitHub repository.

Taking the [v0.49.1](https://github.com/aquasecurity/trivy/releases/tag/v0.49.1) release as an example, we can look at the list of assets and see things with all sorts of different names. We see `.deb` and `.rpm` files for Linux, we see `.tar.gz` files, we see `.zip` files for Windows, and we see `.pem` and `.sig` files for just about everything.

We also see things like `Linux-64bit` (which 64-bit?), `macOS-ARM64` (which is usually `darwin`), and `windows-64bit.zip` (which has different capitalization from `Linux`). Some macOS builds in other repositories have _universal_ builds instead of separate Intel/Apple Silicon builds.

* How do we simplify how we automate downloading the things we need?
* And even more, how do we automate our `Dockerfile`?

### Downloading an archive from GitHub.com

First, set-up the `GITHUB_TOKEN` environment variable. GitHub's [unauthenticated rate limit](https://docs.github.com/en/rest/using-the-rest-api/rate-limits-for-the-rest-api?apiVersion=2022-11-28#primary-rate-limit-for-unauthenticated-users) is 60 requests/hour. However, creating a token with no permissions will raise the [(authenticated) rate limit](https://docs.github.com/en/rest/using-the-rest-api/rate-limits-for-the-rest-api?apiVersion=2022-11-28#primary-rate-limit-for-authenticated-users) to 5,000 requests/hour.

The `download-asset` binary will read this environment variable by default, and use it to make requests.

```bash
download-asset get \
  --owner-repo aquasecurity/trivy \
  --tag v0.49.1 \
  --darwin macOS \
  --arm64 ARM64 \
  --pattern 'trivy_{{.Ver}}_{{.OS}}-{{.Arch}}.{{.Ext}}$' \
  ;
```

Let's break down this set of flags.

<details>
<summary>Read more…</summary>

#### `download-asset get`

This is the binary, and the subcommand `get`. Use the `--help` flag to get more information about additional options.

#### `--owner-repo`

Since we (at the moment) only support GitHub releases, this is the `owner/repository` pattern. In this example, we're going to download from [`aquasecurity/trivy`](https://github.com/aquasecurity/trivy/releases).

#### `--tag v0.49.1`

Here, we've specified a tag in the repository. Since _Assets_ can only be attached to _Releases_, this MUST be a _Tag_ that has a _Release_ attached to it. If you just want to grab the latest release (i.e., the release that is flagged as _latest_, not necessarily the highest version number), then you can either set `--tag latest`, or omit the flag all-together.

It will try the tag with a prepended `v`, then without a prepended `v`, and will respond if either of them match. If the tag doesn't exist, or follows a different format, `download-asset` will throw an error.

#### `--darwin macOS`

This flag only applies when the current system is a macOS system. The same is true for the `--linux`, `--windows`, `--freebsd`, and other OS-specific flags. If the current system is macOS (`darwin`), then this is the string to use for the `{{.OS}}` value in the `--pattern` tag (more in a moment.)

You should set this for each OS you plan to download assets for, with the values matching the strings in the list of assets.

```bash
--darwin macOS \
--linux Linux \
--windows windows \
--freebsd FreeBSD \
--netbsd NetBSD \
```

#### `--arm64 ARM64`

This flag only applies when the current CPU architecture is 64-bit ARM. The same is true for the `--arm32`, `--intel64`, `--intel32`, and other CPU architecture-specific flags. If the current system is 64-bit ARM (`arm64`), then this is the string to use for the `{{.Arch}}` value in the `--pattern` tag (more in a moment.)

You should set this for each CPU architecture you plan to download assets for, with the values matching the strings in the list of assets.

```bash
--arm32 ARM \
--arm64 ARM64 \
--intel32 32bit \
--intel64 64bit \
--ppc32 PPC \
--ppc64 PPC64 \
--s390x s390x \
```

#### `--pattern 'trivy_{{.Ver}}_{{.OS}}-{{.Arch}}.{{.Ext}}$'`

This is the naming pattern to match when looking through the list of _Assets_ attached to the _Release_. We already talked about the `.OS` and `.Arch` values, above.

The `.Ver` value is the tag (or the tag resolved when we selected `latest`) WITHOUT the prepended `v`. If the _Asset_ name contains a `v` before the version, you should add the `v` directly in the `--pattern` value.

The `.Ext` value is a regular expression that matches most common archive file extensions (e.g., `7z`, `xz`, `tar.gz`, `tgz`, `tar.bz2`, `tbz2`, `zip`) WITHOUT the preceding `.`.

Since this is a [regular expression](https://pkg.go.dev/regexp), the `$` at the end means _end of the string_. This helps you avoid matches for `Linux-ARM64.tar.gz.sig` or `windows-64bit.zip.pem` since this tool will download the _first match it finds_. In order to ensure you get what you want, you are advised to make your _pattern_ as specific as possible.

If your pattern (after resolving for `.Ver`, `.OS`, `.Arch`, and `.Ext`) is not a valid Go [regular expression](https://pkg.go.dev/regexp) pattern, the app will panic and exit. The regular expression is passed through [`regexp.MustCompile`](https://pkg.go.dev/regexp#MustCompile).

</details>

### Downloading an archive from GitHub Enterprise Server

Same thing as above, with small changes and a couple of notes:

<details>
<summary>Read more…</summary>

1. The `GITHUB_TOKEN` environment variable should be generated from your GitHub Enterprise Server instance, not public GitHub.com.

1. If your instance has [Subdomain Isolation](https://docs.github.com/en/enterprise-server@latest/admin/configuration/hardening-security-for-your-enterprise/enabling-subdomain-isolation) enabled, then your `--endpoint` flag is likely going to be `api.github.company.com`. Without subdomain isolation, it will likely be `github.company.com`. If you're not sure, ask your GitHub Enterprise Server administrators.

1. For `--endpoint`, pass the scheme+hostname (e.g., `https://github.company.com` or `https://api.github.company.com`). If your instance is running over insecure HTTP (port 80), specify `http://`. If you do not specify a scheme (e.g., `api.github.company.com`), then we will _assume_ HTTPS.

1. Keep in mind that the `--owner-repo` flag will refer to your organization's GitHub Enterprise Server environment, and NOT public GitHub.com.

</details>

```bash
download-asset get \
  --owner-repo myteam/myproject \
  --endpoint github.company.com \
    # other flags... \
  ;
```

### Automating a `Dockerfile`

We'll make a few assumptions here:

1. You are building [multi-platform](https://docs.docker.com/build/building/multi-platform/) Docker/OCI images. (We'll just focus on `linux/amd64` and `linux/arm64` in this example.)
1. You are leveraging [multi-stage](https://docs.docker.com/build/building/multi-stage/) builds.
1. You are leveraging [BuildKit secrets](https://docs.docker.com/build/building/secrets/), or equivalent.
1. You are willing to have one of your build stages be a Go image that will be thrown out after the stage is complete.

```Dockerfile
# syntax=docker/dockerfile:1
FROM --platform=$TARGETPLATFORM golang:1.22-alpine AS go-installer

RUN go install github.com/northwood-labs/download-asset@latest
RUN --mount=type=secret,id=github_token \
    GITHUB_TOKEN="$(cat /run/secrets/github_token)" \
    download-asset get \
        --owner-repo aquasecurity/trivy \
        --tag latest \
        --linux Linux \
        --intel64 64bit \
        --arm64 ARM64 \
        --pattern 'trivy_{{.Ver}}_{{.OS}}-{{.Arch}}.{{.Ext}}$' \
        ;
```

When setting up your final _build stage_, you would use `COPY --from` syntax to pull the downloaded binaries from the temporary Go-based stage into your final build stage.

```Dockerfile
COPY --from=go-installer /usr/local/bin/trivy /usr/local/bin/trivy
```
