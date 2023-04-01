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

"""Utilities for working with Subscription resources"""

import datetime
import time

import boto3
import pytest

DEFAULT_WAIT_UNTIL_EXISTS_TIMEOUT_SECONDS = 90
DEFAULT_WAIT_UNTIL_EXISTS_INTERVAL_SECONDS = 15
DEFAULT_WAIT_UNTIL_DELETED_TIMEOUT_SECONDS = 60*10
DEFAULT_WAIT_UNTIL_DELETED_INTERVAL_SECONDS = 15


def wait_until_exists(
        subscription_arn: str,
        timeout_seconds: int = DEFAULT_WAIT_UNTIL_EXISTS_TIMEOUT_SECONDS,
        interval_seconds: int = DEFAULT_WAIT_UNTIL_EXISTS_INTERVAL_SECONDS,
    ) -> None:
    """Waits until a Subscription with a supplied ARN is returned from SNS
    GetSubscriptionAttributes API.

    Usage:
        from e2e.subscription import wait_until_exists

        wait_until_exists(subscription_arn)

    Raises:
        pytest.fail upon timeout
    """
    now = datetime.datetime.now()
    timeout = now + datetime.timedelta(seconds=timeout_seconds)

    while True:
        if datetime.datetime.now() >= timeout:
            pytest.fail(
                "Timed out waiting for Subscription to exist "
                "in SNS API"
            )
        time.sleep(interval_seconds)

        latest = get_attributes(subscription_arn)
        if latest is not None:
            break


def wait_until_deleted(
        sub_arn: str,
        timeout_seconds: int = DEFAULT_WAIT_UNTIL_DELETED_TIMEOUT_SECONDS,
        interval_seconds: int = DEFAULT_WAIT_UNTIL_DELETED_INTERVAL_SECONDS,
    ) -> None:
    """Waits until a Subscription with a supplied ARN is no longer returned from
    the SNS API.

    Usage:
        from e2e.subscription import wait_until_deleted

        wait_until_deleted(sub_arn)

    Raises:
        pytest.fail upon timeout
    """
    now = datetime.datetime.now()
    timeout = now + datetime.timedelta(seconds=timeout_seconds)

    while True:
        if datetime.datetime.now() >= timeout:
            pytest.fail(
                "Timed out waiting for Subscription to be "
                "deleted in SNS API"
            )
        time.sleep(interval_seconds)

        latest = get_attributes(sub_arn)
        if latest is None:
            break


def get_attributes(subscription_arn):
    """Returns a dict containing the Subscription attributes from the SNS
    GetSubscriptionAttributes API.

    If no such Subscription exists, returns None.
    """
    c = boto3.client('sns')
    try:
        resp = c.get_subscription_attributes(SubscriptionArn=subscription_arn)
        return resp['Attributes']
    except c.exceptions.NotFoundException:
        return None


def get_tags(subscription_arn):
    """Returns the tags for the subscription with a supplied ARN.

    If no such Subscription exists, returns None.
    """
    c = boto3.client('sns')
    try:
        resp = c.list_tags_for_resource(ResourceArn=subscription_arn)
        return resp['Tags']
    except c.exceptions.ResourceNotFoundException:
        return None
