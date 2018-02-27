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
			switch ponum {
			case GetReplicasRequestPID:
				var request GetReplicasMessageRequest
				if obj, ok := po.(bw2.MsgPackPayloadObject); !ok {
					log.Println("Received query was not msgpack")
				} else if err := obj.ValueInto(&request); err != nil {
					log.Println(errors.Wrap(err, "Could not unmarshal received query"))
					return
				}
				log.Println("Got Request", request)
				containers, err := mgr.GetContainersWithLabel(request.Label)
				var resp = request.Response()
				if err != nil {
					log.Println(errors.Wrap(err, "Could not get containers"))
					resp.Error = err.Error()
				} else {
					resp.Containers = containers
				}
				log.Println("Response on", iface.SignalURI("response"))
				if err := client.Publish(&bw2.PublishParams{
					URI:            iface.SignalURI("response"),
					PayloadObjects: []bw2.PayloadObject{resp.PayloadObject()},
				}); err != nil {
					log.Println(errors.Wrap(err, "could not publish respones"))
				}
			case DeployModelRequestPID:
				var request DeployModelRequest
				if obj, ok := po.(bw2.MsgPackPayloadObject); !ok {
					log.Println("Received query was not msgpack")
				} else if err := obj.ValueInto(&request); err != nil {
					log.Println(errors.Wrap(err, "Could not unmarshal received query"))
					return
				}
				log.Println("Got Request", request)
				resp := request.Response()
				err := mgr.DeployModel(request.Name, request.Version, request.Input_type, request.Image, 1)
				if err != nil {
					log.Println(errors.Wrap(err, "Could not deploy model"))
					resp.Error = err.Error()
				}
				log.Println("Response on", iface.SignalURI("response"))
				if err := client.Publish(&bw2.PublishParams{
					URI:            iface.SignalURI("response"),
					PayloadObjects: []bw2.PayloadObject{resp.PayloadObject()},
				}); err != nil {
					log.Println(errors.Wrap(err, "could not publish respones"))
				}
			case RegisterApplicationRequestPID:
				var request RegisterApplicationRequest
				if obj, ok := po.(bw2.MsgPackPayloadObject); !ok {
					log.Println("Received query was not msgpack")
				} else if err := obj.ValueInto(&request); err != nil {
					log.Println(errors.Wrap(err, "Could not unmarshal received query"))
					return
				}
				resp := request.Response()
				log.Println("Got Request", request)
				err := mgr.RegisterApplication(request)
				if err != nil {
					log.Println(errors.Wrap(err, "Could not register app"))
					resp.Error = err.Error()
				}
				log.Println("Response on", iface.SignalURI("response"))
				if err := client.Publish(&bw2.PublishParams{
					URI:            iface.SignalURI("response"),
					PayloadObjects: []bw2.PayloadObject{resp.PayloadObject()},
				}); err != nil {
					log.Println(errors.Wrap(err, "could not publish respones"))
				}

			case LinkModelToAppRequestPID:
				var request LinkModelToAppRequest
				if obj, ok := po.(bw2.MsgPackPayloadObject); !ok {
					log.Println("Received query was not msgpack")
				} else if err := obj.ValueInto(&request); err != nil {
					log.Println(errors.Wrap(err, "Could not unmarshal received query"))
					return
				}
				resp := request.Response()
				log.Println("Got Request", request)
				err := mgr.LinkModelToApp(request)
				if err != nil {
					log.Println(errors.Wrap(err, "Could not link model to app"))
					resp.Error = err.Error()
				}
				log.Println("Response on", iface.SignalURI("response"))
				if err := client.Publish(&bw2.PublishParams{
					URI:            iface.SignalURI("response"),
					PayloadObjects: []bw2.PayloadObject{resp.PayloadObject()},
				}); err != nil {
					log.Println(errors.Wrap(err, "could not publish respones"))
				}

			case RegisterModelRequestPID:
				var request RegisterModelRequest
				if obj, ok := po.(bw2.MsgPackPayloadObject); !ok {
					log.Println("Received query was not msgpack")
				} else if err := obj.ValueInto(&request); err != nil {
					log.Println(errors.Wrap(err, "Could not unmarshal received query"))
					return
				}
				resp := request.Response()
				log.Println("Got Request", request)
				err := mgr.RegisterModel(request)
				if err != nil {
					log.Println(errors.Wrap(err, "Could not register model"))
					resp.Error = err.Error()
				}
				log.Println("Response on", iface.SignalURI("response"))
				if err := client.Publish(&bw2.PublishParams{
					URI:            iface.SignalURI("response"),
					PayloadObjects: []bw2.PayloadObject{resp.PayloadObject()},
				}); err != nil {
					log.Println(errors.Wrap(err, "could not publish respones"))
				}
			case ListModelRequestPID:
				var request ListModelRequest
				if obj, ok := po.(bw2.MsgPackPayloadObject); !ok {
					log.Println("Received query was not msgpack")
				} else if err := obj.ValueInto(&request); err != nil {
					log.Println(errors.Wrap(err, "Could not unmarshal received query"))
					return
				}
				resp := request.Response()
				log.Println("Got Request", request)
				names, models, err := mgr.GetAllModels(request)
				if len(names) > 0 {
					resp.ModelNames = names
				} else {
					resp.ModelDescriptions = models
				}
				log.Println(models)
				if err != nil {
					log.Println(errors.Wrap(err, "Could not register model"))
					resp.Error = err.Error()
				}
				log.Println("Response on", iface.SignalURI("response"))
				if err := client.Publish(&bw2.PublishParams{
					URI:            iface.SignalURI("response"),
					PayloadObjects: []bw2.PayloadObject{resp.PayloadObject()},
				}); err != nil {
					log.Println(errors.Wrap(err, "could not publish respones"))
				}
			case ApplicationListRequestPID:
				var request ApplicationListRequest
				if obj, ok := po.(bw2.MsgPackPayloadObject); !ok {
					log.Println("Received query was not msgpack")
				} else if err := obj.ValueInto(&request); err != nil {
					log.Println(errors.Wrap(err, "Could not unmarshal received query"))
					return
				}
				resp := request.Response()
				log.Println("Got Request", request)
				names, models, err := mgr.GetAllApplications(request)
				if len(names) > 0 {
					resp.ApplicationNames = names
				} else {
					resp.ApplicationDescriptions = models
				}
				log.Println(models)
				if err != nil {
					log.Println(errors.Wrap(err, "Could not list apps"))
					resp.Error = err.Error()
				}
				log.Println("Response on", iface.SignalURI("response"))
				if err := client.Publish(&bw2.PublishParams{
					URI:            iface.SignalURI("response"),
					PayloadObjects: []bw2.PayloadObject{resp.PayloadObject()},
				}); err != nil {
					log.Println(errors.Wrap(err, "could not publish respones"))
				}

			case GetApplicationInfoRequestPID:
				var request GetApplicationInfoRequest
				if obj, ok := po.(bw2.MsgPackPayloadObject); !ok {
					log.Println("Received query was not msgpack")
				} else if err := obj.ValueInto(&request); err != nil {
					log.Println(errors.Wrap(err, "Could not unmarshal received query"))
					return
				}
				resp := request.Response()
				log.Println("Got Request", request)
				info, err := mgr.GetApplicationInfo(request)
				resp.Info = info
				if err != nil {
					log.Println(errors.Wrap(err, "Could not get application info"))
					resp.Error = err.Error()
				}
				log.Println("Response on", iface.SignalURI("response"))
				if err := client.Publish(&bw2.PublishParams{
					URI:            iface.SignalURI("response"),
					PayloadObjects: []bw2.PayloadObject{resp.PayloadObject()},
				}); err != nil {
					log.Println(errors.Wrap(err, "could not publish respones"))
				}

			case GetModelInfoRequestPID:
				var request GetModelInfoRequest
				if obj, ok := po.(bw2.MsgPackPayloadObject); !ok {
					log.Println("Received query was not msgpack")
				} else if err := obj.ValueInto(&request); err != nil {
					log.Println(errors.Wrap(err, "Could not unmarshal received query"))
					return
				}
				resp := request.Response()
				log.Println("Got Request", request)
				info, err := mgr.GetModelInfo(request)
				resp.Info = info
				if err != nil {
					log.Println(errors.Wrap(err, "Could not get model info"))
					resp.Error = err.Error()
				}
				log.Println("Response on", iface.SignalURI("response"))
				if err := client.Publish(&bw2.PublishParams{
					URI:            iface.SignalURI("response"),
					PayloadObjects: []bw2.PayloadObject{resp.PayloadObject()},
				}); err != nil {
					log.Println(errors.Wrap(err, "could not publish respones"))
				}

			case GetLinkedModelsRequestPID:
				var request GetLinkedModelsRequest
				if obj, ok := po.(bw2.MsgPackPayloadObject); !ok {
					log.Println("Received query was not msgpack")
				} else if err := obj.ValueInto(&request); err != nil {
					log.Println(errors.Wrap(err, "Could not unmarshal received query"))
					return
				}
				resp := request.Response()
				log.Println("Got Request", request)
				models, err := mgr.GetLinkedModels(request)
				log.Println("models", models)
				resp.Models = models
				if err != nil {
					log.Println(errors.Wrap(err, "Could not get model info"))
					resp.Error = err.Error()
				}
				log.Println("Response on", iface.SignalURI("response"))
				if err := client.Publish(&bw2.PublishParams{
					URI:            iface.SignalURI("response"),
					PayloadObjects: []bw2.PayloadObject{resp.PayloadObject()},
				}); err != nil {
					log.Println(errors.Wrap(err, "could not publish respones"))
				}

			case GetAllModelReplicasRequestPID:
				var request GetAllModelReplicasRequest
				if obj, ok := po.(bw2.MsgPackPayloadObject); !ok {
					log.Println("Received query was not msgpack")
				} else if err := obj.ValueInto(&request); err != nil {
					log.Println(errors.Wrap(err, "Could not unmarshal received query"))
					return
				}
				resp := request.Response()
				log.Println("Got Request", request)
				err := mgr.GetAllModelReplicas(request)
				if err != nil {
					log.Println(errors.Wrap(err, "Could not get model info"))
					resp.Error = err.Error()
				}
				log.Println("Response on", iface.SignalURI("response"))
				if err := client.Publish(&bw2.PublishParams{
					URI:            iface.SignalURI("response"),
					PayloadObjects: []bw2.PayloadObject{resp.PayloadObject()},
				}); err != nil {
					log.Println(errors.Wrap(err, "could not publish respones"))
				}

			default:
				log.Println(ponum)
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
