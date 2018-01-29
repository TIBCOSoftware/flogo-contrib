package instance2

import (
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

func applyInputMapper(taskData *TaskData) error {

	// get the input mapper
	inputMapper := taskData.task.ActivityConfig().InputMapper()

	pi := taskData.execEnv.Instance

	//if pi.Patch != nil {
	//	// check if the patch has a overriding mapper
	//	mapper := pi.Patch.GetInputMapper(taskData.task.ID())
	//	if mapper != nil {
	//		inputMapper = mapper
	//	}
	//}

	if inputMapper != nil {
		logger.Debug("Applying InputMapper")

		var inputScope data.Scope
		inputScope = taskData.execEnv

		if taskData.workingData != nil {
			inputScope = NewWorkingDataScope(taskData.execEnv, taskData.workingData)
		}

		err := inputMapper.Apply(inputScope, taskData.InputScope())

		if err != nil {
			return err
		}
	}

	return nil
}

func applyInputInterceptor(taskData *TaskData) bool {

	//pi := taskData.execEnv.Instance

	//if pi.Interceptor != nil {
	//
	//	// check if this task as an interceptor
	//	taskInterceptor := pi.Interceptor.GetTaskInterceptor(taskData.task.ID())
	//
	//	if taskInterceptor != nil {
	//
	//		logger.Debug("Applying Interceptor")
	//
	//		if len(taskInterceptor.Inputs) > 0 {
	//			// override input attributes
	//			for _, attribute := range taskInterceptor.Inputs {
	//
	//				logger.Debugf("Overriding Attr: %s = %s", attribute.Name(), attribute.Value())
	//
	//				//todo: validation
	//				taskData.InputScope().SetAttrValue(attribute.Name(), attribute.Value())
	//			}
	//		}
	//
	//		// check if we should not evaluate the task
	//		return !taskInterceptor.Skip
	//	}
	//}

	return true
}

func applyOutputInterceptor(taskData *TaskData) {

	//pi := taskData.execEnv.Instance
	//
	//if pi.Interceptor != nil {
	//
	//	// check if this task as an interceptor and overrides ouputs
	//	taskInterceptor := pi.Interceptor.GetTaskInterceptor(taskData.task.ID())
	//	if taskInterceptor != nil && len(taskInterceptor.Outputs) > 0 {
	//		// override output attributes
	//		for _, attribute := range taskInterceptor.Outputs {
	//
	//			//todo: validation
	//			taskData.OutputScope().SetAttrValue(attribute.Name(), attribute.Value())
	//		}
	//	}
	//}
}

// applyOutputMapper applies the output mapper, returns flag indicating if
// there was an output mapper
func applyOutputMapper(taskData *TaskData) (bool, error) {

	// get the Output Mapper for the TaskOld if one exists
	outputMapper := taskData.task.ActivityConfig().OutputMapper()

	//pi := taskData.execEnv.Instance

	//if pi.Patch != nil {
	//	// check if the patch overrides the Output Mapper
	//	mapper := pi.Patch.GetOutputMapper(taskData.task.ID())
	//	if mapper != nil {
	//		outputMapper = mapper
	//	}
	//}

	if outputMapper != nil {
		logger.Debug("Applying OutputMapper")
		err := outputMapper.Apply(taskData.OutputScope(), taskData.execEnv)

		return true, err
	}

	return false, nil
}
