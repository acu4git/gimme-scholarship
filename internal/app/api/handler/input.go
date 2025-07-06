package handler

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
