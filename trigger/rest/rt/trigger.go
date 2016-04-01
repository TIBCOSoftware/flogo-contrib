package rest

import (
	"encoding/json"
	"net/http"

	"github.com/TIBCOSoftware/flogo-lib/engine/ext/trigger"
	"github.com/TIBCOSoftware/flogo-lib/engine/starter"

	"github.com/julienschmidt/httprouter"
	"github.com/op/go-logging"
)

// log is the default package logger
var log = logging.MustGetLogger("rest-trigger")

// RestTrigger REST trigger struct
type RestTrigger struct {
	metadata       *trigger.Metadata
	processStarter starter.ProcessStarter
	server         *Server
}

func init() {
	md := trigger.NewMetadata(jsonMetadata)
	trigger.Register(&RestTrigger{metadata: md})
}

// Metadata implements trigger.Trigger.Metadata
func (t *RestTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Init implements ext.Trigger.Init
func (t *RestTrigger) Init(processStarter starter.ProcessStarter, config map[string]string) {

	router := httprouter.New()
	router.OPTIONS("/process/start", handleOption)
	router.POST("/process/start", t.StartProcess)

	router.OPTIONS("/process/restart", handleOption)
	router.POST("/process/restart", t.RestartProcess)

	router.OPTIONS("/process/resume", handleOption)
	router.POST("/process/resume", t.ResumeProcess)

	addr := ":" + config["port"]
	t.server = NewServer(addr, router)

	t.processStarter = processStarter
}

// Start implements trigger.Trigger.Start
func (t *RestTrigger) Start() {
	t.server.Start()
}

// Stop implements trigger.Trigger.Start
func (t *RestTrigger) Stop() {
	t.server.Stop()
}

func handleOption(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Add("Access-Control-Allow-Headers", "Origin")
	w.Header().Add("Access-Control-Allow-Headers", "X-Requested-With")
	w.Header().Add("Access-Control-Allow-Headers", "Accept")
	w.Header().Add("Access-Control-Allow-Headers", "Accept-Language")
	w.Header().Set("Content-Type", "application/json")
}

// IDResponse id response object
type IDResponse struct {
	ID string `json:"id"`
}

// StartProcess starts a new Process Instance (POST "/process/start").
//
// To post a start process, try this at a shell:
// $ curl -H "Content-Type: application/json" -X POST -d '{"processUri":"base"}' http://localhost:8080/process/start
func (t *RestTrigger) StartProcess(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	w.Header().Add("Access-Control-Allow-Origin", "*")

	req := &starter.StartRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := t.processStarter.StartProcess(req)

	// If we didn't find it, 404
	//w.WriteHeader(http.StatusNotFound)

	resp := &IDResponse{ID: id}

	encoder := json.NewEncoder(w)
	encoder.Encode(resp)

	w.WriteHeader(http.StatusOK)
}

// RestartProcess restarts a Process Instance (POST "/process/restart").
//
// To post a restart process, try this at a shell:
// $ curl -H "Content-Type: application/json" -X POST -d '{...}' http://localhost:8080/process/restart
func (t *RestTrigger) RestartProcess(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	w.Header().Add("Access-Control-Allow-Origin", "*")

	//defer func() {
	//	if r := recover(); r != nil {
	//		log.Error("Unable to restart process, make sure definition registered")
	//	}
	//}()

	req := &starter.RestartRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := t.processStarter.RestartProcess(req)

	// If we didn't find it, 404
	//w.WriteHeader(http.StatusNotFound)

	resp := &IDResponse{ID: id}

	encoder := json.NewEncoder(w)
	encoder.Encode(resp)

	w.WriteHeader(http.StatusOK)
}

// ResumeProcess resumes a Process Instance (POST "/process/resume").
//
// To post a resume process, try this at a shell:
// $ curl -H "Content-Type: application/json" -X POST -d '{...}' http://localhost:8080/process/resume
func (t *RestTrigger) ResumeProcess(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	w.Header().Add("Access-Control-Allow-Origin", "*")

	defer func() {
		if r := recover(); r != nil {
			log.Error("Unable to resume process, make sure definition registered")
		}
	}()

	req := &starter.ResumeRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	t.processStarter.ResumeProcess(req)

	w.WriteHeader(http.StatusOK)
}
