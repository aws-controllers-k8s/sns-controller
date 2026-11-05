# Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License"). You may
# not use this file except in compliance with the License. A copy of the
# License is located at
#
#     http://aws.amazon.com/apache2.0/
#
# or in the "license" file accompanying this file. This file is distributed
# on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
# express or implied. See the License for the specific language governing
# permissions and limitations under the License.

"""Integration tests for Topic maximum message size and late initialization."""

import time

import pytest

from acktest.k8s import condition
from acktest.k8s import resource as k8s
from acktest.resources import random_suffix_name
from e2e import CRD_GROUP, CRD_VERSION, load_resource, service_marker, topic
from e2e.common.types import TOPIC_RESOURCE_PLURAL
from e2e.replacement_values import REPLACEMENT_VALUES

DEFAULT_MAXIMUM_MESSAGE_SIZE = "262144"
MAXIMUM_MESSAGE_SIZE_WAIT_SECONDS = 90
MAXIMUM_MESSAGE_SIZE_POLL_SECONDS = 3
DELETE_WAIT_AFTER_SECONDS = 10


@pytest.fixture
def maximum_message_size_topic(request):
    topic_name = random_suffix_name("my-size-topic", 24)
    replacements = REPLACEMENT_VALUES.copy()
    replacements["TOPIC_NAME"] = topic_name
    resource_data = load_resource(
        "topic_maximum_message_size",
        additional_replacements=replacements,
    )
    if request.param is not None:
        resource_data["spec"]["maximumMessageSize"] = request.param

    ref = k8s.CustomResourceReference(
        CRD_GROUP, CRD_VERSION, TOPIC_RESOURCE_PLURAL,
        topic_name, namespace="default",
    )
    k8s.create_custom_resource(ref, resource_data)
    topic_arn = None
    try:
        cr = k8s.wait_resource_consumed_by_controller(ref)
        assert cr is not None
        topic_arn = cr["status"]["ackResourceMetadata"]["arn"]
        yield ref, topic_arn
    finally:
        _, deleted = k8s.delete_custom_resource(
            ref, period_length=DELETE_WAIT_AFTER_SECONDS,
        )
        assert deleted
        if topic_arn is not None:
            topic.wait_until_deleted(
                topic_arn,
                timeout_seconds=MAXIMUM_MESSAGE_SIZE_WAIT_SECONDS,
                interval_seconds=MAXIMUM_MESSAGE_SIZE_POLL_SECONDS,
            )


def wait_for_maximum_message_size(ref, topic_arn, expected):
    deadline = time.monotonic() + MAXIMUM_MESSAGE_SIZE_WAIT_SECONDS
    last_spec_value = None
    last_aws_value = None
    last_conditions = {}
    while time.monotonic() < deadline:
        cr = k8s.get_resource(ref)
        attributes = topic.get_attributes(topic_arn)
        if cr is not None and attributes is not None:
            last_spec_value = cr.get("spec", {}).get("maximumMessageSize")
            last_aws_value = attributes.get("MaximumMessageSize")
            last_conditions = {
                entry["type"]: entry["status"]
                for entry in cr.get("status", {}).get("conditions", [])
            }
            if expected is None:
                # SNS may omit the attribute while its default is effective.
                # If SNS reports the default, ACK may late initialize it.
                values_match = (
                    last_spec_value in (None, DEFAULT_MAXIMUM_MESSAGE_SIZE)
                    and attributes.get(
                        "MaximumMessageSize", DEFAULT_MAXIMUM_MESSAGE_SIZE,
                    ) == DEFAULT_MAXIMUM_MESSAGE_SIZE
                )
            else:
                values_match = (
                    last_spec_value == expected and last_aws_value == expected
                )
            if (
                values_match
                and last_conditions.get(condition.CONDITION_TYPE_RESOURCE_SYNCED) == "True"
                and last_conditions.get("ACK.LateInitialized") == "True"
            ):
                return cr, attributes
        time.sleep(MAXIMUM_MESSAGE_SIZE_POLL_SECONDS)
    pytest.fail(
        f"MaximumMessageSize did not converge to {expected!r}: "
        f"CR value={last_spec_value!r}, SNS value={last_aws_value!r}, "
        f"conditions={last_conditions!r}"
    )


@service_marker
@pytest.mark.canary
class TestTopicMaximumMessageSize:
    @pytest.mark.parametrize("maximum_message_size_topic", [None], indirect=True)
    def test_omitted_on_creation(self, maximum_message_size_topic):
        ref, topic_arn = maximum_message_size_topic
        _, attributes = wait_for_maximum_message_size(ref, topic_arn, None)
        assert attributes.get(
            "MaximumMessageSize", DEFAULT_MAXIMUM_MESSAGE_SIZE,
        ) == DEFAULT_MAXIMUM_MESSAGE_SIZE
        condition.assert_synced(ref)

    @pytest.mark.parametrize("maximum_message_size_topic", ["524288"], indirect=True)
    def test_create_update_remove_and_reset(self, maximum_message_size_topic):
        ref, topic_arn = maximum_message_size_topic

        # An explicit larger limit must be accepted by CreateTopic.
        wait_for_maximum_message_size(ref, topic_arn, "524288")

        # Increasing the limit must use the topic attribute update path.
        k8s.patch_custom_resource(ref, {"spec": {"maximumMessageSize": "1048576"}})
        wait_for_maximum_message_size(ref, topic_arn, "1048576")

        # Removing the field preserves AWS's limit and restores it in the CR.
        k8s.patch_custom_resource(ref, {"spec": {"maximumMessageSize": None}})
        wait_for_maximum_message_size(ref, topic_arn, "1048576")

        # Restoring the default requires explicitly requesting its numeric value.
        k8s.patch_custom_resource(
            ref, {"spec": {"maximumMessageSize": DEFAULT_MAXIMUM_MESSAGE_SIZE}},
        )
        wait_for_maximum_message_size(
            ref, topic_arn, DEFAULT_MAXIMUM_MESSAGE_SIZE,
        )
        condition.assert_synced(ref)
