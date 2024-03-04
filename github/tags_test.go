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

const (
	oneOhOh    = "1.0.0"
	complexTag = "0.123.7"
)

func TestInvertTag(t *testing.T) {
	var tests = map[string]struct { // lint:no_dupe
		Input    string
		Expected string
	}{
		"1": {
			Input:    "1",
			Expected: "v1",
		},
		"v1": {
			Input:    "v1",
			Expected: "1",
		},
		oneOhOh: {
			Input:    oneOhOh,
			Expected: "v" + oneOhOh,
		},
		"v" + oneOhOh: {
			Input:    "v" + oneOhOh,
			Expected: oneOhOh,
		},
		complexTag: {
			Input:    complexTag,
			Expected: "v" + complexTag,
		},
		"v" + complexTag: {
			Input:    "v" + complexTag,
			Expected: complexTag,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actualTag := InvertTag(tc.Input)

			if actualTag != tc.Expected {
				t.Errorf("got %q; want %q", actualTag, tc.Expected)
			}
		})
	}
}

func TestRemoveVFromTag(t *testing.T) {
	var tests = map[string]struct { // lint:no_dupe
		Input    string
		Expected string
	}{
		"1": {
			Input:    "1",
			Expected: "1",
		},
		"v1": {
			Input:    "v1",
			Expected: "1",
		},
		oneOhOh: {
			Input:    oneOhOh,
			Expected: oneOhOh,
		},
		"v" + oneOhOh: {
			Input:    "v" + oneOhOh,
			Expected: oneOhOh,
		},
		complexTag: {
			Input:    complexTag,
			Expected: complexTag,
		},
		"v" + complexTag: {
			Input:    "v" + complexTag,
			Expected: complexTag,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actualTag := RemoveVFromTag(tc.Input)

			if actualTag != tc.Expected {
				t.Errorf("got %q; want %q", actualTag, tc.Expected)
			}
		})
	}
}
