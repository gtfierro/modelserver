package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/levigross/grequests"
	"github.com/pkg/errors"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const STOP_TIMEOUT = 30 * time.Second
const CLIPPER_DOCKER_LABEL = "ai.clipper.container.label"
const CLIPPER_MODEL_CONTAINER_LABEL = "ai.clipper.model_container.label"
const CLIPPER_QUERY_FRONTEND_CONTAINER_LABEL = "ai.clipper.query_frontend.label"
const CLIPPER_MGMT_FRONTEND_CONTAINER_LABEL = "ai.clipper.management_frontend.label"
const CONTAINERLESS_MODEL_IMAGE = "NO_CONTAINER"
const _MODEL_CONTAINER_LABEL_DELIMITER = "_"
const CLIPPER_MANAGEMENT_PORT = "1338"
const CLIPPER_QUERY_PORT = "1337"

// implementation of Clipper API client

func createModelContainerLabel(name, version string) string {
	return fmt.Sprintf("%s%s%s", name, _MODEL_CONTAINER_LABEL_DELIMITER, version)
}
func parseModelContainerLabel(label string) (name, version string, err error) {
	parts := strings.Split(label, _MODEL_CONTAINER_LABEL_DELIMITER)
	if len(parts) != 2 {
		return "", "", errors.Errorf("Unable to parse model container label %s", label)
	}
	return parts[0], parts[1], nil
}

type DockerContainerManager struct {
	c             *client.Client
	defaultLabels map[string]string
	networkID     string // clipper network
}

// TODO: extra_container_kwargs
func NewDockerContainerManager(defaultLabels map[string]string) (*DockerContainerManager, error) {
	c, err := client.NewEnvClient()
	if err != nil {
		return nil, errors.Wrap(err, "Could not connect to Docker daemon")
	}
	defaultLabels[CLIPPER_DOCKER_LABEL] = ""
	mgr := &DockerContainerManager{
		c:             c,
		defaultLabels: defaultLabels,
	}

	// create network
	net_cfg, err := c.NetworkCreate(context.TODO(), "clipper_network", types.NetworkCreate{})
	if err != nil {
		return nil, errors.Wrap(err, "Could not create clipper_network network")
	}
	mgr.networkID = net_cfg.ID

	return mgr, nil
}

// TODO: this is really set_num_replicas
func (mgr *DockerContainerManager) DeployModel(name, version, input_type, image string, num_replicas int) error {
	label := createModelContainerLabel(name, version)
	current_replicas, err := mgr.GetContainersWithLabel(label)
	if err != nil {
		return errors.Wrapf(err, "Could not fetch replicas for %s", label)
	}
	if len(current_replicas) < num_replicas {
		missing := num_replicas - len(current_replicas)
		log.Printf("Found %d replicas for %s. Adding %d", len(current_replicas), label, missing)
		for i := 0; i < missing; i++ {
			if err := mgr.StartContainer(name, version, input_type, image); err != nil {
				return errors.Wrapf(err, "Could not start %s when adding replicas", label)
			}
		}
	} else if len(current_replicas) > num_replicas {
		extra := len(current_replicas) - num_replicas
		log.Println("Found %d replicas for %s. Removing %d", len(current_replicas), label, extra)
		// TODO: stop container in loop
		for i := 0; i < extra; i++ {
			if err := mgr.c.ContainerStop(context.TODO(), current_replicas[i].ID, nil); err != nil {
				return errors.Wrapf(err, "Could not stop %s (container %s) when removing replicas", label, current_replicas[i].ID)
			}
		}
	}
	return nil
}

func (mgr *DockerContainerManager) RegisterApplication(req RegisterApplicationRequest) error {
	log.Println("POST to", fmt.Sprintf("http://localhost:%s/admin/add_app", CLIPPER_MANAGEMENT_PORT))
	resp, err := grequests.Post(fmt.Sprintf("http://localhost:%s/admin/add_app", CLIPPER_MANAGEMENT_PORT),
		&grequests.RequestOptions{
			JSON:    req,
			Headers: map[string]string{"Content-type": "application/json"},
		})
	response := resp.String()
	if err != nil && err != io.EOF {
		return errors.Wrap(err, "Could not register app")
	}
	if resp.Ok != true {
		return errors.Errorf("Could not register app: %s", response)
	}
	log.Printf("Application %s was successfully registered", req.Name)
	return nil
}

func (mgr *DockerContainerManager) LinkModelToApp(req LinkModelToAppRequest) error {
	log.Println("POST to", fmt.Sprintf("http://localhost:%s/admin/add_model_links", CLIPPER_MANAGEMENT_PORT))
	resp, err := grequests.Post(fmt.Sprintf("http://localhost:%s/admin/add_model_links", CLIPPER_MANAGEMENT_PORT),
		&grequests.RequestOptions{
			JSON:    req,
			Headers: map[string]string{"Content-type": "application/json"},
		})
	response := resp.String()
	if err != nil && err != io.EOF {
		return errors.Wrap(err, "Could not link model to app")
	}
	if resp.Ok != true {
		return errors.Errorf("Could not link model to app: %s", response)
	}
	log.Printf("Model %s is now linked to application %s", req.Model_names[0], req.App_name)
	return nil
}

func (mgr *DockerContainerManager) GetLinkedModels(req GetLinkedModelsRequest) ([]string, error) {
	log.Println("POST to", fmt.Sprintf("http://localhost:%s/admin/get_linked_models", CLIPPER_MANAGEMENT_PORT))
	var models []string
	resp, err := grequests.Post(fmt.Sprintf("http://localhost:%s/admin/get_linked_models", CLIPPER_MANAGEMENT_PORT),
		&grequests.RequestOptions{
			JSON:    req,
			Headers: map[string]string{"Content-type": "application/json"},
		})
	response := resp.String()
	if err != nil && err != io.EOF {
		return models, errors.Wrap(err, "Could not link model to app")
	}
	if resp.Ok != true {
		return models, errors.Errorf("Could not link model to app: %s", response)
	}
	err = resp.JSON(&models)
	return models, err
}

func (mgr *DockerContainerManager) GetAllModelReplicas(req GetAllModelReplicasRequest) (replicas []string, infos []ReplicaInfo, err error) {
	log.Println("POST to", fmt.Sprintf("http://localhost:%s/admin/get_all_containers", CLIPPER_MANAGEMENT_PORT))
	resp, err := grequests.Post(fmt.Sprintf("http://localhost:%s/admin/get_all_containers", CLIPPER_MANAGEMENT_PORT),
		&grequests.RequestOptions{
			JSON:    req,
			Headers: map[string]string{"Content-type": "application/json"},
		})
	response := resp.String()
	if err != nil && err != io.EOF {
		err = errors.Wrap(err, "Could not link model to app")
		return
	}
	if resp.Ok != true {
		err = errors.Errorf("Could not link model to app: %s", response)
		return
	}

	if req.Verbose {
		err = resp.JSON(&infos)
	} else {
		err = resp.JSON(&replicas)
	}
	return
}

func (mgr *DockerContainerManager) GetModelReplicaInfo(req GetAllModelReplicasRequest) (info ReplicaInfo, err error) {
	log.Println("POST to", fmt.Sprintf("http://localhost:%s/admin/get_containers", CLIPPER_MANAGEMENT_PORT))
	resp, err := grequests.Post(fmt.Sprintf("http://localhost:%s/admin/get_containers", CLIPPER_MANAGEMENT_PORT),
		&grequests.RequestOptions{
			JSON:    req,
			Headers: map[string]string{"Content-type": "application/json"},
		})
	response := resp.String()
	if err != nil && err != io.EOF {
		err = errors.Wrap(err, "Could not link model to app")
		return
	}
	if resp.Ok != true {
		err = errors.Errorf("Could not link model to app: %s", response)
		return
	}

	err = resp.JSON(&info)
	return
}

func (mgr *DockerContainerManager) InspectInstance() (info interface{}, err error) {
	log.Println("GET to", fmt.Sprintf("http://localhost:%s/metrics", CLIPPER_QUERY_PORT))
	resp, err := grequests.Get(fmt.Sprintf("http://localhost:%s/metrics", CLIPPER_QUERY_PORT), nil)
	response := resp.String()
	if err != nil && err != io.EOF {
		err = errors.Wrap(err, "Could not inspect instance")
		return
	}
	if resp.Ok != true {
		err = errors.Errorf("Could not inspect instance: (%s)", response)
		return
	}

	err = resp.JSON(&info)
	return
}

func (mgr *DockerContainerManager) SetModelVersion(req SetModelVersionRequest) (err error) {
	log.Println("GET to", fmt.Sprintf("http://localhost:%s/admin/add_model", CLIPPER_MANAGEMENT_PORT))
	resp, err := grequests.Post(fmt.Sprintf("http://localhost:%s/admin/add_model", CLIPPER_MANAGEMENT_PORT),
		&grequests.RequestOptions{
			JSON:    req,
			Headers: map[string]string{"Content-type": "application/json"},
		})
	response := resp.String()
	if err != nil && err != io.EOF {
		err = errors.Wrap(err, "Could not inspect instance")
		return
	}
	if resp.Ok != true {
		err = errors.Errorf("Could not inspect instance: (%s)", response)
		return
	}

	log.Println(response)
	return
}

func (mgr *DockerContainerManager) GetContainerLogs(req GetContainerLogsRequest) (stdout string, stderr string, err error) {
	rdr, err := mgr.c.ContainerLogs(context.TODO(), req.ContainerID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Details:    true,
	})
	if err != nil {
		return "", "", err
	}
	_stdout := new(bytes.Buffer)
	_stderr := new(bytes.Buffer)
	_, err = stdcopy.StdCopy(_stdout, _stderr, rdr)

	stdout = _stdout.String()
	stderr = _stderr.String()
	log.Println("stdout", len(stdout), stdout)
	log.Println("stderr", len(stderr), stderr)
	return
}

func (mgr *DockerContainerManager) RegisterModel(req RegisterModelRequest) error {
	log.Println("POST to", fmt.Sprintf("http://localhost:%s/admin/add_model", CLIPPER_MANAGEMENT_PORT))
	if req.Labels == nil {
		req.Labels = []string{}
	}
	req.Model_data_path = "DEPRECATED"
	resp, err := grequests.Post(fmt.Sprintf("http://localhost:%s/admin/add_model", CLIPPER_MANAGEMENT_PORT),
		&grequests.RequestOptions{
			JSON:    req,
			Headers: map[string]string{"Content-type": "application/json"},
		})
	response := resp.String()
	if err != nil && err != io.EOF {
		return errors.Wrap(err, "Could not register model")
	}
	if resp.Ok != true {
		return errors.Errorf("Could not register model (%s)", response)
	}
	log.Printf("Successfully registered model %s:%s", req.Model_name, req.Model_version)
	return nil
}

func (mgr *DockerContainerManager) GetCurrentModelVersion(name string) (version string, err error) {
	// get all models
	return "", nil
}

func (mgr *DockerContainerManager) GetAllModels(req ListModelRequest) (names []string, infos []ModelInfo, err error) {
	log.Println("POST to", fmt.Sprintf("http://localhost:%s/admin/get_all_models", CLIPPER_MANAGEMENT_PORT))
	resp, err := grequests.Post(fmt.Sprintf("http://localhost:%s/admin/get_all_models", CLIPPER_MANAGEMENT_PORT),
		&grequests.RequestOptions{
			JSON:    req,
			Headers: map[string]string{"Content-type": "application/json"},
		})
	response := resp.String()
	if err != nil && err != io.EOF {
		err = errors.Wrap(err, "Could not get models")
		return
	}
	if resp.Ok != true {
		err = errors.Errorf("Could not get models (%s)", response)
		return
	}
	if req.Verbose {
		err = resp.JSON(&infos)
	} else {
		err = resp.JSON(&names)
	}
	return
}

func (mgr *DockerContainerManager) GetAllApplications(req ApplicationListRequest) (names []string, infos []ApplicationInfo, err error) {
	log.Println("POST to", fmt.Sprintf("http://localhost:%s/admin/get_all_applications", CLIPPER_MANAGEMENT_PORT))
	resp, err := grequests.Post(fmt.Sprintf("http://localhost:%s/admin/get_all_applications", CLIPPER_MANAGEMENT_PORT),
		&grequests.RequestOptions{
			JSON:    req,
			Headers: map[string]string{"Content-type": "application/json"},
		})
	response := resp.String()
	if err != nil && err != io.EOF {
		err = errors.Wrap(err, "Could not get models")
		return
	}
	if resp.Ok != true {
		err = errors.Errorf("Could not get models (%s)", response)
		return
	}
	if req.Verbose {
		err = resp.JSON(&infos)
	} else {
		err = resp.JSON(&names)
	}
	return
}

func (mgr *DockerContainerManager) GetApplicationInfo(req GetApplicationInfoRequest) (ApplicationInfo, error) {
	log.Println("POST to", fmt.Sprintf("http://localhost:%s/admin/get_application", CLIPPER_MANAGEMENT_PORT))
	var info ApplicationInfo
	resp, err := grequests.Post(fmt.Sprintf("http://localhost:%s/admin/get_application", CLIPPER_MANAGEMENT_PORT),
		&grequests.RequestOptions{
			JSON:    req,
			Headers: map[string]string{"Content-type": "application/json"},
		})
	response := resp.String()
	if err != nil && err != io.EOF {
		return info, errors.Wrap(err, "Could not get models")
	}
	if resp.Ok != true {
		return info, errors.Errorf("Could not get models (%s)", response)
	}
	err = resp.JSON(&info)
	return info, err
}

func (mgr *DockerContainerManager) GetModelInfo(req GetModelInfoRequest) (ModelInfo, error) {
	log.Println("POST to", fmt.Sprintf("http://localhost:%s/admin/get_model", CLIPPER_MANAGEMENT_PORT))
	var info ModelInfo
	resp, err := grequests.Post(fmt.Sprintf("http://localhost:%s/admin/get_model", CLIPPER_MANAGEMENT_PORT),
		&grequests.RequestOptions{
			JSON:    req,
			Headers: map[string]string{"Content-type": "application/json"},
		})
	response := resp.String()
	if err != nil && err != io.EOF {
		return info, errors.Wrap(err, "Could not get models")
	}
	if resp.Ok != true {
		return info, errors.Errorf("Could not get models (%s)", response)
	}
	err = resp.JSON(&info)
	return info, err
	return ModelInfo{}, nil
}

// https://docs.docker.com/develop/sdk/examples/#run-a-container
// https://godoc.org/github.com/docker/docker/client#Client.ContainerList
func (mgr *DockerContainerManager) GetContainersWithLabel(label string) ([]types.Container, error) {
	args := filters.NewArgs()

	ctx := context.TODO()

	args.Add("label", label)
	opts := types.ContainerListOptions{
		Filters: args,
	}
	log.Printf("find %+v", opts)
	return mgr.c.ContainerList(ctx, opts)
}

func (mgr *DockerContainerManager) GetReplicas(name, version string) ([]types.Container, error) {
	label := createModelContainerLabel(name, version)
	return mgr.GetContainersWithLabel(label)
}

func (mgr *DockerContainerManager) StartContainer(name, version, input_type, image string) error {
	containers, err := mgr.GetContainersWithLabel(CLIPPER_QUERY_FRONTEND_CONTAINER_LABEL)
	if err != nil {
		return errors.Wrapf(err, "Could not list replicas for %s:%s", name, version)
	}
	if len(containers) < 1 {
		log.Println("No Clipper query frontend found")
		return errors.New("No Clipper query frontend to attach model container to")
	}
	query_frontend_hostname := strings.TrimPrefix(containers[0].Names[0], "/")
	log.Println("query frontend hostname:", query_frontend_hostname)
	env_vars := map[string]string{
		"CLIPPER_MODEL_NAME":    name,
		"CLIPPER_MODEL_VERSION": version,
		"CLIPPER_IP":            query_frontend_hostname,
		"CLIPPER_INPUT_TYPE":    input_type,
	}
	var transformed_env_vars []string
	for k, v := range env_vars {
		transformed_env_vars = append(transformed_env_vars, fmt.Sprintf("%s=%v", k, v))
	}
	model_container_label := createModelContainerLabel(name, version)
	model_container_name := fmt.Sprintf("%s-%d", model_container_label, rand.Int63n(100000))

	labels := mgr.defaultLabels
	labels[CLIPPER_MODEL_CONTAINER_LABEL] = ""
	ctx := context.TODO()
	config := &container.Config{
		Image:  image,
		Tty:    true,
		Env:    transformed_env_vars,
		Labels: labels,
	}
	var networkCfg *network.NetworkingConfig
	var hostCfg *container.HostConfig

	created, err := mgr.c.ContainerCreate(ctx, config, hostCfg, networkCfg, model_container_name)
	if err != nil {
		return errors.Wrap(err, "Could not create container")
	}

	log.Println("Created container", created.ID)

	// connect to network
	if err = mgr.c.NetworkConnect(ctx, mgr.networkID, created.ID, &network.EndpointSettings{}); err != nil {
		return errors.Wrapf(err, "Could not connect container to network %s", mgr.networkID)
	}

	err = mgr.c.ContainerStart(ctx, created.ID, types.ContainerStartOptions{})
	if err != nil {
		return errors.Wrap(err, "Could not start container")
	}

	// TODO:
	//  add_to_metric_config(model_container_name,
	//                       CLIPPER_INTERNAL_METRIC_PORT)
	log.Println("Started container", created.ID)

	return nil
}

// models is map of name=>list of versions
func (mgr *DockerContainerManager) StopModels(models map[string][]string) error {
	containers, err := mgr.GetContainersWithLabel(CLIPPER_MODEL_CONTAINER_LABEL)
	if err != nil {
		return errors.Wrapf(err, "Could not list containers with label '%s'", CLIPPER_MODEL_CONTAINER_LABEL)
	}
	for _, container := range containers {
		name, version, err := parseModelContainerLabel(container.Labels[CLIPPER_MODEL_CONTAINER_LABEL])
		if err != nil {
			return err
		}
		if versions, found := models[name]; !found {
			continue
		} else {
			for _, delversion := range versions {
				if delversion == version {
					if err := mgr.c.ContainerStop(context.TODO(), container.ID, nil); err != nil {
						return errors.Wrapf(err, "Could not stop (container %s) when removing replicas", container.ID)
					}
				}
			}
		}
	}
	return nil
}

// TODO: implement
//func (mgr *DockerContainerManager) GetLogs(loggingdir string) error {
//    return nil
//}

func (mgr *DockerContainerManager) StopAllModelContainers() error {
	containers, err := mgr.GetContainersWithLabel(CLIPPER_MODEL_CONTAINER_LABEL)
	if err != nil {
		return errors.Wrapf(err, "Could not list containers with label '%s'", CLIPPER_MODEL_CONTAINER_LABEL)
	}
	for _, container := range containers {
		if err := mgr.c.ContainerStop(context.TODO(), container.ID, nil); err != nil {
			return errors.Wrapf(err, "Could not stop (container %s) when removing replicas", container.ID)
		}
	}
	return nil
}

func (mgr *DockerContainerManager) StopAll() error {
	containers, err := mgr.GetContainersWithLabel(CLIPPER_DOCKER_LABEL)
	if err != nil {
		return errors.Wrapf(err, "Could not list containers with label '%s'", CLIPPER_MODEL_CONTAINER_LABEL)
	}
	for _, container := range containers {
		if err := mgr.c.ContainerStop(context.TODO(), container.ID, nil); err != nil {
			return errors.Wrapf(err, "Could not stop (container %s) when removing replicas", container.ID)
		}
	}
	return nil
}
