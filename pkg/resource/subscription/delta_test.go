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

package subscription

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"

	svcapitypes "github.com/aws-controllers-k8s/sns-controller/apis/v1alpha1"
)

// helper to build a *resource with only FilterPolicy set.
func subscriptionWithFilterPolicy(filterPolicy *string) *resource {
	return &resource{
		ko: &svcapitypes.Subscription{
			Spec: svcapitypes.SubscriptionSpec{
				TopicARN:     aws.String("arn:aws:sns:us-east-1:123456789012:my-topic"),
				Protocol:     aws.String("sqs"),
				Endpoint:     aws.String("arn:aws:sqs:us-east-1:123456789012:my-queue"),
				FilterPolicy: filterPolicy,
			},
		},
	}
}

func TestNewResourceDelta_FilterPolicy_JSONComparison(t *testing.T) {
	tests := []struct {
		name     string
		desired  *string
		latest   *string
		wantDiff bool
	}{
		{
			name:     "both nil produces no diff",
			desired:  nil,
			latest:   nil,
			wantDiff: false,
		},
		{
			name:     "desired nil latest non-nil produces diff",
			desired:  nil,
			latest:   aws.String(`{"key":["value"]}`),
			wantDiff: true,
		},
		{
			name:     "desired non-nil latest nil produces diff",
			desired:  aws.String(`{"key":["value"]}`),
			latest:   nil,
			wantDiff: true,
		},
		{
			name:     "identical JSON strings produce no diff",
			desired:  aws.String(`{"notificationOrigin":["PRODUCT"]}`),
			latest:   aws.String(`{"notificationOrigin":["PRODUCT"]}`),
			wantDiff: false,
		},
		{
			name:     "pretty vs compact JSON produces no diff (issue 2877)",
			desired:  aws.String("{\n  \"notificationOrigin\": [\n    \"PRODUCT\"\n  ],\n  \"productType\": [\n    {\n      \"anything-but\": \"simple\"\n    }\n  ]\n}\n"),
			latest:   aws.String(`{"notificationOrigin":["PRODUCT"],"productType":[{"anything-but":"simple"}]}`),
			wantDiff: false,
		},
		{
			name:     "trailing newline from YAML block scalar produces no diff",
			desired:  aws.String("{\"notificationOrigin\":[\"PRODUCT\"]}\n"),
			latest:   aws.String(`{"notificationOrigin":["PRODUCT"]}`),
			wantDiff: false,
		},
		{
			name:     "different key ordering produces no diff",
			desired:  aws.String(`{"productType":[{"anything-but":"simple"}],"notificationOrigin":["PRODUCT"]}`),
			latest:   aws.String(`{"notificationOrigin":["PRODUCT"],"productType":[{"anything-but":"simple"}]}`),
			wantDiff: false,
		},
		{
			name:     "whitespace differences produce no diff",
			desired:  aws.String("{\n  \"notificationOrigin\": [\n    \"PRODUCT\"\n  ]\n}"),
			latest:   aws.String(`{"notificationOrigin":["PRODUCT"]}`),
			wantDiff: false,
		},
		{
			name:     "different filter values produce diff",
			desired:  aws.String(`{"notificationOrigin":["PRODUCT"]}`),
			latest:   aws.String(`{"notificationOrigin":["ORDER"]}`),
			wantDiff: true,
		},
		{
			name:     "extra filter key produces diff",
			desired:  aws.String(`{"notificationOrigin":["PRODUCT"]}`),
			latest:   aws.String(`{"notificationOrigin":["PRODUCT"],"productType":["simple"]}`),
			wantDiff: true,
		},
		{
			name:     "different nested structure produces diff",
			desired:  aws.String(`{"productType":[{"anything-but":"simple"}]}`),
			latest:   aws.String(`{"productType":["simple"]}`),
			wantDiff: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			desired := subscriptionWithFilterPolicy(tt.desired)
			latest := subscriptionWithFilterPolicy(tt.latest)

			delta := newResourceDelta(desired, latest)
			hasDiff := delta.DifferentAt("Spec.FilterPolicy")

			assert.Equal(t, tt.wantDiff, hasDiff,
				"DifferentAt(Spec.FilterPolicy) = %v, want %v", hasDiff, tt.wantDiff)
		})
	}
}
