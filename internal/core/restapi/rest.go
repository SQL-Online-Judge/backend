package restapi

import (
	"net/http"
	"os"
	"time"

	"github.com/SQL-Online-Judge/backend/internal/core/repository"
	"github.com/SQL-Online-Judge/backend/internal/core/service"
	"github.com/SQL-Online-Judge/backend/internal/pkg/db/mongo"
	"github.com/go-chi/jwtauth/v5"
)

var tokenAuth *jwtauth.JWTAuth
var repo *repository.MongoRepository

var (
	userService       *service.UserService
	classService      *service.ClassService
	problemService    *service.ProblemService
	answerService     *service.AnswerService
	taskService       *service.TaskService
	submissionService *service.SubmissionService
)

func init() {
	tokenAuth = jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET")), nil)
	repo = repository.NewMongoRepository(mongo.GetMongoDB())

	userService = service.NewUserService(repo)
	classService = service.NewClassService(repo)
	problemService = service.NewProblemService(repo)
	answerService = service.NewAnswerService(repo)
	taskService = service.NewTaskService(repo)
	submissionService = service.NewSubmissionService(repo)
}

func Serve() {
	r := NewRouter()

	server := &http.Server{
		Addr:         ":3000",
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	server.ListenAndServe()
}
