package instance2

import (
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

func applyInputMapper(taskData *TaskData) error {

	// get the input mapper
	inputMapper := taskData.task.ActivityConfig().InputMapper()

	master := taskData.inst.master

	if master.patch != nil {
		// check if the patch has a overriding mapper
		mapper := master.patch.GetInputMapper(taskData.task.ID())
		if mapper != nil {
			inputMapper = mapper
		}
	}

	if inputMapper != nil {
		logger.Debug("Applying InputMapper")

		var inputScope data.Scope
		inputScope = taskData.inst

		if taskData.workingData != nil {
			inputScope = NewWorkingDataScope(taskData.inst, taskData.workingData)
		}

		err := inputMapper.Apply(inputScope, taskData.InputScope())

		if err != nil {
			return err
		}
	}

	return nil
}

func applyInputInterceptor(taskData *TaskData) bool {

	master := taskData.inst.master

	if master.interceptor != nil {

		// check if this task as an interceptor
		taskInterceptor := master.interceptor.GetTaskInterceptor(taskData.task.ID())

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

func applyOutputInterceptor(taskData *TaskData) {

	master := taskData.inst.master

	if master.interceptor != nil {

		// check if this task as an interceptor and overrides ouputs
		taskInterceptor := master.interceptor.GetTaskInterceptor(taskData.task.ID())
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
func applyOutputMapper(taskData *TaskData) (bool, error) {

	// get the Output Mapper for the TaskOld if one exists
	outputMapper := taskData.task.ActivityConfig().OutputMapper()

	master := taskData.inst.master

	if master.patch != nil {
		// check if the patch overrides the Output Mapper
		mapper := master.patch.GetOutputMapper(taskData.task.ID())
		if mapper != nil {
			outputMapper = mapper
		}
	}

	if outputMapper != nil {
		logger.Debug("Applying OutputMapper")
		err := outputMapper.Apply(taskData.OutputScope(), taskData.inst)

		return true, err
	}

	return false, nil
}
