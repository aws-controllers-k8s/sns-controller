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

"""Integration tests for the SNS Topic resource"""

import time

import pytest
import boto3

from acktest.k8s import condition
from acktest.k8s import resource as k8s
from acktest.resources import random_suffix_name
from acktest import adoption as adoption
from acktest import tags
from e2e import service_marker, CRD_GROUP, CRD_VERSION, load_resource
from e2e.common.types import TOPIC_RESOURCE_KIND, TOPIC_RESOURCE_PLURAL
from e2e.replacement_values import REPLACEMENT_VALUES
from e2e import topic

DELETE_WAIT_AFTER_SECONDS = 10
MODIFY_WAIT_AFTER_SECONDS = 10


@pytest.fixture(scope="module")
def simple_topic():
    topic_name = random_suffix_name("my-simple-topic", 24)
    display_name  = "a simple topic"

    replacements = REPLACEMENT_VALUES.copy()
    replacements['TOPIC_NAME'] = topic_name
    replacements['DISPLAY_NAME'] = display_name

    resource_data = load_resource(
        "topic_simple",
        additional_replacements=replacements,
    )

    ref = k8s.CustomResourceReference(
        CRD_GROUP, CRD_VERSION, TOPIC_RESOURCE_PLURAL,
        topic_name, namespace="default",
    )
    k8s.create_custom_resource(ref, resource_data)
    cr = k8s.wait_resource_consumed_by_controller(ref)

    assert cr is not None
    assert k8s.get_resource_exists(ref)

    yield (ref, cr)

    _, deleted = k8s.delete_custom_resource(
        ref,
        period_length=DELETE_WAIT_AFTER_SECONDS,
    )
    assert deleted


@pytest.fixture(scope="module")
def fifo_topic():
    topic_name = random_suffix_name("my-fifo-topic", 16)
    # NOTE(jaypipes): FIFO Topics must have a name that ends in ".fifo"
    topic_name += ".fifo"
    display_name  = "a fifo topic"

    replacements = REPLACEMENT_VALUES.copy()
    replacements['TOPIC_NAME'] = topic_name
    replacements['DISPLAY_NAME'] = display_name

    resource_data = load_resource(
        "topic_fifo",
        additional_replacements=replacements,
    )

    ref = k8s.CustomResourceReference(
        CRD_GROUP, CRD_VERSION, TOPIC_RESOURCE_PLURAL,
        topic_name, namespace="default",
    )
    k8s.create_custom_resource(ref, resource_data)
    cr = k8s.wait_resource_consumed_by_controller(ref)

    assert cr is not None
    assert k8s.get_resource_exists(ref)

    yield (ref, cr)

    _, deleted = k8s.delete_custom_resource(
        ref,
        period_length=DELETE_WAIT_AFTER_SECONDS,
    )
    assert deleted


@service_marker
@pytest.mark.canary
class TestTopic:
    def test_crud(self, simple_topic):
        ref, res = simple_topic

        condition.assert_synced(ref)

        # Before we update the Topic CR below, let's check to see that the
        # DisplayName field in the CR is still what we set in the original
        # Create call.
        cr = k8s.get_resource(ref)
        assert cr is not None
        assert 'spec' in cr
        assert 'displayName' in cr['spec']
        assert cr['spec']['displayName'] == "a simple topic"

        assert 'status' in cr
        assert 'ackResourceMetadata' in cr['status']
        assert 'arn' in cr['status']['ackResourceMetadata']
        topic_arn = cr['status']['ackResourceMetadata']['arn']

        attrs = topic.get_attributes(topic_arn)
        assert attrs is not None
        assert 'DisplayName' in attrs
        assert attrs['DisplayName'] == "a simple topic"

        new_display_name = "new display name"

        # We're now going to modify the DisplayName field of the Topic, wait
        # some time and verify that the SNS server-side resource shows the new
        # value of the field.
        updates = {
            "spec": {"displayName": new_display_name},
        }
        k8s.patch_custom_resource(ref, updates)
        time.sleep(MODIFY_WAIT_AFTER_SECONDS)

        latest = topic.get_attributes(topic_arn)
        assert latest is not None
        assert 'DisplayName' in latest
        assert latest['DisplayName'] == new_display_name

        # Same update code path check for tags...
        latest_tags = topic.get_tags(topic_arn)
        expect_before_update_tags = [
            {
                "Key": "tag1",
                "Value": "val1"
            }
        ]
        tags.assert_equal_without_ack_tags(
            expect_before_update_tags, latest_tags,
        )
        new_tags = [
            {
                "key": "tag2",
                "value": "val2",
            }
        ]
        updates = {
            "spec": {"tags": new_tags},
        }
        k8s.patch_custom_resource(ref, updates)
        time.sleep(MODIFY_WAIT_AFTER_SECONDS)

        expect_after_update_tags = [
            {
                "Key": "tag2",
                "Value": "val2",
            }
        ]
        latest_tags = topic.get_tags(topic_arn)
        tags.assert_equal_without_ack_tags(
            expect_after_update_tags, latest_tags,
        )

        updates = {
            "spec": {"name": "my-simple-topic-edited"}
        }
        k8s.patch_custom_resource(ref, updates)
        time.sleep(MODIFY_WAIT_AFTER_SECONDS)
        k8s.wait_resource_consumed_by_controller(ref)
        condition.assert_type_status(ref, condition.CONDITION_TYPE_TERMINAL) 

        expected_msg = "Immutable Spec fields have been modified: Name"
        terminal_condition = k8s.get_resource_condition(ref, condition.CONDITION_TYPE_TERMINAL)
        # The name is immutable, testing if we get a terminal error
        assert expected_msg in terminal_condition['message']

    def test_crud_fifo(self, fifo_topic):
        ref, res = fifo_topic

        condition.assert_synced(ref)

        # Before we update the Topic CR below, let's check to see that the
        # DisplayName field in the CR is still what we set in the original
        # Create call.
        cr = k8s.get_resource(ref)
        assert cr is not None
        assert 'spec' in cr
        assert 'displayName' in cr['spec']
        assert cr['spec']['displayName'] == "a fifo topic"
        assert 'fifoTopic' in cr['spec']
        assert bool(cr['spec']['fifoTopic']) == True
        assert bool(cr['spec']['contentBasedDeduplication']) == True

        assert 'status' in cr
        assert 'ackResourceMetadata' in cr['status']
        assert 'arn' in cr['status']['ackResourceMetadata']
        topic_arn = cr['status']['ackResourceMetadata']['arn']

        attrs = topic.get_attributes(topic_arn)
        assert attrs is not None
        assert 'DisplayName' in attrs
        assert attrs['DisplayName'] == "a fifo topic"
        assert 'FifoTopic' in attrs
        assert bool(attrs['FifoTopic']) == True
        assert 'ContentBasedDeduplication' in attrs
        assert bool(attrs['ContentBasedDeduplication']) == True


class TestAdoptTopic(adoption.AbstractAdoptionTest):
    RESOURCE_PLURAL: str = TOPIC_RESOURCE_PLURAL
    RESOURCE_VERSION: str = CRD_VERSION

    _topic_name: str = random_suffix_name("ack-adopted-topic", 24)
    _topic_arn: str

    def bootstrap_resource(self):
        c = boto3.client('sns')
        resp = c.create_topic(Name=self._topic_name)
        self._topic_arn = resp['TopicArn']

    def cleanup_resource(self):
        client = boto3.client('sns')
        client.delete_topic(TopicArn=self._topic_arn)

    def get_resource_spec(self) -> adoption.AdoptedResourceSpec:
        return adoption.AdoptedResourceSpec(
            aws=adoption.AdoptedResourceARNIdentifier(additionalKeys={}, arn=self._topic_arn),
            kubernetes=adoption.AdoptedResourceKubernetesIdentifiers(CRD_GROUP, TOPIC_RESOURCE_KIND),
        )