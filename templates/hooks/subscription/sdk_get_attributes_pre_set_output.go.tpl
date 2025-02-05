
	// If one of the below 3 fields (i.e. Protocol, Endpoint, TopicARN/TopicRef) is empty
	// populate them with their respective attribute values
	// present in the response of the GetSubscriptionAttributes API call
	// Use case: adopting an existing subscription by subcription ARN
	if ko.Spec.Protocol == nil {
		ko.Spec.Protocol = aws.String(resp.Attributes["Protocol"])
	}
	if ko.Spec.Endpoint == nil {
		ko.Spec.Endpoint = aws.String(resp.Attributes["Endpoint"])
	}
	if ko.Spec.TopicARN == nil && ko.Spec.TopicRef == nil {
		ko.Spec.TopicARN = aws.String(resp.Attributes["TopicArn"])
	}
