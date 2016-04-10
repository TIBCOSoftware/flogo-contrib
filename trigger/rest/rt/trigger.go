package rest

import (
	"encoding/json"
	"net/http"

	"github.com/TIBCOSoftware/flogo-lib/core/ext/trigger"
	"github.com/TIBCOSoftware/flogo-lib/core/process"
	"github.com/TIBCOSoftware/flogo-lib/core/processinst"
	"github.com/julienschmidt/httprouter"
	"github.com/op/go-logging"
)

// log is the default package logger
var log = logging.MustGetLogger("rest-trigger")

// todo: switch to use endpoint registration

// RestTrigger REST trigger struct
type RestTrigger struct {
	metadata       *trigger.Metadata
	processStarter processinst.Starter
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
func (t *RestTrigger) Init(processStarter processinst.Starter, config *trigger.Config) {

	router := httprouter.New()
	router.OPTIONS("/process/start", handleOption)
	router.POST("/process/start", t.StartProcess)

	addr := ":" + config.Settings["port"]
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

	req := &StartRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := t.processStarter.StartProcessInstance(req.ProcessURI, req.Data, nil, nil)

	// If we didn't find it, 404
	//w.WriteHeader(http.StatusNotFound)

	resp := &IDResponse{ID: id}

	encoder := json.NewEncoder(w)
	encoder.Encode(resp)

	w.WriteHeader(http.StatusOK)
}

// StartRequest describes a request for starting a ProcessInstance
type StartRequest struct {
	ProcessURI  string               `json:"processUri"`
	Data        map[string]string    `json:"data"`
	Interceptor *process.Interceptor `json:"interceptor"`
	Patch       *process.Patch       `json:"patch"`
}
