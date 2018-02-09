package histocompare

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"fmt"
	"sync"
	"math"
)

// List of input and output variables names
const (
	ivVarName     = "varName"
	ivVarValue = "varValue"
	ivThreshold   = "threshold"
	ivThresholdUnit            = "thresholdUnit"
	ivStoreIfInRange          = "storeIfInRange"
	ivStoreIfExceed            = "storeIfExceed"
	ovPrevStoredValue              = "prevStoredValue"
	ovExceedThreshold         = "exceedThreshold"
)

// log is the default package logger
var log = logger.GetLogger("activity-tibco-histocompare")

// HistoCompareActivity is a stub for your Activity implementation
type HistoCompareActivity struct {
	metadata *activity.Metadata
	sync.Mutex
	storedVars map[string]float64
}

// NewActivity creates a new activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &HistoCompareActivity{metadata: metadata, storedVars: make(map[string]float64)}
}

// Metadata implements activity.Activity.Metadata
func (a *HistoCompareActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *HistoCompareActivity) Eval(context activity.Context) (done bool, err error)  {

	if context.GetInput(ivVarName) == nil || context.GetInput(ivVarValue) == nil || context.GetInput(ivThreshold) == nil {
		log.Error("Required variables have not been set !")
		return false, fmt.Errorf("required variables have not been set")
	}

	
	VarName := context.GetInput(ivVarName).(string)
	VarValue := context.GetInput(ivVarValue).(float64)
	Threshold := context.GetInput(ivThreshold).(float64)
	ThresholdUnit := context.GetInput(ivThresholdUnit).(string)
	StoreIfInRange := context.GetInput(ivStoreIfInRange).(bool)
	StoreIfExceed := context.GetInput(ivStoreIfExceed).(bool)
	
	log.Debugf("Compare Histo [Variable = %s, Value = %s, Threshold = %s, Threshold unit = %s, StoreIfInRange = %s, StoreIfExceed = %s]", VarName, VarValue, Threshold, ThresholdUnit, StoreIfInRange, StoreIfExceed)
	
	storedValue, exceedThreshold, err := a.compareHistoValue(VarName, VarValue, Threshold, ThresholdUnit, StoreIfInRange, StoreIfExceed)
	
	if err != nil {
		return false, err
	}
	context.SetOutput(ovPrevStoredValue, storedValue)
	context.SetOutput(ovExceedThreshold, exceedThreshold)
	return true, nil
}


func (a *HistoCompareActivity) compareHistoValue(varName string, varNewValue float64, threshold float64, thresholdUnit string, StoreIfInRange bool, StoreIfExceed bool) (storedValue float64, exceedThreshold bool, err error)  {
	a.Lock()
	defer a.Unlock()

	exceedThreshold = false
	
	storedValue := varNewValue
	
	if contains(a.storedVars,varName) {
		storedValue = a.storedVars[varName]
		log.Debugf("Variable [%s] is already stored with value [%v]", varName, storedValue)
	} else {
		a.storedVars[varName] = storedValue
		log.Debugf("Variable [%s] didn't exist. Storing it with value [%v]", varName, storedValue)
	}

	if thresholdUnit == "%" {
		threshold = storedValue * (threshold / 100)
	}

	if  math.Abs(varNewValue - storedValue) > threshold {
		log.Debugf("Value [%v] exceed threshold (Stored value = [%v])", varNewValue, storedValue)
		exceedThreshold = true
		if StoreIfExceed {
			a.storedVars[varName] = varNewValue
		} 
	} else {
		log.Debugf("Value [%v] in range.", varNewValue)
	}

	if StoreIfInRange {
		a.storedVars[varName] = varNewValue
	}
	return storedValue, exceedThreshold, nil
}