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
"""Bootstraps the resources required to run the SNS integration tests.
"""
import logging

from acktest.bootstrapping import Resources, BootstrapFailureException
from acktest.bootstrapping.sqs import Queue
from acktest.bootstrapping.sns import Topic
from acktest.aws.identity import get_region, get_account_id

from e2e import bootstrap_directory
from e2e.bootstrap_resources import BootstrapResources

topic = Topic(name_prefix="subscribe-topic")

queue_policy = """{
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "sns.amazonaws.com"
      },
      "Action": "sqs:SendMessage",
      "Resource": "arn:aws:sqs:$REGION:$ACCOUNT_ID:$NAME",
      "Condition": {
        "ArnEquals": {
          "aws:SourceArn": "arn:aws:sns:$REGION:$ACCOUNT_ID:$TOPIC_NAME"
        }
      }
    }
  ]
}
"""

queue_policy_vars = {
    "$TOPIC_NAME": topic.name,
}

def service_bootstrap() -> Resources:
    logging.getLogger().setLevel(logging.INFO)

    resources = BootstrapResources(
        Topic=topic,
        Queue=Queue(
            name_prefix="subscribe-queue",
            policy=queue_policy,
            policy_vars=queue_policy_vars,
        ),
    )

    try:
        resources.bootstrap()
    except BootstrapFailureException as ex:
        exit(254)

    return resources

if __name__ == "__main__":
    config = service_bootstrap()
    # Write config to current directory by default
    config.serialize(bootstrap_directory)
