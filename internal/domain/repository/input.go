package repository

type FilterOption struct {
	Target string
	Type   string
}

type UserInput struct {
	ID          string
	Email       string
	Level       string
	Grade       int64
	AcceptEmail bool
}
