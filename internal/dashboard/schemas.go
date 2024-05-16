package dashboard

type CoursePercentage struct {
    CourseName   string
    TotalPercent float64
	TotalStudents int
    CertifiedCount  int
    UncertifiedCount int
}


type StcCount struct {
    StcName            string
    StcID              uint
    TotalStudents      int
    TotalCertified     int
}