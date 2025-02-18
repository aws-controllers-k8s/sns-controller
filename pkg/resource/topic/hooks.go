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
	"errors"
	"fmt"

	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	"github.com/aws/aws-sdk-go-v2/aws"
	svcsdk "github.com/aws/aws-sdk-go-v2/service/sns"
	svcsdktypes "github.com/aws/aws-sdk-go-v2/service/sns/types"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/smithy-go"

	svcapitypes "github.com/aws-controllers-k8s/sns-controller/apis/v1alpha1"
	commonutil "github.com/aws-controllers-k8s/sns-controller/pkg/util"
)

var (
	// settableAttributesNames is a list of normalized CRD field names that may
	// be set in a SetTopicAttribute API call
	settableAttributeNames = []string{
		"DeliveryPolicy",
		"DisplayName",
		"Policy",
		"TracingConfig",
		"KMSMasterKeyID",
		"SignatureVersion",
		"ContentBasedDeduplication",
	}
)

// customUpdate updates topic attributes and tags. We require a custom update
// method because (unlike other SetAttributes APIs in SNS -- like
// PlatformApplication and PlatformEndpoint -- for topics, you can only update
// a single attribute at a time. Yes, even though the name is
// SetTopicAttributes (plural), you can only update one attribute at a time. :(
// This is why we can't have nice things.
func (rm *resourceManager) customUpdate(
	ctx context.Context,
	desired *resource,
	latest *resource,
	delta *ackcompare.Delta,
) (*resource, error) {
	var err error
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.customUpdate")
	defer func() {
		exit(err)
	}()
	if delta.DifferentAt("Spec.Tags") {
		if err := rm.syncTags(ctx, desired, latest); err != nil {
			return nil, err
		}
	}
	if !delta.DifferentExcept("Spec.Tags") {
		return desired, nil
	}

	// If any required fields in the input shape are missing, AWS resource is
	// not created yet. And sdkUpdate should never be called if this is the
	// case, and it's an error in the generated code if it is...
	if rm.requiredFieldsMissingFromSetAttributesInput(desired) {
		panic("Required field in SetAttributes input shape missing!")
	}

	for _, crdFieldName := range settableAttributeNames {
		if !delta.DifferentAt("Spec." + crdFieldName) {
			continue
		}
		err = rm.setTopicAttribute(ctx, desired, crdFieldName)
		if err != nil {
			return nil, err
		}
	}

	ko := desired.ko.DeepCopy()
	rm.setStatusDefaults(ko)
	return &resource{ko}, nil
}

// requiredFieldsMissingFromSetAtttributesInput returns true if there are any
// fields for the SetAttributes Input shape that are required by not present in
// the resource's Spec or Status
func (rm *resourceManager) requiredFieldsMissingFromSetAttributesInput(
	r *resource,
) bool {
	return (r.ko.Status.ACKResourceMetadata == nil || r.ko.Status.ACKResourceMetadata.ARN == nil)

}

// setTopicAttribute sets a single attribute for a topic.
func (rm *resourceManager) setTopicAttribute(
	ctx context.Context,
	desired *resource,
	crdFieldName string,
) error {
	var err error
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.setTopicAttribute")
	defer func() {
		exit(err)
	}()
	input := rm.newSetAttributesRequestPayload(desired, crdFieldName)
	attrValue := ""
	if input.AttributeValue != nil {
		attrValue = *input.AttributeValue
	}
	rlog.Debug(
		"setting topic attribute",
		"attr_name", crdFieldName,
		"attr_value", attrValue,
	)

	// NOTE(jaypipes): SetAttributes calls return a response but they don't
	// contain any useful information. Instead, below, we'll be returning a
	// DeepCopy of the supplied desired state, which should be fine because
	// that desired state has been constructed from a call to GetAttributes...
	_, respErr := rm.sdkapi.SetTopicAttributes(ctx, input)
	rm.metrics.RecordAPICall("SET_ATTRIBUTES", "SetTopicAttributes", respErr)
	if respErr != nil {
		var awsErr smithy.APIError
		if errors.As(err, &awsErr) && awsErr.ErrorCode() == "NotFound" {
			// Technically, this means someone deleted the backend resource in
			// between the time we got a result back from sdkFind() and here...
			return ackerr.NotFound
		}
		return respErr
	}
	return nil
}

// newSetAttributesRequestPayload returns SDK-specific struct for the HTTP
// request payload of the SetAttributes API call for the resource
func (rm *resourceManager) newSetAttributesRequestPayload(
	r *resource,
	crdFieldName string,
) *svcsdk.SetTopicAttributesInput {
	res := &svcsdk.SetTopicAttributesInput{}
	res.TopicArn = aws.String(string(*r.ko.Status.ACKResourceMetadata.ARN))
	switch crdFieldName {
	case "DeliveryPolicy":
		res.AttributeName = aws.String("DeliveryPolicy")
		res.AttributeValue = r.ko.Spec.DeliveryPolicy
	case "DisplayName":
		res.AttributeName = aws.String("DisplayName")
		res.AttributeValue = r.ko.Spec.DisplayName
	case "Policy":
		res.AttributeName = aws.String("Policy")
		res.AttributeValue = r.ko.Spec.Policy
	case "TracingConfig":
		res.AttributeName = aws.String("TracingConfig")
		res.AttributeValue = r.ko.Spec.TracingConfig
	case "KMSMasterKeyID":
		res.AttributeName = aws.String("KmsMasterKeyId")
		res.AttributeValue = r.ko.Spec.KMSMasterKeyID
	case "SignatureVersion":
		res.AttributeName = aws.String("SignatureVersion")
		res.AttributeValue = r.ko.Spec.SignatureVersion
	case "ContentBasedDeduplication":
		res.AttributeName = aws.String("ContentBasedDeduplication")
		res.AttributeValue = r.ko.Spec.ContentBasedDeduplication
	}
	return res
}

// compareTags is a custom comparison function for comparing lists of Tag
// structs where the order of the structs in the list is not important.
func compareTags(
	delta *ackcompare.Delta,
	a *resource,
	b *resource,
) {
	if len(a.ko.Spec.Tags) != len(b.ko.Spec.Tags) {
		delta.Add("Spec.Tags", a.ko.Spec.Tags, b.ko.Spec.Tags)
	} else if len(a.ko.Spec.Tags) > 0 {
		if !commonutil.EqualTags(a.ko.Spec.Tags, b.ko.Spec.Tags) {
			delta.Add("Spec.Tags", a.ko.Spec.Tags, b.ko.Spec.Tags)
		}
	}
}

// syncTags examines the Tags in the supplied Topic and calls the
// ListTagsForResource, TagResource and UntagResource APIs to ensure that the
// set of associated Tags  stays in sync with the Topic.Spec.Tags
func (rm *resourceManager) syncTags(
	ctx context.Context,
	desired *resource,
	latest *resource,
) (err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.syncTags")
	defer func() { exit(err) }()
	toAdd := []*svcapitypes.Tag{}
	toDelete := []*svcapitypes.Tag{}

	existingTags := latest.ko.Spec.Tags

	for _, t := range desired.ko.Spec.Tags {
		if !inTags(*t.Key, *t.Value, existingTags) {
			toAdd = append(toAdd, t)
		}
	}

	for _, t := range existingTags {
		if !inTags(*t.Key, *t.Value, desired.ko.Spec.Tags) {
			toDelete = append(toDelete, t)
		}
	}

	if len(toAdd) > 0 {
		for _, t := range toAdd {
			rlog.Debug("adding tag to topic", "key", *t.Key, "value", *t.Value)
		}
		if err = rm.addTags(ctx, desired, toAdd); err != nil {
			return err
		}
	}
	if len(toDelete) > 0 {
		for _, t := range toDelete {
			rlog.Debug("removing tag from topic", "key", *t.Key, "value", *t.Value)
		}
		if err = rm.removeTags(ctx, desired, toDelete); err != nil {
			return err
		}
	}

	return nil
}

// inTags returns true if the supplied key and value can be found in the
// supplied list of Tag structs.
//
// TODO(jaypipes): When we finally standardize Tag handling in ACK, move this
// to the ACK common runtime/ or pkg/ repos
func inTags(
	key string,
	value string,
	tags []*svcapitypes.Tag,
) bool {
	for _, t := range tags {
		if *t.Key == key && *t.Value == value {
			return true
		}
	}
	return false
}

// getTags returns the list of tags to the Topic
func (rm *resourceManager) getTags(
	ctx context.Context,
	r *resource,
) ([]*svcapitypes.Tag, error) {
	var err error
	var resp *svcsdk.ListTagsForResourceOutput
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.getTags")
	defer func() { exit(err) }()

	arn := (*string)(r.ko.Status.ACKResourceMetadata.ARN)

	input := &svcsdk.ListTagsForResourceInput{}
	input.ResourceArn = arn
	res := []*svcapitypes.Tag{}

	// NOTE(jaypipes): Unlike other list tag APIs, SNS' ListTagsForResource is
	// not paginated...
	resp, err = rm.sdkapi.ListTagsForResource(ctx, input)
	if err != nil || resp == nil {
		return nil, err
	}
	for _, t := range resp.Tags {
		res = append(res, &svcapitypes.Tag{Key: t.Key, Value: t.Value})
	}
	rm.metrics.RecordAPICall("READ_MANY", "ListTagsForResource", err)
	return res, err
}

// addTags adds the supplied Tags to the supplied Topic resource
func (rm *resourceManager) addTags(
	ctx context.Context,
	r *resource,
	tags []*svcapitypes.Tag,
) (err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.addTag")
	defer func() { exit(err) }()

	arn := (*string)(r.ko.Status.ACKResourceMetadata.ARN)

	input := &svcsdk.TagResourceInput{}
	input.ResourceArn = arn
	inTags := []svcsdktypes.Tag{}
	for _, t := range tags {
		inTags = append(inTags, svcsdktypes.Tag{Key: t.Key, Value: t.Value})
	}
	input.Tags = inTags

	_, err = rm.sdkapi.TagResource(ctx, input)
	rm.metrics.RecordAPICall("UPDATE", "TagResource", err)
	return err
}

// removeTags removes the supplied Tags from the supplied Topic resource
func (rm *resourceManager) removeTags(
	ctx context.Context,
	r *resource,
	tags []*svcapitypes.Tag,
) (err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.removeTag")
	defer func() { exit(err) }()

	arn := (*string)(r.ko.Status.ACKResourceMetadata.ARN)

	input := &svcsdk.UntagResourceInput{}
	input.ResourceArn = arn
	inTagKeys := []string{}
	for _, t := range tags {
		inTagKeys = append(inTagKeys, *t.Key)
	}
	input.TagKeys = inTagKeys

	_, err = rm.sdkapi.UntagResource(ctx, input)
	rm.metrics.RecordAPICall("UPDATE", "UntagResource", err)
	return err
}

// getTopicNameFromARN retrieves the topic name from the topic arn
func (rm *resourceManager) getTopicNameFromARN(tmpARN ackv1alpha1.AWSResourceName) (string, error) {
	topicARN, err := arn.Parse(string(tmpARN))
	if err != nil {
		return "", fmt.Errorf("error parsing topic ARN: %s, error: %w", tmpARN, err)
	}
	return topicARN.Resource, nil
}
