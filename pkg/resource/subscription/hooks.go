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
	"context"

	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	svcsdk "github.com/aws/aws-sdk-go/service/sns"
)

var (
	// settableAttributesNames is a list of normalized CRD field names that may
	// be set in a SetSubscriptionAttribute API call
	settableAttributeNames = []string{
		"DeliveryPolicy",
		"FilterPolicy",
		"FilterPolicyScope",
		"RawMessageDelivery",
		"RedrivePolicy",
		"SubscriptionRoleARN",
	}
)

// customUpdate updates subscription attributes and tags. We require a custom
// update method because (unlike other SetAttributes APIs in SNS -- like
// PlatformApplication and PlatformEndpoint -- for subscriptions, you can only
// update a single attribute at a time. Yes, even though the name is
// SetSubscriptionAttributes (plural), you can only update one attribute at a
// time. :( This is why we can't have nice things.
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
		err = rm.setSubscriptionAttribute(ctx, desired, crdFieldName)
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

// setSubscriptionAttribute sets a single attribute for a subscription.
func (rm *resourceManager) setSubscriptionAttribute(
	ctx context.Context,
	desired *resource,
	crdFieldName string,
) error {
	var err error
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.setSubscriptionAttribute")
	defer func() {
		exit(err)
	}()
	input := rm.newSetAttributesRequestPayload(desired, crdFieldName)
	attrValue := ""
	if input.AttributeValue != nil {
		attrValue = *input.AttributeValue
	}
	rlog.Debug(
		"setting subscription attribute",
		"attr_name", crdFieldName,
		"attr_value", attrValue,
	)

	// NOTE(jaypipes): SetAttributes calls return a response but they don't
	// contain any useful information. Instead, below, we'll be returning a
	// DeepCopy of the supplied desired state, which should be fine because
	// that desired state has been constructed from a call to GetAttributes...
	_, respErr := rm.sdkapi.SetSubscriptionAttributesWithContext(ctx, input)
	rm.metrics.RecordAPICall("SET_ATTRIBUTES", "SetSubscriptionAttributes", respErr)
	if respErr != nil {
		if awsErr, ok := ackerr.AWSError(respErr); ok && awsErr.Code() == "NotFound" {
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
) *svcsdk.SetSubscriptionAttributesInput {
	res := &svcsdk.SetSubscriptionAttributesInput{}
	res.SetSubscriptionArn(string(*r.ko.Status.ACKResourceMetadata.ARN))
	switch crdFieldName {
	case "DeliveryPolicy":
		res.SetAttributeName("DeliveryPolicy")
		res.AttributeValue = r.ko.Spec.DeliveryPolicy
	case "FilterPolicy":
		res.SetAttributeName("FilterPolicy")
		res.AttributeValue = r.ko.Spec.FilterPolicy
	case "FilterPolicyScope":
		res.SetAttributeName("FilterPolicyScope")
		res.AttributeValue = r.ko.Spec.FilterPolicyScope
	case "RawMessageDelivery":
		res.SetAttributeName("RawMessageDelivery")
		res.AttributeValue = r.ko.Spec.RawMessageDelivery
	case "RedrivePolicy":
		res.SetAttributeName("RedrivePolicy")
		res.AttributeValue = r.ko.Spec.RedrivePolicy
	case "SubscriptionRoleARN":
		res.SetAttributeName("SubscriptionRoleArn")
		res.AttributeValue = r.ko.Spec.SubscriptionRoleARN
	}
	return res
}
