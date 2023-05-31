
	// If the Protocol field is empty, populate it with the Protocol attribute value 
	// present in the response of the GetSubscriptionAttributes API call
	// This is a workaround for the fact that the Protocol field is a required field
	// Use case: adopting an existing subscription by subcription ARN
	if ko.Spec.Protocol == nil {
		ko.Spec.Protocol = resp.Attributes["Protocol"]
	}
	