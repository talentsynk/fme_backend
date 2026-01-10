package schemas

import "time"

type MdaCreateSchema struct {
	RegisterName     string
	Email            string
	PhoneNumber      string
	Address          string
	StateOfOperation string
}

type GetAllMdaSchema struct {
	Id               uint
	StateOfOperation string
	Name             string
	Address          string
	IsActive         bool `json:"is_active"`
	STCCount         uint `json:"stc_count"`
	StudentCount     uint `json:"student_count"`
	UserId           uint
	CreatedAt        time.Time
	CourseCount      uint
	Email            string `json:"email"`
}

type GetMdaSchema struct {
	Id           uint
	Name         string
	CreatedAt    time.Time
	CourseCount  uint
	Address      string
	Email        string `json:"email"`
	IsActive     bool   `json:"is_active"`
	STCCount     int    `json:"stc_count"`
	StudentCount int    `json:"student_count"`
	UserId       uint
}

type UpdateMdaSchema struct {
	Id               int
	RegisterName     string
	Address          string
	StateOfOperation string
}
