
	// If the Name field is empty, populate it with the last part of the topic ARN
	// This is a workaround for the fact that the Name field is required by the
	// CreateTopic API call, but not by the GetTopicAttributes API call
	// Use case: adopting an existing topic by topic ARN
	if ko.Spec.Name == nil {
		topicName, err := rm.getTopicNameFromARN(tmpARN)
		if err != nil {
			return nil, err
		}
		ko.Spec.Name = &topicName
	}
	