package handler

type GetScholarshipInput struct {
	Target string `query:"target"`
	Type   string `query:"type"`
}

type PostUserInput struct {
	Email       string `json:"email"`
	Grade       int    `json:"grade"`
	Level       string `json:"level"`
	AcceptEmail bool   `json:"accept_email"`
}
