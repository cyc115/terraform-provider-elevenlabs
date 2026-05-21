package models

type CreateWebhookRequest struct {
	Name        string          `json:"name"`
	WebhookURL  string          `json:"webhook_url"`
	Usage       []string        `json:"usage,omitempty"`
	AuthType    string          `json:"auth_type,omitempty"`
	RetryEnabled bool           `json:"retry_enabled"`
}

type WebhookResponse struct {
	WebhookID                   string  `json:"webhook_id"`
	Name                        string  `json:"name"`
	WebhookURL                  string  `json:"webhook_url"`
	IsDisabled                  bool    `json:"is_disabled"`
	IsAutoDisabled              bool    `json:"is_auto_disabled"`
	CreatedAtUnix               int64   `json:"created_at_unix"`
	AuthType                    string  `json:"auth_type"`
	RetryEnabled                bool    `json:"retry_enabled"`
	Secret                      string  `json:"secret,omitempty"` // only on Create response
	MostRecentFailureErrorCode  *string `json:"most_recent_failure_error_code"`
	MostRecentFailureTimestamp  *string `json:"most_recent_failure_timestamp"`
}

type ListWebhooksResponse struct {
	Webhooks []WebhookResponse `json:"webhooks"`
}

type UpdateWebhookRequest struct {
	Name         *string  `json:"name,omitempty"`
	WebhookURL   *string  `json:"webhook_url,omitempty"`
	AuthType     *string  `json:"auth_type,omitempty"`
	RetryEnabled *bool    `json:"retry_enabled,omitempty"`
}
