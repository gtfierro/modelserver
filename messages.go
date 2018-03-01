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
const ListModelRequestPIDString = "2.2.0.12"
const ListModelResponsePIDString = "2.2.0.13"
const ApplicationListRequestPIDString = "2.2.0.14"
const ApplicationListResponsePIDString = "2.2.0.15"
const GetApplicationInfoRequestPIDString = "2.2.0.16"
const GetApplicationInfoResponsePIDString = "2.2.0.17"
const GetModelInfoRequestPIDString = "2.2.0.18"
const GetModelInfoResponsePIDString = "2.2.0.19"
const GetLinkedModelsRequestPIDString = "2.2.0.20"
const GetLinkedModelsResponsePIDString = "2.2.0.21"
const GetAllModelReplicasRequestPIDString = "2.2.0.22"
const GetAllModelReplicasResponsePIDString = "2.2.0.23"
const GetModelReplicaInfoRequestPIDString = "2.2.0.24"
const GetModelReplicaInfoResponsePIDString = "2.2.0.25"
const GetContainerLogsRequestPIDString = "2.2.0.26"
const GetContainerLogsResponsePIDString = "2.2.0.27"
const InspectInstanceRequestPIDString = "2.2.0.28"
const InspectInstanceResponsePIDString = "2.2.0.29"
const SetModelVersionRequestPIDString = "2.2.0.30"
const SetModelVersionResponsePIDString = "2.2.0.31"

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
var ListModelRequestPID int
var ListModelResponsePID int
var ApplicationListRequestPID int
var ApplicationListResponsePID int
var GetApplicationInfoRequestPID int
var GetApplicationInfoResponsePID int
var GetModelInfoRequestPID int
var GetModelInfoResponsePID int
var GetLinkedModelsRequestPID int
var GetLinkedModelsResponsePID int
var GetAllModelReplicasRequestPID int
var GetAllModelReplicasResponsePID int
var GetModelReplicaInfoRequestPID int
var GetModelReplicaInfoResponsePID int
var GetContainerLogsRequestPID int
var GetContainerLogsResponsePID int
var InspectInstanceRequestPID int
var InspectInstanceResponsePID int
var SetModelVersionRequestPID int
var SetModelVersionResponsePID int

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
	ListModelRequestPID, _ = bw2.PONumFromDotForm(ListModelRequestPIDString)
	ListModelResponsePID, _ = bw2.PONumFromDotForm(ListModelResponsePIDString)
	ApplicationListRequestPID, _ = bw2.PONumFromDotForm(ApplicationListRequestPIDString)
	ApplicationListResponsePID, _ = bw2.PONumFromDotForm(ApplicationListResponsePIDString)
	GetApplicationInfoRequestPID, _ = bw2.PONumFromDotForm(GetApplicationInfoRequestPIDString)
	GetApplicationInfoResponsePID, _ = bw2.PONumFromDotForm(GetApplicationInfoResponsePIDString)
	GetModelInfoRequestPID, _ = bw2.PONumFromDotForm(GetModelInfoRequestPIDString)
	GetModelInfoResponsePID, _ = bw2.PONumFromDotForm(GetModelInfoResponsePIDString)
	GetLinkedModelsRequestPID, _ = bw2.PONumFromDotForm(GetLinkedModelsRequestPIDString)
	GetLinkedModelsResponsePID, _ = bw2.PONumFromDotForm(GetLinkedModelsResponsePIDString)
	GetAllModelReplicasRequestPID, _ = bw2.PONumFromDotForm(GetAllModelReplicasRequestPIDString)
	GetAllModelReplicasResponsePID, _ = bw2.PONumFromDotForm(GetAllModelReplicasResponsePIDString)
	GetModelReplicaInfoRequestPID, _ = bw2.PONumFromDotForm(GetModelReplicaInfoRequestPIDString)
	GetModelReplicaInfoResponsePID, _ = bw2.PONumFromDotForm(GetModelReplicaInfoResponsePIDString)
	GetContainerLogsRequestPID, _ = bw2.PONumFromDotForm(GetContainerLogsRequestPIDString)
	GetContainerLogsResponsePID, _ = bw2.PONumFromDotForm(GetContainerLogsResponsePIDString)
	InspectInstanceRequestPID, _ = bw2.PONumFromDotForm(InspectInstanceRequestPIDString)
	InspectInstanceResponsePID, _ = bw2.PONumFromDotForm(InspectInstanceResponsePIDString)
	SetModelVersionRequestPID, _ = bw2.PONumFromDotForm(SetModelVersionRequestPIDString)
	SetModelVersionResponsePID, _ = bw2.PONumFromDotForm(SetModelVersionResponsePIDString)
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
	MsgID              int64  `json:"-"`
	Name               string `json:"name"`
	Input_type         string `json:"input_type"`
	Default_output     string `json:"default_output"`
	Latency_slo_micros int64  `json:"latency_slo_micros"`
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

type ListModelRequest struct {
	MsgID   int64 `json:"-"`
	Verbose bool  `json:"verbose"`
}

func (msg *ListModelRequest) Response() *ListModelResponse {
	return &ListModelResponse{
		MsgID: msg.MsgID,
	}
}

type ListModelResponse struct {
	MsgID             int64
	Error             string
	ModelNames        []string
	ModelDescriptions []ModelInfo
}

func (msg *ListModelResponse) PayloadObject() bw2.PayloadObject {
	po, _ := bw2.CreateMsgPackPayloadObject(ListModelResponsePID, msg)
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

type GetModelInfoRequest struct {
	MsgID         int64  `json:"-"`
	Model_name    string `json:"model_name"`
	Model_version string `json:"model_version"`
}

func (msg *GetModelInfoRequest) Response() *GetModelInfoResponse {
	return &GetModelInfoResponse{
		MsgID: msg.MsgID,
	}
}

type GetModelInfoResponse struct {
	MsgID int64
	Error string
	Info  ModelInfo
}

func (msg *GetModelInfoResponse) PayloadObject() bw2.PayloadObject {
	po, _ := bw2.CreateMsgPackPayloadObject(GetModelInfoResponsePID, msg)
	return po
}

type GetLinkedModelsRequest struct {
	MsgID    int64  `json:"-"`
	App_name string `json:"app_name"`
}

func (msg *GetLinkedModelsRequest) Response() *GetLinkedModelsResponse {
	return &GetLinkedModelsResponse{
		MsgID: msg.MsgID,
	}
}

type GetLinkedModelsResponse struct {
	MsgID  int64
	Error  string
	Models []string
}

func (msg *GetLinkedModelsResponse) PayloadObject() bw2.PayloadObject {
	po, _ := bw2.CreateMsgPackPayloadObject(GetLinkedModelsResponsePID, msg)
	return po
}

type GetAllModelReplicasRequest struct {
	MsgID   int64 `json:"-"`
	Verbose bool  `json:"verbose"`
}

func (msg *GetAllModelReplicasRequest) Response() *GetAllModelReplicasResponse {
	return &GetAllModelReplicasResponse{
		MsgID: msg.MsgID,
	}
}

type GetAllModelReplicasResponse struct {
	MsgID               int64
	Error               string
	ReplicaNames        []string
	ReplicaDescriptions []ReplicaInfo
}

func (msg *GetAllModelReplicasResponse) PayloadObject() bw2.PayloadObject {
	po, _ := bw2.CreateMsgPackPayloadObject(GetAllModelReplicasResponsePID, msg)
	return po
}

type ReplicaInfo struct {
	Model_id         string `json:"model_id"`
	Model_name       string `json:"model_name"`
	Model_version    string `json:"model_version"`
	Model_replica_id int    `json:"model_replica_id"`
	Input_type       string `json:"input_type"`
}

type GetModelReplicaInfoRequest struct {
	MsgID         int64  `json:"-"`
	Model_name    string `json:"model_name"`
	Model_version string `json:"model_version"`
	Replica_id    int    `json:"replica_id"`
}

func (msg *GetModelReplicaInfoRequest) Response() *GetModelReplicaInfoResponse {
	return &GetModelReplicaInfoResponse{
		MsgID: msg.MsgID,
	}
}

type GetModelReplicaInfoResponse struct {
	MsgID int64
	Error string
}

func (msg *GetModelReplicaInfoResponse) PayloadObject() bw2.PayloadObject {
	po, _ := bw2.CreateMsgPackPayloadObject(GetModelReplicaInfoResponsePID, msg)
	return po
}

type GetContainerLogsRequest struct {
	MsgID       int64  `json:"-"`
	ContainerID string `json:"ContainerID"`
}

func (msg *GetContainerLogsRequest) Response() *GetContainerLogsResponse {
	return &GetContainerLogsResponse{
		MsgID: msg.MsgID,
	}
}

type GetContainerLogsResponse struct {
	MsgID  int64
	Error  string
	Stdout string
	Stderr string
}

func (msg *GetContainerLogsResponse) PayloadObject() bw2.PayloadObject {
	po, _ := bw2.CreateMsgPackPayloadObject(GetContainerLogsResponsePID, msg)
	return po
}

type InspectInstanceRequest struct {
	MsgID int64 `json:"-"`
}

func (msg *InspectInstanceRequest) Response() *InspectInstanceResponse {
	return &InspectInstanceResponse{
		MsgID: msg.MsgID,
	}
}

type InspectInstanceResponse struct {
	MsgID int64
	Error string
	Info  interface{}
}

func (msg *InspectInstanceResponse) PayloadObject() bw2.PayloadObject {
	po, _ := bw2.CreateMsgPackPayloadObject(InspectInstanceResponsePID, msg)
	return po
}

type SetModelVersionRequest struct {
	MsgID         int64  `json:"-"`
	Model_name    string `json:"model_name"`
	Model_version string `json:"model_version"`
}

func (msg *SetModelVersionRequest) Response() *SetModelVersionResponse {
	return &SetModelVersionResponse{
		MsgID: msg.MsgID,
	}
}

type SetModelVersionResponse struct {
	MsgID int64
	Error string
}

func (msg *SetModelVersionResponse) PayloadObject() bw2.PayloadObject {
	po, _ := bw2.CreateMsgPackPayloadObject(SetModelVersionResponsePID, msg)
	return po
}
