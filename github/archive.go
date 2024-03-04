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

package github

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/ulikunitz/xz"
)

func Decompress(archiveStream io.ReadCloser, filename, findPattern, writeToBin string) (string, error) {
	var binPath string

	// .tar.gz or .tgz
	if regexp.MustCompile(`\.tar\.gz$`).MatchString(filename) ||
		regexp.MustCompile(`\.tgz$`).MatchString(filename) {
		g, err := gzip.NewReader(archiveStream)
		if err != nil {
			return "", errors.Wrap(err, "failed to create gzip reader")
		}

		binPath, err = handleTar(g, findPattern, writeToBin)
		if err != nil {
			return binPath, err
		}

		return binPath, nil
	} else if regexp.MustCompile(`\.tar\.xz$`).MatchString(filename) ||
		regexp.MustCompile(`\.txz$`).MatchString(filename) {
		x, err := xz.NewReader(archiveStream)
		if err != nil {
			return "", errors.Wrap(err, "failed to create xz reader")
		}

		binPath, err = handleTar(x, findPattern, writeToBin)
		if err != nil {
			return binPath, err
		}

		return binPath, nil
	} else if regexp.MustCompile(`\.tar\.bz2$`).MatchString(filename) ||
		regexp.MustCompile(`\.tbz2$`).MatchString(filename) {
		b := bzip2.NewReader(archiveStream)

		binPath, err := handleTar(b, findPattern, writeToBin)
		if err != nil {
			return binPath, err
		}

		return binPath, nil
	} else if regexp.MustCompile(`\.zip$`).MatchString(filename) {
		binPath, err := handleZip(archiveStream, findPattern, writeToBin)
		if err != nil {
			return binPath, err
		}

		return binPath, nil
	} else {
		// /usr/local/bin
		binPath := "/" + filepath.Join("usr", "local", "bin", writeToBin)

		f, err := os.Create(binPath) // lint:allow_include_file
		if err != nil {
			// ~/bin
			binPath = filepath.Join(os.Getenv("HOME"), "bin", writeToBin)

			f, err = os.Create(binPath) // lint:allow_include_file
			if err != nil {
				return binPath, errors.Wrap(err, "failed to create file")
			}
		}

		err = f.Chmod(0o0755) // lint:allow_raw_number
		if err != nil {
			return binPath, errors.Wrap(err, "failed to make executable")
		}

		_, err = io.Copy(f, archiveStream) // lint:allow_decompress
		if err != nil {
			return binPath, errors.Wrap(err, "error reading binary file")
		}

		err = f.Close()
		if err != nil {
			return binPath, errors.Wrap(err, "could not close the new file")
		}

		return binPath, nil
	}

	return "", nil
}

func handleTar(g io.Reader, findPattern, writeToBin string) (string, error) {
	t := tar.NewReader(g)

	// /usr/local/bin
	binPath := "/" + filepath.Join("usr", "local", "bin", writeToBin)

	for {
		hdr, err := t.Next()
		if err == io.EOF {
			break // End of archive
		}

		if err != nil {
			return binPath, errors.Wrap(err, "error reading tar header")
		}

		if hdr.Typeflag == tar.TypeReg {
			if !strings.EqualFold(hdr.Name, findPattern) {
				continue
			}

			f, err := os.Create(binPath) // lint:allow_include_file
			if err != nil {
				// ~/bin
				binPath = filepath.Join(os.Getenv("HOME"), "bin", writeToBin)

				f, err = os.Create(binPath) // lint:allow_include_file
				if err != nil {
					return binPath, errors.Wrap(err, "failed to create file")
				}
			}

			err = f.Chmod(0o0755) // lint:allow_raw_number
			if err != nil {
				return binPath, errors.Wrap(err, "failed to make executable")
			}

			_, err = io.Copy(f, t) // lint:allow_decompress
			if err != nil {
				return binPath, errors.Wrap(err, "error reading tar header")
			}

			err = f.Close()
			if err != nil {
				return binPath, errors.Wrap(err, "could not close the new file")
			}
		}
	}

	return binPath, nil
}

func handleZip(z io.ReadCloser, findPattern, writeToBin string) (string, error) {
	b, err := io.ReadAll(z) // The readCloser is the one from the zip-package
	if err != nil {
		return "", errors.Wrap(err, "error reading zip file into memory")
	}

	// bytes.Reader implements io.Reader, io.ReaderAt, etc. All you need!
	readerAt := bytes.NewReader(b)

	r, err := zip.NewReader(readerAt, readerAt.Size())
	if err != nil {
		return "", errors.Wrap(err, "error reading zip header")
	}

	// /usr/local/bin
	binPath := "/" + filepath.Join("usr", "local", "bin", writeToBin)

	for i := range r.File {
		hdr := r.File[i]

		if !strings.EqualFold(hdr.Name, findPattern) {
			continue
		}

		zp, err := hdr.Open()
		if err != nil {
			log.Fatal(err)
		}

		f, err := os.Create(binPath) // lint:allow_include_file
		if err != nil {
			// ~/bin
			binPath = filepath.Join(os.Getenv("HOME"), "bin", writeToBin)

			f, err = os.Create(binPath) // lint:allow_include_file
			if err != nil {
				return binPath, errors.Wrap(err, "failed to create file")
			}
		}

		err = f.Chmod(0o0755) // lint:allow_raw_number
		if err != nil {
			return binPath, errors.Wrap(err, "failed to make executable")
		}

		_, err = io.Copy(f, zp) // lint:allow_decompress
		if err != nil {
			return binPath, errors.Wrap(err, "error reading tar header")
		}

		err = f.Close()
		if err != nil {
			return binPath, errors.Wrap(err, "could not close the new file")
		}
	}

	return binPath, nil
}
