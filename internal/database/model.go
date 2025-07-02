package database

type UserScholarshipNotification struct {
	UserID    string `db:"user_id"`
	UserName  string `db:"user_name"`
	UserEmail string `db:"user_email"`
	scholarship
}
