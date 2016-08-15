package coap

import (
	"bytes"
	"fmt"
	"net"
	"net/url"
	"strings"
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/ext/trigger"
	"github.com/TIBCOSoftware/flogo-lib/core/flowinst"
	"github.com/TIBCOSoftware/flogo-lib/util"
	"github.com/dustin/go-coap"
	"github.com/op/go-logging"
	"encoding/json"
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
	metadata    *trigger.Metadata
	flowStarter flowinst.Starter
	resources   map[string]*CoapResource
	server      *Server
}

type CoapResource struct {
	path      string
	attrs     map[string]string
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

	t.flowStarter = flowStarter
	mux := coap.NewServeMux()
	mux.Handle("/.well-known/core", coap.FuncHandler(t.handleDiscovery))

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
				resource = &CoapResource{path: path, attrs: make(map[string]string), endpoints: make(map[string]*EndpointCfg)}
				t.resources[path] = resource
			}

			resource.endpoints[method] = &EndpointCfg{flowURI: endpoint.FlowURI, autoIdReply: autoIdReply}

			mux.Handle(path, newStartFlowHandler(t, resource))

		} else {
			panic(fmt.Sprintf("Invalid endpoint: %v", endpoint))
		}
	}

	log.Debugf("CoAP Trigger: Configured on port %s", config.Settings["port"])

	t.server = NewServer("udp", ":5683", mux)
}

// Start implements trigger.Trigger.Start
func (t *CoapTrigger) Start() error {

	return t.server.Start()
}

// Stop implements trigger.Trigger.Start
func (t *CoapTrigger) Stop() {
	t.server.Stop()
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

	return coap.FuncHandler(func(conn *net.UDPConn, addr *net.UDPAddr, msg *coap.Message) *coap.Message {

		log.Debugf("CoAP Trigger: Recieved request")

		method := toMethod(msg.Code)
		uriQuery := msg.Option(coap.URIQuery)
		var data map[string]interface{}

		if uriQuery != nil {
			//todo handle error
			queryValues, _ := url.ParseQuery(uriQuery.(string))

			queryParams := make(map[string]string, len(queryValues))

			for key, value := range queryValues {
				queryParams[key] = strings.Join(value, ",")
			}

			data = map[string]interface{}{
				"queryParams": queryParams,
				"payload":     string(msg.Payload),
			}
		} else {
			data = map[string]interface{}{
				"payload": string(msg.Payload),
			}
		}

		endpointCfg, exists := resource.endpoints[method]

		if !exists {
			res := &coap.Message{
				Type:      coap.Reset,
				Code:      coap.MethodNotAllowed,
				MessageID: msg.MessageID,
				Token:     msg.Token,
			}

			return res
		}

		//todo handle error
		startAttrs, _ := rt.metadata.OutputsToAttrs(data, false)

		rh := &AsyncReplyHandler{addr: addr.String(), msg: msg}
		rh.addr2 = addr
		rh.conn = conn

		//todo: implement reply handler?
		id, err := rt.flowStarter.StartFlowInstance(endpointCfg.flowURI, startAttrs, rh, nil)

		if err != nil {
			//todo determing if 404 or 500
			res := &coap.Message{
				Type:      coap.Reset,
				Code:      coap.NotFound,
				MessageID: msg.MessageID,
				Token:     msg.Token,
				Payload:   []byte(fmt.Sprintf("Flow '%s' not found", endpointCfg.flowURI)),
			}

			return res
		}

		log.Debugf("Started Flow: %s", id)

		//var payload []byte

		//if endpointCfg.autoIdReply {
		//	payload = []byte(id)
		//}

		if msg.IsConfirmable() {
			res := &coap.Message{
				Type:      coap.Acknowledgement,
				Code:      0,
				MessageID: msg.MessageID,
				//Token:     msg.Token,
				//Payload:   payload,
			}
			//res.SetOption(coap.ContentFormat, coap.TextPlain)

			log.Debugf("Transmitting %#v", res)
			return res
		}

		return nil
	})
}

type SyncReplyHandler struct {
	wg  sync.WaitGroup
	msg *coap.Message
}

func (rh *SyncReplyHandler) Reply(replyData map[string]string) {
	defer rh.wg.Done()

}

func (rh *SyncReplyHandler) GetReply() *coap.Message {
	rh.wg.Wait()
	return nil
}

type AsyncReplyHandler struct {
	msg  *coap.Message
	addr string

	conn  *net.UDPConn
	addr2 *net.UDPAddr
}

func (rh *AsyncReplyHandler) Reply(replyCode int, replyData interface{}) {

	//c, err := coap.Dial("udp", rh.addr)
	//if err != nil {
	//	log.Errorf("Error dialing: %v", err)
	//	return
	//}

	payload, err := json.Marshal(replyData)

	if err != nil {
		log.Errorf("Unable to marshall replyData: %v", err)
		return
	}

	res := coap.Message{
		Type:      coap.Confirmable,
		Code:      coap.Content,
		MessageID: rh.msg.MessageID,
		Token:     rh.msg.Token,
		Payload:   payload,
	}
	res.SetOption(coap.ContentFormat, coap.TextPlain)

	util.HandlePanic("CoAP Reply", nil)
	err = coap.Transmit(rh.conn, rh.addr2, res)

	//_, err = c.Send(res)
	if err != nil {
		log.Errorf("Error sending response: %v", err)
	}
}

func (rh *AsyncReplyHandler) Release() {

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
