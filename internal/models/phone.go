package models

type ImportPhoneNumberRequest struct {
	PhoneNumber      string  `json:"phone_number"`
	Provider         string  `json:"provider"`
	Label            string  `json:"label,omitempty"`
	TwilioAccountSID *string `json:"twilio_account_sid,omitempty"`
	TwilioAuthToken  *string `json:"twilio_auth_token,omitempty"`
}

type PhoneNumberResponse struct {
	PhoneNumberID    string  `json:"phone_number_id"`
	PhoneNumber      string  `json:"phone_number"`
	Label            string  `json:"label"`
	Provider         string  `json:"provider"`
	SupportsInbound  bool    `json:"supports_inbound"`
	SupportsOutbound bool    `json:"supports_outbound"`
	AssignedAgent    *string `json:"assigned_agent"`
}

type ListPhoneNumbersResponse struct {
	PhoneNumbers []PhoneNumberResponse `json:"phone_numbers"`
}

type UpdatePhoneNumberRequest struct {
	Label         *string `json:"label,omitempty"`
	AssignedAgent *string `json:"assigned_agent,omitempty"`
}
