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

package topic

import (
	"bytes"
	"reflect"

	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	acktags "github.com/aws-controllers-k8s/runtime/pkg/tags"
)

// Hack to avoid import errors during build...
var (
	_ = &bytes.Buffer{}
	_ = &reflect.Method{}
	_ = &acktags.Tags{}
)

// newResourceDelta returns a new `ackcompare.Delta` used to compare two
// resources
func newResourceDelta(
	a *resource,
	b *resource,
) *ackcompare.Delta {
	delta := ackcompare.NewDelta()
	if (a == nil && b != nil) ||
		(a != nil && b == nil) {
		delta.Add("", a, b)
		return delta
	}

	if ackcompare.HasNilDifference(a.ko.Spec.DataProtectionPolicy, b.ko.Spec.DataProtectionPolicy) {
		delta.Add("Spec.DataProtectionPolicy", a.ko.Spec.DataProtectionPolicy, b.ko.Spec.DataProtectionPolicy)
	} else if a.ko.Spec.DataProtectionPolicy != nil && b.ko.Spec.DataProtectionPolicy != nil {
		if *a.ko.Spec.DataProtectionPolicy != *b.ko.Spec.DataProtectionPolicy {
			delta.Add("Spec.DataProtectionPolicy", a.ko.Spec.DataProtectionPolicy, b.ko.Spec.DataProtectionPolicy)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.DeliveryPolicy, b.ko.Spec.DeliveryPolicy) {
		delta.Add("Spec.DeliveryPolicy", a.ko.Spec.DeliveryPolicy, b.ko.Spec.DeliveryPolicy)
	} else if a.ko.Spec.DeliveryPolicy != nil && b.ko.Spec.DeliveryPolicy != nil {
		if *a.ko.Spec.DeliveryPolicy != *b.ko.Spec.DeliveryPolicy {
			delta.Add("Spec.DeliveryPolicy", a.ko.Spec.DeliveryPolicy, b.ko.Spec.DeliveryPolicy)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.DisplayName, b.ko.Spec.DisplayName) {
		delta.Add("Spec.DisplayName", a.ko.Spec.DisplayName, b.ko.Spec.DisplayName)
	} else if a.ko.Spec.DisplayName != nil && b.ko.Spec.DisplayName != nil {
		if *a.ko.Spec.DisplayName != *b.ko.Spec.DisplayName {
			delta.Add("Spec.DisplayName", a.ko.Spec.DisplayName, b.ko.Spec.DisplayName)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.KMSMasterKeyID, b.ko.Spec.KMSMasterKeyID) {
		delta.Add("Spec.KMSMasterKeyID", a.ko.Spec.KMSMasterKeyID, b.ko.Spec.KMSMasterKeyID)
	} else if a.ko.Spec.KMSMasterKeyID != nil && b.ko.Spec.KMSMasterKeyID != nil {
		if *a.ko.Spec.KMSMasterKeyID != *b.ko.Spec.KMSMasterKeyID {
			delta.Add("Spec.KMSMasterKeyID", a.ko.Spec.KMSMasterKeyID, b.ko.Spec.KMSMasterKeyID)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.Name, b.ko.Spec.Name) {
		delta.Add("Spec.Name", a.ko.Spec.Name, b.ko.Spec.Name)
	} else if a.ko.Spec.Name != nil && b.ko.Spec.Name != nil {
		if *a.ko.Spec.Name != *b.ko.Spec.Name {
			delta.Add("Spec.Name", a.ko.Spec.Name, b.ko.Spec.Name)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.Policy, b.ko.Spec.Policy) {
		delta.Add("Spec.Policy", a.ko.Spec.Policy, b.ko.Spec.Policy)
	} else if a.ko.Spec.Policy != nil && b.ko.Spec.Policy != nil {
		if *a.ko.Spec.Policy != *b.ko.Spec.Policy {
			delta.Add("Spec.Policy", a.ko.Spec.Policy, b.ko.Spec.Policy)
		}
	}

	return delta
}
