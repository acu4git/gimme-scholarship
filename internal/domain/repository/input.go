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

// At Mode field,
// specify "REGISTER" if you want to register favorite scholarship,
// specify "DELETE" if you want to delete.
type UserFavoriteInput struct {
	Mode          string
	UserID        string
	ScholarshipID int64
}
