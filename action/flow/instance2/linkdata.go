package instance2

import (
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/definition"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/model"
)

// LinkData represents data associated with an instance of a Link
type LinkData struct {
	inst *Instance
	link    *definition.Link
	status  model.LinkStatus

	changes int
	linkID int //needed for serialization
}

// NewLinkData creates a LinkData for the specified link in the specified task
// environment
func NewLinkData(inst *Instance, link *definition.Link) *LinkData {
	var linkData LinkData

	linkData.inst = inst
	linkData.link = link

	return &linkData
}

// Status returns the current state indicator for the LinkData
func (ld *LinkData) Status() model.LinkStatus {
	return ld.status
}

// SetStatus sets the current state indicator for the LinkData
func (ld *LinkData) SetStatus(status model.LinkStatus) {
	ld.status = status
	ld.inst.ChangeTracker.trackLinkData(&LinkDataChange{ChgType: CtUpd, ID: ld.link.ID(), LinkData: ld})
}

// Link returns the Link associated with ld context
func (ld *LinkData) Link() *definition.Link {
	return ld.link
}
