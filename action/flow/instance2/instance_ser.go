package instance2

import (
	"encoding/json"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/util"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/model"
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
	TaskDatas []*TaskData         `json:"taskDatas"`
	LinkDatas []*LinkData         `json:"linkDatas"`
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

	tds := make([]*TaskData, 0, len(inst.taskDataMap))

	for _, value := range inst.taskDataMap {
		tds = append(tds, value)
	}

	lds := make([]*LinkData, 0, len(inst.linkDataMap))

	for _, value := range inst.linkDataMap {
		lds = append(lds, value)
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
		TaskDatas: tds,
		LinkDatas: lds,
		SubFlows: sfs,
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

	inst.taskDataMap = make(map[string]*TaskData, len(ser.TaskDatas))
	inst.linkDataMap = make(map[int]*LinkData, len(ser.LinkDatas))

	for _, value := range ser.TaskDatas {
		inst.taskDataMap[value.taskID] = value
	}

	for _, value := range ser.LinkDatas {
		inst.linkDataMap[value.linkID] = value
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

		taskDatas := inst.taskDataMap

		if workItem.SubFlowID > 0 {
			taskDatas = inst.subFlows[workItem.SubFlowID].taskDataMap
		}

		workItem.TaskData = taskDatas[workItem.TaskID]
		inst.workItemQueue.Push(workItem)
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
// Embedded Flow Instance Serialization

type serEmbeddedInstance struct {
	SubFlowId int
	Status    model.FlowStatus    `json:"status"`
	FlowURI   string              `json:"flowUri"`
	Attrs     []*data.Attribute   `json:"attrs"`
	TaskDatas []*TaskData         `json:"taskDatas"`
	LinkDatas []*LinkData         `json:"linkDatas"`
}

// MarshalJSON overrides the default MarshalJSON for FlowInstance
func (inst *EmbeddedInstance) MarshalJSON() ([]byte, error) {

	attrs := make([]*data.Attribute, 0, len(inst.attrs))

	for _, value := range inst.attrs {
		attrs = append(attrs, value)
	}

	tds := make([]*TaskData, 0, len(inst.taskDataMap))

	for _, value := range inst.taskDataMap {
		tds = append(tds, value)
	}

	lds := make([]*LinkData, 0, len(inst.linkDataMap))

	for _, value := range inst.linkDataMap {
		lds = append(lds, value)
	}

	//serialize all the subFlows

	return json.Marshal(&serEmbeddedInstance{
		SubFlowId:        inst.subFlowId,
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

	inst.taskDataMap = make(map[string]*TaskData, len(ser.TaskDatas))
	inst.linkDataMap = make(map[int]*LinkData, len(ser.LinkDatas))

	for _, value := range ser.TaskDatas {
		inst.taskDataMap[value.taskID] = value
	}

	for _, value := range ser.LinkDatas {
		inst.linkDataMap[value.linkID] = value
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
// TaskData Serialization

// MarshalJSON overrides the default MarshalJSON for TaskData
func (td *TaskData) MarshalJSON() ([]byte, error) {

	return json.Marshal(&struct {
		TaskID string `json:"taskId"`
		State  int    `json:"state"`
		Status int    `json:"status"`
	}{
		TaskID: td.task.ID(),
		State:  int(td.status),
		Status: int(td.status),
	})
}

// UnmarshalJSON overrides the default UnmarshalJSON for TaskData
func (td *TaskData) UnmarshalJSON(d []byte) error {
	ser := &struct {
		TaskID string `json:"taskId"`
		State  int    `json:"state"`
		Status int    `json:"status"`
	}{}

	if err := json.Unmarshal(d, ser); err != nil {
		return err
	}

	td.status = model.TaskStatus(ser.Status)
	td.taskID = ser.TaskID

	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
// LinkData Serialization

// MarshalJSON overrides the default MarshalJSON for LinkData
func (ld *LinkData) MarshalJSON() ([]byte, error) {

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

// UnmarshalJSON overrides the default UnmarshalJSON for LinkData
func (ld *LinkData) UnmarshalJSON(d []byte) error {
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

	var tdc []*TaskDataChange

	if ict.tdChanges != nil {
		tdc = make([]*TaskDataChange, 0, len(ict.tdChanges))

		for _, value := range ict.tdChanges {
			tdc = append(tdc, value)
		}
	} else {
		tdc = nil
	}

	var ldc []*LinkDataChange

	if ict.ldChanges != nil {
		ldc = make([]*LinkDataChange, 0, len(ict.ldChanges))

		for _, value := range ict.ldChanges {
			ldc = append(ldc, value)
		}
	} else {
		ldc = nil
	}

	return json.Marshal(&struct {
		Status      model.FlowStatus       `json:"status"`
		AttrChanges []*AttributeChange     `json:"attrs"`
		WqChanges   []*WorkItemQueueChange `json:"wqChanges"`
		TdChanges   []*TaskDataChange      `json:"tdChanges"`
		LdChanges   []*LinkDataChange      `json:"ldChanges"`
	}{
		Status:      ict.instChange.Status,
		AttrChanges: ict.instChange.AttrChanges,
		WqChanges:   wqc,
		TdChanges:   tdc,
		LdChanges:   ldc,
	})
}
