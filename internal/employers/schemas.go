package employer

var CreateEmployerSchema struct {
    FirstName   string   
    LastName    string  
    BusinessName    string 
    Email       string   
    PhoneNumber string   
    NIN         string   
    State       string   
    LGA         string   
    Password    string
    CompanyName string
    IsCompany   bool
    CompanyCAC  string
}



type GetEmployerSchema struct {
    Id          uint
    FirstName   string   
    LastName    string   
    PhoneNumber string   
    NIN         string   
    State       string   
    LGA         string   
    UserId      uint  
}    


type RatingFilterSchema struct {
	DaysAgo uint			`form:"days_ago"`
	MaxRating uint 		    `form:"max_ratings"`
	MinRating uint 		    `form:"min_ratings"`

}


type JobFilterSchema struct {
	Status string		   `form:"status"`
	MinBudget float64	   `form:"min_budget"`
	MaxBudget float64	   `form:"max_budget"`
	JobType string		    `form:"job_type"`
	DaysAgo uint			`form:"days_ago"`
	Lga string				`form:"lga"`
	State string			`form:"state"`
}

type GetJobSchema struct {
	Id uint
	JobTitle string
	Description string
	Amount string
	JobType string
	Status string

}