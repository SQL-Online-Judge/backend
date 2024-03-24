package repository

import "github.com/SQL-Online-Judge/backend/internal/model"

type Repository interface {
	UserRepository
}

type UserRepository interface {
	CreateUser(username, password, role string) (int64, error)
	FindByUserID(userID int64) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
	ExistByUserID(userID int64) bool
	ExistByUsername(username string) bool
	DeleteByUserID(userID int64) error
	GetRoleByUserID(userID int64) (string, error)
	UpdateUsernameByUserID(userID int64, username string) error
	IsDeletedByUserID(userID int64) bool
	GetStudents(contains string) ([]*model.User, error)
}

type ClassRepository interface {
	CreateClass(className string, teacherID int64) (int64, error)
	ExistByClassID(classID int64) bool
	IsClassOwner(teacherID, classID int64) bool
	DeleteByClassID(classID int64) error
	IsClassDeleted(classID int64) bool
	UpdateClassNameByClassID(classID int64, className string) error
	FindClassesByTeacherID(teacherID int64) ([]*model.Class, error)
	IsClassMember(classID, studentID int64) bool
	AddStudentToClass(classID, studentID int64) error
	RemoveStudentFromClass(classID, studentID int64) error
	FindStudentsByClassID(classID int64) ([]*model.User, error)
	IsClassTask(classID, taskID int64) bool
	AddTaskToClass(classID, taskID int64) error
	RemoveTaskFromClass(classID, taskID int64) error
	GetTasksInClass(classID int64) ([]*model.Task, error)
	FindByClassID(classID int64) (*model.Class, error)
}

type ProblemRepository interface {
	CreateProblem(p *model.Problem) (int64, error)
	ExistByProblemID(problemID int64) bool
	IsProblemDeleted(problemID int64) bool
	IsProblemAuthor(teacherID, problemID int64) bool
	DeleteByProblemID(problemID int64) error
	UpdateProblem(p *model.Problem) error
	FindByProblemID(problemID int64) (*model.Problem, error)
	FindProblemsByAuthorID(authorID int64) ([]*model.Problem, error)
	FindProblems(contains string) ([]*model.Problem, error)
}

type AnswerRepository interface {
	IsAnswerExist(problemID int64, dbName string) bool
	CreateAnswer(answer *model.Answer) (int64, error)
	ExistByAnswerID(answerID int64) bool
	IsAnswerDeleted(answerID int64) bool
	DeleteByAnswerID(answerID int64) error
	IsAnswerOfProblem(problemID, answerID int64) bool
	UpdateAnswer(answer *model.Answer) error
	FindAnswersByProblemID(problemID int64) ([]*model.Answer, error)
	// FindByAnswerID(answerID int64) (*model.Answer, error)
}

type TaskRepository interface {
	CreateTask(t *model.Task) (int64, error)
	ExistByTaskID(taskID int64) bool
	IsTaskDeleted(taskID int64) bool
	IsTaskAuthor(teacherID, taskID int64) bool
	DeleteByTaskID(taskID int64) error
	UpdateTask(t *model.Task) error
	IsTaskProblem(taskID, problemID int64) bool
	AddTaskProblem(taskID int64, problem *model.TaskProblem) error
	RemoveTaskProblem(taskID, problemID int64) error
	FindByTaskID(taskID int64) (*model.Task, error)
	FindTasks(contains string) ([]*model.Task, error)
	FindTasksByAuthorID(authorID int64) ([]*model.Task, error)
	FindTasksByStudentID(studentID int64) ([]*model.Task, error)
	CanStudentAccessTask(studentID, taskID int64) bool
	FindTaskProblemsByStudentIDAndTaskID(studentID, taskID int64) ([]*model.TaskProblem, error)
	FindProblemsInStudentTask(studentID, taskID int64) ([]*model.Problem, error)
}
