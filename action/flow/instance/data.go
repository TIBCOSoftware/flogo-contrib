package instance

import (
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/definition"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

func applyInputMapper(pi *Instance, taskData *TaskData) error {

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
		logger.Debug("Applying InputMapper")
		err := inputMapper.Apply(pi, taskData.InputScope())

		if err != nil {
			return err
		}
	}

	return nil
}

func applyInputInterceptor(pi *Instance, taskData *TaskData) bool {

	if pi.Interceptor != nil {

		// check if this task as an interceptor
		taskInterceptor := pi.Interceptor.GetTaskInterceptor(taskData.task.ID())

		if taskInterceptor != nil {

			logger.Debug("Applying Interceptor")

			if len(taskInterceptor.Inputs) > 0 {
				// override input attributes
				for _, attribute := range taskInterceptor.Inputs {

					logger.Debugf("Overriding Attr: %s = %s", attribute.Name(), attribute.Value())

					//todo: validation
					taskData.InputScope().SetAttrValue(attribute.Name(), attribute.Value())
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
				taskData.OutputScope().SetAttrValue(attribute.Name(), attribute.Value())
			}
		}
	}
}

// applyOutputMapper applies the output mapper, returns flag indicating if
// there was an output mapper
func applyOutputMapper(pi *Instance, taskData *TaskData) (bool, error) {

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
		logger.Debug("Applying OutputMapper")
		err := outputMapper.Apply(taskData.OutputScope(), pi)

		return true, err
	}

	return false, nil
}

// FixedTaskScope is scope restricted by the set of reference attrs and backed by the specified Task
type FixedTaskScope struct {
	attrs    map[string]*data.Attribute
	refAttrs map[string]*data.Attribute
	task     *definition.Task
	isInput  bool
}

// NewFixedTaskScope creates a FixedTaskScope
func NewFixedTaskScope(refAttrs map[string]*data.Attribute, task *definition.Task, isInput bool) data.Scope {

	scope := &FixedTaskScope{
		refAttrs: refAttrs,
		task:     task,
		isInput:  isInput,
	}

	return scope
}

// GetAttr implements Scope.GetAttr
func (s *FixedTaskScope) GetAttr(attrName string) (attr *data.Attribute, exists bool) {

	if len(s.attrs) > 0 {

		attr, found := s.attrs[attrName]

		if found {
			return attr, true
		}
	}

	if s.task != nil {

		var attr *data.Attribute
		var found bool

		if s.isInput {
			attr, found = s.task.GetInputAttr(attrName)
		} else {
			attr, found = s.task.GetOutputAttr(attrName)
		}

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

	logger.Debugf("SetAttr: %s = %v\n", attrName, value)

	attr, found := s.attrs[attrName]

	var err error
	if found {
		err = attr.SetValue(value)
	} else {
		// look up reference for type
		attr, found = s.refAttrs[attrName]
		if found {
			s.attrs[attrName], err = data.NewAttribute(attrName, attr.Type(), value)
		} else {
			logger.Debugf("SetAttr: Attr %s ref not found\n", attrName)
			logger.Debugf("SetAttr: refs %v\n", s.refAttrs)
		}
		//todo: else error
	}

	return err
}
