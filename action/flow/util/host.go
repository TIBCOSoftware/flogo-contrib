package util

import "github.com/TIBCOSoftware/flogo-lib/core/data"

type FlowHost interface {

	// Reply is used to reply to the Flow Host with the results of the execution
	Reply(replyData map[string]*data.Attribute, err error)

	// Return is used to indicate to the Flow Host that it should complete and return the results of the execution
	Return(returnData map[string]*data.Attribute, err error)

}