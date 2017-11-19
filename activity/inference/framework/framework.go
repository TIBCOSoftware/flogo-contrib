package framework

import (
	"sync"

	"github.com/TIBCOSoftware/flogo-contrib/activity/inference/model"
)

var (
	frameworksMu sync.Mutex
	frameworks   = make(map[string]model.Framework)
)

func init() {
	frameworks = make(map[string]model.Framework)
}

// Register is used by the init() func from each framework to register
func Register(framework model.Framework) {
	frameworksMu.Lock()
	defer frameworksMu.Unlock()

	if framework == nil {
		panic("framework.Register: framework is nil")
	}

	id := framework.FrameworkTyp()

	if _, dup := frameworks[id]; dup {
		panic("framework.Register: framework already registered " + id)
	}

	// copy on write to avoid synchronization on access
	newFrameworks := make(map[string]model.Framework, len(frameworks))

	for k, v := range frameworks {
		newFrameworks[k] = v
	}

	newFrameworks[id] = framework
	frameworks = newFrameworks
}

// Get is used to fecth the specified framework
func Get(id string) model.Framework {
	return frameworks[id]
}
