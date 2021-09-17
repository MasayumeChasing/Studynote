// Copyright (c) Huawei Technologies Co., Ltd. 2020-2020. All rights reserved.
// File Name: amsgrpc.go
// Description: Contains grpc intefaces

// Package amsgrpc implements ams grpc interfaces
package amsgrpc

import (
	"codehub-dg-g.huawei.com/IntelligentVision/IVM/Resources/apis/golang/ams"
	"codehub-dg-g.huawei.com/IntelligentVision/IVM/Resources/utils/zaplog"
	"context"
	"fmt"
	"strings"
	// "codehub-dg-g.huawei.com/IntelligentVision/IVM/Resources/apis/golang/slms"
	"codehub-dg-g.huawei.com/IntelligentVision/IVM/AMS/internal/configmanager"
	"codehub-dg-g.huawei.com/IntelligentVision/IVM/AMS/internal/dbmanager"
	"codehub-dg-g.huawei.com/IntelligentVision/IVM/AMS/internal/taskmanager"
	"codehub-dg-g.huawei.com/IntelligentVision/IVM/Resources/security/sccapi"
	"codehub-dg-g.huawei.com/IntelligentVision/IVM/Resources/utils/jwt"
	"codehub-dg-g.huawei.com/IntelligentVision/IVM/Resources/x/rpc"
	"time"
)

type GRPCInterfaceConf struct {
	GRPC   configmanager.ServerInfo
	JwtAMS configmanager.JwtInfo
}

type AMSgRPCInterface struct {
	TM *taskmanager.TaskManager // task management functions
	ams.UnimplementedAlertManagerServiceServer
	Conf   GRPCInterfaceConf
	inited bool
}

// Create New AMS gRPC Interface instance
func NewAMSgRPCInterface(conf GRPCInterfaceConf, tm *taskmanager.TaskManager) (
	*AMSgRPCInterface, error) {
	var am *AMSgRPCInterface = &AMSgRPCInterface{
		inited: false,
	}

	if am == nil {
		err := fmt.Errorf("ams create instance failed")
		zaplog.Error(err)
		return nil, err
	}

	err := am.Init(conf, tm)
	if err != nil {
		newErr := fmt.Errorf("ams create instance failed, [%s]", err)
		zaplog.Error(newErr)
		return nil, newErr
	}

	return am, nil
}

// AMSgRPCInterface Initialization
func (am *AMSgRPCInterface) Init(conf GRPCInterfaceConf,
	tm *taskmanager.TaskManager) error {
	if am.inited == true {
		return nil
	}

	am.TM = tm
	am.Conf = conf
	am.inited = true

	return nil
}

// AMSgRPCInterface De-Initialization
func (am *AMSgRPCInterface) DeInit() {
	if am.inited == false {
		return
	}

	am.TM.DeInit()
	am.inited = false
	return
}

// Validate the JwtToken
func (am *AMSgRPCInterface) ValidateJwt(ctx context.Context) (map[string]string, error) {

	// Get header from context
	jwtToken := rpc.GetHeaderFromContext(ctx, "authorization")
	if jwtToken == "" {
		err := fmt.Errorf("GetHeaderFromContext failed for JWT token")
		zaplog.Error(err)
		return nil, err
	}

	var newJwtT string
	if strings.Index(jwtToken, "bearer ") == 0 {
		newJwtT = strings.Replace(jwtToken, "bearer ", "", 1)
	} else if strings.Index(jwtToken, "Bearer ") == 0 {
		newJwtT = strings.Replace(jwtToken, "Bearer ", "", 1)
	} else {
		return nil, fmt.Errorf("token illegal")
	}

	plaintext, err := sccapi.Decrypt(am.Conf.JwtAMS.Secret)
	if err != nil {
		err := fmt.Errorf("Decrypt cipher failed [%s]", err)
		zaplog.Error(err)
		return nil, err
	}

	err = jwt.ValidateToken(newJwtT, []byte(plaintext), []byte(plaintext), nil)
	if err != nil {
		newErr := fmt.Errorf("jwt token validation error, [%s]", err)
		zaplog.Error(newErr)
		return nil, newErr
	}

	body, err := jwt.GetBodyFromToken(newJwtT, []byte(plaintext), []byte(plaintext))
	if err != nil {
		newErr := fmt.Errorf("jwt token validation error, [%s]", err)
		return nil, newErr
	}
	
	enterpriseId := body["enterprise_id"]
	if enterpriseId == "" {
		newErr := fmt.Errorf("jwt token:enterprise id error")
		return nil, newErr
	}

	return body, nil
}

// Creating task request
func (am *AMSgRPCInterface) CreateTask(ctx context.Context, req *ams.CreateTaskRequest) (
	*ams.CreateTaskResponse, error) {

	var resp ams.CreateTaskResponse
	resp.TaskId = ""
	var taskId string
	
	// Validate request
	err := req.Validate()
	if err != nil {
		newErr := fmt.Errorf("create task: request error, [%s]", err)
		zaplog.Error(newErr)
		return &resp, newErr
	}

	// Validate req.jwt
	body, ret := am.ValidateJwt(ctx)
	if ret != nil {
		newErr := fmt.Errorf("create task: jwt validation error, [%s]", ret)
		zaplog.Error(newErr)
		return &resp, newErr
	}

	err = am.TM.CreateTask(req, body, &taskId)
	if err != nil || taskId == "" {
		newErr := fmt.Errorf("create task error, [%s]", err)
		zaplog.Error(newErr)
		return &resp, newErr
	}
	resp.TaskId = taskId
	return &resp, nil
}

// Get task info as per per response message format
func (am *AMSgRPCInterface) GetRespTaskStruct(v *dbmanager.AlertTaskInfo) *ams.Task {
	t := new(ams.Task)
	t.TaskId = v.TaskId
	t.TaskName = v.TaskName
	t.State = v.State
	t.Threshold = v.Threshold
	t.LibId = v.LibId

	l := len(v.Devices)
	if l == 0 {
		t.Devices = t.Devices[:0]
	} else {
		for j := 0; j < l; j++ {
			d := new(ams.DeviceInfo)
			d.DeviceId = v.Devices[j].DeviceId
			d.ChannelId = v.Devices[j].ChannelId
			t.Devices = append(t.Devices, d)
		}
	}
	if 0 < len(v.IndexName) {
		t.IndexName = v.IndexName[0]
	}
	t.StartTime = v.StartTime.Format(time.RFC3339)
	t.EndTime = v.EndTime.Format(time.RFC3339)
	t.CreateTime = v.CreateTime.Format(time.RFC3339)
	t.UpdateTime = v.UpdateTime.Format(time.RFC3339)
	return t
}

// Get the List of Tasks
func (am *AMSgRPCInterface) ListTasks(ctx context.Context, req *ams.ListTaskRequest) (
	*ams.ListTaskResponse, error) {
	var resp ams.ListTaskResponse
	resp.Total = 0

	// Validate request
	err := req.Validate()
	if err != nil {
		newErr := fmt.Errorf("list task: request error, [%s]", err)
		zaplog.Error(newErr)
		return &resp, newErr
	}

	// Validate req.jwt
	body, ret := am.ValidateJwt(ctx)
	if ret != nil {
		newErr := fmt.Errorf("list task: jwt validation error, [%s]", ret)
		zaplog.Error(newErr)
		return &resp, newErr
	}

	queryPara := &dbmanager.QueryPara{Offset:req.Offset, Limit:req.Limit, SortDir:req.SortDir, SortKey:req.SortKey}
	total, tasks, err := am.TM.QueryTaskList(req.UserId, body["enterprise_id"], queryPara)
	if err != nil {
		newErr := fmt.Errorf("list task error, [%s]", err)
		zaplog.Error(newErr)
		return &resp, newErr
	}

	resp.Total = int32(total)

	if tasks == nil {
		return &resp, nil
	}

	taskList := []*ams.Task{}
	var empty bool = true

	for _, v := range tasks {
		empty = false
		t := am.GetRespTaskStruct(v)
		taskList = append(taskList, t)
	}

	if empty == false {
		resp.Tasks = taskList
	}

	return &resp, nil
}

// Get Task Details
func (am *AMSgRPCInterface) ShowTask(ctx context.Context, req *ams.ShowTaskRequest) (
	*ams.ShowTaskResponse, error) {
	var resp ams.ShowTaskResponse
	resp.Task = nil

	// Validate request
	err := req.Validate()
	if err != nil {
		newErr := fmt.Errorf("show task: request error, [%s]", err)
		zaplog.Error(newErr)
		return &resp, newErr
	}

	// Validate req.jwt
	body, ret := am.ValidateJwt(ctx)
	if ret != nil {
		newErr := fmt.Errorf("show task: jwt validation error, [%s]", ret)
		zaplog.Error(newErr)
		return &resp, newErr
	}

	task, err := am.TM.QueryTaskById(req.UserId, body["enterprise_id"], req.TaskId)
	if err != nil {
		newErr := fmt.Errorf("show task error, [%s]", err)
		zaplog.Error(newErr)
		return &resp, newErr
	}

	if task == nil {
		return &resp, nil
	}

	t := new(ams.Task)

	t.TaskId = task.TaskId
	t.TaskName = task.TaskName
	t.State = task.State
	t.Threshold = task.Threshold
	t.LibId = task.LibId

	l := len(task.Devices)
	if l == 0 {
		t.Devices = t.Devices[:0]
	} else {
		for i := 0; i < l; i++ {
			d := new(ams.DeviceInfo)
			d.DeviceId = task.Devices[i].DeviceId
			d.ChannelId = task.Devices[i].ChannelId
			t.Devices = append(t.Devices, d)
		}
	}
	if 0 < len(task.IndexName) {
		t.IndexName = task.IndexName[0]
	}
	t.StartTime = task.StartTime.Format(time.RFC3339)
	t.EndTime = task.EndTime.Format(time.RFC3339)
	t.CreateTime = task.CreateTime.Format(time.RFC3339)
	t.UpdateTime = task.UpdateTime.Format(time.RFC3339)
	resp.Task = t
	return &resp, nil
}

// Update Task with new configurations
func (am *AMSgRPCInterface) UpdateTask(ctx context.Context, req *ams.UpdateTaskRequest) (
	*ams.UpdateTaskResponse, error) {

	var resp ams.UpdateTaskResponse
	
	// Validate request
	err := req.Validate()
	if err != nil {
		newErr := fmt.Errorf("update task: request error, [%s]", err)
		zaplog.Error(newErr)
		return &resp, newErr
	}

	// Validate req.jwt
	body, ret := am.ValidateJwt(ctx)
	if ret != nil {
		newErr := fmt.Errorf("update task: jwt validation error, [%s]", ret)
		zaplog.Error(newErr)
		return &resp, newErr
	}

	err = am.TM.UpdateTask(req, body)
	if err != nil {
		newErr := fmt.Errorf("update task error, [%s]", err)
		zaplog.Error(newErr)
		return &resp, newErr
	}

	return &resp, nil
}

// Delete Task by Task ID.
func (am *AMSgRPCInterface) DeleteTask(ctx context.Context, req *ams.DeleteTaskRequest) (
	*ams.DeleteTaskResponse, error) {

	var resp ams.DeleteTaskResponse

	// Validate request
	err := req.Validate()
	if err != nil {
		newErr := fmt.Errorf("delete task: request error, [%s]", err)
		zaplog.Error(newErr)
		return &resp, newErr
	}

	// Validate req.jwt
	body, ret := am.ValidateJwt(ctx)
	if ret != nil {
		newErr := fmt.Errorf("delete task: jwt validation error, [%s]", ret)
		zaplog.Error(newErr)
		return &resp, newErr
	}

	err = am.TM.DeleteTaskById(req.UserId, body["enterprise_id"], req.TaskId)
	if err != nil {
		newErr := fmt.Errorf("delete task error, [%s]", err)
		zaplog.Error(newErr)
		return &resp, newErr
	}

	return &resp, nil
}

// Start the Task
func (am *AMSgRPCInterface) StartTask(ctx context.Context, req *ams.StartTaskRequest) (
	*ams.StartTaskResponse, error) {
	var resp ams.StartTaskResponse

	// Validate request
	err := req.Validate()
	if err != nil {
		newErr := fmt.Errorf("start task: request error, [%s]", err)
		zaplog.Error(newErr)
		return &resp, newErr
	}

	// Validate req.jwt
	body, ret := am.ValidateJwt(ctx)
	if ret != nil {
		newErr := fmt.Errorf("start task: jwt validation error, [%s]", ret)
		zaplog.Error(newErr)
		return &resp, newErr
	}

	err = am.TM.ActivateTaskById(req.UserId, body["enterprise_id"], req.TaskId)
	if err != nil {
		newErr := fmt.Errorf("start task error, [%s]", err)
		zaplog.Error(newErr)
		return &resp, newErr
	}

	return &resp, nil
}

// Stoping the Task.
func (am *AMSgRPCInterface) StopTask(ctx context.Context, req *ams.StopTaskRequest) (
	*ams.StopTaskResponse, error) {

	var resp ams.StopTaskResponse

	// Validate request
	err := req.Validate()
	if err != nil {
		newErr := fmt.Errorf("stop task: request error, [%s]", err)
		zaplog.Error(newErr)
		return &resp, newErr
	}

	// Validate req.jwt
	body, ret := am.ValidateJwt(ctx)
	if ret != nil {
		newErr := fmt.Errorf("stop task: jwt validation error, [%s]", ret)
		zaplog.Error(newErr)
		return &resp, newErr
	}

	err = am.TM.DeActivateTaskById(req.UserId, body["enterprise_id"], req.TaskId)
	if err != nil {
		newErr := fmt.Errorf("stop task error, [%s]", err)
		zaplog.Error(newErr)
		return &resp, newErr
	}

	return &resp, nil
}
