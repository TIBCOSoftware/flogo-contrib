package flowinst

import (
	"strconv"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/flow/activity"
	"github.com/TIBCOSoftware/flogo-lib/flow/flowdef"
)

func applyInputMapper(pi *Instance, taskData *TaskData) {

	// get the input mapper
	inputMapper := taskData.task.InputMapper()

	if pi.Patch != nil {
		// check if the patch has a overriding mapper
		mapper := pi.Patch.GetInputMapper(taskData.task.ID())
		if mapper != nil {
			inputMapper = mapper
		}
	}

	if inputMapper != nil {
		log.Debug("Applying InputMapper")
		inputMapper.Apply(pi, taskData.InputScope())
	}
}

func applyInputInterceptor(pi *Instance, taskData *TaskData) bool {

	if pi.Interceptor != nil {

		// check if this task as an interceptor
		taskInterceptor := pi.Interceptor.GetTaskInterceptor(taskData.task.ID())

		if taskInterceptor != nil {

			log.Debug("Applying Interceptor")

			if len(taskInterceptor.Inputs) > 0 {
				// override input attributes
				for _, attribute := range taskInterceptor.Inputs {

					log.Debugf("Overriding Attr: %s = %s", attribute.Name, attribute.Value)

					//todo: validation
					taskData.InputScope().SetAttrValue(attribute.Name, attribute.Value)
				}
			}

			// check if we should not evaluate the task
			return !taskInterceptor.Skip
		}
	}

	return true
}

func applyOutputInterceptor(pi *Instance, taskData *TaskData) {

	if pi.Interceptor != nil {

		// check if this task as an interceptor and overrides ouputs
		taskInterceptor := pi.Interceptor.GetTaskInterceptor(taskData.task.ID())
		if taskInterceptor != nil && len(taskInterceptor.Outputs) > 0 {
			// override output attributes
			for _, attribute := range taskInterceptor.Outputs {

				//todo: validation
				taskData.OutputScope().SetAttrValue(attribute.Name, attribute.Value)
			}
		}
	}
}

// applyOutputMapper applies the output mapper, returns flag indicating if
// there was an output mapper
func applyOutputMapper(pi *Instance, taskData *TaskData) bool {

	// get the Output Mapper for the Task if one exists
	outputMapper := taskData.task.OutputMapper()

	if pi.Patch != nil {
		// check if the patch overrides the Output Mapper
		mapper := pi.Patch.GetOutputMapper(taskData.task.ID())
		if mapper != nil {
			outputMapper = mapper
		}
	}

	if outputMapper != nil {
		log.Debug("Applying OutputMapper")
		outputMapper.Apply(taskData.OutputScope(), pi)
		return true
	}

	return false
}

func applyDefaultActivityOutputMappings(pi *Instance, taskData *TaskData) {

	activity := activity.Get(taskData.task.ActivityType())

	attrNS := "{A" + strconv.Itoa(taskData.task.ID()) + "."

	for _, attr := range activity.Metadata().Outputs {

		oAttr, _ := taskData.OutputScope().GetAttr(attr.Name)

		if oAttr != nil {
			pi.AddAttr(attrNS+attr.Name+"}", attr.Type, oAttr.Value)
		}
	}
}

func applyDefaultInstanceInputMappings(pi *Instance, attrs []*data.Attribute) {

	if len(attrs) == 0 {
		return
	}

	for _, attr := range attrs {

		attrName := "{T." + attr.Name + "}"
		pi.AddAttr(attrName, attr.Type, attr.Value)
	}
}

// FixedTaskScope is scope restricted by the set of reference attrs and backed by the specified Task
type FixedTaskScope struct {
	attrs    map[string]*data.Attribute
	refAttrs map[string]*data.Attribute
	task     *flowdef.Task
}

// NewFixedTaskScope creates a FixedTaskScope
func NewFixedTaskScope(refAttrs map[string]*data.Attribute, task *flowdef.Task) data.Scope {

	scope := &FixedTaskScope{
		refAttrs: refAttrs,
		task:     task,
	}

	return scope
}

//// GetAttrType implements Scope.GetAttrType
//func (s *FixedTaskScope) GetAttrType(attrName string) (attrType data.Type, exists bool) {
//
//	attr, found := s.refAttrs[attrName]
//
//	if found {
//		return attr.Type, true
//	}
//
//	return 0, false
//}

// GetAttr implements Scope.GetAttr
func (s *FixedTaskScope) GetAttr(attrName string) (attr *data.Attribute, exists bool) {

	if len(s.attrs) > 0 {

		attr, found := s.attrs[attrName]

		if found {
			return attr, true
		}
	}

	if s.task != nil {

		attr, found := s.task.GetAttr(attrName)

		if !found {
			attr, found = s.refAttrs[attrName]
		}

		return attr, found
	}

	return nil, false
}

// SetAttrValue implements Scope.SetAttrValue
func (s *FixedTaskScope) SetAttrValue(attrName string, value interface{}) error {

	if len(s.attrs) == 0 {
		s.attrs = make(map[string]*data.Attribute)
	}

	log.Debugf("SetAttr: %s = %v\n", attrName, value)

	attr, found := s.attrs[attrName]

	if found {
		//todo handle errors
		coercedVal, _ := data.CoerceToValue(value, attr.Type)
		attr.Value = coercedVal
	} else {
		// look up reference for type
		attr, found = s.refAttrs[attrName]
		if found {
			coercedVal, _ := data.CoerceToValue(value, attr.Type)
			s.attrs[attrName] = data.NewAttribute(attrName, attr.Type, coercedVal)
		} else {
			log.Debugf("SetAttr: Attr %s ref not found\n", attrName)
			log.Debugf("SetAttr: refs %v\n", s.refAttrs)
		}
		//todo: else error
	}

	return nil
}
