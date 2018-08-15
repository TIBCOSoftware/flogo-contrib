package filter

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"fmt"
)

// activityLogger is the default logger for the Filter Activity
var activityLogger = logger.GetLogger("activity-filter")

const (
	sType              = "type"
	sProceedOnlyOnEmit = "proceedOnlyOnEmit"
	ivValue            = "value"
	ovFiltered         = "filtered"
	ovValue            = "value"
)

//we can generate json from this! - we could also create a "validate-able" object from this
type Settings struct {
	Type              string `md:"required,allowed(non-zero)"`
	ProceedOnlyOnEmit bool
}

func init() {
	activityLogger.SetLogLevel(logger.InfoLevel)
}

var metadata *activity.Metadata

func New(config *activity.Config) (activity.Activity, error) {
	act := &FilterActivity{}

	filterType, _ := config.Settings[sType]

	if filterType == "non-zero" {
		act.filter = &NonZeroFilter{}
	} else {
		return nil, fmt.Errorf("unsupported filter: '%s'", filterType)
	}

	if proceedOnlyOnEmit, ok := config.Settings[sProceedOnlyOnEmit]; ok {
		act.proceedOnlyOnEmit = proceedOnlyOnEmit.(bool)
	}

	return act, nil
}

// FilterActivity is an Activity that is used to Filter a message to the console
type FilterActivity struct {
	filter            Filter
	proceedOnlyOnEmit bool
}

// NewActivity creates a new AppActivity
func NewActivity(md *activity.Metadata) activity.Activity {
	metadata = md
	activity.RegisterFactory(md.ID, New)
	return &FilterActivity{}
}

// Metadata returns the activity's metadata
func (a *FilterActivity) Metadata() *activity.Metadata {
	return metadata
}

// Eval implements api.Activity.Eval - Filters the Message
func (a *FilterActivity) Eval(ctx activity.Context) (done bool, err error) {

	filter := a.filter
	proceedOnlyOnEmit := a.proceedOnlyOnEmit

	if filter == nil {
		//backwards compatibility support

		settings, err := getSettings(ctx)
		if err != nil {
			return false, err
		}

		if settings.Type == "non-zero" {
			filter = &NonZeroFilter{}
		}

		proceedOnlyOnEmit = settings.ProceedOnlyOnEmit
	}

	in := ctx.GetInput(ivValue)

	filteredOut := filter.FilterOut(in)

	done = !(proceedOnlyOnEmit && filteredOut)

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
