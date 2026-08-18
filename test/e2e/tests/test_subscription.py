# Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License"). You may
# not use this file except in compliance with the License. A copy of the
# License is located at
#
#	 http://aws.amazon.com/apache2.0/
#
# or in the "license" file accompanying this file. This file is distributed
# on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
# express or implied. See the License for the specific language governing
# permissions and limitations under the License.

"""Integration tests for the SNS Subscription resource"""

import json
import time

import pytest
import boto3

from acktest.k8s import condition
from acktest.k8s import resource as k8s
from acktest.resources import random_suffix_name
from acktest import adoption
from e2e import service_marker, CRD_GROUP, CRD_VERSION, load_resource
from e2e.bootstrap_resources import get_bootstrap_resources
from e2e.common.types import SUBSCRIPTION_RESOURCE_KIND, SUBSCRIPTION_RESOURCE_PLURAL
from e2e.replacement_values import REPLACEMENT_VALUES
from e2e import subscription

MODIFY_WAIT_AFTER_SECONDS = 10
CHECK_WAIT_AFTER_REF_RESOLVE_SECONDS = 10
DELETE_SUBSCRIPTION_TIMEOUT_SECONDS = 10

# PendingConfirmation is unknown until the first read-back, so every
# Subscription needs one requeue before it can report Synced=True.
SYNCED_WAIT_PERIODS = 6
SYNCED_WAIT_PERIOD_LENGTH_SECONDS = 15

# Spans at least one requeue of a not-synced resource.
REQUEUE_OBSERVE_WAIT_SECONDS = 45


@pytest.fixture(scope="module")
def subscription_sqs():
    subscription_name = random_suffix_name("subscription-sqs", 24)
    display_name  = "a subscription to a queue"

    boot_resources = get_bootstrap_resources()
    q = boot_resources.Queue1
    topic = boot_resources.Topic1

    replacements = REPLACEMENT_VALUES.copy()
    replacements['SUBSCRIPTION_NAME'] = subscription_name
    replacements['TOPIC_ARN'] = topic.arn
    replacements['ENDPOINT'] = q.arn

    resource_data = load_resource(
        "subscription_with_refs",
        additional_replacements=replacements,
    )

    ref = k8s.CustomResourceReference(
        CRD_GROUP, CRD_VERSION, SUBSCRIPTION_RESOURCE_PLURAL,
        subscription_name, namespace="default",
    )
    k8s.create_custom_resource(ref, resource_data)
    cr = k8s.wait_resource_consumed_by_controller(ref)

    assert cr is not None
    assert k8s.get_resource_exists(ref)
    # NOTE(jaypipes): This works because we manually override the
    # ReturnSubscriptionArn field in SubscribeInput to "true"
    assert 'status' in cr
    assert 'ackResourceMetadata' in cr['status']
    assert 'arn' in cr['status']['ackResourceMetadata']
    sub_arn = cr['status']['ackResourceMetadata']['arn']

    yield (ref, cr, sub_arn)

    _, deleted = k8s.delete_custom_resource(
        ref,
        period_length=DELETE_SUBSCRIPTION_TIMEOUT_SECONDS,
    )
    assert deleted

    subscription.wait_until_deleted(sub_arn)


@service_marker
@pytest.mark.canary
class TestSubscription:
    def test_crud(self, subscription_sqs):
        sub_ref, sub_cr, sub_arn = subscription_sqs

        subscription.wait_until_exists(sub_arn)

        # sqs is confirmed by SNS on our behalf, but only reaches Synced=True
        # once PendingConfirmation has been read back.
        assert k8s.wait_on_condition(
            sub_ref,
            condition.CONDITION_TYPE_RESOURCE_SYNCED,
            "True",
            wait_periods=SYNCED_WAIT_PERIODS,
            period_length=SYNCED_WAIT_PERIOD_LENGTH_SECONDS,
        )

        cr = k8s.get_resource(sub_ref)
        assert cr['status']['pendingConfirmation'] == "false"

        # Before we update the Topic CR below, let's check to see that the
        # DisplayName field in the CR is still what we set in the original
        # Create call.
        cr = k8s.get_resource(sub_ref)
        assert cr is not None
        assert 'spec' in cr
        assert 'deliveryPolicy' not in cr['spec']

        attrs = subscription.get_attributes(sub_arn)
        assert attrs is not None
        assert 'DeliveryPolicy' not in attrs

        delivery_policy = {
            "healthyRetryPolicy": {
                "minDelayTarget": 1,
                "maxDelayTarget": 60,
                "numRetries": 50,
                "numNoDelayRetries": 3,
                "numMinDelayRetries": 2,
                "numMaxDelayRetries": 35,
                "backoffFunction": "exponential"
            }
        }
        new_delivery_policy = json.dumps(delivery_policy)

        # We're now going to modify the DeliveryPolicy field of the
        # Subscription, wait some time and verify that the SNS server-side
        # resource shows the new value of the field.
        updates = {
            "spec": {"deliveryPolicy": new_delivery_policy},
        }
        k8s.patch_custom_resource(sub_ref, updates)
        time.sleep(MODIFY_WAIT_AFTER_SECONDS)

        latest = subscription.get_attributes(sub_arn)
        assert latest is not None
        assert 'DeliveryPolicy' in latest

        # NOTE(jaypipes): SNS adds some default field values to the
        # DeliveryPolicy JSON object on the server-side, including things like
        # `"guaranteed": false` and `"requestPolicy": null`. We will simply
        # verify that the healthRetryPolicy segment we updated is correct.
        got_delivery_policy= json.loads(latest['DeliveryPolicy'])
        assert 'healthyRetryPolicy' in got_delivery_policy
        exp_healthy_retry_policy = delivery_policy['healthyRetryPolicy']
        assert exp_healthy_retry_policy == got_delivery_policy['healthyRetryPolicy']

        # Verify semantic JSON comparison (is_document): patch with the same
        # logical JSON but different key ordering. With DocumentEqual, this
        # should NOT trigger a reconciliation loop or unnecessary update.
        reordered_policy = {
            "healthyRetryPolicy": {
                "backoffFunction": "exponential",
                "numMaxDelayRetries": 35,
                "numMinDelayRetries": 2,
                "numNoDelayRetries": 3,
                "numRetries": 50,
                "maxDelayTarget": 60,
                "minDelayTarget": 1
            }
        }
        updates = {
            "spec": {"deliveryPolicy": json.dumps(reordered_policy)},
        }
        k8s.patch_custom_resource(sub_ref, updates)
        time.sleep(MODIFY_WAIT_AFTER_SECONDS)

        # Should still be synced — no unnecessary update triggered
        condition.assert_synced(sub_ref)


@pytest.fixture(scope="module")
def subscription_pending():
    """Creates a Subscription that stays unconfirmed.

    An email endpoint is not reachability-checked at Subscribe time (unlike
    http/https, which SNS rejects outright with InvalidParameter if the endpoint
    does not answer its confirmation request), and example.com is reserved by
    RFC 2606, so nobody ever confirms it.
    """
    subscription_name = random_suffix_name("subscription-pending", 28)
    boot_resources = get_bootstrap_resources()
    topic = boot_resources.Topic2

    replacements = REPLACEMENT_VALUES.copy()
    replacements['SUBSCRIPTION_NAME'] = subscription_name
    replacements['TOPIC_ARN'] = topic.arn
    replacements['ENDPOINT'] = f"{subscription_name}@example.com"

    resource_data = load_resource(
        "subscription_pending",
        additional_replacements=replacements,
    )

    ref = k8s.CustomResourceReference(
        CRD_GROUP, CRD_VERSION, SUBSCRIPTION_RESOURCE_PLURAL,
        subscription_name, namespace="default",
    )
    k8s.create_custom_resource(ref, resource_data)
    cr = k8s.wait_resource_consumed_by_controller(ref)

    assert cr is not None
    # A real ARN comes back despite being unconfirmed because
    # ReturnSubscriptionArn is overridden to true in SubscribeInput.
    assert 'status' in cr
    assert 'ackResourceMetadata' in cr['status']
    assert 'arn' in cr['status']['ackResourceMetadata']
    sub_arn = cr['status']['ackResourceMetadata']['arn']

    yield (ref, sub_arn)

    # The resource carries deletion-policy: retain, so the CR is removed without
    # calling Unsubscribe -- SNS rejects that for a pending subscription. The
    # unconfirmed subscription is left for SNS to reap after 48h, so there is no
    # wait_until_deleted here.
    _, deleted = k8s.delete_custom_resource(
        ref,
        period_length=DELETE_SUBSCRIPTION_TIMEOUT_SECONDS,
    )
    assert deleted


@service_marker
class TestSubscriptionPendingConfirmation:
    def test_pending_confirmation(self, subscription_pending):
        """An unconfirmed Subscription reports Synced=False and keeps requeueing
        so a later confirmation is noticed.
        """
        sub_ref, sub_arn = subscription_pending

        subscription.wait_until_exists(sub_arn)

        attrs = subscription.get_attributes(sub_arn)
        assert attrs is not None
        assert attrs['PendingConfirmation'] == "true"

        assert k8s.wait_on_condition(
            sub_ref,
            condition.CONDITION_TYPE_RESOURCE_SYNCED,
            "False",
            wait_periods=SYNCED_WAIT_PERIODS,
            period_length=SYNCED_WAIT_PERIOD_LENGTH_SECONDS,
        )

        cr = k8s.get_resource(sub_ref)
        assert cr is not None
        assert cr['status']['pendingConfirmation'] == "true"

        # Awaiting confirmation is a recoverable wait, not terminal.
        terminal = k8s.get_resource_condition(
            sub_ref, condition.CONDITION_TYPE_TERMINAL,
        )
        assert terminal is None or terminal['status'] != "True"

        # ACK rewrites lastTransitionTime on every reconcile, so a newer
        # timestamp proves the resource is still being requeued.
        before = condition.get_synced_last_transition_time(sub_ref)
        assert before is not None
        time.sleep(REQUEUE_OBSERVE_WAIT_SECONDS)
        after = condition.get_synced_last_transition_time(sub_ref)
        assert after is not None
        assert after > before, (
            "expected the controller to requeue and re-reconcile the pending "
            f"subscription, but ACK.ResourceSynced still reads {before}"
        )
        condition.assert_not_synced(sub_ref)
