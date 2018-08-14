package filter

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

// activityLogger is the default logger for the Filter Activity
var activityLogger = logger.GetLogger("activity-filter")

const (
	sType              = "type"
	sProceedOnlyOnEmit = "proceedOnlyOnEmit"

	ivValue = "value"

	ovFiltered = "filtered"
	ovValue    = "value"
)

//we can generate json from this! - we could also create a "validate-able" object from this
type Settings struct {
	Type              string `md:"required,allowed(non-zero)"`
	ProceedOnlyOnEmit bool
}

func init() {
	activityLogger.SetLogLevel(logger.InfoLevel)
}

// FilterActivity is an Activity that is used to Filter a message to the console
type FilterActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new AppActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &FilterActivity{metadata: metadata}
}

// Metadata returns the activity's metadata
func (a *FilterActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Filters the Message
func (a *FilterActivity) Eval(ctx activity.Context) (done bool, err error) {

	//todo move to Activity instance creation
	settings, err := getSettings(ctx)
	if err != nil {
		return false, err
	}

	var filter Filter
	if settings.Type == "non-zero" {
		filter = &NonZeroFilter{}
	}

	in := ctx.GetInput(ivValue)

	filteredOut := filter.FilterOut(in)

	done = !(settings.ProceedOnlyOnEmit && filteredOut)

	ctx.SetOutput(ovFiltered, filteredOut)
	ctx.SetOutput(ovValue, in)

	return done, nil
}

func getSettings(ctx activity.Context) (*Settings, error) {

	settings := &Settings{}

	settings.Type = "non-zero" // default function
	setting, exists := ctx.GetSetting(sType)
	if exists {
		val, err := data.CoerceToString(setting)
		if err == nil {
			settings.Type = val
		}
	}

	settings.ProceedOnlyOnEmit = true // by default only proceed on emit
	setting, exists = ctx.GetSetting(sProceedOnlyOnEmit)
	if exists {
		val, err := data.CoerceToBoolean(setting)
		if err == nil {
			settings.ProceedOnlyOnEmit = val
		}
	}

	// settings validation can be done here once activities are created on configuration instead of
	// setting up during runtime

	return settings, nil
}

type Filter interface {
	FilterOut(val interface{}) bool
}
