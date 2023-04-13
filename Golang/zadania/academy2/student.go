package academy

//go:generate mockery --name Student
type Student interface {
	FinalGrade() int
	Name() string
	Year() uint8
}
