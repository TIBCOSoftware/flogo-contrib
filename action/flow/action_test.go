package flow

import (
	"github.com/TIBCOSoftware/flogo-lib/types"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

//TestInitNoFlavorError
func TestInitNoFlavorError(t *testing.T) {
	flowAction := NewFlowAction()
	assert.NotNil(t, flowAction)

	mockConfig := &types.ActionConfig{Id: "myMockConfig", Ref: "github.com/my/mock/ref"}

	// Recover from expected panic
	defer func() {
		if r := recover(); r != nil {
			// Expected error
			if !strings.HasPrefix(r.(string), "Error while loading flow") {
				t.Fail()
			}
		}
	}()

	flowAction.Init(*mockConfig)

	// If reaches here it should fail, as it should panic before
	t.Fail()
}

//TestInitUnCompressedFlowFlavorError
func TestInitUnCompressedFlowFlavorError(t *testing.T) {
	flowAction := NewFlowAction()
	assert.NotNil(t, flowAction)

	mockFlowData := []byte(`{"flow":{}}`)

	mockConfig := &types.ActionConfig{Id: "myMockConfig", Ref: "github.com/my/mock/ref", Data: mockFlowData}

	// Recover from expected panic
	defer func() {
		if r := recover(); r != nil {
			// Expected error
			if !strings.HasPrefix(r.(string), "Error while loading uncompressed flow") {
				t.Fail()
			}
		}
	}()

	flowAction.Init(*mockConfig)

	// If reaches here it should fail, as it should panic before
	t.Fail()
}

//TestInitCompressedFlowFlavorError
func TestInitCompressedFlowFlavorError(t *testing.T) {
	flowAction := NewFlowAction()
	assert.NotNil(t, flowAction)

	mockFlowData := []byte(`{"flowCompressed":""}`)

	mockConfig := &types.ActionConfig{Id: "myMockConfig", Ref: "github.com/my/mock/ref", Data: mockFlowData}

	// Recover from expected panic
	defer func() {
		if r := recover(); r != nil {
			// Expected error
			if !strings.HasPrefix(r.(string), "Error while loading compressed flow") {
				t.Fail()
			}
		}
	}()

	flowAction.Init(*mockConfig)

	// If reaches here it should fail, as it should panic before
	t.Fail()
}

//TestInitURIFlowFlavorError
func TestInitURIFlowFlavorError(t *testing.T) {
	flowAction := NewFlowAction()
	assert.NotNil(t, flowAction)

	mockFlowData := []byte(`{"flowURI":""}`)

	mockConfig := &types.ActionConfig{Id: "myMockConfig", Ref: "github.com/my/mock/ref", Data: mockFlowData}

	// Recover from expected panic
	defer func() {
		if r := recover(); r != nil {
			// Expected error
			if !strings.HasPrefix(r.(string), "Error while loading flow URI") {
				t.Fail()
			}
		}
	}()

	flowAction.Init(*mockConfig)

	// If reaches here it should fail, as it should panic before
	t.Fail()
}
