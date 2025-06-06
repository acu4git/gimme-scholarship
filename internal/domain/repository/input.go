package repository

type FilterOption struct {
	UserID *string
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

type UserFavoriteInput struct {
	UserID        string
	ScholarshipID int64
}
