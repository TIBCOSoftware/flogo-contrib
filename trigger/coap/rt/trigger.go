package coap

import (
	"bytes"
	"fmt"
	"net"
	"strings"

	"github.com/dustin/go-coap"
	"github.com/TIBCOSoftware/flogo-lib/core/ext/trigger"
	"github.com/TIBCOSoftware/flogo-lib/core/flowinst"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/op/go-logging"
)

const (
	methodGET    = "GET"
	methodPOST   = "POST"
	methodPUT    = "PUT"
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
	resources   map[string]*CoapResource
	mux         *coap.ServeMux
}

type CoapResource struct {
	path string
	attrs map[string]string
	endpoints map[string]*EndpointCfg
}

type EndpointCfg struct {
	method      string
	flowURI     string
	autoIdReply bool
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

	t.mux = coap.NewServeMux()
	t.mux.Handle("/.well-known/core", coap.FuncHandler(t.handleDiscovery))

	endpoints := config.Endpoints
	t.resources = make(map[string]*CoapResource)


	for _, endpoint := range endpoints {

		if endpointIsValid(endpoint) {
			method := strings.ToUpper(endpoint.Settings["method"])
			path := endpoint.Settings["path"]
			autoIdReply, _ := data.CoerceToBoolean(endpoint.Settings["autoIdReply"])

			log.Debugf("CoAP Trigger: Registering endpoint [%s: %s] for Flow: %s", method, path, endpoint.FlowURI)
			if autoIdReply {
				log.Debug("CoAP Trigger: AutoIdReply Enabled")
			}

			resource, exists := t.resources[path]

			if !exists {
				resource = &CoapResource{path:path, attrs:make(map[string]string), endpoints:make(map[string]*EndpointCfg)}
				t.resources[path] = resource
			}

			resource.endpoints[method] = &EndpointCfg{flowURI:endpoint.FlowURI, autoIdReply:autoIdReply}

			t.mux.Handle(path, newStartFlowHandler(t, resource))

		} else {
			panic(fmt.Sprintf("Invalid endpoint: %v", endpoint))
		}
	}

	log.Debugf("CoAP Trigger: Configured on port %s", config.Settings["port"])
}

// Start implements trigger.Trigger.Start
func (t *CoapTrigger) Start() error {

	err := coap.ListenAndServe("udp", ":5683", t.mux)

	return err
}

// Stop implements trigger.Trigger.Start
func (t *CoapTrigger) Stop() {
}

// IDResponse id response object
type IDResponse struct {
	ID string `json:"id"`
}

func (t *CoapTrigger) handleDiscovery(conn *net.UDPConn, addr *net.UDPAddr, msg *coap.Message) *coap.Message {

	//path := msg.PathString() //handle queries

	//todo add filter support

	var buffer bytes.Buffer

	numResources := len(t.resources)

	i := 0
	for _, resource := range t.resources {

		i++

		buffer.WriteString("<")
		buffer.WriteString(resource.path)
		buffer.WriteString(">")

		if len(resource.attrs) > 0 {
			for k, v := range resource.attrs {
				buffer.WriteString(";")
				buffer.WriteString(k)
				buffer.WriteString("=")
				buffer.WriteString(v)
			}
		}

		if i < numResources {
			buffer.WriteString(",\n")
		} else {
			buffer.WriteString("\n")
		}
	}

	payloadStr := buffer.String()

	res := &coap.Message{
		Type:      msg.Type,
		Code:      coap.Content,
		MessageID: msg.MessageID,
		Token:     msg.Token,
		Payload:   []byte(payloadStr),
	}
	res.SetOption(coap.ContentFormat, coap.AppLinkFormat)

	log.Debugf("Transmitting %#v", res)

	return res
}


func newStartFlowHandler(rt *CoapTrigger, resource *CoapResource) coap.Handler {

	return coap.FuncHandler(func (conn *net.UDPConn, addr *net.UDPAddr, msg *coap.Message) *coap.Message {

		log.Debugf("CoAP Trigger: Recieved request")

		data := map[string]interface{}{
			"payload":  string(msg.Payload),
		}

		method := toMethod(msg.Code);

		endpointCfg, exists := resource.endpoints[method]

		if !exists {
			res := &coap.Message{
				Type:      msg.Type,
				Code:      coap.BadRequest,
				MessageID: msg.MessageID,
				Token:     msg.Token,
				Payload:   []byte("Unknown Endpoint"),
			}

			return res
		}

		//todo handle error
		startAttrs, _ := rt.metadata.OutputsToAttrs(data, false)


		//todo: implement reply handler?
		id, _ := rt.flowStarter.StartFlowInstance(endpointCfg.flowURI, startAttrs, nil, nil)

		if msg.IsConfirmable() && endpointCfg.autoIdReply {
			res := &coap.Message{
				Type:      coap.Acknowledgement,
				Code:      coap.Content,
				MessageID: msg.MessageID,
				Token:     msg.Token,
				Payload:   []byte(id),
			}
			res.SetOption(coap.ContentFormat, coap.TextPlain)

			log.Debugf("Transmitting %#v", res)
			return res
		}

		return nil
	})
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
