	// SNS does not allow deleting a FIFO topic that still has an active
	// message archive policy: DeleteTopic returns
	// "InvalidState: Cannot delete a topic with an ArchivePolicy".
	// The archive policy must be deactivated first by setting the
	// ArchivePolicy attribute to an empty value. Deactivating the policy
	// also deletes any archived messages, which is the expected behavior
	// when the topic itself is being deleted.
	if r.ko.Spec.ArchivePolicy != nil {
		_, err = rm.sdkapi.SetTopicAttributes(ctx, &svcsdk.SetTopicAttributesInput{
			TopicArn:       (*string)(r.ko.Status.ACKResourceMetadata.ARN),
			AttributeName:  aws.String("ArchivePolicy"),
			AttributeValue: aws.String("{}"),
		})
		rm.metrics.RecordAPICall("SET_ATTRIBUTES", "SetTopicAttributes", err)
		if err != nil {
			return nil, err
		}
	}
