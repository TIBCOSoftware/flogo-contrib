package support

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/op/go-logging"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/flow/flowdef"
	"github.com/TIBCOSoftware/flogo-lib/util"
)

const (
	uriSchemeFile = "file://"
)

var log = logging.MustGetLogger("manager")

// FlowManager is a simple manager for flows
type FlowManager struct {
	mu    *sync.Mutex // protects the flow maps
	flows map[string]*FlowEntry
}

// FlowEntry will contain either a compressed flow, an uncompressed flow or a flow uri
type FlowEntry struct {
	compressed   string
	uncompressed []byte
	uri          string
}

// NewFlowManager creates a new FlowManager
func NewFlowManager() *FlowManager {
	return &FlowManager{}
}

// AddCompressed adds a compressed flow to the map of flow entries
func (mgr *FlowManager) AddCompressed(id string, newFlow string) error {
	if len(newFlow) == 0 {
		return fmt.Errorf("Empty Flow with id '%s'", id)
	}
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	_, ok := mgr.flows[id]
	if ok {
		return fmt.Errorf("Flow with id '%s' already exists", id)
	}
	// Add the flow
	mgr.flows[id] = &FlowEntry{compressed: newFlow}
	return nil
}

// AddUncompressed adds an uncompressed flow to the map of flow entries
func (mgr *FlowManager) AddUncompressed(id string, newFlow []byte) error {
	if len(newFlow) == 0 {
		return fmt.Errorf("Empty Flow with id '%s'", id)
	}
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	_, ok := mgr.flows[id]
	if ok {
		return fmt.Errorf("Flow with id '%s' already exists", id)
	}
	// Add the flow
	mgr.flows[id] = &FlowEntry{uncompressed: newFlow}
	return nil
}

// AddURI adds a uri flow to the map of flow entries
func (mgr *FlowManager) AddURI(id string, newUri string) error {
	if len(newUri) == 0 {
		return fmt.Errorf("Empty Flow URI with id '%s'", id)
	}
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	_, ok := mgr.flows[id]
	if ok {
		return fmt.Errorf("Flow with id '%s' already exists", id)
	}
	// Add the flow
	mgr.flows[id] = &FlowEntry{uri: newUri}
	return nil
}

// GetFlow gets the specified embedded flow
func (mgr *FlowManager) GetFlow(id string) (*flowdef.DefinitionRep, error) {

	entry, ok := mgr.flows[id]

	if !ok {
		err := fmt.Errorf("No flow found for id '%s'", id)
		log.Errorf(err.Error())
		return nil, err
	}

	var flowDefBytes []byte

	// Uncompressed Flow condition
	if len(entry.uncompressed) > 0 {
		// Uncompressed flow
		flowDefBytes = entry.uncompressed
	}

	// Compressed Flow condition
	if len(entry.compressed) > 0 {

		decodedBytes, err := decodeAndUnzip(entry.compressed)
		if err != nil {
			decodeErr := fmt.Errorf("Error decoding compressed flow with id '%s', %s", id, err.Error())
			log.Errorf(decodeErr.Error())
			return nil, decodeErr
		}
		flowDefBytes = decodedBytes
	}

	// URI Flow condition
	if len(entry.uri) > 0 {
		if strings.HasPrefix(entry.uri, uriSchemeFile) {
			// File URI
			log.Infof("Loading Local Flow: %s\n", entry.uri)
			flowFilePath, _ := util.URLStringToFilePath(entry.uri)

			readBytes, err := ioutil.ReadFile(flowFilePath)
			if err != nil {
				readErr := fmt.Errorf("Error reading flow file with id '%s', %s", id, err.Error())
				log.Errorf(readErr.Error())
				return nil, readErr
			}
			flowDefBytes = readBytes
		} else {
			// URI
			req, err := http.NewRequest("GET", entry.uri, nil)
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				getErr := fmt.Errorf("Error getting flow uri with id '%s', %s", id, err.Error())
				log.Errorf(getErr.Error())
				return nil, getErr
			}
			defer resp.Body.Close()

			log.Infof("response Status:", resp.Status)

			if resp.StatusCode >= 300 {
				//not found
				getErr := fmt.Errorf("Error getting flow uri with id '%s', status code %d", id, resp.StatusCode)
				log.Errorf(getErr.Error())
				return nil, getErr
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				readErr := fmt.Errorf("Error reading flow uri response body with id '%s', %s", id, err.Error())
				log.Errorf(readErr.Error())
				return nil, readErr
			}
			flowDefBytes = body
		}
	}

	var flow *flowdef.DefinitionRep
	err := json.Unmarshal(flowDefBytes, &flow)
	if err != nil {
		log.Errorf(err.Error())
		return nil, fmt.Errorf("Error marshalling flow with id '%s', %s", id, err.Error())
	}
	return flow, nil
}

func decodeAndUnzip(encoded string) ([]byte, error) {

	decoded, _ := base64.StdEncoding.DecodeString(encoded)
	buf := bytes.NewBuffer(decoded)
	r, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}
	jsonAsBytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return jsonAsBytes, nil
}
