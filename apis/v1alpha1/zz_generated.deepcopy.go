//go:build !ignore_autogenerated

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

// Code generated by ack-generate. DO NOT EDIT.

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	corev1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BatchResultErrorEntry) DeepCopyInto(out *BatchResultErrorEntry) {
	*out = *in
	if in.Code != nil {
		in, out := &in.Code, &out.Code
		*out = new(string)
		**out = **in
	}
	if in.ID != nil {
		in, out := &in.ID, &out.ID
		*out = new(string)
		**out = **in
	}
	if in.Message != nil {
		in, out := &in.Message, &out.Message
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BatchResultErrorEntry.
func (in *BatchResultErrorEntry) DeepCopy() *BatchResultErrorEntry {
	if in == nil {
		return nil
	}
	out := new(BatchResultErrorEntry)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Endpoint) DeepCopyInto(out *Endpoint) {
	*out = *in
	if in.Attributes != nil {
		in, out := &in.Attributes, &out.Attributes
		*out = make(map[string]*string, len(*in))
		for key, val := range *in {
			var outVal *string
			if val == nil {
				(*out)[key] = nil
			} else {
				inVal := (*in)[key]
				in, out := &inVal, &outVal
				*out = new(string)
				**out = **in
			}
			(*out)[key] = outVal
		}
	}
	if in.EndpointARN != nil {
		in, out := &in.EndpointARN, &out.EndpointARN
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Endpoint.
func (in *Endpoint) DeepCopy() *Endpoint {
	if in == nil {
		return nil
	}
	out := new(Endpoint)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MessageAttributeValue) DeepCopyInto(out *MessageAttributeValue) {
	*out = *in
	if in.DataType != nil {
		in, out := &in.DataType, &out.DataType
		*out = new(string)
		**out = **in
	}
	if in.StringValue != nil {
		in, out := &in.StringValue, &out.StringValue
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MessageAttributeValue.
func (in *MessageAttributeValue) DeepCopy() *MessageAttributeValue {
	if in == nil {
		return nil
	}
	out := new(MessageAttributeValue)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PhoneNumberInformation) DeepCopyInto(out *PhoneNumberInformation) {
	*out = *in
	if in.Status != nil {
		in, out := &in.Status, &out.Status
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PhoneNumberInformation.
func (in *PhoneNumberInformation) DeepCopy() *PhoneNumberInformation {
	if in == nil {
		return nil
	}
	out := new(PhoneNumberInformation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PlatformApplication) DeepCopyInto(out *PlatformApplication) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PlatformApplication.
func (in *PlatformApplication) DeepCopy() *PlatformApplication {
	if in == nil {
		return nil
	}
	out := new(PlatformApplication)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PlatformApplication) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PlatformApplicationList) DeepCopyInto(out *PlatformApplicationList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]PlatformApplication, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PlatformApplicationList.
func (in *PlatformApplicationList) DeepCopy() *PlatformApplicationList {
	if in == nil {
		return nil
	}
	out := new(PlatformApplicationList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PlatformApplicationList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PlatformApplicationSpec) DeepCopyInto(out *PlatformApplicationSpec) {
	*out = *in
	if in.EventDeliveryFailure != nil {
		in, out := &in.EventDeliveryFailure, &out.EventDeliveryFailure
		*out = new(string)
		**out = **in
	}
	if in.EventEndpointCreated != nil {
		in, out := &in.EventEndpointCreated, &out.EventEndpointCreated
		*out = new(string)
		**out = **in
	}
	if in.EventEndpointCreatedRef != nil {
		in, out := &in.EventEndpointCreatedRef, &out.EventEndpointCreatedRef
		*out = new(corev1alpha1.AWSResourceReferenceWrapper)
		(*in).DeepCopyInto(*out)
	}
	if in.EventEndpointDeleted != nil {
		in, out := &in.EventEndpointDeleted, &out.EventEndpointDeleted
		*out = new(string)
		**out = **in
	}
	if in.EventEndpointDeletedRef != nil {
		in, out := &in.EventEndpointDeletedRef, &out.EventEndpointDeletedRef
		*out = new(corev1alpha1.AWSResourceReferenceWrapper)
		(*in).DeepCopyInto(*out)
	}
	if in.EventEndpointUpdated != nil {
		in, out := &in.EventEndpointUpdated, &out.EventEndpointUpdated
		*out = new(string)
		**out = **in
	}
	if in.EventEndpointUpdatedRef != nil {
		in, out := &in.EventEndpointUpdatedRef, &out.EventEndpointUpdatedRef
		*out = new(corev1alpha1.AWSResourceReferenceWrapper)
		(*in).DeepCopyInto(*out)
	}
	if in.FailureFeedbackRoleARN != nil {
		in, out := &in.FailureFeedbackRoleARN, &out.FailureFeedbackRoleARN
		*out = new(string)
		**out = **in
	}
	if in.FailureFeedbackRoleRef != nil {
		in, out := &in.FailureFeedbackRoleRef, &out.FailureFeedbackRoleRef
		*out = new(corev1alpha1.AWSResourceReferenceWrapper)
		(*in).DeepCopyInto(*out)
	}
	if in.Name != nil {
		in, out := &in.Name, &out.Name
		*out = new(string)
		**out = **in
	}
	if in.Platform != nil {
		in, out := &in.Platform, &out.Platform
		*out = new(string)
		**out = **in
	}
	if in.PlatformCredential != nil {
		in, out := &in.PlatformCredential, &out.PlatformCredential
		*out = new(string)
		**out = **in
	}
	if in.PlatformPrincipal != nil {
		in, out := &in.PlatformPrincipal, &out.PlatformPrincipal
		*out = new(string)
		**out = **in
	}
	if in.SuccessFeedbackRoleARN != nil {
		in, out := &in.SuccessFeedbackRoleARN, &out.SuccessFeedbackRoleARN
		*out = new(string)
		**out = **in
	}
	if in.SuccessFeedbackRoleRef != nil {
		in, out := &in.SuccessFeedbackRoleRef, &out.SuccessFeedbackRoleRef
		*out = new(corev1alpha1.AWSResourceReferenceWrapper)
		(*in).DeepCopyInto(*out)
	}
	if in.SuccessFeedbackSampleRate != nil {
		in, out := &in.SuccessFeedbackSampleRate, &out.SuccessFeedbackSampleRate
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PlatformApplicationSpec.
func (in *PlatformApplicationSpec) DeepCopy() *PlatformApplicationSpec {
	if in == nil {
		return nil
	}
	out := new(PlatformApplicationSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PlatformApplicationStatus) DeepCopyInto(out *PlatformApplicationStatus) {
	*out = *in
	if in.ACKResourceMetadata != nil {
		in, out := &in.ACKResourceMetadata, &out.ACKResourceMetadata
		*out = new(corev1alpha1.ResourceMetadata)
		(*in).DeepCopyInto(*out)
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]*corev1alpha1.Condition, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(corev1alpha1.Condition)
				(*in).DeepCopyInto(*out)
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PlatformApplicationStatus.
func (in *PlatformApplicationStatus) DeepCopy() *PlatformApplicationStatus {
	if in == nil {
		return nil
	}
	out := new(PlatformApplicationStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PlatformApplication_SDK) DeepCopyInto(out *PlatformApplication_SDK) {
	*out = *in
	if in.Attributes != nil {
		in, out := &in.Attributes, &out.Attributes
		*out = make(map[string]*string, len(*in))
		for key, val := range *in {
			var outVal *string
			if val == nil {
				(*out)[key] = nil
			} else {
				inVal := (*in)[key]
				in, out := &inVal, &outVal
				*out = new(string)
				**out = **in
			}
			(*out)[key] = outVal
		}
	}
	if in.PlatformApplicationARN != nil {
		in, out := &in.PlatformApplicationARN, &out.PlatformApplicationARN
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PlatformApplication_SDK.
func (in *PlatformApplication_SDK) DeepCopy() *PlatformApplication_SDK {
	if in == nil {
		return nil
	}
	out := new(PlatformApplication_SDK)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PlatformEndpoint) DeepCopyInto(out *PlatformEndpoint) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PlatformEndpoint.
func (in *PlatformEndpoint) DeepCopy() *PlatformEndpoint {
	if in == nil {
		return nil
	}
	out := new(PlatformEndpoint)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PlatformEndpoint) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PlatformEndpointList) DeepCopyInto(out *PlatformEndpointList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]PlatformEndpoint, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PlatformEndpointList.
func (in *PlatformEndpointList) DeepCopy() *PlatformEndpointList {
	if in == nil {
		return nil
	}
	out := new(PlatformEndpointList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PlatformEndpointList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PlatformEndpointSpec) DeepCopyInto(out *PlatformEndpointSpec) {
	*out = *in
	if in.CustomUserData != nil {
		in, out := &in.CustomUserData, &out.CustomUserData
		*out = new(string)
		**out = **in
	}
	if in.Enabled != nil {
		in, out := &in.Enabled, &out.Enabled
		*out = new(string)
		**out = **in
	}
	if in.PlatformApplicationARN != nil {
		in, out := &in.PlatformApplicationARN, &out.PlatformApplicationARN
		*out = new(string)
		**out = **in
	}
	if in.Token != nil {
		in, out := &in.Token, &out.Token
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PlatformEndpointSpec.
func (in *PlatformEndpointSpec) DeepCopy() *PlatformEndpointSpec {
	if in == nil {
		return nil
	}
	out := new(PlatformEndpointSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PlatformEndpointStatus) DeepCopyInto(out *PlatformEndpointStatus) {
	*out = *in
	if in.ACKResourceMetadata != nil {
		in, out := &in.ACKResourceMetadata, &out.ACKResourceMetadata
		*out = new(corev1alpha1.ResourceMetadata)
		(*in).DeepCopyInto(*out)
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]*corev1alpha1.Condition, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(corev1alpha1.Condition)
				(*in).DeepCopyInto(*out)
			}
		}
	}
	if in.EndpointARN != nil {
		in, out := &in.EndpointARN, &out.EndpointARN
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PlatformEndpointStatus.
func (in *PlatformEndpointStatus) DeepCopy() *PlatformEndpointStatus {
	if in == nil {
		return nil
	}
	out := new(PlatformEndpointStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PublishBatchRequestEntry) DeepCopyInto(out *PublishBatchRequestEntry) {
	*out = *in
	if in.ID != nil {
		in, out := &in.ID, &out.ID
		*out = new(string)
		**out = **in
	}
	if in.MessageDeduplicationID != nil {
		in, out := &in.MessageDeduplicationID, &out.MessageDeduplicationID
		*out = new(string)
		**out = **in
	}
	if in.MessageGroupID != nil {
		in, out := &in.MessageGroupID, &out.MessageGroupID
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PublishBatchRequestEntry.
func (in *PublishBatchRequestEntry) DeepCopy() *PublishBatchRequestEntry {
	if in == nil {
		return nil
	}
	out := new(PublishBatchRequestEntry)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PublishBatchResultEntry) DeepCopyInto(out *PublishBatchResultEntry) {
	*out = *in
	if in.ID != nil {
		in, out := &in.ID, &out.ID
		*out = new(string)
		**out = **in
	}
	if in.SequenceNumber != nil {
		in, out := &in.SequenceNumber, &out.SequenceNumber
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PublishBatchResultEntry.
func (in *PublishBatchResultEntry) DeepCopy() *PublishBatchResultEntry {
	if in == nil {
		return nil
	}
	out := new(PublishBatchResultEntry)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Subscription) DeepCopyInto(out *Subscription) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Subscription.
func (in *Subscription) DeepCopy() *Subscription {
	if in == nil {
		return nil
	}
	out := new(Subscription)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Subscription) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SubscriptionList) DeepCopyInto(out *SubscriptionList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Subscription, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SubscriptionList.
func (in *SubscriptionList) DeepCopy() *SubscriptionList {
	if in == nil {
		return nil
	}
	out := new(SubscriptionList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SubscriptionList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SubscriptionSpec) DeepCopyInto(out *SubscriptionSpec) {
	*out = *in
	if in.DeliveryPolicy != nil {
		in, out := &in.DeliveryPolicy, &out.DeliveryPolicy
		*out = new(string)
		**out = **in
	}
	if in.Endpoint != nil {
		in, out := &in.Endpoint, &out.Endpoint
		*out = new(string)
		**out = **in
	}
	if in.FilterPolicy != nil {
		in, out := &in.FilterPolicy, &out.FilterPolicy
		*out = new(string)
		**out = **in
	}
	if in.FilterPolicyScope != nil {
		in, out := &in.FilterPolicyScope, &out.FilterPolicyScope
		*out = new(string)
		**out = **in
	}
	if in.Protocol != nil {
		in, out := &in.Protocol, &out.Protocol
		*out = new(string)
		**out = **in
	}
	if in.RawMessageDelivery != nil {
		in, out := &in.RawMessageDelivery, &out.RawMessageDelivery
		*out = new(string)
		**out = **in
	}
	if in.RedrivePolicy != nil {
		in, out := &in.RedrivePolicy, &out.RedrivePolicy
		*out = new(string)
		**out = **in
	}
	if in.SubscriptionRoleARN != nil {
		in, out := &in.SubscriptionRoleARN, &out.SubscriptionRoleARN
		*out = new(string)
		**out = **in
	}
	if in.TopicARN != nil {
		in, out := &in.TopicARN, &out.TopicARN
		*out = new(string)
		**out = **in
	}
	if in.TopicRef != nil {
		in, out := &in.TopicRef, &out.TopicRef
		*out = new(corev1alpha1.AWSResourceReferenceWrapper)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SubscriptionSpec.
func (in *SubscriptionSpec) DeepCopy() *SubscriptionSpec {
	if in == nil {
		return nil
	}
	out := new(SubscriptionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SubscriptionStatus) DeepCopyInto(out *SubscriptionStatus) {
	*out = *in
	if in.ACKResourceMetadata != nil {
		in, out := &in.ACKResourceMetadata, &out.ACKResourceMetadata
		*out = new(corev1alpha1.ResourceMetadata)
		(*in).DeepCopyInto(*out)
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]*corev1alpha1.Condition, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(corev1alpha1.Condition)
				(*in).DeepCopyInto(*out)
			}
		}
	}
	if in.ConfirmationWasAuthenticated != nil {
		in, out := &in.ConfirmationWasAuthenticated, &out.ConfirmationWasAuthenticated
		*out = new(string)
		**out = **in
	}
	if in.EffectiveDeliveryPolicy != nil {
		in, out := &in.EffectiveDeliveryPolicy, &out.EffectiveDeliveryPolicy
		*out = new(string)
		**out = **in
	}
	if in.Owner != nil {
		in, out := &in.Owner, &out.Owner
		*out = new(string)
		**out = **in
	}
	if in.PendingConfirmation != nil {
		in, out := &in.PendingConfirmation, &out.PendingConfirmation
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SubscriptionStatus.
func (in *SubscriptionStatus) DeepCopy() *SubscriptionStatus {
	if in == nil {
		return nil
	}
	out := new(SubscriptionStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Subscription_SDK) DeepCopyInto(out *Subscription_SDK) {
	*out = *in
	if in.Endpoint != nil {
		in, out := &in.Endpoint, &out.Endpoint
		*out = new(string)
		**out = **in
	}
	if in.Owner != nil {
		in, out := &in.Owner, &out.Owner
		*out = new(string)
		**out = **in
	}
	if in.Protocol != nil {
		in, out := &in.Protocol, &out.Protocol
		*out = new(string)
		**out = **in
	}
	if in.SubscriptionARN != nil {
		in, out := &in.SubscriptionARN, &out.SubscriptionARN
		*out = new(string)
		**out = **in
	}
	if in.TopicARN != nil {
		in, out := &in.TopicARN, &out.TopicARN
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Subscription_SDK.
func (in *Subscription_SDK) DeepCopy() *Subscription_SDK {
	if in == nil {
		return nil
	}
	out := new(Subscription_SDK)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Tag) DeepCopyInto(out *Tag) {
	*out = *in
	if in.Key != nil {
		in, out := &in.Key, &out.Key
		*out = new(string)
		**out = **in
	}
	if in.Value != nil {
		in, out := &in.Value, &out.Value
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Tag.
func (in *Tag) DeepCopy() *Tag {
	if in == nil {
		return nil
	}
	out := new(Tag)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Topic) DeepCopyInto(out *Topic) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Topic.
func (in *Topic) DeepCopy() *Topic {
	if in == nil {
		return nil
	}
	out := new(Topic)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Topic) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TopicList) DeepCopyInto(out *TopicList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Topic, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TopicList.
func (in *TopicList) DeepCopy() *TopicList {
	if in == nil {
		return nil
	}
	out := new(TopicList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *TopicList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TopicSpec) DeepCopyInto(out *TopicSpec) {
	*out = *in
	if in.ApplicationFailureFeedbackRoleARN != nil {
		in, out := &in.ApplicationFailureFeedbackRoleARN, &out.ApplicationFailureFeedbackRoleARN
		*out = new(string)
		**out = **in
	}
	if in.ApplicationSuccessFeedbackRoleARN != nil {
		in, out := &in.ApplicationSuccessFeedbackRoleARN, &out.ApplicationSuccessFeedbackRoleARN
		*out = new(string)
		**out = **in
	}
	if in.ApplicationSuccessFeedbackSampleRate != nil {
		in, out := &in.ApplicationSuccessFeedbackSampleRate, &out.ApplicationSuccessFeedbackSampleRate
		*out = new(string)
		**out = **in
	}
	if in.ContentBasedDeduplication != nil {
		in, out := &in.ContentBasedDeduplication, &out.ContentBasedDeduplication
		*out = new(string)
		**out = **in
	}
	if in.DataProtectionPolicy != nil {
		in, out := &in.DataProtectionPolicy, &out.DataProtectionPolicy
		*out = new(string)
		**out = **in
	}
	if in.DeliveryPolicy != nil {
		in, out := &in.DeliveryPolicy, &out.DeliveryPolicy
		*out = new(string)
		**out = **in
	}
	if in.DisplayName != nil {
		in, out := &in.DisplayName, &out.DisplayName
		*out = new(string)
		**out = **in
	}
	if in.FIFOTopic != nil {
		in, out := &in.FIFOTopic, &out.FIFOTopic
		*out = new(string)
		**out = **in
	}
	if in.FirehoseFailureFeedbackRoleARN != nil {
		in, out := &in.FirehoseFailureFeedbackRoleARN, &out.FirehoseFailureFeedbackRoleARN
		*out = new(string)
		**out = **in
	}
	if in.FirehoseSuccessFeedbackRoleARN != nil {
		in, out := &in.FirehoseSuccessFeedbackRoleARN, &out.FirehoseSuccessFeedbackRoleARN
		*out = new(string)
		**out = **in
	}
	if in.FirehoseSuccessFeedbackSampleRate != nil {
		in, out := &in.FirehoseSuccessFeedbackSampleRate, &out.FirehoseSuccessFeedbackSampleRate
		*out = new(string)
		**out = **in
	}
	if in.HTTPFailureFeedbackRoleARN != nil {
		in, out := &in.HTTPFailureFeedbackRoleARN, &out.HTTPFailureFeedbackRoleARN
		*out = new(string)
		**out = **in
	}
	if in.HTTPSuccessFeedbackRoleARN != nil {
		in, out := &in.HTTPSuccessFeedbackRoleARN, &out.HTTPSuccessFeedbackRoleARN
		*out = new(string)
		**out = **in
	}
	if in.HTTPSuccessFeedbackSampleRate != nil {
		in, out := &in.HTTPSuccessFeedbackSampleRate, &out.HTTPSuccessFeedbackSampleRate
		*out = new(string)
		**out = **in
	}
	if in.KMSMasterKeyID != nil {
		in, out := &in.KMSMasterKeyID, &out.KMSMasterKeyID
		*out = new(string)
		**out = **in
	}
	if in.KMSMasterKeyRef != nil {
		in, out := &in.KMSMasterKeyRef, &out.KMSMasterKeyRef
		*out = new(corev1alpha1.AWSResourceReferenceWrapper)
		(*in).DeepCopyInto(*out)
	}
	if in.LambdaFailureFeedbackRoleARN != nil {
		in, out := &in.LambdaFailureFeedbackRoleARN, &out.LambdaFailureFeedbackRoleARN
		*out = new(string)
		**out = **in
	}
	if in.LambdaSuccessFeedbackRoleARN != nil {
		in, out := &in.LambdaSuccessFeedbackRoleARN, &out.LambdaSuccessFeedbackRoleARN
		*out = new(string)
		**out = **in
	}
	if in.LambdaSuccessFeedbackSampleRate != nil {
		in, out := &in.LambdaSuccessFeedbackSampleRate, &out.LambdaSuccessFeedbackSampleRate
		*out = new(string)
		**out = **in
	}
	if in.Name != nil {
		in, out := &in.Name, &out.Name
		*out = new(string)
		**out = **in
	}
	if in.Policy != nil {
		in, out := &in.Policy, &out.Policy
		*out = new(string)
		**out = **in
	}
	if in.PolicyRef != nil {
		in, out := &in.PolicyRef, &out.PolicyRef
		*out = new(corev1alpha1.AWSResourceReferenceWrapper)
		(*in).DeepCopyInto(*out)
	}
	if in.SQSFailureFeedbackRoleARN != nil {
		in, out := &in.SQSFailureFeedbackRoleARN, &out.SQSFailureFeedbackRoleARN
		*out = new(string)
		**out = **in
	}
	if in.SQSSuccessFeedbackRoleARN != nil {
		in, out := &in.SQSSuccessFeedbackRoleARN, &out.SQSSuccessFeedbackRoleARN
		*out = new(string)
		**out = **in
	}
	if in.SQSSuccessFeedbackSampleRate != nil {
		in, out := &in.SQSSuccessFeedbackSampleRate, &out.SQSSuccessFeedbackSampleRate
		*out = new(string)
		**out = **in
	}
	if in.SignatureVersion != nil {
		in, out := &in.SignatureVersion, &out.SignatureVersion
		*out = new(string)
		**out = **in
	}
	if in.Tags != nil {
		in, out := &in.Tags, &out.Tags
		*out = make([]*Tag, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(Tag)
				(*in).DeepCopyInto(*out)
			}
		}
	}
	if in.TracingConfig != nil {
		in, out := &in.TracingConfig, &out.TracingConfig
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TopicSpec.
func (in *TopicSpec) DeepCopy() *TopicSpec {
	if in == nil {
		return nil
	}
	out := new(TopicSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TopicStatus) DeepCopyInto(out *TopicStatus) {
	*out = *in
	if in.ACKResourceMetadata != nil {
		in, out := &in.ACKResourceMetadata, &out.ACKResourceMetadata
		*out = new(corev1alpha1.ResourceMetadata)
		(*in).DeepCopyInto(*out)
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]*corev1alpha1.Condition, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(corev1alpha1.Condition)
				(*in).DeepCopyInto(*out)
			}
		}
	}
	if in.EffectiveDeliveryPolicy != nil {
		in, out := &in.EffectiveDeliveryPolicy, &out.EffectiveDeliveryPolicy
		*out = new(string)
		**out = **in
	}
	if in.Owner != nil {
		in, out := &in.Owner, &out.Owner
		*out = new(string)
		**out = **in
	}
	if in.TopicARN != nil {
		in, out := &in.TopicARN, &out.TopicARN
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TopicStatus.
func (in *TopicStatus) DeepCopy() *TopicStatus {
	if in == nil {
		return nil
	}
	out := new(TopicStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Topic_SDK) DeepCopyInto(out *Topic_SDK) {
	*out = *in
	if in.TopicARN != nil {
		in, out := &in.TopicARN, &out.TopicARN
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Topic_SDK.
func (in *Topic_SDK) DeepCopy() *Topic_SDK {
	if in == nil {
		return nil
	}
	out := new(Topic_SDK)
	in.DeepCopyInto(out)
	return out
}
