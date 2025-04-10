package model

type Course struct {
	ID           int64
	TimetableID  string
	Class        string
	Type         string
	Credits      int
	Instructors  string
	Title        string
	Year         string
	Semester     string
	Day          string
	SyllabusYear int
}

type Courses []Course
