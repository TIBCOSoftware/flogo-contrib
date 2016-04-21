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

			log.Debugf("REST Trigger: Registering endpoint [%s: %s] for Flow: %s", method, path, endpoint.FlowURI)

			router.OPTIONS(path, handleOption) // for CORS
			router.Handle(method, path, newStartFlowHandler(t, endpoint.FlowURI))

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

func newStartFlowHandler(rt *RestTrigger, flowURI string) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		log.Debugf("REST Trigger: Recieved request")

		w.Header().Add("Access-Control-Allow-Origin", "*")

		params := make(map[string]string)
		for _, param := range ps {
			params[param.Key] = param.Value
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

		data := map[string]interface{}{
			"params":  params,
			"content": content,
		}

		//todo: implement reply handler?
		id := rt.flowStarter.StartFlowInstance(flowURI, data, nil, nil)

		// paramsJSON, _ := json.Marshal(params)
		// contentJSON, _ := json.Marshal(content)
		//
		// dataJSON := map[string]string{
		// 	"params":  string(paramsJSON),
		// 	"content": string(contentJSON),
		// }
		//
		// if log.IsEnabledFor(logging.DEBUG) {
		// 	log.Debugf("REST Trigger: Starting Flow [%s] - data: %v", flowURI, dataJSON)
		// }
		//
		// id := rt.flowStarter.StartFlowInstance(flowURI, dataJSON, nil, nil)

		// If we didn't find it, 404
		//w.WriteHeader(http.StatusNotFound)

		resp := &IDResponse{ID: id}

		encoder := json.NewEncoder(w)
		encoder.Encode(resp)

		w.WriteHeader(http.StatusOK)
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
