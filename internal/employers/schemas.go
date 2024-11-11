package employer



var CreateEmployerSchema struct {
    FirstName    string   
    LastName     string  
    Email        string   
    PhoneNumber  string   
    NIN          string   
    State        string   
    LGA          string   
    Password     string
}




type GetEmployerSchema struct {
    Id           uint
    FirstName    string   
    LastName     string   
    Email        string
    EmployerType string   
    PhoneNumber  string   
    NIN          string   
    State        string   
    LGA          string   
    UserId       uint      
}
