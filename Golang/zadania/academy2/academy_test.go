package academy

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type repositoryMock struct {
	mapStudents map[uint8][]Sophomore
}

func (r *repositoryMock) List(year uint8) (names []string, err error) {
	var nameSlice []string
	for _, student := range r.mapStudents[year] {
		nameSlice = append(nameSlice, student.Name())
	}
	if len(nameSlice) == 0 {
		return nameSlice, errors.New("Wrong year")
	}
	return nameSlice, nil
}

func (r *repositoryMock) Get(name string) (Student, error) {
	for _, yearSlice := range r.mapStudents {
		for _, student := range yearSlice {
			if student.Name() == name {
				return student, nil
			}
		}
	}
	return Sophomore{}, ErrStudentNotFound
}

func (r *repositoryMock) Save(name string, year uint8) error {
	err := r.Graduate(name)
	if err != nil {
		return err
	}
	yearSlice, ok := r.mapStudents[year]
	if !ok {
		r.mapStudents[year] = []Sophomore{{name: name}}
	} else {
		r.mapStudents[year] = append(yearSlice, Sophomore{name: name})
	}

	return nil
}

func (r *repositoryMock) Graduate(name string) error {
	for key, yearSlice := range r.mapStudents {
		for index, student := range yearSlice {
			if student.Name() == name {
				switch index {
				case 0:
					r.mapStudents[key] = yearSlice[1:]
				case len(yearSlice):
					r.mapStudents[key] = yearSlice[:index-1]
				default:
					r.mapStudents[key] = append(yearSlice[:index], yearSlice[index+1:]...)
				}
				return nil
			}
		}
	}
	return ErrStudentNotFound
}

var sophomoreSLice = []Sophomore{
	{
		name:       "John Doe",
		grades:     []int{5, 4, 5, 5, 5},
		project:    5,
		attendance: []bool{true, true, true, true, true},
	},
	{
		name:       "Jane Smith",
		grades:     []int{5, 4, 5, 5, 5},
		project:    5,
		attendance: []bool{true, false, true, true, false},
	},
	{
		name:       "Jane Forge",
		grades:     []int{5, 4, 5, 5, 5},
		project:    4,
		attendance: []bool{true, false, false, true, false},
	},
	{
		name:       "Jason Bourne",
		grades:     []int{2, 2, 3, 2, 2},
		project:    2,
		attendance: []bool{true, false, true, true, false},
	},
	{
		name:       "John",
		grades:     []int{4, 4, 4},
		project:    5,
		attendance: []bool{true, true},
	},
	{
		name:       "John Foe",
		grades:     []int{5, 4, 5, 5, 5},
		project:    2,
		attendance: []bool{true, true, true, true, true},
	},
	{
		name:       "John Boe",
		grades:     []int{5, 4, 5, 5, 5},
		project:    1,
		attendance: []bool{true, true, true, true, true},
	},
	{
		name:       "John Goe",
		grades:     []int{1, 2, 1, 1, 2},
		project:    5,
		attendance: []bool{true, true, true, true, true},
	},
	{
		name:       "Jane Farmer",
		grades:     []int{1, 1, 1, 1, 1},
		project:    1,
		attendance: []bool{true, false, true, true, false},
	},
}

func TestGradeStudent(t *testing.T) {
	testCases := []struct {
		name        string
		r           repositoryMock
		studentName string
		wantErorr   error
		expectGrade int
	}{
		{
			name:        "Empty repository",
			r:           repositoryMock{},
			studentName: "",
			wantErorr:   nil,
			expectGrade: 0,
		},
		{
			name:        "Raise a year",
			r:           repositoryMock{mapStudents: map[uint8][]Sophomore{2: sophomoreSLice}},
			studentName: "John Doe",
			wantErorr:   nil,
			expectGrade: 5,
		},
		{
			name:        "Didn't pass",
			r:           repositoryMock{mapStudents: map[uint8][]Sophomore{2: sophomoreSLice}},
			studentName: "Jane Forge",
			wantErorr:   nil,
			expectGrade: 1,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			studentToGrade, err := tc.r.Get(tc.studentName)
			assert.Equal(t, tc.wantErorr, GradeStudent(&tc.r, tc.studentName))
			if err != nil {
				t.SkipNow()
			}
			switch {
			case tc.expectGrade == 1:
				got, _ := tc.r.List(studentToGrade.Year())
				assert.Contains(t, got, tc.studentName)
			case studentToGrade.Year() == 3:
				got, _ := tc.r.List(studentToGrade.Year())
				assert.NotContains(t, got, tc.studentName)
			default:
				got, _ := tc.r.List(studentToGrade.Year() + 1)
				assert.Contains(t, got, tc.studentName)
			}
		})
	}
}

func TestGradeYear(t *testing.T) {
	testCases := []struct {
		name      string
		r         repositoryMock
		giveYear  uint8
		wantErorr error
	}{
		{
			name:      "Empty repository",
			r:         repositoryMock{},
			giveYear:  2,
			wantErorr: errors.New("Wrong year"),
		},
		{
			name:      "Correct Grades",
			r:         repositoryMock{mapStudents: map[uint8][]Sophomore{1: {}, 2: sophomoreSLice, 3: {}}},
			giveYear:  2,
			wantErorr: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			wantMap := map[uint8][]Sophomore{1: {}, 2: {}, 3: {}}
			for _, student := range tc.r.mapStudents[tc.giveYear] {
				grade := student.FinalGrade()
				switch {
				case grade == 1:
					wantMap[student.Year()] = append(wantMap[student.Year()], Sophomore{name: student.Name()})
				default:
					wantMap[student.Year()+1] = append(wantMap[student.Year()+1], Sophomore{name: student.Name()})
				}
			}
			err := GradeYear(&tc.r, tc.giveYear)
			assert.Equal(t, tc.wantErorr, err)
			if err != nil {
				t.SkipNow()
			}
			assert.Equal(t, wantMap, tc.r.mapStudents)
		})
	}
}
