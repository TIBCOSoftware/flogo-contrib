package simple

import (
	"github.com/TIBCOSoftware/flogo-lib/flow/flowdef"
	"github.com/TIBCOSoftware/flogo-lib/flow/model"
	"github.com/op/go-logging"
)

// log is the default package logger
var log = logging.MustGetLogger("model-tibco-simple")

const (
	MODEL_NAME = "tibco-simple"
)

func init() {
	model.Register(New())
}

////////////////////////////////////////////////////////////////////////////////////////////////////////

func New() *model.FlowModel {
	m := model.New(MODEL_NAME)
	m.RegisterFlowBehavior(&SimpleFlowBehavior{})
	m.RegisterTaskBehavior(1, &SimpleTaskBehavior{})
	return m
}

// SimpleFlowBehavior implements model.FlowBehavior
type SimpleFlowBehavior struct {
}

// Start implements model.FlowBehavior.Start
func (pb *SimpleFlowBehavior) Start(context model.FlowContext) (start bool, evalCode int) {
	// just schedule the root task
	return true, 0
}

// Resume implements model.FlowBehavior.Resume
func (pb *SimpleFlowBehavior) Resume(context model.FlowContext) bool {
	return true
}

// TasksDone implements model.FlowBehavior.TasksDone
func (pb *SimpleFlowBehavior) TasksDone(context model.FlowContext, doneCode int) {
	// all tasks are done
}

// Done implements model.FlowBehavior.Done
func (pb *SimpleFlowBehavior) Done(context model.FlowContext) {
	log.Debugf("Flow Done\n")
}

////////////////////////////////////////////////////////////////////////////////////////////////////////

// SimpleTaskBehavior implements model.TaskBehavior
type SimpleTaskBehavior struct {
}

// Enter implements model.TaskBehavior.Enter
func (tb *SimpleTaskBehavior) Enter(context model.TaskContext, enterCode int) (eval bool, evalCode int) {

	task := context.Task()
	log.Debugf("Task Enter: %s\n", task.Name())

	context.SetState(STATE_ENTERED)

	//check if all predecessor links are done

	linkContexts := context.FromInstLinks()

	ready := true

	if len(linkContexts) == 0 {
		// has no predecessor links, so task is ready
		ready = true
	} else {

		log.Debugf("Num Links: %d\n", len(linkContexts))
		for _, linkContext := range linkContexts {

			log.Debugf("Task: %s, linkData: %v\n", task.Name(), linkContext)
			if linkContext.State() != STATE_LINK_TRUE {
				ready = false
				break
			}
		}
	}

	if ready {
		log.Debugf("Task Ready\n")
		context.SetState(STATE_READY)
	} else {
		log.Debugf("Task Not Ready\n")
	}

	return ready, 0
}

// Eval implements model.TaskBehavior.Eval
func (tb *SimpleTaskBehavior) Eval(context model.TaskContext, evalCode int) (done bool, doneCode int, err error) {

	task := context.Task()
	log.Debugf("Task Eval: %s\n", task)

	if len(task.ChildTasks()) > 0 {
		log.Debugf("Has Children\n")

		//has children, so set to waiting
		context.SetState(STATE_WAITING)

		context.EnterLeadingChildren(0)

		return false, 0, nil

	} else {

		if context.HasActivity() {

			done, err := context.EvalActivity()

			// todo handle error transition
			if err != nil {
				log.Errorf("Error evaluating activity '%s'[%s] - %s", context.Task().Name(), context.Task().ActivityType(), err.Error())
				context.SetState(STATE_FAILED)

				//we don't have an error transition, so we'll return it so the global error handler can deal with it
				return false, 0, err
			}

			return done, 0, nil
		}

		//no-op
		return true, 0, nil
	}
}

// PostEval implements model.TaskBehavior.PostEval
func (tb *SimpleTaskBehavior) PostEval(context model.TaskContext, evalCode int, data interface{}) (done bool, doneCode int, err error) {

	log.Debugf("Task PostEval\n")

	if context.HasActivity() { //if activity is async

		//done := activity.PostEval(activityContext, data)
		done := true
		return done, 0, nil
	}

	//no-op
	return true, 0, nil
}

// Done implements model.TaskBehavior.Done
func (tb *SimpleTaskBehavior) Done(context model.TaskContext, doneCode int) (notifyParent bool, childDoneCode int, taskEntries []*model.TaskEntry) {

	task := context.Task()
	log.Debugf("done task:%s\n", task.Name())

	context.SetState(STATE_DONE)
	//context.SetTaskDone() for task garbage collection

	linkInsts := context.ToInstLinks()
	numLinks := len(linkInsts)

	// process outgoing links
	if numLinks > 0 {

		taskEntries = make([]*model.TaskEntry, 0, numLinks)

		for _, linkInst := range linkInsts {

			follow := true

			if linkInst.Link().Type() == flowdef.LtExpression {
				//todo handle error
				follow, _ = context.EvalLink(linkInst.Link())
			}

			if follow {
				linkInst.SetState(STATE_LINK_TRUE)

				taskEntry := &model.TaskEntry{Task: linkInst.Link().ToTask(), EnterCode: 0}
				taskEntries = append(taskEntries, taskEntry)
			}
		}

		//continue on to successor tasks
		return false, 0, taskEntries
	}

	// there are no outgoing links, so just notify parent that we are done
	return true, 0, nil
}

// ChildDone implements model.TaskBehavior.ChildDone
func (tb *SimpleTaskBehavior) ChildDone(context model.TaskContext, childTask *flowdef.Task, childDoneCode int) (done bool, doneCode int) {

	childTasks, hasChildren := context.ChildTaskInsts()

	if !hasChildren {
		log.Debugf("Task ChildDone - No Children\n")
		return true, 0
	}

	for _, taskInst := range childTasks {

		if taskInst.State() != STATE_DONE {
			return false, 0
		}
	}

	// our children are done, so just transition ourselves to done
	return true, 0
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
// State
const (
	STATE_NOT_STARTED int = 0

	STATE_LINK_FALSE int = 1
	STATE_LINK_TRUE  int = 2

	STATE_ENTERED int = 10
	STATE_READY   int = 20
	STATE_WAITING int = 30
	STATE_DONE    int = 40
	STATE_FAILED  int = 100
)
