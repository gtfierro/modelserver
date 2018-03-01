package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"

	bw2 "github.com/immesys/bw2bind"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

type server struct {
}

func newServer() *server {
	return &server{}
}

func (srv *server) newModel(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	defer req.Body.Close()

	dec := json.NewDecoder(req.Body)
	var model Model
	if err := dec.Decode(&model); err != nil {
		http.Error(rw, err.Error(), 500)
		return
	}
	log.Printf("%+v", model)
	rw.WriteHeader(200)
}

func main() {
	mgr, err := NewDockerContainerManager(map[string]string{"boo": "ya"})
	if err != nil {
		panic(err)
	}
	//err = mgr.StartContainer("model", "1", "inp", "gtfierro/hod:latest")
	//if err != nil {
	//	panic(err)
	//}

	client := bw2.ConnectOrExit("")
	client.OverrideAutoChainTo(true)
	client.SetEntityFromEnvironOrExit()
	svc := client.RegisterService("scratch.ns", "s.modelserver")
	iface := svc.RegisterInterface("_", "i.modelserver")
	requests, err := client.Subscribe(&bw2.SubscribeParams{
		URI: iface.SlotURI("request"),
	})
	if err != nil {
		panic(err)
	}
	log.Println("listening on", iface.SlotURI("request"))
	//TODO: use worker pattern from mdal
	for msg := range requests {
		msg.Dump()
		for _, po := range msg.POs {
			ponum := po.GetPONum()
			var resp Response
			switch ponum {
			case GetReplicasRequestPID:
				var request GetReplicasMessageRequest
				if obj, ok := po.(bw2.MsgPackPayloadObject); !ok {
					log.Println("Received query was not msgpack")
				} else if err := obj.ValueInto(&request); err != nil {
					log.Println(errors.Wrap(err, "Could not unmarshal received query"))
					return
				}
				log.Printf("Got Request %+v", request)
				containers, err := mgr.GetContainersWithLabel(request.Label)
				_resp := request.Response()
				if err != nil {
					log.Println(errors.Wrap(err, "Could not get containers"))
					_resp.Error = err.Error()
				} else {
					_resp.Containers = containers
				}
				resp = _resp
			case DeployModelRequestPID:
				var request DeployModelRequest
				if obj, ok := po.(bw2.MsgPackPayloadObject); !ok {
					log.Println("Received query was not msgpack")
				} else if err := obj.ValueInto(&request); err != nil {
					log.Println(errors.Wrap(err, "Could not unmarshal received query"))
					return
				}
				log.Printf("Got Request %+v", request)
				_resp := request.Response()
				err := mgr.DeployModel(request.Name, request.Version, request.Input_type, request.Image, 1)
				if err != nil {
					log.Println(errors.Wrap(err, "Could not deploy model"))
					_resp.Error = err.Error()
				}
				resp = _resp
			case RegisterApplicationRequestPID:
				var request RegisterApplicationRequest
				if obj, ok := po.(bw2.MsgPackPayloadObject); !ok {
					log.Println("Received query was not msgpack")
				} else if err := obj.ValueInto(&request); err != nil {
					log.Println(errors.Wrap(err, "Could not unmarshal received query"))
					return
				}
				log.Printf("Got Request %+v", request)
				err := mgr.RegisterApplication(request)
				_resp := request.Response()
				if err != nil {
					log.Println(errors.Wrap(err, "Could not register app"))
					_resp.Error = err.Error()
				}
				resp = _resp

			case LinkModelToAppRequestPID:
				var request LinkModelToAppRequest
				if obj, ok := po.(bw2.MsgPackPayloadObject); !ok {
					log.Println("Received query was not msgpack")
				} else if err := obj.ValueInto(&request); err != nil {
					log.Println(errors.Wrap(err, "Could not unmarshal received query"))
					return
				}
				_resp := request.Response()
				log.Printf("Got Request %+v", request)
				err := mgr.LinkModelToApp(request)
				if err != nil {
					log.Println(errors.Wrap(err, "Could not link model to app"))
					_resp.Error = err.Error()
				}
				resp = _resp

			case RegisterModelRequestPID:
				var request RegisterModelRequest
				if obj, ok := po.(bw2.MsgPackPayloadObject); !ok {
					log.Println("Received query was not msgpack")
				} else if err := obj.ValueInto(&request); err != nil {
					log.Println(errors.Wrap(err, "Could not unmarshal received query"))
					return
				}
				_resp := request.Response()
				log.Printf("Got Request %+v", request)
				err := mgr.RegisterModel(request)
				if err != nil {
					log.Println(errors.Wrap(err, "Could not register model"))
					_resp.Error = err.Error()
				}
				resp = _resp

			case ListModelRequestPID:
				var request ListModelRequest
				if obj, ok := po.(bw2.MsgPackPayloadObject); !ok {
					log.Println("Received query was not msgpack")
				} else if err := obj.ValueInto(&request); err != nil {
					log.Println(errors.Wrap(err, "Could not unmarshal received query"))
					return
				}
				_resp := request.Response()
				log.Printf("Got Request %+v", request)
				names, models, err := mgr.GetAllModels(request)
				if len(names) > 0 {
					_resp.ModelNames = names
				} else {
					_resp.ModelDescriptions = models
				}
				if err != nil {
					log.Println(errors.Wrap(err, "Could not register model"))
					_resp.Error = err.Error()
				}
				resp = _resp

			case ApplicationListRequestPID:
				var request ApplicationListRequest
				if obj, ok := po.(bw2.MsgPackPayloadObject); !ok {
					log.Println("Received query was not msgpack")
				} else if err := obj.ValueInto(&request); err != nil {
					log.Println(errors.Wrap(err, "Could not unmarshal received query"))
					return
				}
				_resp := request.Response()
				log.Printf("Got Request %+v", request)
				names, models, err := mgr.GetAllApplications(request)
				if len(names) > 0 {
					_resp.ApplicationNames = names
				} else {
					_resp.ApplicationDescriptions = models
				}
				if err != nil {
					log.Println(errors.Wrap(err, "Could not list apps"))
					_resp.Error = err.Error()
				}
				resp = _resp

			case GetApplicationInfoRequestPID:
				var request GetApplicationInfoRequest
				if obj, ok := po.(bw2.MsgPackPayloadObject); !ok {
					log.Println("Received query was not msgpack")
				} else if err := obj.ValueInto(&request); err != nil {
					log.Println(errors.Wrap(err, "Could not unmarshal received query"))
					return
				}
				_resp := request.Response()
				log.Printf("Got Request %+v", request)
				info, err := mgr.GetApplicationInfo(request)
				_resp.Info = info
				if err != nil {
					log.Println(errors.Wrap(err, "Could not get application info"))
					_resp.Error = err.Error()
				}
				resp = _resp

			case GetModelInfoRequestPID:
				var request GetModelInfoRequest
				if obj, ok := po.(bw2.MsgPackPayloadObject); !ok {
					log.Println("Received query was not msgpack")
				} else if err := obj.ValueInto(&request); err != nil {
					log.Println(errors.Wrap(err, "Could not unmarshal received query"))
					return
				}
				_resp := request.Response()
				log.Printf("Got Request %+v", request)
				info, err := mgr.GetModelInfo(request)
				_resp.Info = info
				if err != nil {
					log.Println(errors.Wrap(err, "Could not get model info"))
					_resp.Error = err.Error()
				}
				resp = _resp

			case GetLinkedModelsRequestPID:
				var request GetLinkedModelsRequest
				if obj, ok := po.(bw2.MsgPackPayloadObject); !ok {
					log.Println("Received query was not msgpack")
				} else if err := obj.ValueInto(&request); err != nil {
					log.Println(errors.Wrap(err, "Could not unmarshal received query"))
					return
				}
				_resp := request.Response()
				log.Printf("Got Request %+v", request)
				models, err := mgr.GetLinkedModels(request)
				_resp.Models = models
				if err != nil {
					log.Println(errors.Wrap(err, "Could not get linked models"))
					_resp.Error = err.Error()
				}
				resp = _resp

			case GetAllModelReplicasRequestPID:
				var request GetAllModelReplicasRequest
				if obj, ok := po.(bw2.MsgPackPayloadObject); !ok {
					log.Println("Received query was not msgpack")
				} else if err := obj.ValueInto(&request); err != nil {
					log.Println(errors.Wrap(err, "Could not unmarshal received query"))
					return
				}
				_resp := request.Response()
				log.Printf("Got Request %+v", request)
				replicas, infos, err := mgr.GetAllModelReplicas(request)
				_resp.ReplicaNames = replicas
				_resp.ReplicaDescriptions = infos
				if err != nil {
					log.Println(errors.Wrap(err, "Could not get all model replicas"))
					_resp.Error = err.Error()
				}
				resp = _resp

			case GetModelReplicaInfoRequestPID:
				var request GetModelInfoRequest
				if obj, ok := po.(bw2.MsgPackPayloadObject); !ok {
					log.Println("Received query was not msgpack")
				} else if err := obj.ValueInto(&request); err != nil {
					log.Println(errors.Wrap(err, "Could not unmarshal received query"))
					return
				}
				_resp := request.Response()
				log.Printf("Got Request %+v", request)
				info, err := mgr.GetModelInfo(request)
				_resp.Info = info
				if err != nil {
					log.Println(errors.Wrap(err, "Could not get model replica info"))
					_resp.Error = err.Error()
				}
				resp = _resp

			case GetContainerLogsRequestPID:
				var request GetContainerLogsRequest
				if obj, ok := po.(bw2.MsgPackPayloadObject); !ok {
					log.Println("Received query was not msgpack")
				} else if err := obj.ValueInto(&request); err != nil {
					log.Println(errors.Wrap(err, "Could not unmarshal received query"))
					return
				}
				_resp := request.Response()
				log.Printf("Got Request %+v", request)
				stdout, stderr, err := mgr.GetContainerLogs(request)
				_resp.Stdout = stdout
				_resp.Stderr = stderr
				if err != nil {
					log.Println(errors.Wrap(err, "Could not get container logs"))
					_resp.Error = err.Error()
				}

			case InspectInstanceRequestPID:
				var request InspectInstanceRequest
				if obj, ok := po.(bw2.MsgPackPayloadObject); !ok {
					log.Println("Received query was not msgpack")
				} else if err := obj.ValueInto(&request); err != nil {
					log.Println(errors.Wrap(err, "Could not unmarshal received query"))
					return
				}
				_resp := request.Response()
				log.Printf("Got Request %+v", request)
				info, err := mgr.InspectInstance()
				_resp.Info = info
				if err != nil {
					log.Println(errors.Wrap(err, "Could not inspect instance"))
					_resp.Error = err.Error()
				}
				resp = _resp

			case SetModelVersionRequestPID:
				var request SetModelVersionRequest
				if obj, ok := po.(bw2.MsgPackPayloadObject); !ok {
					log.Println("Received query was not msgpack")
				} else if err := obj.ValueInto(&request); err != nil {
					log.Println(errors.Wrap(err, "Could not unmarshal received query"))
					return
				}
				_resp := request.Response()
				log.Printf("Got Request %+v", request)
				err := mgr.SetModelVersion(request)
				if err != nil {
					log.Println(errors.Wrap(err, "Could not set model version"))
					_resp.Error = err.Error()
				}
				resp = _resp
			default:
				log.Println(ponum)
				continue
			}
			log.Println("Response on", iface.SignalURI("response"))
			if err := client.Publish(&bw2.PublishParams{
				URI:            iface.SignalURI("response"),
				PayloadObjects: []bw2.PayloadObject{resp.PayloadObject()},
			}); err != nil {
				log.Println(errors.Wrap(err, "could not publish respones"))
			}
		}
	}

	addrString := "127.0.0.1:5555"
	r := httprouter.New()
	srv := newServer()
	r.POST("/api/model", srv.newModel)
	http.Handle("/", r)
	address, err := net.ResolveTCPAddr("tcp4", addrString)
	if err != nil {
		log.Fatalf("Error resolving address %s (%s)", addrString, err.Error())
	}
	_srv := &http.Server{
		Addr: address.String(),
	}
	log.Println("Starting HTTP Server on ", addrString)
	log.Fatal(_srv.ListenAndServe())
}
