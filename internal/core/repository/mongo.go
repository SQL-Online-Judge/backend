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

func (mr *MongoRepository) getAnswerCollection() *mongo.Collection {
	return mr.db.Collection("answer")
}

func (mr *MongoRepository) getTaskCollection() *mongo.Collection {
	return mr.db.Collection("task")
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
	defer cursor.Close(ctx)

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

func (mr *MongoRepository) ExistByProblemID(problemID int64) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "problemID", Value: problemID}}
	count, err := mr.getProblemCollection().CountDocuments(ctx, filter)
	if err != nil {
		logger.Logger.Error("failed to count documents", zap.Error(err))
		return false
	}

	return count > 0
}

func (mr *MongoRepository) IsProblemDeleted(problemID int64) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "problemID", Value: problemID}}
	var problem model.Problem
	err := mr.getProblemCollection().FindOne(ctx, filter).Decode(&problem)
	if err != nil {
		logger.Logger.Error("failed to find problem by problemID", zap.Int64("problemID", problemID), zap.Error(err))
		return true
	}

	return problem.Deleted
}

func (mr *MongoRepository) IsProblemAuthor(teacherID, problemID int64) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{
		{Key: "problemID", Value: problemID},
		{Key: "authorID", Value: teacherID},
	}
	count, err := mr.getProblemCollection().CountDocuments(ctx, filter)
	if err != nil {
		logger.Logger.Error("failed to count documents", zap.Error(err))
		return false
	}

	return count > 0
}

func (mr *MongoRepository) DeleteByProblemID(problemID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "problemID", Value: problemID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "deleted", Value: true}}}}
	_, err := mr.getProblemCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Logger.Error("failed to delete problem", zap.Int64("problemID", problemID), zap.Error(err))
		return fmt.Errorf("failed to delete problem: %w", err)
	}

	return nil
}

func (mr *MongoRepository) UpdateProblem(p *model.Problem) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "problemID", Value: p.ProblemID}}
	update := bson.D{{Key: "$set", Value: p}}
	_, err := mr.getProblemCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Logger.Error("failed to update problem", zap.Int64("problemID", p.ProblemID), zap.Error(err))
		return fmt.Errorf("failed to update problem: %w", err)
	}

	return nil
}

func (mr *MongoRepository) FindByProblemID(problemID int64) (*model.Problem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "problemID", Value: problemID}}
	var problem model.Problem
	err := mr.getProblemCollection().FindOne(ctx, filter).Decode(&problem)
	if err != nil {
		logger.Logger.Error("failed to find problem by problemID", zap.Int64("problemID", problemID), zap.Error(err))
		return nil, fmt.Errorf("failed to find problem by problemID: %w", err)
	}

	return &problem, nil
}

func (mr *MongoRepository) FindProblems(contains string) ([]*model.Problem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{
		{Key: "title", Value: bson.D{{Key: "$regex", Value: primitive.Regex{Pattern: contains, Options: "i"}}}},
		{Key: "deleted", Value: false},
	}
	cursor, err := mr.getProblemCollection().Find(ctx, filter)
	if err != nil {
		logger.Logger.Error("failed to get problems", zap.Error(err))
		return nil, fmt.Errorf("failed to get problems: %w", err)
	}
	defer cursor.Close(ctx)

	var problems []*model.Problem
	for cursor.Next(ctx) {
		var problem model.Problem
		err := cursor.Decode(&problem)
		if err != nil {
			logger.Logger.Error("failed to decode problem", zap.Error(err))
			return nil, fmt.Errorf("failed to decode problem: %w", err)
		}
		problems = append(problems, &problem)
	}

	return problems, nil
}

func (mr *MongoRepository) FindProblemsByAuthorID(authorID int64) ([]*model.Problem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{
		{Key: "authorID", Value: authorID},
		{Key: "deleted", Value: false},
	}
	cursor, err := mr.getProblemCollection().Find(ctx, filter)
	if err != nil {
		logger.Logger.Error("failed to get problems", zap.Error(err))
		return nil, fmt.Errorf("failed to get problems: %w", err)
	}
	defer cursor.Close(ctx)

	var problems []*model.Problem
	for cursor.Next(ctx) {
		var problem model.Problem
		err := cursor.Decode(&problem)
		if err != nil {
			logger.Logger.Error("failed to decode problem", zap.Error(err))
			return nil, fmt.Errorf("failed to decode problem: %w", err)
		}
		problems = append(problems, &problem)
	}

	return problems, nil
}

func (mr *MongoRepository) IsAnswerExist(problemID int64, dbName string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{
		{Key: "problemID", Value: problemID},
		{Key: "dbName", Value: dbName},
		{Key: "deleted", Value: false},
	}
	count, err := mr.getAnswerCollection().CountDocuments(ctx, filter)
	if err != nil {
		logger.Logger.Error("failed to count documents", zap.Error(err))
		return false
	}

	return count > 0
}

func (mr *MongoRepository) CreateAnswer(answer *model.Answer) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := mr.getAnswerCollection().InsertOne(ctx, answer)
	if err != nil {
		logger.Logger.Error("failed to create answer", zap.Error(err))
		return 0, fmt.Errorf("failed to create answer: %w", err)
	}

	return answer.AnswerID, nil
}

func (mr *MongoRepository) ExistByAnswerID(answerID int64) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "answerID", Value: answerID}}
	count, err := mr.getAnswerCollection().CountDocuments(ctx, filter)
	if err != nil {
		logger.Logger.Error("failed to count documents", zap.Error(err))
		return false
	}

	return count > 0
}

func (mr *MongoRepository) IsAnswerDeleted(answerID int64) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "answerID", Value: answerID}}
	var answer model.Answer
	err := mr.getAnswerCollection().FindOne(ctx, filter).Decode(&answer)
	if err != nil {
		logger.Logger.Error("failed to find answer by answerID", zap.Int64("answerID", answerID), zap.Error(err))
		return true
	}

	return answer.Deleted
}

func (mr *MongoRepository) DeleteByAnswerID(answerID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "answerID", Value: answerID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "deleted", Value: true}}}}
	_, err := mr.getAnswerCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Logger.Error("failed to delete answer", zap.Int64("answerID", answerID), zap.Error(err))
		return fmt.Errorf("failed to delete answer: %w", err)
	}

	return nil
}

func (mr *MongoRepository) IsAnswerOfProblem(problemID, answerID int64) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{
		{Key: "answerID", Value: answerID},
		{Key: "problemID", Value: problemID},
	}
	count, err := mr.getAnswerCollection().CountDocuments(ctx, filter)
	if err != nil {
		logger.Logger.Error("failed to count documents", zap.Error(err))
		return false
	}

	return count > 0
}

func (mr *MongoRepository) UpdateAnswer(answer *model.Answer) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "answerID", Value: answer.AnswerID}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "prepareSQL", Value: answer.PrepareSQL},
		{Key: "answerSQL", Value: answer.AnswerSQL},
		{Key: "judgeSQL", Value: answer.JudgeSQL},
		{Key: "isReady", Value: false},
	}}}
	_, err := mr.getAnswerCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Logger.Error("failed to update answer", zap.Int64("answerID", answer.AnswerID), zap.Error(err))
		return fmt.Errorf("failed to update answer: %w", err)
	}

	return nil
}

func (mr *MongoRepository) FindAnswersByProblemID(problemID int64) ([]*model.Answer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{
		{Key: "problemID", Value: problemID},
		{Key: "deleted", Value: false},
	}
	cursor, err := mr.getAnswerCollection().Find(ctx, filter)
	if err != nil {
		logger.Logger.Error("failed to get answers", zap.Error(err))
		return nil, fmt.Errorf("failed to get answers: %w", err)
	}
	defer cursor.Close(ctx)

	var answers []*model.Answer
	for cursor.Next(ctx) {
		var answer model.Answer
		err := cursor.Decode(&answer)
		if err != nil {
			logger.Logger.Error("failed to decode answer", zap.Error(err))
			return nil, fmt.Errorf("failed to decode answer: %w", err)
		}
		answers = append(answers, &answer)
	}

	return answers, nil
}

func (mr *MongoRepository) CreateTask(task *model.Task) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := mr.getTaskCollection().InsertOne(ctx, task)
	if err != nil {
		logger.Logger.Error("failed to create task", zap.Error(err))
		return 0, fmt.Errorf("failed to create task: %w", err)
	}

	return task.TaskID, nil
}

func (mr *MongoRepository) ExistByTaskID(taskID int64) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "taskID", Value: taskID}}
	count, err := mr.getTaskCollection().CountDocuments(ctx, filter)
	if err != nil {
		logger.Logger.Error("failed to count documents", zap.Error(err))
		return false
	}

	return count > 0
}

func (mr *MongoRepository) IsTaskDeleted(taskID int64) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "taskID", Value: taskID}}
	var task model.Task
	err := mr.getTaskCollection().FindOne(ctx, filter).Decode(&task)
	if err != nil {
		logger.Logger.Error("failed to find task by taskID", zap.Int64("taskID", taskID), zap.Error(err))
		return true
	}

	return task.Deleted
}

func (mr *MongoRepository) IsTaskAuthor(teacherID, taskID int64) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{
		{Key: "taskID", Value: taskID},
		{Key: "authorID", Value: teacherID},
	}
	count, err := mr.getTaskCollection().CountDocuments(ctx, filter)
	if err != nil {
		logger.Logger.Error("failed to count documents", zap.Error(err))
		return false
	}

	return count > 0
}

func (mr *MongoRepository) DeleteByTaskID(taskID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "taskID", Value: taskID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "deleted", Value: true}}}}
	_, err := mr.getTaskCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Logger.Error("failed to delete task", zap.Int64("taskID", taskID), zap.Error(err))
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}

func (mr *MongoRepository) UpdateTask(task *model.Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "taskID", Value: task.TaskID}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "taskName", Value: task.TaskName},
		{Key: "isTimeLimited", Value: task.IsTimeLimited},
		{Key: "beginTime", Value: task.BeginTime},
		{Key: "endTime", Value: task.EndTime},
	}}}
	_, err := mr.getTaskCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Logger.Error("failed to update task", zap.Int64("taskID", task.TaskID), zap.Error(err))
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}

func (mr *MongoRepository) IsTaskProblem(taskID, problemID int64) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{
		{Key: "taskID", Value: taskID},
		{Key: "problems.problemID", Value: problemID},
	}
	count, err := mr.getTaskCollection().CountDocuments(ctx, filter)
	if err != nil {
		logger.Logger.Error("failed to count documents", zap.Error(err))
		return false
	}

	return count > 0
}

func (mr *MongoRepository) AddTaskProblem(taskID int64, problem *model.TaskProblem) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "taskID", Value: taskID}}
	update := bson.D{{Key: "$addToSet", Value: bson.D{{Key: "problems", Value: problem}}}}
	_, err := mr.getTaskCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Logger.Error("failed to add task problem", zap.Int64("taskID", taskID), zap.Int64("problemID", problem.ProblemID), zap.Error(err))
		return fmt.Errorf("failed to add task problem: %w", err)
	}

	return nil
}

func (mr *MongoRepository) RemoveTaskProblem(taskID, problemID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "taskID", Value: taskID}}
	update := bson.D{{Key: "$pull", Value: bson.D{{Key: "problems", Value: bson.D{{Key: "problemID", Value: problemID}}}}}}
	_, err := mr.getTaskCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Logger.Error("failed to remove task problem", zap.Int64("taskID", taskID), zap.Int64("problemID", problemID), zap.Error(err))
		return fmt.Errorf("failed to remove task problem: %w", err)
	}

	return nil
}

func (mr *MongoRepository) FindByTaskID(taskID int64) (*model.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "taskID", Value: taskID}}
	var task model.Task
	err := mr.getTaskCollection().FindOne(ctx, filter).Decode(&task)
	if err != nil {
		logger.Logger.Error("failed to find task by taskID", zap.Int64("taskID", taskID), zap.Error(err))
		return nil, fmt.Errorf("failed to find task by taskID: %w", err)
	}

	return &task, nil
}

func (mr *MongoRepository) FindTasks(contains string) ([]*model.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{
		{Key: "taskName", Value: bson.D{{Key: "$regex", Value: primitive.Regex{Pattern: contains, Options: "i"}}}},
		{Key: "deleted", Value: false},
	}
	cursor, err := mr.getTaskCollection().Find(ctx, filter)
	if err != nil {
		logger.Logger.Error("failed to get tasks", zap.Error(err))
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}
	defer cursor.Close(ctx)

	var tasks []*model.Task
	for cursor.Next(ctx) {
		var task model.Task
		err := cursor.Decode(&task)
		if err != nil {
			logger.Logger.Error("failed to decode task", zap.Error(err))
			return nil, fmt.Errorf("failed to decode task: %w", err)
		}
		tasks = append(tasks, &task)
	}

	return tasks, nil
}

func (mr *MongoRepository) FindTasksByAuthorID(authorID int64) ([]*model.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{
		{Key: "authorID", Value: authorID},
		{Key: "deleted", Value: false},
	}
	cursor, err := mr.getTaskCollection().Find(ctx, filter)
	if err != nil {
		logger.Logger.Error("failed to get tasks", zap.Error(err))
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}
	defer cursor.Close(ctx)

	var tasks []*model.Task
	for cursor.Next(ctx) {
		var task model.Task
		err := cursor.Decode(&task)
		if err != nil {
			logger.Logger.Error("failed to decode task", zap.Error(err))
			return nil, fmt.Errorf("failed to decode task: %w", err)
		}
		tasks = append(tasks, &task)
	}

	return tasks, nil
}

func (mr *MongoRepository) IsClassTask(classID, taskID int64) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{
		{Key: "classID", Value: classID},
		{Key: "tasks", Value: bson.D{{Key: "$in", Value: []int64{taskID}}}},
	}
	count, err := mr.getClassCollection().CountDocuments(ctx, filter)
	if err != nil {
		logger.Logger.Error("failed to count documents", zap.Error(err))
		return false
	}

	return count > 0
}

func (mr *MongoRepository) AddTaskToClass(classID, taskID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "classID", Value: classID}}
	update := bson.D{{Key: "$addToSet", Value: bson.D{{Key: "tasks", Value: taskID}}}}
	_, err := mr.getClassCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Logger.Error("failed to add task to class", zap.Int64("taskID", taskID), zap.Int64("classID", classID), zap.Error(err))
		return fmt.Errorf("failed to add task to class: %w", err)
	}

	return nil
}

func (mr *MongoRepository) RemoveTaskFromClass(classID, taskID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "classID", Value: classID}}
	update := bson.D{{Key: "$pull", Value: bson.D{{Key: "tasks", Value: taskID}}}}
	_, err := mr.getClassCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Logger.Error("failed to remove task from class", zap.Int64("taskID", taskID), zap.Int64("classID", classID), zap.Error(err))
		return fmt.Errorf("failed to remove task from class: %w", err)
	}

	return nil
}

func (mr *MongoRepository) GetTasksInClass(classID int64) ([]*model.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "classID", Value: classID}}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "task"},
			{Key: "localField", Value: "tasks"},
			{Key: "foreignField", Value: "taskID"},
			{Key: "as", Value: "tasks"},
		}}},
		{{Key: "$unwind", Value: "$tasks"}},
		{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "taskID", Value: "$tasks.taskID"},
			{Key: "taskName", Value: "$tasks.taskName"},
			{Key: "isTimeLimited", Value: "$tasks.isTimeLimited"},
			{Key: "beginTime", Value: "$tasks.beginTime"},
			{Key: "endTime", Value: "$tasks.endTime"},
		}}},
	}

	cursor, err := mr.getClassCollection().Aggregate(ctx, pipeline)
	if err != nil {
		logger.Logger.Error("failed to aggregate", zap.Error(err))
		return nil, fmt.Errorf("failed to aggregate: %w", err)
	}
	defer cursor.Close(ctx)

	var tasks []*model.Task
	err = cursor.All(ctx, &tasks)
	if err != nil {
		logger.Logger.Error("failed to decode tasks", zap.Error(err))
		return nil, fmt.Errorf("failed to decode tasks: %w", err)
	}

	return tasks, nil
}

func (mr *MongoRepository) FindByClassID(classID int64) (*model.Class, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "classID", Value: classID}}
	var class model.Class
	err := mr.getClassCollection().FindOne(ctx, filter).Decode(&class)
	if err != nil {
		logger.Logger.Error("failed to find class by classID", zap.Int64("classID", classID), zap.Error(err))
		return nil, fmt.Errorf("failed to find class by classID: %w", err)
	}

	return &class, nil
}

func (mr *MongoRepository) getStudentTasksPipeline(studentID int64) mongo.Pipeline {
	return mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "students", Value: studentID}}}},
		{{Key: "$unwind", Value: "$tasks"}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$tasks"},
		}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "task"},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "taskID"},
			{Key: "as", Value: "tasks"},
		}}},
		{{Key: "$match", Value: bson.D{{Key: "tasks.deleted", Value: false}}}},
		{{Key: "$unwind", Value: "$tasks"}},
		{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "taskID", Value: "$tasks.taskID"},
			{Key: "taskName", Value: "$tasks.taskName"},
			{Key: "problems", Value: "$tasks.problems"},
			{Key: "isTimeLimited", Value: "$tasks.isTimeLimited"},
			{Key: "beginTime", Value: "$tasks.beginTime"},
			{Key: "endTime", Value: "$tasks.endTime"},
		}}},
	}
}

func (mr *MongoRepository) FindTasksByStudentID(studentID int64) ([]*model.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipeline := mr.getStudentTasksPipeline(studentID)
	cursor, err := mr.getClassCollection().Aggregate(ctx, pipeline)
	if err != nil {
		logger.Logger.Error("failed to aggregate", zap.Error(err))
		return nil, fmt.Errorf("failed to aggregate: %w", err)
	}
	defer cursor.Close(ctx)

	var tasks []*model.Task
	err = cursor.All(ctx, &tasks)
	if err != nil {
		logger.Logger.Error("failed to decode tasks", zap.Error(err))
		return nil, fmt.Errorf("failed to decode tasks: %w", err)
	}

	return tasks, nil
}

func (mr *MongoRepository) CanStudentAccessTask(studentID, taskID int64) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{
		{Key: "students", Value: studentID},
		{Key: "tasks", Value: taskID},
		{Key: "deleted", Value: false},
	}
	count, err := mr.getClassCollection().CountDocuments(ctx, filter)
	if err != nil {
		logger.Logger.Error("failed to count documents", zap.Error(err))
		return false
	}

	return count > 0
}

func (mr *MongoRepository) getStudentTaskProblemsPipeline(studentID, taskID int64) mongo.Pipeline {
	pipeline := mr.getStudentTasksPipeline(studentID)
	pipeline = append(pipeline, mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "taskID", Value: taskID}}}},
		{{Key: "$unwind", Value: "$problems"}},
		{{Key: "$project", Value: bson.D{
			{Key: "problemID", Value: "$problems.problemID"},
			{Key: "score", Value: "$problems.score"},
		}}},
	}...)
	return pipeline
}

func (mr *MongoRepository) FindTaskProblemsByStudentIDAndTaskID(studentID, taskID int64) ([]*model.TaskProblem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipeline := mr.getStudentTaskProblemsPipeline(studentID, taskID)

	cursor, err := mr.getClassCollection().Aggregate(ctx, pipeline)
	if err != nil {
		logger.Logger.Error("failed to aggregate", zap.Error(err))
		return nil, fmt.Errorf("failed to aggregate: %w", err)
	}
	defer cursor.Close(ctx)

	var taskProblems []*model.TaskProblem
	err = cursor.All(ctx, &taskProblems)
	if err != nil {
		logger.Logger.Error("failed to decode task", zap.Error(err))
		return nil, fmt.Errorf("failed to decode task: %w", err)
	}

	return taskProblems, nil
}

func (mr *MongoRepository) getProblemsInStudentTaskPipeline(studentID, taskID int64) mongo.Pipeline {
	pipeline := mr.getStudentTaskProblemsPipeline(studentID, taskID)
	pipeline = append(pipeline, mongo.Pipeline{
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "problem"},
			{Key: "localField", Value: "problemID"},
			{Key: "foreignField", Value: "problemID"},
			{Key: "as", Value: "problems"},
		}}},
		{{Key: "$unwind", Value: "$problems"}},
		{{Key: "$match", Value: bson.D{{Key: "problems.deleted", Value: false}}}},
		{{Key: "$project", Value: bson.D{
			{Key: "problemID", Value: "$problems.problemID"},
			{Key: "title", Value: "$problems.title"},
			{Key: "tags", Value: "$problems.tags"},
			{Key: "content", Value: "$problems.content"},
			{Key: "timeLimit", Value: "$problems.timeLimit"},
			{Key: "memoryLimit", Value: "$problems.memoryLimit"},
		}}},
	}...)
	return pipeline
}

func (mr *MongoRepository) FindProblemsInStudentTask(studentID, taskID int64) ([]*model.Problem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipeline := mr.getProblemsInStudentTaskPipeline(studentID, taskID)

	cursor, err := mr.getClassCollection().Aggregate(ctx, pipeline)
	if err != nil {
		logger.Logger.Error("failed to aggregate", zap.Error(err))
		return nil, fmt.Errorf("failed to aggregate: %w", err)
	}
	defer cursor.Close(ctx)

	var problems []*model.Problem
	err = cursor.All(ctx, &problems)
	if err != nil {
		logger.Logger.Error("failed to decode problems", zap.Error(err))
		return nil, fmt.Errorf("failed to decode problems: %w", err)
	}

	return problems, nil
}
