// Copyright (c) 2026 Daniel Montanari. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitutil

import (
	"testing"
	"time"

	"github.com/go-git/go-git/v6/plumbing"
)

// TestVersionIncrementing tests the Major, Minor, and Patch incrementing logic.
func TestVersionIncrementing(t *testing.T) {
	// 1. Setup
	initialTag := TagInfo{
		Name:         "v1.2.3",
		Date:         time.Now(),
		Hash:         plumbing.ZeroHash,
		MajorVersion: 1,
		MinorVersion: 2,
		PatchVersion: 3,
	}

	// A GitTags struct with only one tag, which will be the "Newer" one.
	gitTags := GitTags{
		Tags: []TagInfo{initialTag},
	}

	// 2. Test Cases
	testCases := []struct {
		name     string
		action   func() string
		expected string
	}{
		{
			name:     "IncrementMajor",
			action:   gitTags.IncrementMajor,
			expected: "v2.0.0",
		},
		{
			name:     "IncrementMinor",
			action:   gitTags.IncrementMinor,
			expected: "v1.3.0",
		},
		{
			name:     "IncrementPatch",
			action:   gitTags.IncrementPatch,
			expected: "v1.2.4",
		},
	}

	// 3. Execution and Assertion
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.action()
			if got != tc.expected {
				t.Errorf("expected %q, but got %q", tc.expected, got)
			}
		})
	}
}

// TestNewerOlder ensures the correct tag is returned from a sorted list.
func TestNewerOlder(t *testing.T) {
	tagOlder := TagInfo{Name: "v1.0.0", Date: time.Now().Add(-1 * time.Hour)}
	tagNewer := TagInfo{Name: "v2.0.0", Date: time.Now()}

	// The NewGitTags function sorts tags by date descending (newest first).
	// So the first element should be the newest.
	gitTags := GitTags{
		Tags: []TagInfo{tagNewer, tagOlder},
	}

	if newer := gitTags.Newer(); newer.Name != "v2.0.0" {
		t.Errorf("expected Newer() to be 'v2.0.0', but got %q", newer.Name)
	}

	if older := gitTags.Older(); older.Name != "v1.0.0" {
		t.Errorf("expected Older() to be 'v1.0.0', but got %q", older.Name)
	}
}
