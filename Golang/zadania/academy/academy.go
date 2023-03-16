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
	var averageGrade int
	var sumGrades int
	for _, grade := range grades {
		sumGrades += grade
	}
	averageGrade = int(math.Round(float64(sumGrades) / float64(len(grades))))
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
	if len(s.Grades) == 0 {
		return 1
	}
	if s.Project == 1 {
		return 1
	}

	averageGrade := AverageGrade(s.Grades)

	if averageGrade == 1 {
		return 1
	}

	percentAttended := AttendancePercentage(s.Attendace)
	var finalGrade int
	switch {
	case percentAttended >= 0.8:
		finalGrade = int(math.Round(float64(averageGrade+s.Project) / 2.0))
	case percentAttended < 0.6:
		finalGrade = 1
	case percentAttended >= 0.6:
		finalGrade = int(math.Round(float64(averageGrade+s.Project)/2.0)) - 1
	}
	return finalGrade
}

// GradeStudents returns a map of final grades for a given slice of
// Student structs. The key is a student's name and the value is a
// final grade.
func GradeStudents(students []Student) map[string]uint8 {
	if len(students) == 0 {
		return map[string]uint8{}
	}
	gradeStudents := make(map[string]uint8, len(students))
	for _, student := range students {
		gradeStudents[student.Name] = uint8(FinalGrade(student))
	}
	return gradeStudents
}
