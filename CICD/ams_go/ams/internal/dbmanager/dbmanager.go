// Copyright (c) Huawei Technologies Co., Ltd. 2020-2020. All rights reserved.
// File Name: dbmanager.go
// Description: Contians DB manager functions

// Package dbmanager implements DB management functions
package dbmanager

import (
	"codehub-dg-g.huawei.com/IntelligentVision/IVM/AMS/common/constants"
	"codehub-dg-g.huawei.com/IntelligentVision/IVM/AMS/internal/configmanager"
	"codehub-dg-g.huawei.com/IntelligentVision/IVM/Resources/security/sccapi"
	"codehub-dg-g.huawei.com/IntelligentVision/IVM/Resources/utils/zaplog"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// Default page size
const (
	DefaultLimit uint32 = 10
)

type DBConnectionInfo struct {
	collection *mongo.Collection
	client     *mongo.Client
	ctx        context.Context
}
type DBManager struct {
	inited bool
	DBConn DBConnectionInfo
}

type DeviceInfo struct {
	DeviceId  string `bson:"device_id"`
	ChannelId string `bson:"channel_id"`
}

// Alert task info db-table structure
type AlertTaskInfo struct {
	UserId       string       `bson:"user_id"`
	EnterpriseId string       `bson:"enterprise_id"`
	TaskId       string       `bson:"task_id"`
	TaskName     string       `bson:"task_name"`
	State        uint32       `bson:"state"`
	Threshold    float64      `bson:"threshold"`
	LibId        string       `bson:"lib_id"`
	Devices      []DeviceInfo `bson:"devices"`
	StartTime    time.Time    `bson:"start_time"`
	EndTime      time.Time    `bson:"end_time"`
	CreateTime   time.Time    `bson:"create_time"`
	UpdateTime   time.Time    `bson:"update_time"`
	IndexName    []string     `bson:"index_name"`
}

// Query db based on these parameters
type QueryPara struct {
	Offset uint32
	Limit uint32
	SortDir string
	SortKey string
}

// Get DB manager instance
func NewDBManager(conf *configmanager.MongoDBInfo) (*DBManager, error) {
	var db *DBManager = &DBManager{}

	if db == nil {
		err := fmt.Errorf("db manager create instance failed")
		zaplog.Error(err)
		return nil, err
	}
	err := db.Init(conf)
	if err != nil {
		newErr := fmt.Errorf("db manager create instance failed, [%s]", err)
		zaplog.Error(newErr)
		return nil, newErr
	}

	return db, nil
}

// Init DB manager instance
func (db *DBManager) Init(conf *configmanager.MongoDBInfo) error {
	if db.inited == true {
		return nil
	}

	err := db.Connect(conf)
	if err != nil {
		newErr := fmt.Errorf("db connect failed, [%s]", err)
		zaplog.Error(newErr)
		return newErr
	}

	db.inited = true
	return nil
}

// Delete DB manager instance
func (db *DBManager) DeInit() {
	if db.inited == false {
		return
	}

	err := db.Disconnect()
	if err != nil {
		newErr := fmt.Errorf("db disconnect failed, [%s]", err)
		zaplog.Error(newErr)
	}

	db.inited = false
	return
}

// Connect to DB.
func (db *DBManager) Connect(conf *configmanager.MongoDBInfo) error {
	var err error

	if len(conf.UserName) == 0 || len(conf.Password) == 0 {
		newErr := fmt.Errorf("mongo db credential not set")
		zaplog.Error(newErr)
		return newErr
	}

	port := fmt.Sprint(conf.Port)
	applyURI := "mongodb://" + conf.IP + ":" + port

	clientOptions := options.Client().ApplyURI(applyURI)
	plainText, err := sccapi.Decrypt(conf.Password)
	if err != nil {
		newErr := fmt.Errorf("mongo decrypt failed, [%s]", err)
		zaplog.Error(newErr)
		return newErr
	}
	credential := options.Credential{
		Username: conf.UserName,
		Password: plainText,
	}

	clientOptions = clientOptions.SetAuth(credential)
	if constants.EnvConfig.Mongo.EnableSsl {
		clientOptions, err = setSslOpt(clientOptions)
		if err != nil {
			return errors.New("set ssl error")
		}
	}

	db.DBConn.client, err = mongo.Connect(db.DBConn.ctx, clientOptions)
	if err != nil {
		newErr := fmt.Errorf("db connect operation failed")
		zaplog.Error(newErr)
		return newErr
	}

	err = db.DBConn.client.Ping(db.DBConn.ctx, nil)
	if err != nil {
		newErr := fmt.Errorf("db ping operation failed")
		zaplog.Error(newErr)
		return newErr
	}

	zaplog.Info("Connected to MongoDB!")

	db.DBConn.collection = db.DBConn.client.Database("alert_manager").
		Collection("alert_tasks_info")

	return nil
}

func setSslOpt(clientOptions *options.ClientOptions) (*options.ClientOptions, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	clientOptions = clientOptions.SetTLSConfig(tlsConfig)
	return clientOptions, nil
}

// Disconnect from DB
func (db *DBManager) Disconnect() error {
	err := db.DBConn.client.Disconnect(db.DBConn.ctx)
	if err != nil {
		newErr := fmt.Errorf("db disconnect operation failed, [%s]", err)
		zaplog.Error(newErr)
		return newErr
	}

	zaplog.Info("Connection to MongoDB closed.")

	return nil
}

// Save Task
func (db *DBManager) SaveTask(task *AlertTaskInfo) error {
	insertResult, err := db.DBConn.collection.InsertOne(db.DBConn.ctx, task)
	if err != nil {
		newErr := fmt.Errorf("db insert operation failed, [%s]", err)
		zaplog.Error(newErr)
		return newErr
	}

	zaplog.Info("Inserted a single document: ", insertResult.InsertedID)

	return nil
}

// Query Task List
func (db *DBManager) QueryTaskList(userId string, enterpriseId string, queryPara *QueryPara) (
	[]*AlertTaskInfo, error) {

	// filter by user id and enterprise id
	filter := bson.M{
		"user_id": userId,
		"enterprise_id": enterpriseId,
	}

	tasks, err := db.QueryTaskListInternal(filter, queryPara)
	if err != nil {
		newErr := fmt.Errorf("query task list function error, [%s]", err)
		zaplog.Error(newErr)
		return nil, newErr
	}

	return tasks, err
}

// Query TaskList By State
func (db *DBManager) QueryTaskListByState(state uint32,
	queryPara *QueryPara) ([]*AlertTaskInfo, error) {


	filter := bson.D{
		primitive.E{Key: "state", Value: state},
	}

	tasks, err := db.QueryTaskListInternal(filter, queryPara)
	if err != nil {
		newErr := fmt.Errorf("query task list by state function error, [%s]", err)
		zaplog.Error(newErr)
		return nil, newErr
	}

	return tasks, err
}

// Query TaskList By Start time
func (db *DBManager) QueryTaskListByStartTime(userId string, enterpriseId string, 
	startTime time.Time, queryPara *QueryPara) ([]*AlertTaskInfo, error) {

	filter := bson.M{
		"user_id": userId,
		"enterprise_id": enterpriseId,
		"start_time": startTime,
	}


	tasks, err := db.QueryTaskListInternal(filter, queryPara)
	if err != nil {
		newErr := fmt.Errorf("query task list by start time function error, [%s]", err)
		zaplog.Error(newErr)
		return nil, newErr
	}

	return tasks, err
}

// Internal Query Task List
func (db *DBManager) QueryTaskListInternal(filter interface{}, 
	queryPara *QueryPara) ([]*AlertTaskInfo, error) {

	var sortKey string = queryPara.SortKey
	var limit uint32 = queryPara.Limit
	
	if limit == 0 {
		// default limit
		limit = DefaultLimit
	}

	options := options.Find()

	// Skip the documents till offset
	options.SetSkip(int64(queryPara.Offset * limit))

	// Set limit
	options.SetLimit(int64(limit))
	
	if sortKey == "" {
		sortKey = "create_time"
	}

	// Sort result in ascending or descending by sort key
	sortMap := make(map[string]interface{}, 1)
	if queryPara.SortDir == "ascending" {
		sortMap[sortKey] = 1
	} else if queryPara.SortDir == "descending" {
		sortMap[sortKey] = -1
	} else {
		sortMap[sortKey] = -1 // default sort in descending order
	}

	options.SetSort(sortMap)

	tasks, err := db.filterTasks(filter, options)
	if err != nil {
		newErr := fmt.Errorf("db filtertask function error, [%s]", err)
		zaplog.Error(newErr)
		return nil, newErr
	}

	return tasks, err
}

func (db *DBManager) filterTasks(filter interface{}, options *options.FindOptions) ([]*AlertTaskInfo, error) {
	// A slice of tasks for storing the decoded documents
	var tasks []*AlertTaskInfo

	cur, err := db.DBConn.collection.Find(db.DBConn.ctx, filter, options)
	if err != nil {
		newErr := fmt.Errorf("db find operation failed, [%s]", err)
		zaplog.Error(newErr)
		return nil, newErr
	}

	for cur.Next(db.DBConn.ctx) {
		var t AlertTaskInfo
		err := cur.Decode(&t)
		if err != nil {
			cur.Close(db.DBConn.ctx)
			newErr := fmt.Errorf("decode operation failed, [%s]", err)
			zaplog.Error(newErr)
			return nil, newErr
		}

		tasks = append(tasks, &t)
	}

	if err := cur.Err(); err != nil {
		cur.Close(db.DBConn.ctx)
		newErr := fmt.Errorf("db find operation failed, cursor error, [%s]", err)
		zaplog.Error(newErr)
		return nil, newErr
	}

	// once exhausted, close the cursor
	cur.Close(db.DBConn.ctx)

	if len(tasks) == 0 {
		zaplog.Info("filterTasks: No matching record found")
		return nil, nil
	}

	return tasks, nil
}

// Query Task By Id
func (db *DBManager) QueryTaskById(userId string, enterpriseId string,
	taskId string) (*AlertTaskInfo, error) {
	var task AlertTaskInfo

	filter := bson.M{
		"user_id": userId,
		"enterprise_id": enterpriseId,
		"task_id": taskId,
	}

	// find a single document
	err := db.DBConn.collection.FindOne(db.DBConn.ctx, filter).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			newErr := fmt.Errorf("QueryTaskById: No matching record found")
			return nil, newErr
		}
		newErr := fmt.Errorf("db findone operation failed, [%s]", err)
		zaplog.Error(newErr)
		return nil, newErr
	}

	return &task, nil
}

// Update Task State
func (db *DBManager) UpdateTaskState(userId string, enterpriseId string, taskId string, state uint32) error {

	filter := bson.M{
		"user_id": userId,
		"enterprise_id": enterpriseId,
		"task_id": taskId,
	}

	update := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "state", Value: state},
	}}}

	t := &AlertTaskInfo{}
	err := db.DBConn.collection.FindOneAndUpdate(db.DBConn.ctx, filter,
		update).Decode(t)
	if err != nil {
		newErr := fmt.Errorf("db findoneandupdate operation failed, [%s]", err)
		zaplog.Error(newErr)
		return newErr
	}

	return nil
}

// Update new task details
func (db *DBManager) UpdateTask(userId string, enterpriseId string, task *AlertTaskInfo) error {

	filter := bson.M{
		"user_id": userId,
		"enterprise_id": enterpriseId,
		"task_id": task.TaskId,
	}

	_, err := db.DBConn.collection.ReplaceOne(db.DBConn.ctx, filter, task)
	if err != nil {
		newErr := fmt.Errorf("db replaceone operation failed, [%s]", err)
		zaplog.Error(newErr)
		return newErr
	}

	return nil
}

// Delete DB Manager Task
func (db *DBManager) DeleteTask(userId string, enterpriseId string, taskId string) error {
	filter := bson.M{
		"user_id": userId,
		"enterprise_id": enterpriseId,
		"task_id": taskId,
	}

	result, err := db.DBConn.collection.DeleteOne(db.DBConn.ctx, filter)
	if err != nil {
		newErr := fmt.Errorf("db deleteone operation failed, [%s]", err)
		zaplog.Error(newErr)
		return newErr
	}

	if result.DeletedCount == 0 {
		zaplog.Info("DeleteTask: No matching record found")
	}

	return nil
}

// Get total number of tasks matching to user id and enterprise id
func (db *DBManager) GetTotalTaskCount(userId string, enterpriseId string) (int64, error) {
	total, err := db.DBConn.collection.CountDocuments(db.DBConn.ctx, bson.M{
		"user_id": userId,
		"enterprise_id": enterpriseId,
	})
	if err != nil {
		newErr := fmt.Errorf("db countocuments operation failed, [%s]", err)
		zaplog.Error(newErr)
		return -1, newErr
	}
	return total, nil
}
