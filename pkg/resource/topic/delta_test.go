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

func TestNewResourceDelta_ArchivePolicy_SemanticJSON(t *testing.T) {
	tests := []struct {
		name            string
		desiredPolicy   *string
		latestPolicy    *string
		expectDifferent bool
	}{
		{
			name:            "both nil - no drift",
			desiredPolicy:   nil,
			latestPolicy:    nil,
			expectDifferent: false,
		},
		{
			name:            "desired set, latest nil - drift detected",
			desiredPolicy:   strPtr(`{"MessageRetentionPeriod":"30"}`),
			latestPolicy:    nil,
			expectDifferent: true,
		},
		{
			name:            "identical string - no drift",
			desiredPolicy:   strPtr(`{"MessageRetentionPeriod":"30"}`),
			latestPolicy:    strPtr(`{"MessageRetentionPeriod":"30"}`),
			expectDifferent: false,
		},
		{
			// The key perpetual-reconcile scenario: SNS echoes the policy back
			// with extra whitespace. Must NOT be treated as drift.
			name:            "same JSON, different whitespace - no drift",
			desiredPolicy:   strPtr(`{"MessageRetentionPeriod":"30"}`),
			latestPolicy:    strPtr(`{ "MessageRetentionPeriod": "30" }`),
			expectDifferent: false,
		},
		{
			name:            "same JSON, different key ordering - no drift",
			desiredPolicy:   strPtr(`{"a":"1","b":"2"}`),
			latestPolicy:    strPtr(`{"b":"2","a":"1"}`),
			expectDifferent: false,
		},
		{
			name:            "semantically different JSON - drift detected",
			desiredPolicy:   strPtr(`{"MessageRetentionPeriod":"30"}`),
			latestPolicy:    strPtr(`{"MessageRetentionPeriod":"60"}`),
			expectDifferent: true,
		},
		{
			name:            "invalid JSON, identical strings - no drift",
			desiredPolicy:   strPtr("not-json"),
			latestPolicy:    strPtr("not-json"),
			expectDifferent: false,
		},
		{
			name:            "invalid JSON, different strings - drift detected",
			desiredPolicy:   strPtr("not-json-a"),
			latestPolicy:    strPtr("not-json-b"),
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
						ArchivePolicy: tc.desiredPolicy,
					},
				},
			}
			latest := &resource{
				ko: &svcapitypes.Topic{
					Spec: svcapitypes.TopicSpec{
						ArchivePolicy: tc.latestPolicy,
					},
				},
			}

			delta := newResourceDelta(desired, latest)
			require.NotNil(delta)

			if tc.expectDifferent {
				assert.True(
					delta.DifferentAt("Spec.ArchivePolicy"),
					"expected drift at Spec.ArchivePolicy but none detected",
				)
			} else {
				assert.False(
					delta.DifferentAt("Spec.ArchivePolicy"),
					"expected no drift at Spec.ArchivePolicy but drift was detected",
				)
			}
		})
	}
}
