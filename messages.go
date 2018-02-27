package main

import (
	"math/rand"
	"time"

	"github.com/docker/docker/api/types"
	bw2 "github.com/immesys/bw2bind"
)

const GetReplicasRequestPIDString = "2.2.0.0"
const GetReplicasResponsePIDString = "2.2.0.1"
const DeployModelRequestPIDString = "2.2.0.2"
const DeployModelResponsePIDString = "2.2.0.3"
const RegisterApplicationRequestPIDString = "2.2.0.4"
const RegisterApplicationResponsePIDString = "2.2.0.5"
const LinkModelToAppRequestPIDString = "2.2.0.6"
const LinkModelToAppResponsePIDString = "2.2.0.7"
const BuildAndDeployModelRequestPIDString = "2.2.0.8"
const BuildAndDeployModelResponsePIDString = "2.2.0.9"
const RegisterModelRequestPIDString = "2.2.0.10"
const RegisterModelResponsePIDString = "2.2.0.11"
const ModelInfoRequestPIDString = "2.2.0.12"
const ModelInfoResponsePIDString = "2.2.0.13"
const ApplicationListRequestPIDString = "2.2.0.14"
const ApplicationListResponsePIDString = "2.2.0.15"
const GetApplicationInfoRequestPIDString = "2.2.0.16"
const GetApplicationInfoResponsePIDString = "2.2.0.17"

var GetReplicasRequestPID int
var GetReplicasResponsePID int
var DeployModelRequestPID int
var DeployModelResponsePID int
var RegisterApplicationRequestPID int
var RegisterApplicationResponsePID int
var LinkModelToAppRequestPID int
var LinkModelToAppResponsePID int
var BuildAndDeployModelRequestPID int
var BuildAndDeployModelResponsePID int
var RegisterModelRequestPID int
var RegisterModelResponsePID int
var ModelInfoRequestPID int
var ModelInfoResponsePID int
var ApplicationListRequestPID int
var ApplicationListResponsePID int
var GetApplicationInfoRequestPID int
var GetApplicationInfoResponsePID int

func init() {
	rand.Seed(time.Now().UnixNano())
	GetReplicasRequestPID, _ = bw2.PONumFromDotForm(GetReplicasRequestPIDString)
	GetReplicasResponsePID, _ = bw2.PONumFromDotForm(GetReplicasResponsePIDString)
	DeployModelRequestPID, _ = bw2.PONumFromDotForm(DeployModelRequestPIDString)
	DeployModelResponsePID, _ = bw2.PONumFromDotForm(DeployModelResponsePIDString)
	RegisterApplicationRequestPID, _ = bw2.PONumFromDotForm(RegisterApplicationRequestPIDString)
	RegisterApplicationResponsePID, _ = bw2.PONumFromDotForm(RegisterApplicationResponsePIDString)
	LinkModelToAppRequestPID, _ = bw2.PONumFromDotForm(LinkModelToAppRequestPIDString)
	LinkModelToAppResponsePID, _ = bw2.PONumFromDotForm(LinkModelToAppResponsePIDString)
	BuildAndDeployModelRequestPID, _ = bw2.PONumFromDotForm(BuildAndDeployModelRequestPIDString)
	BuildAndDeployModelResponsePID, _ = bw2.PONumFromDotForm(BuildAndDeployModelResponsePIDString)
	RegisterModelRequestPID, _ = bw2.PONumFromDotForm(RegisterModelRequestPIDString)
	RegisterModelResponsePID, _ = bw2.PONumFromDotForm(RegisterModelResponsePIDString)
	ModelInfoRequestPID, _ = bw2.PONumFromDotForm(ModelInfoRequestPIDString)
	ModelInfoResponsePID, _ = bw2.PONumFromDotForm(ModelInfoResponsePIDString)
	ApplicationListRequestPID, _ = bw2.PONumFromDotForm(ApplicationListRequestPIDString)
	ApplicationListResponsePID, _ = bw2.PONumFromDotForm(ApplicationListResponsePIDString)
	GetApplicationInfoRequestPID, _ = bw2.PONumFromDotForm(GetApplicationInfoRequestPIDString)
	GetApplicationInfoResponsePID, _ = bw2.PONumFromDotForm(GetApplicationInfoResponsePIDString)
}

type GetReplicasMessageRequest struct {
	MsgID int64
	// container label to search
	Label string
}

func NewGetReplicasRequest(label string) GetReplicasMessageRequest {
	return GetReplicasMessageRequest{
		MsgID: rand.Int63(),
		Label: label,
	}
}

func (msg *GetReplicasMessageRequest) PayloadObject() bw2.PayloadObject {
	po, _ := bw2.CreateMsgPackPayloadObject(GetReplicasRequestPID, msg)
	return po
}

func (msg *GetReplicasMessageRequest) Response() *GetReplicasMessageResponse {
	return &GetReplicasMessageResponse{
		MsgID: msg.MsgID,
	}
}

type GetReplicasMessageResponse struct {
	MsgID      int64
	Containers []types.Container
	Error      string
}

func (msg *GetReplicasMessageResponse) PayloadObject() bw2.PayloadObject {
	po, _ := bw2.CreateMsgPackPayloadObject(GetReplicasResponsePID, msg)
	return po
}

type DeployModelRequest struct {
	MsgID      int64
	Name       string
	Version    string
	Input_type string
	Image      string
}

func (msg *DeployModelRequest) Response() *DeployModelResponse {
	return &DeployModelResponse{
		MsgID: msg.MsgID,
	}
}

type DeployModelResponse struct {
	MsgID int64
	Error string
}

func (msg *DeployModelResponse) PayloadObject() bw2.PayloadObject {
	po, _ := bw2.CreateMsgPackPayloadObject(DeployModelResponsePID, msg)
	return po
}

type RegisterApplicationRequest struct {
	MsgID int64 `json:"-"`
	ApplicationInfo
}

func (msg *RegisterApplicationRequest) Response() *RegisterApplicationResponse {
	return &RegisterApplicationResponse{
		MsgID: msg.MsgID,
	}
}

type RegisterApplicationResponse struct {
	MsgID int64
	Error string
}

func (msg *RegisterApplicationResponse) PayloadObject() bw2.PayloadObject {
	po, _ := bw2.CreateMsgPackPayloadObject(RegisterApplicationResponsePID, msg)
	return po
}

type LinkModelToAppRequest struct {
	MsgID       int64    `json:"-"`
	App_name    string   `json:"app_name"`
	Model_names []string `json:"model_names"`
}

func (msg *LinkModelToAppRequest) Response() *LinkModelToAppResponse {
	return &LinkModelToAppResponse{
		MsgID: msg.MsgID,
	}
}

type LinkModelToAppResponse struct {
	MsgID int64
	Error string
}

func (msg *LinkModelToAppResponse) PayloadObject() bw2.PayloadObject {
	po, _ := bw2.CreateMsgPackPayloadObject(LinkModelToAppResponsePID, msg)
	return po
}

type BuildAndDeployModelRequest struct {
	MsgID              int64 `json:"-"`
	Name               string
	Version            string
	Input_type         string
	Model_data_path    string
	Base_image         string
	Labels             map[string]string
	Container_registry string
	Num_replicas       int
	Batch_size         int
}

func (msg *BuildAndDeployModelRequest) Response() *BuildAndDeployModelResponse {
	return &BuildAndDeployModelResponse{
		MsgID: msg.MsgID,
	}
}

type BuildAndDeployModelResponse struct {
	MsgID int64
	Error string
}

func (msg *BuildAndDeployModelResponse) PayloadObject() bw2.PayloadObject {
	po, _ := bw2.CreateMsgPackPayloadObject(BuildAndDeployModelResponsePID, msg)
	return po
}

type RegisterModelRequest struct {
	MsgID int64 `json:"-"`
	// image
	Batch_size      int      `json:"batch_size"`
	Model_data_path string   `json:"model_data_path"`
	Input_type      string   `json:"input_type"`
	Labels          []string `json:"labels"`
	Container_name  string   `json:"container_name"`
	Model_version   string   `json:"model_version"`
	Model_name      string   `json:"model_name"`
}

func (msg *RegisterModelRequest) Response() *RegisterModelResponse {
	return &RegisterModelResponse{
		MsgID: msg.MsgID,
	}
}

type RegisterModelResponse struct {
	MsgID int64
	Error string
}

func (msg *RegisterModelResponse) PayloadObject() bw2.PayloadObject {
	po, _ := bw2.CreateMsgPackPayloadObject(RegisterModelResponsePID, msg)
	return po
}

type ModelInfoRequest struct {
	MsgID   int64 `json:"-"`
	Verbose bool  `json:"verbose"`
}

func (msg *ModelInfoRequest) Response() *ModelInfoResponse {
	return &ModelInfoResponse{
		MsgID: msg.MsgID,
	}
}

type ModelInfoResponse struct {
	MsgID             int64
	Error             string
	ModelNames        []string
	ModelDescriptions []ModelInfo
}

func (msg *ModelInfoResponse) PayloadObject() bw2.PayloadObject {
	po, _ := bw2.CreateMsgPackPayloadObject(ModelInfoResponsePID, msg)
	return po
}

type ModelInfo struct {
	Is_current_version bool     `msgpack:"is_current_version",json:"is_current_version"`
	Model_data_path    string   `msgpack:"model_data_path",json:"model_data_path"`
	Input_type         string   `msgpack:"input_type",json:"input_type"`
	Labels             []string `msgpack:"labels",json:"labels"`
	Container_name     string   `msgpack:"container_name",json:"container_name"`
	Model_version      string   `msgpack:"model_version",json:"model_version"`
	Model_name         string   `msgpack:"model_name",json:"model_name"`
}

type ApplicationListRequest struct {
	MsgID   int64 `json:"-"`
	Verbose bool  `json:"verbose"`
}

func (msg *ApplicationListRequest) Response() *ApplicationListResponse {
	return &ApplicationListResponse{
		MsgID: msg.MsgID,
	}
}

type ApplicationListResponse struct {
	MsgID                   int64
	Error                   string
	ApplicationNames        []string
	ApplicationDescriptions []ApplicationInfo
}

func (msg *ApplicationListResponse) PayloadObject() bw2.PayloadObject {
	po, _ := bw2.CreateMsgPackPayloadObject(ApplicationListResponsePID, msg)
	return po
}

type GetApplicationInfoRequest struct {
	MsgID int64  `json:"-"`
	Name  string `json:"name"`
}

func (msg *GetApplicationInfoRequest) Response() *GetApplicationInfoResponse {
	return &GetApplicationInfoResponse{
		MsgID: msg.MsgID,
	}
}

type GetApplicationInfoResponse struct {
	MsgID int64
	Error string
	Info  ApplicationInfo
}

func (msg *GetApplicationInfoResponse) PayloadObject() bw2.PayloadObject {
	po, _ := bw2.CreateMsgPackPayloadObject(GetApplicationInfoResponsePID, msg)
	return po
}

type ApplicationInfo struct {
	Name               string `msgpack:"name",json:"name"`
	Input_type         string `msgpack:"input_type",json:"input_type"`
	Default_output     string `msgpack:"default_output",json:"default_output"`
	Latency_slo_micros int64  `msgpack:"latency_slo_micros",json:"latency_slo_micros"`
}
