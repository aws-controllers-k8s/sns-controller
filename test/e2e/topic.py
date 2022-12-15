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

"""Utilities for working with Topic resources"""

import datetime
import time

import boto3
import pytest

DEFAULT_WAIT_UNTIL_EXISTS_TIMEOUT_SECONDS = 90
DEFAULT_WAIT_UNTIL_EXISTS_INTERVAL_SECONDS = 15


def wait_until_exists(
        topic_arn: str,
        timeout_seconds: int = DEFAULT_WAIT_UNTIL_EXISTS_TIMEOUT_SECONDS,
        interval_seconds: int = DEFAULT_WAIT_UNTIL_EXISTS_INTERVAL_SECONDS,
    ) -> None:
    """Waits until a Topic with a supplied ARN is returned from SNS
    GetTopicAttributes API.

    Usage:
        from e2e.topic import wait_until_exists

        wait_until_exists(topic_arn)

    Raises:
        pytest.fail upon timeout
    """
    now = datetime.datetime.now()
    timeout = now + datetime.timedelta(seconds=timeout_seconds)

    while True:
        if datetime.datetime.now() >= timeout:
            pytest.fail(
                "Timed out waiting for Topic to exist "
                "in SNS API"
            )
        time.sleep(interval_seconds)

        latest = get_attributes(topic_arn)
        if latest is not None:
            break


def get_attributes(topic_arn):
    """Returns a dict containing the Topic attributes from the SNS
    GetTopicAttributes API.

    If no such Topic exists, returns None.
    """
    c = boto3.client('sns')
    try:
        resp = c.get_topic_attributes(TopicArn=topic_arn)
        return resp['Attributes']
    except c.exceptions.NotFoundException:
        return None


def get_tags(topic_arn):
    """Returns the tags for the topic with a supplied ARN.

    If no such Topic exists, returns None.
    """
    c = boto3.client('sns')
    try:
        resp = c.list_tags_for_resource(ResourceArn=topic_arn)
        return resp['Tags']
    except c.exceptions.ResourceNotFoundException:
        return None
