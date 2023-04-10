package tests

import (
	"errors"
	"testing"

	academy "github.com/grupawp/akademia-programowania/Golang/zadania/academy2"
	"github.com/grupawp/akademia-programowania/Golang/zadania/academy2/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGradeStudent(t *testing.T) {
	testCases := []struct {
		name              string
		studentName       string
		studentYear       uint8
		studentFinalGrade int
		wantError         error
	}{
		{
			name:              "Empty student",
			studentName:       "",
			studentYear:       0,
			studentFinalGrade: 2,
			wantError:         nil,
		},
		{
			name:              "Rais a Year",
			studentName:       "John",
			studentYear:       2,
			studentFinalGrade: 5,
			wantError:         nil,
		},
		{
			name:              "Student didn't pass",
			studentName:       "John",
			studentYear:       2,
			studentFinalGrade: 1,
			wantError:         nil,
		},
		{
			name:              "Student graduated",
			studentName:       "John",
			studentYear:       3,
			studentFinalGrade: 5,
			wantError:         nil,
		},
		{
			name:              "Invalid Grade",
			studentName:       "John",
			studentYear:       3,
			studentFinalGrade: 0,
			wantError:         academy.ErrInvalidGrade,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			student1 := mocks.Student{}
			student1.On("Name").Return(tc.studentName)
			student1.On("Year").Return(tc.studentYear)
			student1.On("FinalGrade").Return(tc.studentFinalGrade)

			studentSlice := []mocks.Student{student1}

			r := repositoryMock{mapStudents: map[uint8][]mocks.Student{tc.studentYear: studentSlice}}

			studentToGrade, err := r.Get(tc.studentName)
			assert.Equal(t, tc.wantError, academy.GradeStudent(&r, tc.studentName))
			if err != nil || tc.wantError != nil {
				t.SkipNow()
			}

			switch {
			case tc.studentFinalGrade == 1:
				got, _ := r.List(studentToGrade.Year())
				assert.Contains(t, got, tc.studentName)
			case studentToGrade.Year() == 3:
				got, _ := r.List(studentToGrade.Year())
				assert.NotContains(t, got, tc.studentName)
			default:
				got, _ := r.List(studentToGrade.Year() + 1)
				assert.Contains(t, got, tc.studentName)
			}
		})
	}

}

func TestGradeYear(t *testing.T) {
	studentToCreate := []struct {
		studentName       string
		studentYear       uint8
		studentFinalGrade int
	}{
		{
			studentName:       "John Filler",
			studentYear:       2,
			studentFinalGrade: 5,
		},
		{
			studentName:       "Angelina Jolie",
			studentYear:       3,
			studentFinalGrade: 2,
		},
		{
			studentName:       "John Paul 2",
			studentYear:       2,
			studentFinalGrade: 1,
		},
		{
			studentName:       "Lady Gaga",
			studentYear:       2,
			studentFinalGrade: 2,
		},
		{
			studentName:       "Aragorn",
			studentYear:       1,
			studentFinalGrade: 4,
		},
		{
			studentName:       "Gandalf",
			studentYear:       3,
			studentFinalGrade: 4,
		},
	}
	studentMap := map[uint8][]mocks.Student{1: {}, 2: {}, 3: {}}
	for _, toCreate := range studentToCreate {
		student1 := mocks.Student{}
		student1.On("Name").Return(toCreate.studentName)
		student1.On("Year").Return(toCreate.studentYear)
		student1.On("FinalGrade").Return(toCreate.studentFinalGrade)
		studentMap[toCreate.studentYear] = append(studentMap[toCreate.studentYear], student1)
	}
	testCases := []struct {
		name      string
		r         repositoryMock
		giveYear  uint8
		wantErorr error
	}{
		{
			name:      "Empty map",
			r:         repositoryMock{},
			giveYear:  0,
			wantErorr: errors.New("Wrong year"),
		},
		{
			name:      "Empty slice",
			r:         repositoryMock{map[uint8][]mocks.Student{1: {}, 2: {}, 3: {}}},
			giveYear:  0,
			wantErorr: errors.New("Wrong year"),
		},
		{
			name:      "Correct Grades",
			r:         repositoryMock{studentMap},
			giveYear:  2,
			wantErorr: nil,
		},
		{
			name:      "A Last year",
			r:         repositoryMock{studentMap},
			giveYear:  3,
			wantErorr: nil,
		},
		{
			name:      "First year",
			r:         repositoryMock{studentMap},
			giveYear:  1,
			wantErorr: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var wantYear1 []string
			wantYear2, _ := tc.r.List(tc.giveYear + 1)
			for _, studentWant := range tc.r.mapStudents[tc.giveYear] {
				grade := studentWant.FinalGrade()
				switch {
				case grade == 1:
					wantYear1 = append(wantYear1, studentWant.Name())
				case studentWant.Year() == 3:
					continue
				default:
					wantYear2 = append(wantYear2, studentWant.Name())
				}
			}
			err := academy.GradeYear(&tc.r, tc.giveYear)
			assert.Equal(t, tc.wantErorr, err)
			if err != nil {
				t.SkipNow()
			}
			gotSlice1, _ := tc.r.List(tc.giveYear)

			assert.Equal(t, wantYear1, gotSlice1)

			if tc.giveYear != 3 {
				gotSlice2, _ := tc.r.List(tc.giveYear + 1)
				assert.Equal(t, wantYear2, gotSlice2)
			}
		})
	}
}

type repositoryMock struct {
	mapStudents map[uint8][]mocks.Student
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

func (r *repositoryMock) Get(name string) (academy.Student, error) {
	for _, yearSlice := range r.mapStudents {
		for _, student := range yearSlice {
			if student.Name() == name {
				return &student, nil
			}
		}
	}
	return &mocks.Student{}, academy.ErrStudentNotFound
}

func (r *repositoryMock) Save(name string, year uint8) error {
	err := r.Graduate(name)
	studentNew := mocks.Student{}
	studentNew.On("Name").Return(name)
	studentNew.On("Year").Return(year)
	studentNew.On("FinalGrade").Return(0)
	if err != nil {
		return err
	}
	yearSlice, ok := r.mapStudents[year]
	if !ok {
		r.mapStudents[year] = []mocks.Student{studentNew}
	} else {
		r.mapStudents[year] = append(yearSlice, studentNew)
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
	return academy.ErrStudentNotFound
}
