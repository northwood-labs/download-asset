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

/*
Package github provides a library for downloading release assets from GitHub.
*/
package github

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"

	gh "github.com/google/go-github/v60/github"
	"github.com/mailgun/errors"
	"golang.org/x/oauth2"
)

var (
	oauthConf oauth2.Config

	ctx = context.Background()
)

type (
	NewClientInput struct {
		Endpoint string
		Token    string
	}
)

func NewClient(input *NewClientInput) (*gh.Client, error) {
	oauthClient := oauthConf.Client(ctx, &oauth2.Token{
		AccessToken: input.Token,
		TokenType:   "Bearer",
	})

	var (
		client *gh.Client
		err    error
	)

	apiEndpoint, uploadEndpoint, _ := ParseDomain(input.Endpoint)

	if input.Endpoint != "" {
		client, err = gh.NewClient(oauthClient).WithEnterpriseURLs(apiEndpoint, uploadEndpoint)
		if err != nil {
			return client, errors.Wrap(err, "failed to create new GitHub client")
		}
	} else {
		client = gh.NewClient(oauthClient)
	}

	return client, nil
}

func GetLatestRelease(client *gh.Client, owner, repo string) (*gh.RepositoryRelease, error) {
	release, _, err := client.Repositories.GetLatestRelease(ctx, owner, repo)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get latest release")
	}

	return release, nil
}

func GetReleaseVersion(client *gh.Client, owner, repo, tag string) (*gh.RepositoryRelease, error) {
	release, _, err := client.Repositories.GetReleaseByTag(ctx, owner, repo, tag)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get release by tag")
	}

	return release, nil
}

func GetAssetStream( // lint:allow_named_returns
	client *gh.Client,
	ownerRepo []string,
	release *gh.RepositoryRelease,
	pattern string,
) (archiveStream io.ReadCloser, name string, err error) {
	var i int

	for i = range release.Assets {
		asset := release.Assets[i]
		rePattern := regexp.MustCompile(pattern)

		if rePattern.MatchString(*asset.Name) {
			break
		}
	}

	asset := release.Assets[i]

	rc, _, err := client.Repositories.DownloadReleaseAsset(
		ctx,
		ownerRepo[0],
		ownerRepo[1],
		*asset.ID,
		http.DefaultClient,
	)

	return rc, asset.GetName(), err
}

func DownloadStream(archiveStream io.ReadCloser, filename, findPattern, writeToBin string) (string, error) {
	tmpDir, err := os.MkdirTemp("", filename+"-*")
	if err != nil {
		return tmpDir, errors.Wrap(err, "failed to create temp dir into which to download")
	}

	binPath, err := Decompress(archiveStream, filename, findPattern, writeToBin)
	if err != nil {
		log.Fatal(err)
	}

	return binPath, nil
}
