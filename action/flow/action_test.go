package flow

import (
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

//TestInitNoFlavorError
func TestInitNoFlavorError(t *testing.T) {
	// Recover from expected panic
	defer func() {
		if r := recover(); r != nil {
			// Expected error
			if !strings.HasPrefix(r.(string), "No flow found in action data") {
				t.Fail()
			}
		}
	}()

	mockConfig := &action.Config{Id: "myMockConfig", Ref: "github.com/my/mock/ref", Data: []byte(`{}`)}
	f := &FlowFactory{}
	flowAction := f.New(mockConfig)
	assert.NotNil(t, flowAction)

	// If reaches here it should fail, as it should panic before
	t.Fail()
}

//TestInitUnCompressedFlowFlavorError
func TestInitUnCompressedFlowFlavorError(t *testing.T) {
	// Recover from expected panic
	defer func() {
		if r := recover(); r != nil {
			// Expected error
			if !strings.HasPrefix(r.(string), "Error while loading uncompressed flow") {
				t.Fail()
			}
		}
	}()

	mockFlowData := []byte(`{"flow":{}}`)

	mockConfig := &action.Config{Id: "myMockConfig", Ref: "github.com/my/mock/ref", Data: mockFlowData}

	f := &FlowFactory{}
	flowAction := f.New(mockConfig)
	assert.NotNil(t, flowAction)

	// If reaches here it should fail, as it should panic before
	t.Fail()
}

//TestInitCompressedFlowFlavorError
func TestInitCompressedFlowFlavorError(t *testing.T) {
	// Recover from expected panic
	defer func() {
		if r := recover(); r != nil {
			// Expected error
			if !strings.HasPrefix(r.(string), "Error while loading compressed flow") {
				t.Fail()
			}
		}
	}()
	mockFlowData := []byte(`{"flowCompressed":""}`)

	mockConfig := &action.Config{Id: "myMockConfig", Ref: "github.com/my/mock/ref", Data: mockFlowData}

	f := &FlowFactory{}
	flowAction := f.New(mockConfig)
	assert.NotNil(t, flowAction)

	// If reaches here it should fail, as it should panic before
	t.Fail()
}

//TestInitURIFlowFlavorError
func TestInitURIFlowFlavorError(t *testing.T) {
	// Recover from expected panic
	defer func() {
		if r := recover(); r != nil {
			// Expected error
			if !strings.HasPrefix(r.(string), "Error while loading flow URI") {
				t.Fail()
			}
		}
	}()
	mockFlowData := []byte(`{"flowURI":""}`)

	mockConfig := &action.Config{Id: "myMockConfig", Ref: "github.com/my/mock/ref", Data: mockFlowData}

	f := &FlowFactory{}
	flowAction := f.New(mockConfig)
	assert.NotNil(t, flowAction)

	// If reaches here it should fail, as it should panic before
	t.Fail()
}
