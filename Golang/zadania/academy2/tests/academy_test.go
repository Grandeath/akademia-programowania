package tests

import (
	"errors"
	"testing"

	academy "github.com/grupawp/akademia-programowania/Golang/zadania/academy2"
	"github.com/grupawp/akademia-programowania/Golang/zadania/academy2/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
		{
			name:              "Student is missing",
			studentName:       "John",
			studentYear:       3,
			studentFinalGrade: 0,
			wantError:         academy.ErrStudentNotFound,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockStudent := mocks.NewStudent(t)

			mockRepository := mocks.NewRepository(t)
			if errors.Is(tc.wantError, academy.ErrStudentNotFound) {
				mockRepository.On("Get", tc.studentName).Return(mockStudent, tc.wantError)
			} else if tc.studentFinalGrade < 1 || tc.studentFinalGrade > 5 {
				mockStudent.On("FinalGrade").Return(tc.studentFinalGrade)
				mockRepository.On("Get", tc.studentName).Return(mockStudent, nil)
			} else {
				mockStudent.On("Name").Return(tc.studentName)
				mockStudent.On("Year").Return(tc.studentYear)
				mockStudent.On("FinalGrade").Return(tc.studentFinalGrade)
				mockRepository.On("Get", tc.studentName).Return(mockStudent, nil)
				switch {
				case tc.studentFinalGrade == 1:
					mockRepository.On("Save", tc.studentName, tc.studentYear).Return(nil)
				case tc.studentYear == 3:
					mockRepository.On("Graduate", tc.studentName, mock.Anything).Return(nil)
				default:
					mockRepository.On("Save", tc.studentName, tc.studentYear+1).Return(nil)
				}
			}
			err := academy.GradeStudent(mockRepository, tc.studentName)
			if errors.Is(tc.wantError, academy.ErrStudentNotFound) {
				assert.Equal(t, nil, err)
			} else {
				assert.Equal(t, tc.wantError, err)
			}
			mockRepository.AssertExpectations(t)
			mockStudent.AssertExpectations(t)
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
	studentMap := map[uint8][]*mocks.Student{1: {}, 2: {}, 3: {}}
	for _, toCreate := range studentToCreate {
		student1 := mocks.NewStudent(t)
		student1.On("Name").Return(toCreate.studentName)
		student1.On("Year").Return(toCreate.studentYear)
		student1.On("FinalGrade").Return(toCreate.studentFinalGrade)
		studentMap[toCreate.studentYear] = append(studentMap[toCreate.studentYear], student1)
	}
	testCases := []struct {
		name      string
		r         []string
		giveYear  uint8
		wantErorr error
	}{
		{
			name:      "Empty map",
			giveYear:  0,
			wantErorr: errors.New("Wrong year"),
		},
		{
			name:      "Empty slice",
			giveYear:  4,
			wantErorr: nil,
		},
		{
			name:      "Correct Grades",
			giveYear:  2,
			wantErorr: nil,
		},
		{
			name:      "A Last year",
			giveYear:  3,
			wantErorr: nil,
		},
		{
			name:      "First year",
			giveYear:  1,
			wantErorr: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, mockStudent := range studentMap[tc.giveYear] {
				mockRepository := mocks.NewRepository(t)
				var gotStudentSlice []string
				gotStudentSlice = append(gotStudentSlice, mockStudent.Name())
				mockRepository.On("Get", mockStudent.Name()).Return(mockStudent, nil)
				switch {
				case mockStudent.FinalGrade() == 1:
					mockRepository.On("Save", mockStudent.Name(), mockStudent.Year()).Return(nil)
				case mockStudent.Year() == 3:
					mockRepository.On("Graduate", mockStudent.Name()).Return(nil)
				default:
					mockRepository.On("Save", mockStudent.Name(), mockStudent.Year()+1).Return(nil)
				}
				mockRepository.On("List", tc.giveYear).Return(gotStudentSlice, tc.wantErorr)
				err := academy.GradeYear(mockRepository, tc.giveYear)
				mockRepository.AssertExpectations(t)
				assert.Equal(t, tc.wantErorr, err)
				mockStudent.AssertExpectations(t)
			}
		})
	}
}
