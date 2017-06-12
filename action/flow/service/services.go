package service


const (
	// ServiceStateRecorder is the name of the StateRecorder service used in configuration
	ServiceStateRecorder string = "stateRecorder"

	// ServiceFlowProvider is the name of the FlowProvider service used in configuration
	ServiceFlowProvider string = "flowProvider"

	// ServiceEngineTester is the name of the EngineTester service used in configuration
	ServiceEngineTester string = "engineTester"
)

/*
// StateRecorderService is the flowinst.StateRecorder wrapped as a service
type StateRecorderService interface {
	util.Service
	flowinst.StateRecorder
}

// FlowProviderService is the flow.Provider wrapped as a service
type FlowProviderService interface {
	util.Service
	flowdef.Provider
}

// EngineTesterService is an engine service to assist in testing flows
type EngineTesterService interface {
	util.Service
}
*/