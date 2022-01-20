package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/E_learning/controllers"
	"github.com/E_learning/models"
	sess "github.com/E_learning/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Coursereq struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `json:"name" binding:"required" bson:"Name,omitempty"`
	Author      string             `json:"author" bson:"Author"`
	Description string             `json:"description" binding:"required" bson:"Description,omitempty"`
	CreatedAt   time.Time          `json:"created_at" bson:"Created_at"`
}

func (server *Server) createCourse(ctx *gin.Context) {
	var req Coursereq
	//var x primitive.ObjectID
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	username := sess.SessionStart().Get("username", ctx)
	instructor, err := controllers.FindInstructor(ctx, username.(string))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorized"})
		return
	}
	args := models.Course{
		ID:          primitive.NewObjectID(),
		Name:        req.Name,
		Author:      instructor.UserName,
		Description: req.Description,
		CreatedAt:   time.Now(),
	}
	course, err := controllers.CreateCourse(ctx, &args)
	if err != nil {
		if we, ok := err.(mongo.WriteException); ok {
			for _, e := range we.WriteErrors {
				if e.Index == 0 {
					ctx.JSON(http.StatusBadRequest, gin.H{"error": "A course with this title already exists"})
					return
				}
			}
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Item not written"})
		return
	}
	ctx.JSON(http.StatusOK, course)
	log.Println("Remove data from Redis")
	server.redisClient.Del("Courses")
}

type getCourseRequest struct {
	ID string `uri:"id"  binding:"required"`
}

func (server *Server) deleteCourse(ctx *gin.Context) {
	var req getCourseRequest
	if err := ctx.BindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username := sess.SessionStart().Get("username", ctx)
	instructor, err := controllers.FindInstructor(ctx, username.(string))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorized"})
		return
	}
	course, err := controllers.FindCourse(ctx, req.ID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Not Found!"})
			return
		}

	}
	if course.Author != instructor.UserName {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "account doesn't belong to the authenticated user"})
		return
	} else {
		err = controllers.DeleteCourse(ctx, req.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong couldn't delete"})
			return
		}
		log.Println("Remove data from Redis")
		server.redisClient.Del("Courses")
		ctx.JSON(http.StatusOK, "Delete Course Successfull!")
	}
}

func (server *Server) findCourse(ctx *gin.Context) {
	var req getCourseRequest
	if err := ctx.BindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	course, err := controllers.FindCourse(ctx, req.ID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Not Found!"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong couldn't fetch data"})
		return
	}
	ctx.JSON(http.StatusOK, course)
}

type updateCourseRequest struct {
	ID          string `uri:"id"  binding:"required"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (server *Server) updateCourse(ctx *gin.Context) {
	var req updateCourseRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	arg := controllers.UpdateCourseParams{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
	}
	username := sess.SessionStart().Get("username", ctx)
	instructor, err := controllers.FindInstructor(ctx, username.(string))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorized"})
		return
	}
	course, err := controllers.FindCourse(ctx, arg.ID)
	if instructor.UserName != course.Author {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "account doesn't belong to the authenticated user"})
		return
	}
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Not Found!"})
			return
		}
	} else {
		results, err := controllers.UpdateCourse(ctx, arg)
		if err != nil {
			if we, ok := err.(mongo.WriteException); ok {
				for _, e := range we.WriteErrors {
					if e.Index == 0 {
						ctx.JSON(http.StatusBadRequest, gin.H{"error": "A course with this title already exists"})
						return
					}
				}
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong!"})
			return
		}
		log.Println("Remove data from Redis")
		server.redisClient.Del("Courses")
		ctx.JSON(http.StatusOK, results)
	}
}

type listCoursesRequest struct {
	Owner    string `json:"Instructor"`
	PageID   int64  `form:"page_id" binding:"required,min=0"`
	PageSize int64  `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listCourses(ctx *gin.Context) {
	var req listCoursesRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	username := sess.SessionStart().Get("username", ctx)
	instructor, err := controllers.FindInstructor(ctx, username.(string))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorized"})
		return
	}
	arg := controllers.ListCourseParams{
		Owner: instructor.UserName,
		Limit: req.PageSize,
		Skip:  (req.PageID - 1) * req.PageSize,
	}
	val, err := server.redisClient.Get("Courses").Result()
	if err == redis.Nil {
		log.Printf("Request to MongoDB")
		cacheresults, err := controllers.ListCourses(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		data, _ := json.Marshal(cacheresults)
		server.redisClient.Set("Courses", string(data), 0)
		ctx.JSON(http.StatusOK, cacheresults)
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		log.Printf("Request to Redis")
		courses := make([]models.Course, 0)
		json.Unmarshal([]byte(val), &courses)
		ctx.JSON(http.StatusOK, courses)
	}

}
