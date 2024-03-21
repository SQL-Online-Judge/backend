package restapi

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(contentTypeJSON)

	r.Group(func(r chi.Router) {
		r.Get("/", sayHello("SQL-Online-Judge"))
		r.Post("/login", login)
	})

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator(tokenAuth))
		r.Use(getRole)

		r.Route("/admin", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(checkRole("admin"))
				r.Post("/teacher", createTeacher)
			})

			r.Group(func(r chi.Router) {
				r.Use(checkRole("teacher"))
				r.Post("/students", createStudents)
				r.Delete("/students", deleteStudents)
				r.Patch("/students/{userID}", updateStudentUsername)
				r.Get("/students/{userID}", getStudent)
				r.Get("/students", getStudents)

				r.Post("/classes", createClass)
				r.Delete("/classes/{classID}", deleteClass)
				r.Patch("/classes/{classID}", updateClassName)
				r.Get("/classes", getClasses)

				r.Post("/classes/{classID}/students", addStudentsToClass)
				r.Delete("/classes/{classID}/students", removeStudentsFromClass)
				r.Get("/classes/{classID}/students", getStudentsInClass)

				r.Post("/problems", createProblem)
				r.Delete("/problems/{problemID}", deleteProblem)
				r.Put("/problems/{problemID}", updateProblem)
				r.Get("/problems/{problemID}", getProblem)
				r.Get("/problems", getProblems)
				r.Get("/my/problems", getTeacherProblems)

				r.Post("/problems/{problemID}/answers", createAnswer)
				r.Delete("/problems/{problemID}/answers/{answerID}", deleteAnswer)
				r.Put("/problems/{problemID}/answers/{answerID}", updateAnswer)
				r.Get("/problems/{problemID}/answers", getAnswers)

				r.Post("/tasks", createTask)
				r.Delete("/tasks/{taskID}", deleteTask)
				r.Put("/tasks/{taskID}", updateTask)
				r.Post("/tasks/{taskID}/problems", addProblemsToTask)
				r.Delete("/tasks/{taskID}/problems", removeProblemsFromTask)
				r.Get("/tasks/{taskID}", getTask)
				r.Get("/tasks", getTasks)
				r.Get("/my/tasks", getTeacherTasks)

				r.Post("/classes/{classID}/tasks", addTasksToClass)
				r.Delete("/classes/{classID}/tasks", removeTasksFromClass)
			})
		})
	})

	return r
}
