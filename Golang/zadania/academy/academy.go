package academy

import (
	"math"
)

type Student struct {
	Name       string
	Grades     []int
	Project    int
	Attendance []bool
}

// AverageGrade returns an average grade given a
// slice containing all grades received during a
// semester, rounded to the nearest integer.
func AverageGrade(grades []int) int {
	if len(grades) == 0 {
		return 0
	}
	var sumGrades int
	for _, grade := range grades {
		sumGrades += grade
	}
	averageGrade := int(math.Round(float64(sumGrades) / float64(len(grades))))
	return averageGrade
}

// AttendancePercentage returns a percentage of class
// attendance, given a slice containing information
// whether a student was present (true) of absent (false).
//
// The percentage of attendance is represented as a
// floating-point number ranging from  0 to 1,
// with 2 digits of precision.
func AttendancePercentage(attendance []bool) float64 {
	if len(attendance) == 0 {
		return 0
	}
	var peopleAttended int
	for _, attend := range attendance {
		if attend {
			peopleAttended += 1
		}
	}
	percentAttended := float64(peopleAttended) / float64(len(attendance))
	rounded := math.Round(percentAttended*1000) / 1000
	return rounded
}

// FinalGrade returns a final grade achieved by a student,
// ranging from 1 to 5.
//
// The final grade is calculated as the average of a project grade
// and an average grade from the semester, with adjustments based
// on the student's attendance. The final grade is rounded
// to the nearest integer.

// If the student's attendance is below 80%, the final grade is
// decreased by 1. If the student's attendance is below 60%, average
// grade is 1 or project grade is 1, the final grade is 1.
func FinalGrade(s Student) int {
	averageGrade := AverageGrade(s.Grades)
	percentAttended := AttendancePercentage(s.Attendance)

	if len(s.Grades) == 0 || s.Project == 1 || averageGrade == 1 || percentAttended < 0.6 {
		return 1
	}

	switch {
	case percentAttended >= 0.8:
		return int(math.Round(float64(averageGrade+s.Project) / 2.0))
	case percentAttended >= 0.6:
		return int(math.Round(float64(averageGrade+s.Project)/2.0)) - 1
	}
	return 1
}

// GradeStudents returns a map of final grades for a given slice of
// Student structs. The key is a student's name and the value is a
// final grade.
func GradeStudents(students []Student) map[string]uint8 {
	gradeStudents := make(map[string]uint8, len(students))
	for _, student := range students {
		gradeStudents[student.Name] = uint8(FinalGrade(student))
	}
	return gradeStudents
}
