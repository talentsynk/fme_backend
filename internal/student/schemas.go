package student

import "time"

// import "fme_backend/internal/course"



type GetAllStudentSchema struct {
	ID               uint
	FirstName        string
	LastName         string
	IsActive         bool
	Email            string
	CoursesTaken     string
	StateOfResidence string
	UserID           int
	PhoneNumber      string
	CreatedAt        time.Time
	Gender           string
	Address          string
	IsGraduated      bool
}

type GetStudentSchema struct {
	ID               uint
	FirstName        string
	LastName         string
	IsActive         bool
	Email            string
	CoursesTaken     string
	Gender           string
	StateOfResidence string
	Address          string
	CreatedAt        time.Time
	PhoneNumber      string
	UserId           int
	IsGraduated      bool
}

type TotalStudentInfo struct {
	TotalStudents         int
	TotalActiveStudents   int
	TotalInactiveStudents int
}

type ResidenceInfo struct {
	StatOfResidence string
	TotalEnroll     int
	TotalGraduated  int
}

type CreateStudentSchematype struct {
	Firstname              string
	Lastname               string
	Gender                 string
	PhoneNumber            string
	StateOfOrigin          string
	StateOfResidence       string
	Email                  string
	DOBstring              string
	SID                    string
	CourseID               uint
	Address                string
	NationalIdentityNumber string
	LocalGovernment        string
	IsDisabled             bool
	DisabilityName         string
}

type Disability struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type GraduateStudentSchemaType struct {
	DateOfGrad string
	NsqLevel   string
}
