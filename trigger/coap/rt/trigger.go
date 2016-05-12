package coap

import (
	"fmt"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/ext/trigger"
	"github.com/TIBCOSoftware/flogo-lib/core/flowinst"
	"github.com/op/go-logging"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"net"
	"github.com/dustin/go-coap"
	"strconv"
)

const (
	methodGET = "GET"
	methodPOST = "POST"
	methodPUT = "PUT"
	methodDELETE = "DELETE"
)

// log is the default package logger
var log = logging.MustGetLogger("trigger-tibco-coap")

var validMethods = []string{methodGET, methodPOST, methodPUT, methodDELETE}

type StartFunc func(payload string) (string, bool)

// CoapTrigger CoAP trigger struct
type CoapTrigger struct {
	addr        string
	metadata    *trigger.Metadata
	flowStarter flowinst.Starter
	starters    map[string]StartFunc
}

func init() {
	md := trigger.NewMetadata(jsonMetadata)
	trigger.Register(&CoapTrigger{metadata: md})
}

// Metadata implements trigger.Trigger.Metadata
func (t *CoapTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Init implements ext.Trigger.Init
func (t *CoapTrigger) Init(flowStarter flowinst.Starter, config *trigger.Config) {

	t.starters = make(map[string]StartFunc)
	t.addr = ":" + config.Settings["port"]

	endpoints := config.Endpoints

	for _, endpoint := range endpoints {

		if endpointIsValid(endpoint) {
			method := strings.ToUpper(endpoint.Settings["method"])
			path := endpoint.Settings["path"]
			autoIdReply, _ := data.CoerceToBoolean(endpoint.Settings["autoIdReply"])

			log.Debugf("CoAP Trigger: Registering endpoint [%s: %s] for Flow: %s", method, path, endpoint.FlowURI)
			if autoIdReply {
				log.Debug("CoAP Trigger: AutoIdReply Enabled")
			}

			t.starters[method + ":" + path] = newStartFlowHandler(t, endpoint.FlowURI, autoIdReply)

		} else {
			panic(fmt.Sprintf("Invalid endpoint: %v", endpoint))
		}
	}

	log.Debugf("CoAP Trigger: Configured on port %s", config.Settings["port"])
}

// Start implements trigger.Trigger.Start
func (t *CoapTrigger) Start() error {

	err := (coap.ListenAndServe("udp", t.addr,
		coap.FuncHandler(func(l *net.UDPConn, a *net.UDPAddr, m *coap.Message) *coap.Message {
			log.Debugf("Got message path=%q: %#v from %v", m.Path(), m, a)

			starter := t.starters[toMethod(m.Code) + ":" + m.PathString()]

			if starter != nil {

				id, reply := starter()

				if m.IsConfirmable() && reply {
					res := &coap.Message{
						Type:      coap.Acknowledgement,
						Code:      coap.Content,
						MessageID: m.MessageID,
						Token:     m.Token,
						Payload:   []byte(id),
					}
					res.SetOption(coap.ContentFormat, coap.TextPlain)

					log.Debugf("Transmitting %#v", res)
					return res
				}
			}

			return nil
		})))

	return err
}

// Stop implements trigger.Trigger.Start
func (t *CoapTrigger) Stop() {
}

// IDResponse id response object
type IDResponse struct {
	ID string `json:"id"`
}

func newStartFlowHandler(rt *CoapTrigger, flowURI string, autoIdReply bool) StartFunc {

	return func(payload string) (string, bool) {

		log.Debugf("CoAP Trigger: Recieved request")

		data := map[string]interface{}{
			"payload":  payload,
		}

		//todo handle error
		startAttrs, _ := rt.metadata.OutputsToAttrs(data, false)

		//todo: implement reply handler?
		id, _ := rt.flowStarter.StartFlowInstance(flowURI, startAttrs, nil, nil)

		if autoIdReply {
			return id, true
		}

		return "", false
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

func toMethod(code coap.COAPCode) string {

	var method string

	switch code {
	case coap.GET:
		method = methodGET
	case coap.POST:
		method = methodPOST
	case coap.PUT:
		method = methodPUT
	case coap.DELETE:
		method = methodDELETE
	}

	return method
}
