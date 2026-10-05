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
	"context"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"

	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	ackmetrics "github.com/aws-controllers-k8s/runtime/pkg/metrics"
	"github.com/aws/aws-sdk-go-v2/aws"
	svcsdk "github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	svcapitypes "github.com/aws-controllers-k8s/sns-controller/apis/v1alpha1"
)

// recordingHTTPClient records the form-encoded body of every SNS API call and
// answers with a successful response with an empty result.
type recordingHTTPClient struct {
	calls []url.Values
}

func (c *recordingHTTPClient) Do(req *http.Request) (*http.Response, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	form, err := url.ParseQuery(string(body))
	if err != nil {
		return nil, err
	}
	c.calls = append(c.calls, form)
	action := form.Get("Action")
	resp := "<" + action + "Response><" + action + "Result></" + action + "Result></" + action + "Response>"
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/xml"}},
		Body:       io.NopCloser(strings.NewReader(resp)),
		Request:    req,
	}, nil
}

func newTagsTestManager(httpClient *recordingHTTPClient) *resourceManager {
	return &resourceManager{
		metrics: ackmetrics.NewMetrics("sns"),
		sdkapi: svcsdk.New(svcsdk.Options{
			Region:      "us-west-2",
			Credentials: aws.AnonymousCredentials{},
			HTTPClient:  httpClient,
		}),
	}
}

func newTopicWithTags(tags map[string]string) *resource {
	arn := ackv1alpha1.AWSResourceName("arn:aws:sns:us-west-2:123456789012:topic")
	ko := &svcapitypes.Topic{}
	ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{ARN: &arn}
	for k, v := range tags {
		ko.Spec.Tags = append(ko.Spec.Tags, &svcapitypes.Tag{Key: strPtr(k), Value: strPtr(v)})
	}
	return &resource{ko: ko}
}

// tagsFromCall returns the Key->Value pairs of a TagResource call.
func tagsFromCall(call url.Values) map[string]string {
	res := map[string]string{}
	for i := 1; ; i++ {
		k := call.Get("Tags.member." + strconv.Itoa(i) + ".Key")
		if k == "" {
			return res
		}
		res[k] = call.Get("Tags.member." + strconv.Itoa(i) + ".Value")
	}
}

func TestSyncTags(t *testing.T) {
	tests := []struct {
		name           string
		desired        map[string]string
		latest         map[string]string
		expectTagged   map[string]string
		expectUntagged []string
	}{
		{
			name:    "no difference makes no API calls",
			desired: map[string]string{"a": "1"},
			latest:  map[string]string{"a": "1"},
		},
		{
			name:         "new tag is added",
			desired:      map[string]string{"a": "1", "b": "2"},
			latest:       map[string]string{"a": "1"},
			expectTagged: map[string]string{"b": "2"},
		},
		{
			name:           "removed tag is untagged",
			desired:        map[string]string{"a": "1"},
			latest:         map[string]string{"a": "1", "b": "2"},
			expectUntagged: []string{"b"},
		},
		{
			// Tag APIs upsert, so a changed value must only be re-tagged. If
			// it were also untagged (after being tagged), the topic would lose
			// the tag entirely.
			name:         "changed value is upserted and not untagged",
			desired:      map[string]string{"a": "2"},
			latest:       map[string]string{"a": "1"},
			expectTagged: map[string]string{"a": "2"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			httpClient := &recordingHTTPClient{}
			rm := newTagsTestManager(httpClient)

			err := rm.syncTags(
				context.Background(),
				newTopicWithTags(tc.desired),
				newTopicWithTags(tc.latest),
			)
			require.NoError(t, err)

			tagged := map[string]string{}
			untagged := []string{}
			for _, call := range httpClient.calls {
				switch call.Get("Action") {
				case "TagResource":
					for k, v := range tagsFromCall(call) {
						tagged[k] = v
					}
				case "UntagResource":
					for i := 1; call.Get("TagKeys.member."+strconv.Itoa(i)) != ""; i++ {
						untagged = append(untagged, call.Get("TagKeys.member."+strconv.Itoa(i)))
					}
				}
			}
			if tc.expectTagged == nil {
				tc.expectTagged = map[string]string{}
			}
			if tc.expectUntagged == nil {
				tc.expectUntagged = []string{}
			}
			assert.Equal(t, tc.expectTagged, tagged)
			assert.ElementsMatch(t, tc.expectUntagged, untagged)
		})
	}
}

func TestCompareTags(t *testing.T) {
	tests := []struct {
		name            string
		a, b            map[string]string
		expectDifferent bool
	}{
		{name: "both empty", expectDifferent: false},
		{name: "same tags", a: map[string]string{"a": "1"}, b: map[string]string{"a": "1"}},
		{name: "different value", a: map[string]string{"a": "1"}, b: map[string]string{"a": "2"}, expectDifferent: true},
		{name: "different key", a: map[string]string{"a": "1"}, b: map[string]string{"b": "1"}, expectDifferent: true},
		{name: "extra tag", a: map[string]string{"a": "1"}, b: map[string]string{"a": "1", "b": "2"}, expectDifferent: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			delta := newResourceDelta(newTopicWithTags(tc.a), newTopicWithTags(tc.b))
			assert.Equal(t, tc.expectDifferent, delta.DifferentAt("Spec.Tags"))
		})
	}
}
