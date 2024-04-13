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
