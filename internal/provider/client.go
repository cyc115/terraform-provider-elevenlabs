package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cyc115/terraform-provider-elevenlabs/internal/models"
)

const baseURL = "https://api.elevenlabs.io"

type Client struct {
	apiKey string
	http   *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{apiKey: apiKey, http: &http.Client{Timeout: 30 * time.Second}}
}

func (c *Client) do(ctx context.Context, method, path string, body any) ([]byte, int, error) {
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, 0, err
		}
		bodyReader = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, baseURL+path, bodyReader)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("xi-api-key", c.apiKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	if resp.StatusCode >= 400 {
		return nil, resp.StatusCode, fmt.Errorf("ElevenLabs API %s %s → %d: %s", method, path, resp.StatusCode, string(data))
	}
	return data, resp.StatusCode, nil
}

// --- Agent CRUD ---

func (c *Client) CreateAgent(ctx context.Context, req *models.CreateAgentRequest) (*models.AgentResponse, error) {
	data, _, err := c.do(ctx, http.MethodPost, "/v1/convai/agents/create", req)
	if err != nil {
		return nil, err
	}
	var resp models.AgentResponse
	return &resp, json.Unmarshal(data, &resp)
}

func (c *Client) GetAgent(ctx context.Context, agentID string) (*models.AgentResponse, error) {
	data, _, err := c.do(ctx, http.MethodGet, "/v1/convai/agents/"+agentID, nil)
	if err != nil {
		return nil, err
	}
	var resp models.AgentResponse
	return &resp, json.Unmarshal(data, &resp)
}

func (c *Client) UpdateAgent(ctx context.Context, agentID string, req *models.UpdateAgentRequest) (*models.AgentResponse, error) {
	data, _, err := c.do(ctx, http.MethodPatch, "/v1/convai/agents/"+agentID, req)
	if err != nil {
		return nil, err
	}
	var resp models.AgentResponse
	return &resp, json.Unmarshal(data, &resp)
}

func (c *Client) DeleteAgent(ctx context.Context, agentID string) error {
	_, _, err := c.do(ctx, http.MethodDelete, "/v1/convai/agents/"+agentID, nil)
	return err
}

// --- Webhook CRUD (list+filter for Read — no individual GET endpoint) ---

func (c *Client) CreateWebhook(ctx context.Context, req *models.CreateWebhookRequest) (*models.WebhookResponse, error) {
	data, _, err := c.do(ctx, http.MethodPost, "/v1/workspace/webhooks", req)
	if err != nil {
		return nil, err
	}
	var resp models.WebhookResponse
	return &resp, json.Unmarshal(data, &resp)
}

func (c *Client) ListWebhooks(ctx context.Context) ([]models.WebhookResponse, error) {
	data, _, err := c.do(ctx, http.MethodGet, "/v1/workspace/webhooks", nil)
	if err != nil {
		return nil, err
	}
	var resp models.ListWebhooksResponse
	return resp.Webhooks, json.Unmarshal(data, &resp)
}

func (c *Client) GetWebhook(ctx context.Context, webhookID string) (*models.WebhookResponse, error) {
	list, err := c.ListWebhooks(ctx)
	if err != nil {
		return nil, err
	}
	for i := range list {
		if list[i].WebhookID == webhookID {
			return &list[i], nil
		}
	}
	return nil, fmt.Errorf("webhook %s not found", webhookID)
}

func (c *Client) UpdateWebhook(ctx context.Context, webhookID string, req *models.UpdateWebhookRequest) (*models.WebhookResponse, error) {
	data, _, err := c.do(ctx, http.MethodPatch, "/v1/workspace/webhooks/"+webhookID, req)
	if err != nil {
		return nil, err
	}
	var resp models.WebhookResponse
	return &resp, json.Unmarshal(data, &resp)
}

func (c *Client) DeleteWebhook(ctx context.Context, webhookID string) error {
	_, _, err := c.do(ctx, http.MethodDelete, "/v1/workspace/webhooks/"+webhookID, nil)
	return err
}

// --- Phone number CRUD (list+filter for Read) ---

func (c *Client) ImportPhone(ctx context.Context, req *models.ImportPhoneNumberRequest) (*models.PhoneNumberResponse, error) {
	data, _, err := c.do(ctx, http.MethodPost, "/v1/convai/phone-numbers/import", req)
	if err != nil {
		return nil, err
	}
	var resp models.PhoneNumberResponse
	return &resp, json.Unmarshal(data, &resp)
}

func (c *Client) ListPhones(ctx context.Context) ([]models.PhoneNumberResponse, error) {
	data, _, err := c.do(ctx, http.MethodGet, "/v1/convai/phone-numbers", nil)
	if err != nil {
		return nil, err
	}
	// API returns a bare JSON array, not a wrapped object.
	var phones []models.PhoneNumberResponse
	return phones, json.Unmarshal(data, &phones)
}

func (c *Client) GetPhone(ctx context.Context, phoneID string) (*models.PhoneNumberResponse, error) {
	list, err := c.ListPhones(ctx)
	if err != nil {
		return nil, err
	}
	for i := range list {
		if list[i].PhoneNumberID == phoneID {
			return &list[i], nil
		}
	}
	return nil, fmt.Errorf("phone number %s not found", phoneID)
}

func (c *Client) UpdatePhone(ctx context.Context, phoneID string, req *models.UpdatePhoneNumberRequest) (*models.PhoneNumberResponse, error) {
	data, _, err := c.do(ctx, http.MethodPatch, "/v1/convai/phone-numbers/"+phoneID, req)
	if err != nil {
		return nil, err
	}
	var resp models.PhoneNumberResponse
	return &resp, json.Unmarshal(data, &resp)
}

func (c *Client) DeletePhone(ctx context.Context, phoneID string) error {
	_, _, err := c.do(ctx, http.MethodDelete, "/v1/convai/phone-numbers/"+phoneID, nil)
	return err
}
