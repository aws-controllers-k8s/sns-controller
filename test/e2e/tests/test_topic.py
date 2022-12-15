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

from acktest.k8s import condition
from acktest.k8s import resource as k8s
from acktest.resources import random_suffix_name
from e2e import service_marker, CRD_GROUP, CRD_VERSION, load_resource
from e2e.common.types import TOPIC_RESOURCE_PLURAL
from e2e.replacement_values import REPLACEMENT_VALUES
from e2e import topic
from e2e import tag

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

#        new_display_name = "new display name"
#
#        # We're now going to modify the DisplayName field of the Role, wait
#        # some time and verify that the SNS server-side resource shows the new
#        # value of the field.
#        updates = {
#            "spec": {"displayName": new_display_name},
#        }
#        k8s.patch_custom_resource(ref, updates)
#        time.sleep(MODIFY_WAIT_AFTER_SECONDS)
#
#        latest = topic.get_attributes(topic_name)
#        assert latest is not None
#        assert 'Attributes' in latest
#        assert 'DisplayName' in latest
#        assert latest['DisplayName'] == new_display_name
#
#        # Same update code path check for tags...
#        latest_tags = topic.get_tags(topic_name)
#        before_update_expected_tags = [
#            {
#                "Key": "tag1",
#                "Value": "val1"
#            }
#        ]
#        assert tag.cleaned(latest_tags) == before_update_expected_tags
#        new_tags = [
#            {
#                "key": "tag2",
#                "value": "val2",
#            }
#        ]
#        updates = {
#            "spec": {"tags": new_tags},
#        }
#        k8s.patch_custom_resource(ref, updates)
#        time.sleep(MODIFY_WAIT_AFTER_SECONDS)
#
#        after_update_expected_tags = [
#            {
#                "Key": "tag2",
#                "Value": "val2",
#            }
#        ]
#        latest_tags = role.get_tags(role_name)
#        assert tag.cleaned(latest_tags) == after_update_expected_tags
