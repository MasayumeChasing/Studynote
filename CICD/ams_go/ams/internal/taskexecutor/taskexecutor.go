// Copyright (c) Huawei Technologies Co., Ltd. 2020-2020. All rights reserved.
// File Name:  taskexecutor.go
// Description: Execution of Alert Task to perform 1:N search

// Package taskexecutor implements the execution of Alert Task to perform 1:N search
package taskexecutor

import (
	"bytes"
	"codehub-dg-g.huawei.com/IntelligentVision/IVM/AMS/internal/configmanager"
	"codehub-dg-g.huawei.com/IntelligentVision/IVM/AMS/internal/dbmanager"
	"codehub-dg-g.huawei.com/IntelligentVision/IVM/AMS/internal/taskmanager"
	"codehub-dg-g.huawei.com/IntelligentVision/IVM/Resources/apis/golang/message"
	"codehub-dg-g.huawei.com/IntelligentVision/IVM/Resources/search/engine"
	"codehub-dg-g.huawei.com/IntelligentVision/IVM/Resources/security/sccapi"
	"codehub-dg-g.huawei.com/IntelligentVision/IVM/Resources/utils/zaplog"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"
)

const MatchedAlarm = "matched_alarm"

type TaskExecutor struct {
	eng    engine.Engine
	inited bool
}

// New Task Executor
func NewTaskExecutor(info *configmanager.ESEngineInfo) (*TaskExecutor, error) {
	zaplog.Info("Initializing Elastic Search Engine")
	te := TaskExecutor{
		inited: true,
	}
	secret, err := sccapi.Decrypt(info.Password)
	if err != nil {
		newErr := fmt.Errorf("Sccapi decrypt failed in Engine init, %s", err)
		zaplog.Error(newErr)
		return nil, newErr
	}
	engineEndpoint := engine.SearchEndpoint{
		Endpoints: info.Endpoints,
		Username:  info.Username,
		Password:  secret,
	}
	te.eng, err = engine.NewEngine(engineEndpoint)
	if err != nil {
		newErr := fmt.Errorf("Engine init failed, %s", err)
		zaplog.Error(newErr)
		return nil, newErr
	}
	return &te, nil
}

// To perform 1:N img search
func (te *TaskExecutor) DoImgSearch(tm *taskmanager.TaskManager, msg string) (
	string, error) {
	searchQuery, err := parseJsonImgSearchMsg(msg)
	if err != nil {
		return "", err
	}
	deviceId := searchQuery.Data.Param.DeviceID
	channelId := searchQuery.Data.Param.ChannelID
	if len(deviceId) == 0 {
		newErr := fmt.Errorf("No device id on Img Search msg")
		zaplog.Error(newErr)
		return "", newErr
	}
	zaplog.Infof("Img Search operation for device ID [%s]", deviceId)
	imgFeature, err := convertFeatureByteToFloat(searchQuery.Data.Param.ImgFeatureVector)
	if err != nil {
		newErr := fmt.Errorf("Feature byte to float conversion failed, %s", err)
		zaplog.Error(newErr)
		return "", newErr
	}
	matchedDatas, err := getMatchedData(te, tm, deviceId, channelId, imgFeature)
	if err != nil {
		newErr := fmt.Errorf("getMatchedData failed, %s", err)
		zaplog.Error(newErr)
		return "", newErr
	}
	if matchedDatas != nil && len(matchedDatas) > 0 {
		return encodeImgSearchResult(searchQuery, matchedDatas)
	}
	return "", nil
}

// Do 1 to N comparison and get matched data
func getMatchedData(te *TaskExecutor, tm *taskmanager.TaskManager, deviceId string, channelId string,
	imgFeature []float32) ([]message.MatchedDataField, error) {
	var offset uint32 = 0
	matchedDatas := make([]message.MatchedDataField, 0, 1)

	for {
		tasks, err := tm.QueryTaskListByState(taskmanager.TaskActive, 
		&dbmanager.QueryPara{Offset:offset, Limit:dbmanager.DefaultLimit, SortDir:"ascending", SortKey:"create_time"})
		if err != nil {
			newErr := fmt.Errorf("Query task list failed, %s", err)
			zaplog.Error(newErr)
			return matchedDatas, newErr
		}
		// if tasks is nil or 0 task length, indicates that there is no active task
		if tasks == nil || len(tasks) == 0 {
			zaplog.Info("All Tasks are checked for img search")
			break
		}
		zaplog.Infof("Received [%d] tasks for img search", len(tasks))
		offset++ // next page
		for _, task := range tasks {
			if isValidTaskForSearch(deviceId, channelId, task) != 1 {
				continue
			}
			zaplog.Infof("Executing Task [taskId=%s] with device id [%s]",
				task.TaskId, deviceId)
			task.State = taskmanager.TaskRunning
			searchResult, err := execute1NImgCompare(te, task, imgFeature)
			task.State = taskmanager.TaskActive
			if err != nil {
				newErr := fmt.Errorf("Img Search failed deviceId [%s], %s", deviceId, err)
				zaplog.Info(newErr)
				continue
			}
			matchedData := message.MatchedDataField{
				TargetDataId: searchResult.Result[0].Entity.DataID,
				TaskId:       task.TaskId, TaskName: task.TaskName,
				Score:      searchResult.Result[0].Score,
				Attributes: searchResult.Result[0].Entity.Attributes,
			}
			matchedDatas = append(matchedDatas, matchedData)
		}
	}
	return matchedDatas, nil
}

func isValidTaskForSearch(deviceId string, channelId string,
	task *dbmanager.AlertTaskInfo) int {
	now := time.Now()
	// -1 in StartTime or EndTime means no time restriction
	if now.Before(task.StartTime)|| now.After(task.EndTime) {
		// Expired task
		return 0
	}
	return isDeviceIdPresent(deviceId, channelId, task.Devices)
}

// To check whether device id and channel id are matching with one present in task

func isDeviceIdPresent(deviceId string, channelId string,
	devices []dbmanager.DeviceInfo) int {
	for _, device := range devices {
		if device.DeviceId == deviceId && device.ChannelId == channelId {
			return 1
		}
	}
	// TODO we need to check current time is in start and stop time of AMS
	return 0
}

// To parse json format of search msg received from IDS via RabbitMQ
func parseJsonImgSearchMsg(msg string) (*message.IdsAlarm, error) {
	var searchQuery message.IdsAlarm
	err := json.Unmarshal([]byte(msg), &searchQuery)
	if err != nil {
		newErr := fmt.Errorf("Parsing Json Img Search Query failed, %s", err)
		zaplog.Error(newErr)
		return nil, newErr
	}
	return &searchQuery, nil
}

// To convert from byte to int32
func byteToInt32(b []byte) (int, error) {
	bytesBuffer := bytes.NewBuffer(b)
	var x int32
	err := binary.Read(bytesBuffer, binary.LittleEndian, &x)
	if err != nil {
		newErr := fmt.Errorf("Binary read failed, %s", err)
		zaplog.Error(newErr)
		return (-1), newErr // -1 indicates error
	}
	return int(x), nil
}

// To convert feature vector from byte to float
func convertFeatureByteToFloat(feature []byte) ([]float32, error) {
	var byteToFloatFactor int = 4
	var ignoredDataLen int = 8
	caps := make([]float32, 0, 5)
	for i := ignoredDataLen; i < len(feature); i += byteToFloatFactor {
		j := i + byteToFloatFactor
		if j > len(feature) {
			j = len(feature)
		}
		val, err := byteToInt32(feature[i:j])
		if err != nil {
			newErr := fmt.Errorf("byte to int32 conversion failed, %s", err)
			zaplog.Error(newErr)
			return caps, newErr
		}
		var f float32
		f = float32(val) / float32(4096.0) // capConversionFactor
		caps = append(caps, f)
	}
	return caps, nil
}

// To parse json format of search msg received from IDS via RabbitMQ

func execute1NImgCompare(te *TaskExecutor, task *dbmanager.AlertTaskInfo,
	imgFeature []float32) (
	*engine.StaticSearchResult, error) {
	staticDocQuery := &engine.StaticDocQuery{
		DataID: &engine.QueryData{
			DataIDs: nil,
		},
		Threshold: task.Threshold,
		Feature: &engine.QueryFeature{
			FeatureVector: imgFeature,
			TopK:          1,
		},
	}
	searchResult, err := te.eng.SearchStatic(task.IndexName, staticDocQuery)
	if err != nil {
		newErr := fmt.Errorf("Img Search call failed,", err)
		zaplog.Error(newErr)
		return nil, newErr
	}
	if len(searchResult.Result) == 0 {
		newErr := fmt.Errorf("Img Search result is empty")
		zaplog.Error(newErr)
		return nil, newErr
	}
	return searchResult, nil
}

// Encode the Img search result in json format

func encodeImgSearchResult(searchQuery *message.IdsAlarm,
	matchedDatas []message.MatchedDataField) (string, error) {
	searchResult := message.UserAlertMsg{
		MsgID:       searchQuery.MsgID,
		MsgType:     searchQuery.MsgType,
		MsgResource: searchQuery.MsgResource,
		PushTime:    searchQuery.PushTime,
		UserAlertData: message.UserAlertDataField{
			DataID:      searchQuery.Data.Param.DataID,
			UserIDs:     searchQuery.Data.UserIDs,
			DeviceID:    searchQuery.Data.Param.DeviceID,
			ChannelID:   searchQuery.Data.Param.ChannelID,
			CaptureTime: searchQuery.Data.Param.CaptureTime,
			AlarmType:   MatchedAlarm,
			MatchedData: matchedDatas,
		},
	}
	jsonBytes, err := json.Marshal(searchResult)
	if err != nil {
		newErr := fmt.Errorf("Img Search result json encoding failed, %s", err)
		zaplog.Error(newErr)
		return "", newErr
	}
	return string(jsonBytes), nil
}
