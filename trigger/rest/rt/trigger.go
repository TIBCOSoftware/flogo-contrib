package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/ext/trigger"
	"github.com/TIBCOSoftware/flogo-lib/core/flowinst"
	"github.com/julienschmidt/httprouter"
	"github.com/op/go-logging"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

// log is the default package logger
var log = logging.MustGetLogger("trigger-tibco-rest")

var validMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE"}

// RestTrigger REST trigger struct
type RestTrigger struct {
	metadata    *trigger.Metadata
	flowStarter flowinst.Starter
	server      *Server
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
func (t *RestTrigger) Init(flowStarter flowinst.Starter, config *trigger.Config) {

	router := httprouter.New()

	addr := ":" + config.Settings["port"]
	t.flowStarter = flowStarter

	endpoints := config.Endpoints

	for _, endpoint := range endpoints {

		if endpointIsValid(endpoint) {
			method := strings.ToUpper(endpoint.Settings["method"])
			path := endpoint.Settings["path"]
			autoIdReply, _ := data.CoerceToBoolean(endpoint.Settings["autoIdReply"])

			log.Debugf("REST Trigger: Registering endpoint [%s: %s] for Flow: %s", method, path, endpoint.FlowURI)
			if autoIdReply {
				log.Debug("REST Trigger: AutoIdReply Enabled")
			}

			router.OPTIONS(path, handleOption) // for CORS
			router.Handle(method, path, newStartFlowHandler(t, endpoint.FlowURI, autoIdReply))

		} else {
			panic(fmt.Sprintf("Invalid endpoint: %v", endpoint))
		}
	}

	log.Debugf("REST Trigger: Configured on port %s", config.Settings["port"])
	t.server = NewServer(addr, router)
}

// Start implements trigger.Trigger.Start
func (t *RestTrigger) Start() error {
	err := t.server.Start()
	return err
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

func newStartFlowHandler(rt *RestTrigger, flowURI string, autoIdReply bool) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		log.Debugf("REST Trigger: Recieved request")

		w.Header().Add("Access-Control-Allow-Origin", "*")

		pathParams := make(map[string]string)
		for _, param := range ps {
			pathParams[param.Key] = param.Value
		}

		var content interface{}
		err := json.NewDecoder(r.Body).Decode(&content)
		if err != nil {
			switch {
			case err == io.EOF:
				// empty body
				//todo should endpoint say if content is expected?
			case err != nil:
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		queryValues := r.URL.Query()
		queryParams := make(map[string]string, len(queryValues))

		for key, value := range queryValues {
			queryParams[key] = strings.Join(value, ",")
		}

		data := map[string]interface{}{
			"params":  pathParams,
			"pathParams" : pathParams,
			"queryParams" : queryParams,
			"content": content,
		}

		//todo handle error
		startAttrs,_ := rt.metadata.OutputsToAttrs(data, false)

		//todo: implement reply handler?
		id, err := rt.flowStarter.StartFlowInstance(flowURI, startAttrs, nil, nil)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if autoIdReply {
			resp := &IDResponse{ID: id}

			encoder := json.NewEncoder(w)
			encoder.Encode(resp)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}

////////////////////////////////////////////////////////////////////////////////////////
// Utils

func endpointIsValid(endpoint *trigger.EndpointConfig) bool {

	if !stringInList(strings.ToUpper(endpoint.Settings["method"]), validMethods) {
		return false
	}

	//validate path

	return true
}

func stringInList(str string, list []string) bool {
	for _, value := range list {
		if value == str {
			return true
		}
	}
	return false
}
