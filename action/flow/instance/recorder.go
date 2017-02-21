package instance

// StateRecorder is the interface that describes a service that can record
// snapshots and steps of a Flow Instance
type StateRecorder interface {

	// RecordSnapshot records a Snapshot of the FlowInstance
	RecordSnapshot(instance *Instance)

	// RecordStep records the changes for the current Step of the Flow Instance
	RecordStep(instance *Instance)
}
