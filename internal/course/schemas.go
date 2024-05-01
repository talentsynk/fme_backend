package course

var CreateCourseSchema struct {
	Name        string
	Description string
	CategoryID uint
}

var CreateCategorySchema struct {
	Name        string
	Description string
}

type GetCourseSchema struct {
	Id				uint
	TotalStudents	uint	
	TotalMda		uint
	TotalStc		uint
	Description		string
	Name			string
	}


type GetAllCoursesSchema struct {
	Id uint
	Name string
	Description string
}

type GetAllCategoriesSchema struct {
	Id uint
	Name string
	Description string
}