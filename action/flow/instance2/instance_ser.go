package instance2

import (
	"encoding/json"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/model"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/util"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////
// Flow Instance Serialization

type serIndependentInstance struct {
	ID        string              `json:"id"`
	Status    model.FlowStatus    `json:"status"`
	State     int                 `json:"state"`
	FlowURI   string              `json:"flowUri"`
	Attrs     []*data.Attribute   `json:"attrs"`
	WorkQueue []*WorkItem         `json:"workQueue"`
	TaskInsts []*TaskInst         `json:"tasks"`
	LinkInsts []*LinkInst         `json:"links"`
	SubFlows  []*EmbeddedInstance `json:"subFlows,omitempty"`
}

// MarshalJSON overrides the default MarshalJSON for FlowInstance
func (inst *IndependentInstance) MarshalJSON() ([]byte, error) {

	queue := make([]*WorkItem, inst.workItemQueue.List.Len())

	for i, e := 0, inst.workItemQueue.List.Front(); e != nil; i, e = i+1, e.Next() {
		queue[i], _ = e.Value.(*WorkItem)
	}

	attrs := make([]*data.Attribute, 0, len(inst.attrs))

	for _, value := range inst.attrs {
		attrs = append(attrs, value)
	}

	tis := make([]*TaskInst, 0, len(inst.taskInsts))

	for _, value := range inst.taskInsts {
		tis = append(tis, value)
	}

	lis := make([]*LinkInst, 0, len(inst.linkInsts))

	for _, value := range inst.linkInsts {
		lis = append(lis, value)
	}

	sfs := make([]*EmbeddedInstance, 0, len(inst.subFlows))

	for _, value := range inst.subFlows {
		sfs = append(sfs, value)
	}

	//serialize all the subFlows

	return json.Marshal(&serIndependentInstance{
		ID:        inst.id,
		Status:    inst.status,
		Attrs:     attrs,
		FlowURI:   inst.flowURI,
		WorkQueue: queue,
		TaskInsts: tis,
		LinkInsts: lis,
		SubFlows:  sfs,
	})
}

// UnmarshalJSON overrides the default UnmarshalJSON for FlowInstance
func (inst *IndependentInstance) UnmarshalJSON(d []byte) error {

	ser := &serIndependentInstance{}
	if err := json.Unmarshal(d, ser); err != nil {
		return err
	}

	inst.id = ser.ID
	inst.status = ser.Status
	inst.flowURI = ser.FlowURI

	inst.attrs = make(map[string]*data.Attribute)

	for _, value := range ser.Attrs {
		inst.attrs[value.Name()] = value
	}

	inst.ChangeTracker = NewInstanceChangeTracker()

	inst.taskInsts = make(map[string]*TaskInst, len(ser.TaskInsts))
	inst.linkInsts = make(map[int]*LinkInst, len(ser.LinkInsts))

	for _, value := range ser.TaskInsts {
		inst.taskInsts[value.taskID] = value
	}

	for _, value := range ser.LinkInsts {
		inst.linkInsts[value.linkID] = value
	}

	subFlowCtr := 0

	if len(ser.SubFlows) > 0 {

		inst.subFlows = make(map[int]*EmbeddedInstance, len(ser.SubFlows))

		for _, value := range ser.SubFlows {
			inst.subFlows[value.subFlowId] = value

			if value.subFlowId > subFlowCtr {
				subFlowCtr = value.subFlowId
			}
		}

		inst.subFlowCtr = subFlowCtr
	}

	inst.workItemQueue = util.NewSyncQueue()

	for _, workItem := range ser.WorkQueue {

		taskInsts := inst.taskInsts

		if workItem.SubFlowID > 0 {
			taskInsts = inst.subFlows[workItem.SubFlowID].taskInsts
		}

		workItem.taskInst = taskInsts[workItem.TaskID]
		inst.workItemQueue.Push(workItem)
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
// Embedded Flow Instance Serialization

type serEmbeddedInstance struct {
	SubFlowId int
	Status    model.FlowStatus  `json:"status"`
	FlowURI   string            `json:"flowUri"`
	Attrs     []*data.Attribute `json:"attrs"`
	TaskDatas []*TaskInst       `json:"taskDatas"`
	LinkDatas []*LinkInst       `json:"linkDatas"`
}

// MarshalJSON overrides the default MarshalJSON for FlowInstance
func (inst *EmbeddedInstance) MarshalJSON() ([]byte, error) {

	attrs := make([]*data.Attribute, 0, len(inst.attrs))

	for _, value := range inst.attrs {
		attrs = append(attrs, value)
	}

	tds := make([]*TaskInst, 0, len(inst.taskInsts))

	for _, value := range inst.taskInsts {
		tds = append(tds, value)
	}

	lds := make([]*LinkInst, 0, len(inst.linkInsts))

	for _, value := range inst.linkInsts {
		lds = append(lds, value)
	}

	//serialize all the subFlows

	return json.Marshal(&serEmbeddedInstance{
		SubFlowId: inst.subFlowId,
		Status:    inst.status,
		Attrs:     attrs,
		FlowURI:   inst.flowURI,
		TaskDatas: tds,
		LinkDatas: lds,
	})
}

// UnmarshalJSON overrides the default UnmarshalJSON for FlowInstance
func (inst *EmbeddedInstance) UnmarshalJSON(d []byte) error {

	ser := &serEmbeddedInstance{}
	if err := json.Unmarshal(d, ser); err != nil {
		return err
	}

	inst.subFlowId = ser.SubFlowId
	inst.status = ser.Status
	inst.flowURI = ser.FlowURI

	inst.attrs = make(map[string]*data.Attribute)

	for _, value := range ser.Attrs {
		inst.attrs[value.Name()] = value
	}

	inst.taskInsts = make(map[string]*TaskInst, len(ser.TaskDatas))
	inst.linkInsts = make(map[int]*LinkInst, len(ser.LinkDatas))

	for _, value := range ser.TaskDatas {
		inst.taskInsts[value.taskID] = value
	}

	for _, value := range ser.LinkDatas {
		inst.linkInsts[value.linkID] = value
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
// TaskInst Serialization

// MarshalJSON overrides the default MarshalJSON for TaskInst
func (ti *TaskInst) MarshalJSON() ([]byte, error) {

	return json.Marshal(&struct {
		TaskID string `json:"taskId"`
		State  int    `json:"state"`
		Status int    `json:"status"`
	}{
		TaskID: ti.task.ID(),
		State:  int(ti.status),
		Status: int(ti.status),
	})
}

// UnmarshalJSON overrides the default UnmarshalJSON for TaskInst
func (ti *TaskInst) UnmarshalJSON(d []byte) error {
	ser := &struct {
		TaskID string `json:"taskId"`
		State  int    `json:"state"`
		Status int    `json:"status"`
	}{}

	if err := json.Unmarshal(d, ser); err != nil {
		return err
	}

	ti.status = model.TaskStatus(ser.Status)
	ti.taskID = ser.TaskID

	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
// LinkInst Serialization

// MarshalJSON overrides the default MarshalJSON for LinkInst
func (ld *LinkInst) MarshalJSON() ([]byte, error) {

	return json.Marshal(&struct {
		LinkID int `json:"linkId"`
		State  int `json:"state"`
		Status int `json:"status"`
	}{
		LinkID: ld.link.ID(),
		State:  int(ld.status),
		Status: int(ld.status),
	})
}

// UnmarshalJSON overrides the default UnmarshalJSON for LinkInst
func (ld *LinkInst) UnmarshalJSON(d []byte) error {
	ser := &struct {
		LinkID int `json:"linkId"`
		State  int `json:"state"`
		Status int `json:"status"`
	}{}

	if err := json.Unmarshal(d, ser); err != nil {
		return err
	}

	ld.status = model.LinkStatus(ser.Status)
	ld.linkID = ser.LinkID

	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
// Flow Instance Changes Serialization

// MarshalJSON overrides the default MarshalJSON for InstanceChangeTracker
func (ict *InstanceChangeTracker) MarshalJSON() ([]byte, error) {

	var wqc []*WorkItemQueueChange

	if ict.wiqChanges != nil {
		wqc = make([]*WorkItemQueueChange, 0, len(ict.wiqChanges))

		for _, value := range ict.wiqChanges {
			wqc = append(wqc, value)
		}

	} else {
		wqc = nil
	}

	// for backwards compatibility //
	var tdc []*TaskInstChange
	var ldc []*LinkInstChange
	var attrs []*AttributeChange
	var status model.FlowStatus
	/////////////////////////////////

	var instChanges []*InstanceChange

	if ict.instChanges != nil {
		instChanges = make([]*InstanceChange, 0, len(ict.instChanges))

		for _, value := range ict.instChanges {
			instChanges = append(instChanges, value)
		}

		// for backwards compatibility
		masterChg, ok := ict.instChanges[0]
		if ok {
			status = masterChg.Status
			attrs = masterChg.AttrChanges
		}

		if masterChg.tiChanges != nil {
			tdc = make([]*TaskInstChange, 0, len(masterChg.tiChanges))

			for _, value := range masterChg.tiChanges {
				tdc = append(tdc, value)
			}
		} else {
			tdc = nil
		}

		if masterChg.liChanges != nil {
			ldc = make([]*LinkInstChange, 0, len(masterChg.liChanges))

			for _, value := range masterChg.liChanges {
				ldc = append(ldc, value)
			}
		} else {
			ldc = nil
		}
	} else {
		instChanges = nil
	}

	//backwards compatibility
	return json.Marshal(&struct {
		Status      model.FlowStatus       `json:"status,omitempty"`
		AttrChanges []*AttributeChange     `json:"attrs,omitempty"`
		WqChanges   []*WorkItemQueueChange `json:"wqChanges,omitempty"`
		TdChanges   []*TaskInstChange      `json:"tdChanges,omitempty"`
		LdChanges   []*LinkInstChange      `json:"ldChanges,omitempty"`
		InstChanges []*InstanceChange      `json:"ldChanges,omitempty"`
	}{
		Status:      status,
		AttrChanges: attrs,
		WqChanges:   wqc,
		TdChanges:   tdc,
		LdChanges:   ldc,
		InstChanges: instChanges,
	})
}

// MarshalJSON overrides the default MarshalJSON for InstanceChangeTracker
func (ic *InstanceChange) MarshalJSON() ([]byte, error) {

	var tdc []*TaskInstChange

	if ic.tiChanges != nil {
		tdc = make([]*TaskInstChange, 0, len(ic.tiChanges))

		for _, value := range ic.tiChanges {
			tdc = append(tdc, value)
		}
	} else {
		tdc = nil
	}

	var ldc []*LinkInstChange

	if ic.liChanges != nil {
		ldc = make([]*LinkInstChange, 0, len(ic.liChanges))

		for _, value := range ic.liChanges {
			ldc = append(ldc, value)
		}
	} else {
		ldc = nil
	}

	if ic.Status > -1 {
		//backwards compatibility
		return json.Marshal(&struct {
			FlowID      int                `json:"flowId"`
			Status      model.FlowStatus   `json:"status"`
			AttrChanges []*AttributeChange `json:"attrs,omitempty"`
			TdChanges   []*TaskInstChange  `json:"tasks,omitempty"`
			LdChanges   []*LinkInstChange  `json:"links,omitempty"`
			SfChange    *SubFlowChange     `json:"subFlow,omitempty"`
		}{
			FlowID:      ic.SubFlowID,
			Status:      ic.Status,
			AttrChanges: ic.AttrChanges,
			TdChanges:   tdc,
			LdChanges:   ldc,
			SfChange:    ic.SubFlowChg,
		})
	}

	return json.Marshal(&struct {
		FlowID      int                `json:"flowId"`
		AttrChanges []*AttributeChange `json:"attrs,omitempty"`
		TdChanges   []*TaskInstChange  `json:"tasks,omitempty"`
		LdChanges   []*LinkInstChange  `json:"links,omitempty"`
		SfChange    *SubFlowChange     `json:"subFlow,omitempty"`
	}{
		FlowID:      ic.SubFlowID,
		AttrChanges: ic.AttrChanges,
		TdChanges:   tdc,
		LdChanges:   ldc,
		SfChange:    ic.SubFlowChg,
	})
}
