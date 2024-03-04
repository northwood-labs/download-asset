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
	"strings"

	"github.com/nlnwa/whatwg-url/canonicalizer"
)

func ParseDomain(domain string) ( // lint:allow_named_returns
	apiEndpoint,
	uploadEndpoint string,
	subdomainIsolation bool,
) {
	url, err := canonicalizer.GoogleSafeBrowsing.Parse(domain)
	if err != nil {
		return "", "", false
	}

	// Remove.
	url.SetHash("")
	url.SetPassword("")
	url.SetUsername("")
	url.SetSearch("")

	// Normalize to HTTPS if missing.
	if !strings.EqualFold(domain[0:4], "http") {
		url.SetProtocol("https")
	}

	// Handle subdomain isolation switching.
	if url.Hostname()[0:4] == "api." {
		subdomainIsolation = true

		url.SetPathname("")
		apiEndpoint = url.String()

		url.SetHostname("uploads." + url.Hostname()[4:])
		uploadEndpoint = url.String()
	} else {
		url.SetPathname("/api/v3")
		apiEndpoint = url.String()

		url.SetPathname("/uploads")
		uploadEndpoint = url.String()
	}

	// Remove trailing slashes.
	if apiEndpoint[len(apiEndpoint)-1:] == "/" {
		apiEndpoint = apiEndpoint[0 : len(apiEndpoint)-1]
	}

	if uploadEndpoint[len(uploadEndpoint)-1:] == "/" {
		uploadEndpoint = uploadEndpoint[0 : len(uploadEndpoint)-1]
	}

	return apiEndpoint, uploadEndpoint, subdomainIsolation
}
