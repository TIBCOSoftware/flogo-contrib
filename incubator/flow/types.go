package flow

import (
	"encoding/json"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

type Flavor struct {
	// The flow is embedded and uncompressed
	Flow json.RawMessage `json:"flow"`
	// The flow is a URI
	FlowCompressed json.RawMessage `json:"flowCompressed"`
	// The flow is a URI
	FlowURI json.RawMessage `json:"flowURI"`
}

type FlowConfig struct {
	Attrs    []*data.Attribute `json:"attributes"`
	RootTask *Task             `json:"rootTask"`
}

// Task is the object that describes the definition of
// a task.  It contains its data (attributes) and its
// nested structure (child tasks & child links).
type Task struct {
}
