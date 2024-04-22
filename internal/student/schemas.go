package student

import "time"

// import "fme_backend/internal/course"

var CreateStudentSchema struct {
	Firstname string
	Lastname string
	Gender string
	PhoneNumber	string
	StateOfOrigin string
	StateOfResidence string
	Email string
	DOBstring string
	SID		string
	NsqLevel	string
	CourseID	uint
	Address string
}

type GetAllStudentSchema struct {
    StudentID    uint
    FirstName    string
    LastName     string
    IsActive     bool
    Email        string
	CoursesTaken 	 string
	StateOfResidence	string
}

type GetStudentSchema struct {
    StudentID    uint
    FirstName    string
    LastName     string
    IsActive     bool
    Email        string
	CoursesTaken 	 string
	Gender			string
	StateOfResidence	string
	Address			string
	CreatedAt		time.Time
	PhoneNumber		string
}
