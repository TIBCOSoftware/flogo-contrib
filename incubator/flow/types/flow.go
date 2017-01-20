package types

import (
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

type FlowConfig struct {
	Attrs    []*data.Attribute `json:"attributes"`
	RootTask *Task             `json:"rootTask"`
}

// Task is the object that describes the definition of
// a task.  It contains its data (attributes) and its
// nested structure (child tasks & child links).
type Task struct {
}
