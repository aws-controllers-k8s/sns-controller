// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package topic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	svcapitypes "github.com/aws-controllers-k8s/sns-controller/apis/v1alpha1"
)

func strPtr(s string) *string {
	return &s
}

func TestNewResourceDelta_DisplayName_NilEqualsZeroValue(t *testing.T) {
	tests := []struct {
		name            string
		desiredDisplay  *string
		latestDisplay   *string
		expectDifferent bool
	}{
		{
			name:            "both nil - no drift",
			desiredDisplay:  nil,
			latestDisplay:   nil,
			expectDifferent: false,
		},
		{
			name:            "desired nil, latest empty string - no drift (nil equals zero value)",
			desiredDisplay:  nil,
			latestDisplay:   strPtr(""),
			expectDifferent: false,
		},
		{
			name:            "desired empty string, latest nil - drift detected (one-way only)",
			desiredDisplay:  strPtr(""),
			latestDisplay:   nil,
			expectDifferent: true,
		},
		{
			name:            "desired nil, latest non-empty - drift detected",
			desiredDisplay:  nil,
			latestDisplay:   strPtr("my-topic"),
			expectDifferent: true,
		},
		{
			name:            "both same non-empty value - no drift",
			desiredDisplay:  strPtr("my-topic"),
			latestDisplay:   strPtr("my-topic"),
			expectDifferent: false,
		},
		{
			name:            "both empty string - no drift",
			desiredDisplay:  strPtr(""),
			latestDisplay:   strPtr(""),
			expectDifferent: false,
		},
		{
			name:            "different non-empty values - drift detected",
			desiredDisplay:  strPtr("old-name"),
			latestDisplay:   strPtr("new-name"),
			expectDifferent: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			assert := assert.New(t)

			desired := &resource{
				ko: &svcapitypes.Topic{
					Spec: svcapitypes.TopicSpec{
						DisplayName: tc.desiredDisplay,
					},
				},
			}
			latest := &resource{
				ko: &svcapitypes.Topic{
					Spec: svcapitypes.TopicSpec{
						DisplayName: tc.latestDisplay,
					},
				},
			}

			delta := newResourceDelta(desired, latest)
			require.NotNil(delta)

			if tc.expectDifferent {
				assert.True(
					delta.DifferentAt("Spec.DisplayName"),
					"expected drift at Spec.DisplayName but none detected",
				)
			} else {
				assert.False(
					delta.DifferentAt("Spec.DisplayName"),
					"expected no drift at Spec.DisplayName but drift was detected",
				)
			}
		})
	}
}
