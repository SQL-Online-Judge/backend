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
}
