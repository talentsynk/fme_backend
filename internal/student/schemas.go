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
    ID    uint
    FirstName    string
    LastName     string
    IsActive     bool
    Email        string
	CoursesTaken 	 string
	StateOfResidence	string
	UserID			int
	PhoneNumber string
	CreatedAt		time.Time
	Gender string
	Address	string


}

type GetStudentSchema struct {
    ID    uint
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
	UserId 			int
}

type TotalStudentInfo struct {
	TotalStudents int
	TotalActiveStudents	int
	TotalInactiveStudents	int
}