package simple

import (
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/definition"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/model"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// log is the default package logger
var log = logger.GetLogger("model-tibco-simple")

const (
	MODEL_NAME = "tibco-simple"
)

func init() {
	model.Register(New())
}

func New() *model.FlowModel {
	m := model.New(MODEL_NAME)
	m.RegisterFlowBehavior(&SimpleFlowBehavior{})
	m.RegisterDefaultTaskBehavior(&SimpleTaskBehavior{})
	m.RegisterTaskBehavior(2, "iterator", &SimpleIteratorTaskBehavior{})
	return m
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
// Flow Behavior

// SimpleFlowBehavior implements model.FlowBehavior
type SimpleFlowBehavior struct {
}

// Start implements model.FlowBehavior.Start
func (pb *SimpleFlowBehavior) Start(ctx model.FlowContext) (start bool, taskEntries []*model.TaskEntry) {
	// just schedule the root task
	return true, GetFlowTaskEntries(ctx, false)
}

// Resume implements model.FlowBehavior.Resume
func (pb *SimpleFlowBehavior) Resume(ctx model.FlowContext) bool {
	return true
}

// TasksDone implements model.FlowBehavior.TasksDone
//todo handler error flow
func (pb *SimpleFlowBehavior) TaskDone(ctx model.FlowContext, doneCode int) (done bool) {
	tasks := ctx.TaskInsts()

	for _, taskInst := range tasks {

		if taskInst.State() < STATE_DONE {

			log.Debugf("task %s not done or skipped", taskInst.Task().Name())
			return false
		}
	}

	log.Debug("all tasks done or skipped")

	// our tasks are done, so the flow is done
	return true
}

// Done implements model.FlowBehavior.Done
func (pb *SimpleFlowBehavior) Done(ctx model.FlowContext) {
	log.Debugf("Flow Done\n")
}

func GetFlowTaskEntries(ctx model.FlowContext, leadingOnly bool) []*model.TaskEntry {

	var taskEntries []*model.TaskEntry

	for _, task := range ctx.FlowDefinition().GetTasks() {

		if len(task.FromLinks()) == 0 || !leadingOnly {

			taskEntry := &model.TaskEntry{Task: task, EnterCode: 0}
			taskEntries = append(taskEntries, taskEntry)
		}
	}

	return taskEntries
}

func (pb *SimpleFlowBehavior) StartEmbeddedFlow(context model.FlowContext, flow *definition.Definition) (start bool, taskEntries []*model.TaskEntry) {
	return false, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////

// SimpleTaskBehavior implements model.TaskBehavior
type SimpleTaskBehavior struct {
}

// Enter implements model.TaskBehavior.Enter
func (tb *SimpleTaskBehavior) Enter(ctx model.TaskContext, enterCode int) (eval bool, evalCode int) {

	task := ctx.Task()
	log.Debugf("Task Enter: %s\n", task.Name())

	ctx.SetStatus(STATE_ENTERED)

	//check if all predecessor links are done

	linkContexts := ctx.FromInstLinks()

	ready := true
	skipped := false

	if len(linkContexts) == 0 {
		// has no predecessor links, so task is ready
		ready = true
	} else {
		skipped = true

		log.Debugf("Num Links: %d\n", len(linkContexts))
		for _, linkContext := range linkContexts {

			log.Debugf("Task: %s, linkData: %v\n", task.Name(), linkContext)
			if linkContext.Status() < STATE_LINK_FALSE {
				ready = false
				break
			} else if linkContext.Status() == STATE_LINK_TRUE {
				skipped = false
			}
		}
	}

	if ready {

		if skipped {
			log.Debugf("Task Skipped\n")
			ctx.SetStatus(STATE_SKIPPED)
			//todo hack, wait for explicit skip support from engine
			return ready, -666
		} else {
			log.Debugf("Task Ready\n")
			ctx.SetStatus(STATE_READY)
		}

	} else {
		log.Debugf("Task Not Ready\n")
	}

	return ready, 0
}

// Eval implements model.TaskBehavior.Eval
func (tb *SimpleTaskBehavior) Eval(ctx model.TaskContext, evalCode int) (evalResult model.EvalResult, doneCode int, err error) {

	if ctx.Status() == STATE_SKIPPED {
		//todo introduce EVAL_SKIPPED?
		return model.EVAL_DONE, 0, nil
	}

	task := ctx.Task()
	log.Debugf("Task Eval: %v\n", task)

	if ctx.HasActivity() {

		done, err := ctx.EvalActivity()

		if err != nil {
			log.Errorf("Error evaluating activity '%s'[%s] - %s", ctx.Task().Name(), ctx.Task().ActivityConfig().Ref(), err.Error())
			ctx.SetStatus(STATE_FAILED)
			return model.EVAL_FAIL, 0, err
		}

		if done {
			evalResult = model.EVAL_DONE
		} else {
			evalResult = model.EVAL_WAIT
		}

		return evalResult, 0, nil
	}

	//no-op
	return model.EVAL_DONE, 0, nil
}

// PostEval implements model.TaskBehavior.PostEval
func (tb *SimpleTaskBehavior) PostEval(ctx model.TaskContext, evalCode int, data interface{}) (done bool, doneCode int, err error) {

	log.Debugf("Task PostEval\n")

	if ctx.HasActivity() { //if activity is async

		//done := activity.PostEval(activityContext, data)
		done := true
		return done, 0, nil
	}

	//no-op
	return true, 0, nil
}

// Done implements model.TaskBehavior.Done
func (tb *SimpleTaskBehavior) Done(ctx model.TaskContext, doneCode int) (notifiyFlow bool, notifyCode int, taskEntries []*model.TaskEntry, err error) {

	task := ctx.Task()

	linkInsts := ctx.ToInstLinks()
	numLinks := len(linkInsts)

	if ctx.Status() == STATE_SKIPPED {
		log.Debugf("skipped task: %s\n", task.Name())

		// skip outgoing links
		if numLinks > 0 {

			taskEntries = make([]*model.TaskEntry, 0, numLinks)
			for _, linkInst := range linkInsts {

				linkInst.SetStatus(STATE_LINK_SKIPPED)

				//todo: engine should not eval mappings for skipped tasks, skip
				//todo: needs to be a state/op understood by the engine
				taskEntry := &model.TaskEntry{Task: linkInst.Link().ToTask(), EnterCode: EC_SKIP}
				taskEntries = append(taskEntries, taskEntry)
			}

			//continue on to successor tasks
			return false, 0, taskEntries, nil
		}
	} else {
		log.Debugf("done task: %s", task.Name())

		ctx.SetStatus(STATE_DONE)
		//ctx.SetTaskDone() for task garbage collection

		// process outgoing links
		if numLinks > 0 {

			taskEntries = make([]*model.TaskEntry, 0, numLinks)

			for _, linkInst := range linkInsts {

				follow := true

				if linkInst.Link().Type() == definition.LtError {
					//todo should we skip or ignore?
					continue
				}

				if linkInst.Link().Type() == definition.LtExpression {
					//todo handle error
					follow, err = ctx.EvalLink(linkInst.Link())

					if err != nil {
						return false, 0, nil, err
					}
				}

				if follow {
					linkInst.SetStatus(STATE_LINK_TRUE)

					taskEntry := &model.TaskEntry{Task: linkInst.Link().ToTask(), EnterCode: 0}
					taskEntries = append(taskEntries, taskEntry)
				} else {
					linkInst.SetStatus(STATE_LINK_FALSE)

					taskEntry := &model.TaskEntry{Task: linkInst.Link().ToTask(), EnterCode: EC_SKIP}
					taskEntries = append(taskEntries, taskEntry)
				}
			}

			//continue on to successor tasks
			return false, 0, taskEntries, nil
		}
	}

	log.Debug("notifying flow that task is done")

	// there are no outgoing links, so just notify parent that we are done
	return true, 0, nil, nil
}

// Done implements model.TaskBehavior.Error
func (tb *SimpleTaskBehavior) Error(ctx model.TaskContext) (handled bool, taskEntry *model.TaskEntry) {

	linkInsts := ctx.ToInstLinks()
	numLinks := len(linkInsts)

	// process outgoing links
	if numLinks > 0 {

		for _, linkInst := range linkInsts {

			if linkInst.Link().Type() == definition.LtError {
				linkInst.SetStatus(STATE_LINK_TRUE)
				taskEntry := &model.TaskEntry{Task: linkInst.Link().ToTask(), EnterCode: 0}
				return true, taskEntry
			}
		}
	}

	// there are no outgoing error links, so just return false
	return false, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
// Status
const (
	EC_SKIP = 1

	STATE_NOT_STARTED int = 0

	STATE_LINK_FALSE   int = 1
	STATE_LINK_TRUE    int = 2
	STATE_LINK_SKIPPED int = 3

	STATE_ENTERED int = 10
	STATE_READY   int = 20
	STATE_WAITING int = 30
	STATE_DONE    int = 40
	STATE_SKIPPED int = 50
	STATE_FAILED  int = 100
)
