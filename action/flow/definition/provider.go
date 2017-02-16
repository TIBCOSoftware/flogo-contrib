package definition

import (
	"fmt"
	"sync"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/support"
	"github.com/TIBCOSoftware/flogo-lib/flow/flowdef"
	"github.com/TIBCOSoftware/flogo-lib/flow/script/fggos"
	"github.com/TIBCOSoftware/flogo-lib/flow/service"
	"github.com/TIBCOSoftware/flogo-lib/util"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("provider")

// Provider is the interface that describes an object
// that can provide flow definitions from a URI
type Provider interface {

	// GetFlow retrieves the flow definition for the specified id
	GetFlow(flowId string) (*flowdef.Definition, error)
	// AddCompressedFlow adds the flow for a specified id
	AddCompressedFlow(id string, flow string) error
	// AddUnCompressedFlow adds the flow for a specified id
	AddUncompressedFlow(id string, flow []byte) error
	// AddFlowURI adds the flow for a specified uri
	AddFlowURI(id string, uri string) error
}

// RemoteFlowProvider is an implementation of FlowProvider service
// that can access flowes via URI
type RemoteFlowProvider struct {
	//todo: switch to LRU cache
	mutex     *sync.Mutex
	flowCache map[string]*flowdef.Definition
	flowMgr   *support.FlowManager
}

// NewRemoteFlowProvider creates a RemoteFlowProvider
func NewRemoteFlowProvider() *RemoteFlowProvider {
	var service RemoteFlowProvider
	service.flowCache = make(map[string]*flowdef.Definition)
	service.mutex = &sync.Mutex{}
	service.flowMgr = support.NewFlowManager()
	return &service
}

func (pps *RemoteFlowProvider) Name() string {
	return service.ServiceFlowProvider
}

// Start implements util.Managed.Start()
func (pps *RemoteFlowProvider) Start() error {
	// no-op
	return nil
}

// Stop implements util.Managed.Stop()
func (pps *RemoteFlowProvider) Stop() error {
	// no-op
	return nil
}

// GetFlow implements flow.Provider.GetFlow
func (pps *RemoteFlowProvider) GetFlow(id string) (*flowdef.Definition, error) {

	// todo turn pps.flowCache to real cache
	if flow, ok := pps.flowCache[id]; ok {
		log.Debugf("Accessing cached Flow: %s\n")
		return flow, nil
	}

	log.Debugf("Getting Flow: %s\n", id)

	flowRep, err := pps.flowMgr.GetFlow(id)
	if err != nil {
		return nil, err
	}

	def, err := flowdef.NewDefinition(flowRep)
	if err != nil {
		errorMsg := fmt.Sprintf("Error unmarshalling flow '%s': %s", id, err.Error())
		log.Errorf(errorMsg)
		return nil, fmt.Errorf(errorMsg)
	}

	//todo optimize this - not needed if flow doesn't have expressions
	//todo have a registry for this?
	def.SetLinkExprManager(fggos.NewGosLinkExprManager(def))
	//def.SetLinkExprManager(fglua.NewLuaLinkExprManager(def))

	//synchronize
	pps.mutex.Lock()
	pps.flowCache[id] = def
	pps.mutex.Unlock()

	return def, nil

}

func (pps *RemoteFlowProvider) AddCompressedFlow(id string, flow string) error {
	return pps.flowMgr.AddCompressed(id, flow)
}

func (pps *RemoteFlowProvider) AddUncompressedFlow(id string, flow []byte) error {
	return pps.flowMgr.AddUncompressed(id, flow)
}

func (pps *RemoteFlowProvider) AddFlowURI(id string, uri string) error {
	return pps.flowMgr.AddURI(id, uri)
}

func DefaultConfig() *util.ServiceConfig {
	return &util.ServiceConfig{Name: service.ServiceFlowProvider, Enabled: true}
}
