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

import "testing"

func TestSubdomain(t *testing.T) {
	// See https://bit.ly/3P1O9Rt for more information.
	// https://docs.github.com/en/enterprise-cloud@latest/rest/enterprise-admin
	// https://docs.github.com/en/enterprise-server@latest/rest/enterprise-admin
	var tests = map[string]struct { // lint:no_dupe
		InputURL           string
		APIEndpoint        string
		UploadEndpoint     string
		SubdomainIsolation bool
	}{
		"api.github.com": {
			InputURL:           "api.github.com",
			APIEndpoint:        "https://api.github.com",
			UploadEndpoint:     "https://uploads.github.com",
			SubdomainIsolation: true,
		},
		"https-api.github.com": {
			InputURL:           "api.github.com",
			APIEndpoint:        "https://api.github.com",
			UploadEndpoint:     "https://uploads.github.com",
			SubdomainIsolation: true,
		},
		"github.company.com": {
			InputURL:           "github.company.com",
			APIEndpoint:        "https://github.company.com/api/v3",
			UploadEndpoint:     "https://github.company.com/uploads",
			SubdomainIsolation: false,
		},
		"http-github.company.com": {
			InputURL:           "HTTP://github.company.com",
			APIEndpoint:        "http://github.company.com/api/v3",
			UploadEndpoint:     "http://github.company.com/uploads",
			SubdomainIsolation: false,
		},
		"https-github.company.com": {
			InputURL:           "HTTPs://github.company.com",
			APIEndpoint:        "https://github.company.com/api/v3",
			UploadEndpoint:     "https://github.company.com/uploads",
			SubdomainIsolation: false,
		},
		"https-github.company.com2": {
			InputURL:           "HTTPs://GITHUB.company.com",
			APIEndpoint:        "https://github.company.com/api/v3",
			UploadEndpoint:     "https://github.company.com/uploads",
			SubdomainIsolation: false,
		},
		"https-github.prod.company.com": {
			InputURL:           "HTTPs://GITHUB.prod.company.com",
			APIEndpoint:        "https://github.prod.company.com/api/v3",
			UploadEndpoint:     "https://github.prod.company.com/uploads",
			SubdomainIsolation: false,
		},
		"api.github.company.com": {
			InputURL:           "api.github.company.com",
			APIEndpoint:        "https://api.github.company.com",
			UploadEndpoint:     "https://uploads.github.company.com",
			SubdomainIsolation: true,
		},
		"http-api.github.company.com": {
			InputURL:           "HTTP://api.github.company.com",
			APIEndpoint:        "http://api.github.company.com",
			UploadEndpoint:     "http://uploads.github.company.com",
			SubdomainIsolation: true,
		},
		"https-api.github.company.com": {
			InputURL:           "HTTPs://api.github.company.com",
			APIEndpoint:        "https://api.github.company.com",
			UploadEndpoint:     "https://uploads.github.company.com",
			SubdomainIsolation: true,
		},
		"https-api.github.company.com2": {
			InputURL:           "HTTPs://API.GITHUB.company.com",
			APIEndpoint:        "https://api.github.company.com",
			UploadEndpoint:     "https://uploads.github.company.com",
			SubdomainIsolation: true,
		},
		"https-api.github.prod.company.com": {
			InputURL:           "HTTPs://api.GITHUB.prod.company.com",
			APIEndpoint:        "https://api.github.prod.company.com",
			UploadEndpoint:     "https://uploads.github.prod.company.com",
			SubdomainIsolation: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			apiEndpoint, uploadEndpoint, subdomainIsolation := ParseDomain(tc.InputURL)

			if apiEndpoint != tc.APIEndpoint {
				t.Errorf("APIEndpoint: got %q; want %q", apiEndpoint, tc.APIEndpoint)
			}

			if uploadEndpoint != tc.UploadEndpoint {
				t.Errorf("UploadEndpoint: got %q; want %q", uploadEndpoint, tc.UploadEndpoint)
			}

			if subdomainIsolation != tc.SubdomainIsolation {
				t.Errorf("SubdomainIsolation: got %v; want %v", subdomainIsolation, tc.SubdomainIsolation)
			}
		})
	}
}
