package handler

import "encoding/json"

type GetScholarshipInput struct {
	Target string `query:"target"`
	Type   string `query:"type"`
}

type PostUserInput struct {
	Email       string `form:"email"`
	Grade       int    `form:"grade"`
	Level       string `form:"level"`
	AcceptEmail bool   `form:"accept_email"`
}

type PutUserInput struct {
	Grade       int    `form:"grade"`
	Level       string `form:"level"`
	AcceptEmail bool   `form:"accept_email"`
}

type ClerkWebhookEvent struct {
	Object string          `json:"object"`
	Type   string          `json:"type"`
	Data   json.RawMessage `json:"data"`
}

type ClerkUser struct {
	ID             string `json:"id"`
	Username       string `json:"username"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	EmailAddresses []struct {
		ID           string `json:"id"`
		EmailAddress string `json:"email_address"`
	} `json:"email_addresses"`
	PublicMetadata  map[string]any `json:"public_metadata"`
	PrivateMetadata map[string]any `json:"private_metadata"`
}
