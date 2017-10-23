package response

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// log is the default package logger
var log = logger.GetLogger("activity-tibco-response")

const (
	ivMappings = "mappings"
	ivOperation = "operation"
)

// ResponseActivity is an Activity that is used to reply/return via the trigger
// inputs : {method,uri,params}
// outputs: {result}
type ResponseActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new ResponseActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &ResponseActivity{metadata: metadata}
}

// Metadata returns the activity's metadata
func (a *ResponseActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Invokes a REST Operation
func (a *ResponseActivity) Eval(context activity.Context) (done bool, err error) {

	mappings := context.GetInput(ivMappings).([]interface{})
	operation := context.GetInput(ivOperation).(string)

	log.Debugf("Operation :'%s', Mappings: %+v", operation, mappings)

	//todo move this to a action instance level initialization, need the notion of static inputs or config
	replyMapper, err := createMapper(mappings)

	if err != nil {
		return false, nil
	}

	actionCtx := context.ActionContext()

	outAttrs := actionCtx.InstanceMetadata().Output
	attrs := make([]*data.Attribute, 0, len(outAttrs))

	for _, outAttr := range outAttrs {
		attrs = append(attrs, outAttr)
	}

	//create a fixed scope using the output metadata
	outputScope := data.NewFixedScope(attrs)
	inputScope  :=  actionCtx.WorkingData() //flow data

	err = replyMapper.Apply(inputScope, outputScope)

	if err != nil {
		return false, nil
	}

	actionCtx.ReplyWithAttrs(outputScope.GetAttrs(), nil)
	//actionCtx.Reply()

	return true, nil
}

func createMapper(mappings []interface{}) (data.Mapper, error) {

	var mappingDefs []*data.MappingDef

	for _, mapping := range mappings {

		mappingObject := mapping.(map[string]interface{})

		mappingType := int(mappingObject["type"].(float64))
		value := mappingObject["value"]
		mapTo := mappingObject["mapTo"].(string)

		mappingDef := &data.MappingDef{Type:data.MappingType(mappingType), MapTo:mapTo, Value:value}
		mappingDefs = append(mappingDefs, mappingDef)
	}

	mapperDef := &data.MapperDef{Mappings:mappingDefs}
	basicMapper := mapper.NewBasicMapper(mapperDef)

	return basicMapper, nil
}
