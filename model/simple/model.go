package simple

import (

	"github.com/TIBCOSoftware/flogo-lib/core/ext/model"
	"github.com/TIBCOSoftware/flogo-lib/core/process"
	"github.com/op/go-logging"
)

// log is the default package logger
var log = logging.MustGetLogger("simple-model")

func init() {
	m := model.New("simple")
	m.RegisterProcessBehavior(1, &SimpleProcessBehavior{})
	m.RegisterTaskBehavior(1, &SimpleTaskBehavior{})
	m.RegisterLinkBehavior(1, &SimpleLinkBehavior{})

	model.Register(m)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////

// SimpleProcessBehavior implements model.ProcessBehavior
type SimpleProcessBehavior struct {
}

// Start implements model.ProcessBehavior.Start
func (pb *SimpleProcessBehavior) Start(context model.ProcessContext, data interface{}) (start bool, evalCode int) {
	// just schedule the root task
	return true, 0
}

// Resume implements model.ProcessBehavior.Resume
func (pb *SimpleProcessBehavior) Resume(context model.ProcessContext, data interface{}) bool {
	return true
}

// TasksDone implements model.ProcessBehavior.TasksDone
func (pb *SimpleProcessBehavior) TasksDone(context model.ProcessContext, doneCode int) {
	// all tasks are done
}

// Done implements model.ProcessBehavior.Done
func (pb *SimpleProcessBehavior) Done(context model.ProcessContext) {
	log.Debugf("Process Done\n")
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

	linkContexts := context.FromLinks()

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
func (tb *SimpleTaskBehavior) Eval(context model.TaskContext, evalCode int) (done bool, doneCode int) {

	task := context.Task()
	log.Debugf("Task Eval: %s\n", task)

	if len(task.ChildTasks()) > 0 {
		log.Debugf("Has Children\n")

		//has children, so set to waiting
		context.SetState(STATE_WAITING)

		//for now enter all children (bpel style) - todo: change to enter leading chlidren
		context.EnterChildren(nil)

		return false, 0

	} else {

		activity, activityContext := context.Activity()

		if activity != nil {

			//log.Debug("Evaluating Activity: ", activity.GetType())
			done := activity.Eval(activityContext)
			return done, 0
		} else {

			//no-op
			return true, 0
		}
	}
}

// PostEval implements model.TaskBehavior.PostEval
func (tb *SimpleTaskBehavior) PostEval(context model.TaskContext, evalCode int, data interface{}) (done bool, doneCode int) {

	log.Debugf("Task PostEval\n")

	//activity, activityContext := context.Activity()
	activity, _ := context.Activity()

	if activity != nil { //if activity is async

		//done := activity.PostEval(activityContext, data)
		done := true
		return done, 0
	} else {

		//no-op
		return true, 0
	}
}

// Done implements model.TaskBehavior.Done
func (tb *SimpleTaskBehavior) Done(context model.TaskContext, doneCode int) (notifyParent bool, childDoneCode int, taskEntries []*model.TaskEntry) {

	task := context.Task()
	log.Debugf("done task:%s\n", task.Name())

	context.SetState(STATE_DONE)
	//context.SetTaskDone() for task garbage collection

	links := task.ToLinks()
	numLinks := len(links)

	// process outgoing links
	if numLinks > 0 {

		taskEntries := make([]*model.TaskEntry, 0, numLinks)

		for _, link := range links {

			linkContext := context.EvalLink(link, 0)
			if linkContext.State() == STATE_LINK_TRUE {

				taskEntry := &model.TaskEntry{Task: link.ToTask(), EnterCode: 0}
				taskEntries = append(taskEntries, taskEntry)
			}
		}

		//continue on to successor tasks
		return false, 0, taskEntries

	} else {
		// there are no outgoing links, so just notify parent that we are done
		return true, 0, nil
	}
}

// ChildDone implements model.TaskBehavior.ChildDone
func (tb *SimpleTaskBehavior) ChildDone(context model.TaskContext, childTask *process.Task, childDoneCode int) (done bool, doneCode int) {

	log.Debugf("Task ChildDone\n")

	// our children are done, so just transition ourselves to done
	return true, 0
}

////////////////////////////////////////////////////////////////////////////////////////////////////////

// SimpleLinkBehavior implements model.LinkBehavior
type SimpleLinkBehavior struct {
}

// Eval implements model.LinkBehavior.Eval
func (lb *SimpleLinkBehavior) Eval(context model.LinkContext, evalCode int) {

	context.SetState(STATE_LINK_TRUE)
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
)
