package controllers

import (
	"context"
	"testing"
	"time"

	"github.com/E_learning/models"
	"github.com/E_learning/util"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func createcourse() models.Course {

	//section := NewSection()
	arg := models.Course{
		ID:          primitive.NewObjectID(),
		Name:        util.RandomAuthor(),
		Author:      util.RandomAuthor(),
		Description: util.RandomString(40),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	return arg
}
func NewSection() []*models.Section {
	result := []*models.Section{}
	content := models.Content{
		ID:       primitive.NewObjectID(),
		SubTitle: util.RandomString(10),
		Content:  util.RandomString(80),
	}
	section := models.Section{
		ID:      primitive.NewObjectID(),
		Title:   util.RandomString(10),
		Content: []*models.Content{&content},
	}
	result = append(result, &section)
	return result
}

func TestCourseCollection(t *testing.T) {
	collection := CourseCollection(context.Background())
	require.NotEmpty(t, collection)
}

func TestCreateCourse(t *testing.T) {
	arg := createcourse()
	course, err := CreateCourse(context.Background(), &arg)
	require.NoError(t, err)
	require.NotEmpty(t, &course)
	require.Equal(t, course.Author, arg.Author)
	require.Equal(t, course.Description, arg.Description)

}

func TestFindCourse(t *testing.T) {
	course := createcourse()
	arg, err := CreateCourse(context.Background(), &course)
	require.NoError(t, err)
	course, err = FindCourse(context.Background(), arg.ID.Hex())
	require.NoError(t, err)
	require.NotNil(t, course)
	require.Equal(t, course.Author, arg.Author)
	require.Equal(t, course.Description, arg.Description)
}
func TestUpdateCourse(t *testing.T) {
	course := createcourse()
	arg, err := CreateCourse(context.Background(), &course)
	require.NoError(t, err)
	updateargs := UpdateCourseParams{
		ID:          arg.ID.Hex(),
		Name:        util.RandomString(6),
		Description: util.RandomString(100),
	}
	results, err := UpdateCourse(context.Background(), updateargs)
	require.NoError(t, err)
	require.NotNil(t, results)
}

func TestDeleteCourse(t *testing.T) {
	course := createcourse()
	arg, err := CreateCourse(context.Background(), &course)
	require.NoError(t, err)
	err = DeleteCourse(context.Background(), arg.ID.Hex())
	require.NoError(t, err)
	result, err := FindCourse(context.Background(), arg.ID.Hex())
	require.Error(t, err)
	require.Empty(t, result)
}

func TestListCourse(t *testing.T) {
	course := createcourse()
	arg, err := CreateCourse(context.Background(), &course)
	require.NoError(t, err)
	arg2 := ListCourseParams{
		Owner: arg.Author,
		Limit: 1,
		Skip:  0,
	}
	results, err := ListCourses(context.Background(), arg2)
	require.NoError(t, err)
	require.NotNil(t, results)
	require.NotEmpty(t, results)
}
