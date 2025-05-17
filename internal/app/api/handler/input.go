package handler

type GetScholarshipInput struct {
	Target string `query:"target"`
	Type   string `query:"type"`
}
