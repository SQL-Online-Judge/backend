package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/SQL-Online-Judge/backend/internal/model"
	"github.com/SQL-Online-Judge/backend/internal/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var ErrUserIsNil = fmt.Errorf("user is nil")

type MongoRepository struct {
	db *mongo.Database
}

func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{
		db: db,
	}
}

func (mr *MongoRepository) getUserCollection() *mongo.Collection {
	return mr.db.Collection("user")
}

func (mr *MongoRepository) getClassCollection() *mongo.Collection {
	return mr.db.Collection("class")
}

func (mr *MongoRepository) getProblemCollection() *mongo.Collection {
	return mr.db.Collection("problem")
}

func (mr *MongoRepository) ExistByUserID(userID int64) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.D{{Key: "userID", Value: userID}}
	count, err := mr.getUserCollection().CountDocuments(ctx, filter)
	if err != nil {
		logger.Logger.Error("failed to count documents", zap.Error(err))
		return false
	}

	return count > 0
}

func (mr *MongoRepository) ExistByUsername(username string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.D{{Key: "username", Value: username}}
	count, err := mr.getUserCollection().CountDocuments(ctx, filter)
	if err != nil {
		logger.Logger.Error("failed to count documents", zap.Error(err))
		return false
	}

	return count > 0
}

func (mr *MongoRepository) CreateUser(username, password, role string) (int64, error) {
	user := model.NewUser(username, password, role)
	hashedPassword, err := user.GetHashedPassword()
	if err != nil {
		return 0, fmt.Errorf("failed to get hashed password: %w", err)
	}
	user.Password = hashedPassword

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = mr.getUserCollection().InsertOne(ctx, user)
	if err != nil {
		logger.Logger.Error("failed to create user", zap.String("username", user.Username), zap.Error(err))
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return user.UserID, nil
}

func (mr *MongoRepository) FindByUserID(userID int64) (*model.User, error) {
	var user model.User
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.D{
		{Key: "userID", Value: userID},
		{Key: "deleted", Value: false},
	}
	err := mr.getUserCollection().FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by user id: %w", err)
	}

	return &user, nil
}

func (mr *MongoRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.D{
		{Key: "username", Value: username},
		{Key: "deleted", Value: false},
	}
	err := mr.getUserCollection().FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by username: %w", err)
	}

	return &user, nil
}

func (mr *MongoRepository) DeleteByUserID(userID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.D{{Key: "userID", Value: userID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "deleted", Value: true}}}}
	_, err := mr.getUserCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Logger.Error("failed to delete user", zap.Int64("userID", userID), zap.Error(err))
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (mr *MongoRepository) GetRoleByUserID(userID int64) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.D{{Key: "userID", Value: userID}}
	var user model.User
	err := mr.getUserCollection().FindOne(ctx, filter).Decode(&user)
	if err != nil {
		logger.Logger.Error("failed to get role by userID", zap.Int64("userID", userID), zap.Error(err))
		return "", fmt.Errorf("failed to get role by userID: %w", err)
	}

	return user.Role, nil
}

func (mr *MongoRepository) UpdateUsernameByUserID(userID int64, username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.D{{Key: "userID", Value: userID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "username", Value: username}}}}
	_, err := mr.getUserCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Logger.Error("failed to update username", zap.Int64("userID", userID), zap.String("username", username), zap.Error(err))
		return fmt.Errorf("failed to update username: %w", err)
	}

	return nil
}

func (mr *MongoRepository) IsDeletedByUserID(userID int64) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.D{{Key: "userID", Value: userID}}
	var user model.User
	err := mr.getUserCollection().FindOne(ctx, filter).Decode(&user)
	if err != nil {
		logger.Logger.Error("failed to find user by userID", zap.Int64("userID", userID), zap.Error(err))
		return true
	}

	return user.Deleted
}

func (mr *MongoRepository) GetStudents(contains string) ([]*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.D{
		{Key: "role", Value: "student"},
		{Key: "username", Value: bson.D{{Key: "$regex", Value: primitive.Regex{Pattern: contains, Options: "i"}}}},
		{Key: "deleted", Value: false},
	}
	options := &options.FindOptions{
		Projection: bson.D{
			{Key: "password", Value: 0},
		},
	}
	cursor, err := mr.getUserCollection().Find(ctx, filter, options)
	if err != nil {
		logger.Logger.Error("failed to get students", zap.Error(err))
		return nil, fmt.Errorf("failed to get students: %w", err)
	}
	defer cursor.Close(ctx)

	var students []*model.User
	for cursor.Next(ctx) {
		var student model.User
		err := cursor.Decode(&student)
		if err != nil {
			logger.Logger.Error("failed to decode student", zap.Error(err))
			return nil, fmt.Errorf("failed to decode student: %w", err)
		}
		students = append(students, &student)
	}

	return students, nil
}

func (mr *MongoRepository) CreateClass(className string, teacherID int64) (int64, error) {
	class := model.NewClass(className, teacherID)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := mr.getClassCollection().InsertOne(ctx, class)
	if err != nil {
		logger.Logger.Error("failed to create class", zap.String("className", className), zap.Int64("teacherID", teacherID), zap.Error(err))
		return 0, fmt.Errorf("failed to create class: %w", err)
	}

	return class.ClassID, nil
}

func (mr *MongoRepository) ExistByClassID(classID int64) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.D{{Key: "classID", Value: classID}}
	count, err := mr.getClassCollection().CountDocuments(ctx, filter)
	if err != nil {
		logger.Logger.Error("failed to count documents", zap.Error(err))
		return false
	}

	return count > 0
}

func (mr *MongoRepository) IsClassOwner(teacherID, classID int64) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.D{
		{Key: "classID", Value: classID},
		{Key: "teacherID", Value: teacherID},
	}
	count, err := mr.getClassCollection().CountDocuments(ctx, filter)
	if err != nil {
		logger.Logger.Error("failed to count documents", zap.Error(err))
		return false
	}

	return count > 0
}

func (mr *MongoRepository) DeleteByClassID(classID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.D{{Key: "classID", Value: classID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "deleted", Value: true}}}}
	_, err := mr.getClassCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Logger.Error("failed to delete class", zap.Int64("classID", classID), zap.Error(err))
		return fmt.Errorf("failed to delete class: %w", err)
	}

	return nil
}

func (mr *MongoRepository) IsClassDeleted(classID int64) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.D{{Key: "classID", Value: classID}}
	var class model.Class
	err := mr.getClassCollection().FindOne(ctx, filter).Decode(&class)
	if err != nil {
		logger.Logger.Error("failed to find class by classID", zap.Int64("classID", classID), zap.Error(err))
		return true
	}

	return class.Deleted
}

func (mr *MongoRepository) UpdateClassNameByClassID(classID int64, className string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.D{{Key: "classID", Value: classID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "className", Value: className}}}}
	_, err := mr.getClassCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Logger.Error("failed to update class name", zap.Int64("classID", classID), zap.String("className", className), zap.Error(err))
		return fmt.Errorf("failed to update class name: %w", err)
	}

	return nil
}

func (mr *MongoRepository) FindClassesByTeacherID(teacherID int64) ([]*model.Class, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.D{
		{Key: "teacherID", Value: teacherID},
		{Key: "deleted", Value: false},
	}
	cursor, err := mr.getClassCollection().Find(ctx, filter)
	if err != nil {
		logger.Logger.Error("failed to get classes", zap.Error(err))
		return nil, fmt.Errorf("failed to get classes: %w", err)
	}
	defer cursor.Close(ctx)

	var classes []*model.Class
	for cursor.Next(ctx) {
		var class model.Class
		err := cursor.Decode(&class)
		if err != nil {
			logger.Logger.Error("failed to decode class", zap.Error(err))
			return nil, fmt.Errorf("failed to decode class: %w", err)
		}
		classes = append(classes, &class)
	}

	return classes, nil
}

func (mr *MongoRepository) IsClassMember(classID, studentID int64) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{
		{Key: "classID", Value: classID},
		{Key: "students", Value: bson.D{{Key: "$in", Value: []int64{studentID}}}},
	}
	count, err := mr.getClassCollection().CountDocuments(ctx, filter)
	if err != nil {
		logger.Logger.Error("failed to count documents", zap.Error(err))
		return false
	}

	return count > 0
}

func (mr *MongoRepository) AddStudentToClass(classID int64, studentID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "classID", Value: classID}}
	update := bson.D{{Key: "$addToSet", Value: bson.D{{Key: "students", Value: studentID}}}}
	_, err := mr.getClassCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Logger.Error("failed to add students to class", zap.Int64("classID", classID), zap.Int64("studentID", studentID), zap.Error(err))
		return fmt.Errorf("failed to add students to class: %w", err)
	}

	return nil
}

func (mr *MongoRepository) RemoveStudentFromClass(classID int64, studentID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "classID", Value: classID}}
	update := bson.D{{Key: "$pull", Value: bson.D{{Key: "students", Value: studentID}}}}
	_, err := mr.getClassCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Logger.Error("failed to remove students from class", zap.Int64("classID", classID), zap.Int64("studentID", studentID), zap.Error(err))
		return fmt.Errorf("failed to remove students from class: %w", err)
	}

	return nil
}

func (mr *MongoRepository) FindStudentsByClassID(classID int64) ([]*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "classID", Value: classID}}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "user"},
			{Key: "localField", Value: "students"},
			{Key: "foreignField", Value: "userID"},
			{Key: "as", Value: "students"},
		}}},
		{{Key: "$unwind", Value: "$students"}},
		{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "userID", Value: "$students.userID"},
			{Key: "username", Value: "$students.username"},
		}}},
	}

	cursor, err := mr.getClassCollection().Aggregate(ctx, pipeline)
	if err != nil {
		logger.Logger.Error("failed to aggregate", zap.Error(err))
		return nil, fmt.Errorf("failed to aggregate: %w", err)
	}

	var students []*model.User
	err = cursor.All(ctx, &students)
	if err != nil {
		logger.Logger.Error("failed to decode students", zap.Error(err))
		return nil, fmt.Errorf("failed to decode students: %w", err)
	}

	return students, nil
}

func (mr *MongoRepository) CreateProblem(p *model.Problem) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := mr.getProblemCollection().InsertOne(ctx, p)
	if err != nil {
		logger.Logger.Error("failed to create problem", zap.Error(err))
		return 0, fmt.Errorf("failed to create problem: %w", err)
	}

	return p.ProblemID, nil
}
