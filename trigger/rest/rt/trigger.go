package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/ext/trigger"
	"github.com/TIBCOSoftware/flogo-lib/core/processinst"
	"github.com/julienschmidt/httprouter"
	"github.com/op/go-logging"
)

// log is the default package logger
var log = logging.MustGetLogger("rest-trigger")

var validMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE"}

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

	addr := ":" + config.Settings["port"]
	t.processStarter = processStarter

	endpoints := config.Endpoints

	for _, endpoint := range endpoints {

		if endpointIsValid(endpoint) {
			method := strings.ToUpper(endpoint.Settings["method"])
			path := endpoint.Settings["path"]

			log.Debugf("REST Trigger: Registering endpoint [%s: %s] for Process: %s", method, path, endpoint.ProcessURI)

			router.OPTIONS(path, handleOption) // for CORS
			router.Handle(method, path, newStartProcessHandler(t, endpoint.ProcessURI))

		} else {
			panic(fmt.Sprintf("Invalid endpoint: %v", endpoint))
		}
	}

	log.Debugf("REST Trigger: Configured on port %s", config.Settings["port"])
	t.server = NewServer(addr, router)
}

// Start implements trigger.Trigger.Start
func (t *RestTrigger) Start() {
	err := t.server.Start()

	if err != nil {
		log.Errorf("REST Trigger: Error starting - %v", err)
	}

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

func newStartProcessHandler(rt *RestTrigger, processURI string) httprouter.Handle {

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

		//data := map[string]interface{}{
		//	"params": params,
		//	"content": content,
		//}
		//
		////todo: implement reply handler?
		//id := rt.processStarter.StartProcessInstance(processURI, data, nil, nil)

		//todo: fix StartProcessInstance to use map[string]interface{} and remove this
		paramsJSON, _ := json.Marshal(params)
		contentJSON, _ := json.Marshal(content)

		dataJSON := map[string]string{
			"params":  string(paramsJSON),
			"content": string(contentJSON),
		}

		if log.IsEnabledFor(logging.DEBUG) {
			log.Debugf("REST Trigger: Starting Process [%s] - data: %v", processURI, dataJSON)
		}

		id := rt.processStarter.StartProcessInstance(processURI, dataJSON, nil, nil)

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
