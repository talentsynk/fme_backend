package artisans


type GetAllArtisansSchema struct {
	
}

type RatingFilterSchema struct {
	DaysAgo uint			`form:"days_ago"`
	MaxRating uint 		`form:"max_ratings"`
	MinRating uint 		`form:"min_ratings"`

}


type ArtisanFilterSchema struct {
	MinRating float64	`form:"min_rating"`
	MaxRating float64	`form:"max_rating"`
	RatingSort bool 		`form:"rating_sort"`

}


type ArtisanDataDownload struct {
	FirstName string
	LastName	string
	ArtisanID string
	JobTitle	string
	Budget	float64
	JobType	string
	JobLocation	string
	JobDescription	string
	ApplicationStatus string
	JobStatus		string
	ArtisanRating	string
	ArtisanRatingDescription	string

}