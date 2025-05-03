package database

import "time"

type educationLevel struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

type user struct {
	UUID             string    `db:"uuid"`
	Email            string    `db:"email"`
	Name             string    `db:"name"`
	EducationLevelID int64     `db:"department_id"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}

type scholarship struct {
	ID             int64     `db:"id"`
	Name           string    `db:"name"`
	Address        string    `db:"name"`
	TargetDetail   string    `db:"target_detail"`
	AmountDetail   string    `db:"amount_detail"`
	TypeDetail     string    `db:"type_detail"`
	CapacityDetail string    `db:"capacity_detail"`
	Deadline       time.Time `db:"deadline"`
	DeadlineDetail string    `db:"deadline_detail"`
	ContactPoint   string    `db:"contact_point"`
	Remark         string    `db:"remark"`
	PostingDate    time.Time `db:"posting_date"`
	CreatedAt      time.Time `db:"created_at"`
}

type scholarshipTarget struct {
	ID               int64 `db:"id"`
	ScholarshipID    int64 `db:"scholarship_id"`
	EducationLevelID int64 `db:"education_level_id"`
}
