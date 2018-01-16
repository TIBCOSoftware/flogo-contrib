package simple

import (
	"fmt"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/definition"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/model"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	//"github.com/TIBCOSoftware/flogo-lib/logger"
)

// SimpleIteratorTaskBehavior implements model.TaskBehavior
type SimpleIteratorTaskBehavior struct {
}

// Enter implements model.TaskBehavior.Enter
func (tb *SimpleIteratorTaskBehavior) Enter(ctx model.TaskContext, enterCode int) (eval bool, evalCode int) {

	//todo inherit this code from base task

	task := ctx.Task()
	log.Debugf("Task Enter: %s\n", task.Name())

	ctx.SetState(STATE_ENTERED)

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

			if linkContext.State() < STATE_LINK_FALSE {
				ready = false
				break
			} else if linkContext.State() == STATE_LINK_TRUE {
				skipped = false
			}
		}
	}

	if ready {

		if skipped {
			log.Debugf("Task Skipped\n")
			ctx.SetState(STATE_SKIPPED)
			//todo hack, wait for explicit skip support from engine
			return ready, -666
		} else {
			log.Debugf("Task Ready\n")
			ctx.SetState(STATE_READY)
		}

	} else {
		log.Debugf("Task Not Ready\n")
	}

	return ready, 0
}

type Iteration struct {
	Key interface{}
	Value interface{}
}

// Eval implements model.TaskBehavior.Eval
func (tb *SimpleIteratorTaskBehavior) Eval(ctx model.TaskContext, evalCode int) (evalResult model.EvalResult, doneCode int, err error) {

	if ctx.State() == STATE_SKIPPED {
		return model.EVAL_DONE, EC_SKIP, nil
	}

	task := ctx.Task()
	log.Debugf("Task Eval: %v\n", task)

	if ctx.HasActivity() {

		var itx Iterator

		itxAttr, ok := ctx.GetWorkingData("_iterator")
		iterationAttr, _ := ctx.GetWorkingData("iteration")

		if ok {
			itx = itxAttr.Value().(Iterator)
		} else {

			iterateOn, ok := ctx.GetSetting("iterate")

			if !ok {
				//todo if iterateOn is not defined, what should we do?
				//just skip for now
				return model.EVAL_DONE, 0, nil
			}

			switch t := iterateOn.(type) {
			case string:
				count, err := data.CoerceToInteger(iterateOn)
				if err != nil {
					return model.EVAL_FAIL, 0, err
				}
				itx = NewIntIterator(count)
			case int:
				count := iterateOn.(int)
				itx = NewIntIterator(count)
			case map[string]interface{}:
				itx = NewObjectIterator(t)
			case []interface{}:
				itx = NewArrayIterator(t)
			default:
				return model.EVAL_FAIL, 0, fmt.Errorf("unsupported type '%s' for iterateOn", t)
			}

			itxAttr, _ = data.NewAttribute("_iterator", data.ANY, itx)
			ctx.AddWorkingData(itxAttr)

			iteration := map[string]interface{}{
				"key": nil,
				"value":   nil,
			}

			iterationAttr, _ = data.NewAttribute("iteration", data.OBJECT, iteration)
			ctx.AddWorkingData(iterationAttr)
		}

		repeat := itx.next()

		if repeat {
			log.Debugf("Repeat:%s, Key:%s, Value:%v", repeat, itx.Key(), itx.Value())

			iteration,_ := iterationAttr.Value().(map[string]interface{})
			iteration["key"] = itx.Key()
			iteration["value"] = itx.Value()

			_, err := ctx.EvalActivity()

			//what to do if eval isn't "done"?
			if err != nil {
				log.Errorf("Error evaluating activity '%s'[%s] - %s", ctx.Task().Name(), ctx.Task().ActivityConfig().Ref(), err.Error())
				ctx.SetState(STATE_FAILED)
				return model.EVAL_FAIL, 0, err
			}

			evalResult = model.EVAL_REPEAT
		} else {
			evalResult = model.EVAL_DONE
		}

		return evalResult, 0, nil
	}

	//no-op
	return model.EVAL_DONE, 0, nil
}

// PostEval implements model.TaskBehavior.PostEval
func (tb *SimpleIteratorTaskBehavior) PostEval(ctx model.TaskContext, evalCode int, data interface{}) (done bool, doneCode int, err error) {

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
func (tb *SimpleIteratorTaskBehavior) Done(ctx model.TaskContext, doneCode int) (notifyParent bool, childDoneCode int, taskEntries []*model.TaskEntry, err error) {

	task := ctx.Task()

	linkInsts := ctx.ToInstLinks()
	numLinks := len(linkInsts)

	if ctx.State() == STATE_SKIPPED {
		log.Debugf("skipped task: %s\n", task.Name())

		// skip outgoing links
		if numLinks > 0 {

			taskEntries = make([]*model.TaskEntry, 0, numLinks)
			for _, linkInst := range linkInsts {

				linkInst.SetState(STATE_LINK_SKIPPED)

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

		ctx.SetState(STATE_DONE)
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
					linkInst.SetState(STATE_LINK_TRUE)

					taskEntry := &model.TaskEntry{Task: linkInst.Link().ToTask(), EnterCode: 0}
					taskEntries = append(taskEntries, taskEntry)
				} else {
					linkInst.SetState(STATE_LINK_FALSE)

					taskEntry := &model.TaskEntry{Task: linkInst.Link().ToTask(), EnterCode: EC_SKIP}
					taskEntries = append(taskEntries, taskEntry)
				}
			}

			//continue on to successor tasks
			return false, 0, taskEntries, nil
		}
	}

	log.Debug("notifying parent that task is done")

	// there are no outgoing links, so just notify parent that we are done
	return true, 0, nil, nil
}

// Done implements model.TaskBehavior.Error
func (tb *SimpleIteratorTaskBehavior) Error(ctx model.TaskContext) (handled bool, taskEntry *model.TaskEntry) {

	linkInsts := ctx.ToInstLinks()
	numLinks := len(linkInsts)

	// process outgoing links
	if numLinks > 0 {

		for _, linkInst := range linkInsts {

			if linkInst.Link().Type() == definition.LtError {
				linkInst.SetState(STATE_LINK_TRUE)
				taskEntry := &model.TaskEntry{Task: linkInst.Link().ToTask(), EnterCode: 0}
				return true, taskEntry
			}
		}
	}

	// there are no outgoing error links, so just return false
	return false, nil
}

///////////////////////////////////
// Iterators

type Iterator interface {
	Key() interface{}
	Value() interface{}
	next() bool
}

type ArrayIterator struct {
	current int
	data    []interface{}
}

func (itx *ArrayIterator) Key() interface{} {
	return itx.current
}

func (itx *ArrayIterator) Value() interface{} {
	return itx.data[itx.current]
}
func (itx *ArrayIterator) next() bool {
	itx.current++
	if itx.current >= len(itx.data) {
		return false
	}
	return true
}

func NewArrayIterator(data []interface{}) *ArrayIterator {
	return &ArrayIterator{data: data, current: -1}
}

type IntIterator struct {
	current int
	count   int
}

func (itx *IntIterator) Key() interface{} {
	return itx.current
}

func (itx *IntIterator) Value() interface{} {
	return itx.current
}

func (itx *IntIterator) next() bool {
	itx.current++
	if itx.current >= itx.count {
		return false
	}
	return true
}

func NewIntIterator(count int) *IntIterator {
	return &IntIterator{count: count, current: -1}
}

type ObjectIterator struct {
	current int
	keyMap  map[int]string
	data    map[string]interface{}
}

func (itx *ObjectIterator) Key() interface{} {
	return itx.keyMap[itx.current]
}

func (itx *ObjectIterator) Value() interface{} {
	key := itx.keyMap[itx.current]
	return itx.data[key]
}

func (itx *ObjectIterator) next() bool {
	itx.current++
	if itx.current >= len(itx.data) {
		return false
	}
	return true
}

func NewObjectIterator(data map[string]interface{}) *ObjectIterator {
	keyMap := make(map[int]string, len(data))
	i := 0
	for key := range data {
		keyMap[i] = key
		i++
	}

	return &ObjectIterator{keyMap: keyMap, data: data, current: -1}
}

