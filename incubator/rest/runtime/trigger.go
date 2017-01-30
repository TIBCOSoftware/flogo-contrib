package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/TIBCOSoftware/flogo-contrib/incubator/rest/runtime/cors"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/types"
	"github.com/julienschmidt/httprouter"
	"github.com/op/go-logging"
)

const (
	REST_CORS_PREFIX = "REST_TRIGGER"
	TRIGGER_REF      = "github.com/TIBCOSoftware/flogo-contrib/incubator/rest/runtime"
)

// log is the default package logger
var log = logging.MustGetLogger("trigger-tibco-rest")

var validMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}

var md = trigger.NewMetadata(jsonMetadata)

// RestTrigger REST trigger struct
type RestTrigger struct {
	Md     *trigger.Metadata
	runner action.Runner
	server *Server
	myId   string
}

type RestFactory struct{}

func init() {
	trigger.RegisterFactory(TRIGGER_REF, &RestFactory{})
}

//New Creates a new trigger instance for a given id
func (t *RestFactory) New(id string) trigger.Trigger2 {
	return &RestTrigger{Md: md, myId: id}
}

// Metadata implements trigger.Trigger.Metadata
func (t *RestTrigger) Metadata() *trigger.Metadata {
	return t.Md
}

func (t *RestTrigger) Init(config types.TriggerConfig, runner action.Runner) {
	log.Debugf("In init, id '%s'", t.myId)

	// Get Trigger Settings
	var triggerSettings *Settings
	err := json.Unmarshal(config.Settings, &triggerSettings)
	if err != nil {
		panic(err.Error())
	}

	t.InitHandlers(triggerSettings, config.Handlers, runner)
}

// Initialize the handlers
func (t *RestTrigger) InitHandlers(triggerSettings *Settings, handlers []*types.TriggerHandler, runner action.Runner) {

	router := httprouter.New()

	addr := ":" + triggerSettings.Port
	t.runner = runner

	for _, handler := range handlers {

		// Get Handler Settings
		var handlerSettings *HandlerSettings
		err := json.Unmarshal(handler.Settings, &handlerSettings)
		if err != nil {
			panic(err.Error())
		}

		if handlerIsValid(handlerSettings) {
			method := strings.ToUpper(handlerSettings.Method)
			path := handlerSettings.Path

			log.Infof("REST Trigger: Registering handler [%s: %s] for Action Id: [%s]", method, path, handler.ActionId)

			router.OPTIONS(path, handleCorsPreflight) // for CORS
			router.Handle(method, path, newActionHandler(t, handler.ActionId, handlerSettings))

		} else {
			panic(fmt.Sprintf("Invalid handler: %v", handlerSettings))
		}
	}

	log.Debugf("REST Trigger: Configured on port %s", triggerSettings.Port)
	t.server = NewServer(addr, router)
}

// Start implements util.Managed.Start
func (t *RestTrigger) Start2() error {
	log.Infof("Start CALLED in Trigger")
	return nil
	//	return t.server.Start()
}

func (t *RestTrigger) Start() error {
	return t.server.Start()
}

// Stop implements util.Managed.Stop
func (t *RestTrigger) Stop() error {
	return t.server.Stop()
}

// Handles the cors preflight request
func handleCorsPreflight(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	log.Infof("Received [OPTIONS] request to CorsPreFlight: %+v", r)

	c := cors.New(REST_CORS_PREFIX, log)
	c.HandlePreflight(w, r)
}

// IDResponse id response object
type IDResponse struct {
	ID string `json:"id"`
}

func newActionHandler(rt *RestTrigger, actionId string, handlerSettings *HandlerSettings) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		log.Infof("REST Trigger: Received request for id '%s'", rt.myId)

		c := cors.New(REST_CORS_PREFIX, log)
		c.WriteCorsActualRequestHeaders(w)

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
			//todo should handler say if content is expected?
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
			"params":      pathParams,
			"pathParams":  pathParams,
			"queryParams": queryParams,
			"content":     content,
		}

		//todo handle error
		startAttrs, _ := rt.Md.OutputsToAttrs(data, false)

		action := action.Get2(actionId)
		log.Debugf("Found action' %+x'", action)

		context := trigger.NewContext(context.Background(), startAttrs)
		replyCode, replyData, err := rt.runner.Run(context, action, actionId, nil)

		if err != nil {
			log.Debugf("REST Trigger Error: %s", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if replyData != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(replyCode)
			if err := json.NewEncoder(w).Encode(replyData); err != nil {
				log.Error(err)
			}
		}

		if replyCode > 0 {
			w.WriteHeader(replyCode)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}

////////////////////////////////////////////////////////////////////////////////////////
// Utils

func handlerIsValid(handlerSettings *HandlerSettings) bool {

	if !stringInList(strings.ToUpper(handlerSettings.Method), validMethods) {
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
