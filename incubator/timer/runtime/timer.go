package timer

import (
	"context"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/types"
	"github.com/carlescere/scheduler"
	"github.com/op/go-logging"
	"math"
	"strconv"
	"strings"
	"time"
)

const (
	TRIGGER_REF = "github.com/TIBCOSoftware/flogo-contrib/incubator/timer/runtime"
)

// log is the default package logger
var log = logging.MustGetLogger("trigger-tibco-timer")
var md = trigger.NewMetadata(jsonMetadata)

type TimerTrigger struct {
	metadata *trigger.Metadata
	runner   action.Runner
	config   types.TriggerConfig
	timers   map[string]*scheduler.Job
	myId     string
}

type TimerFactory struct {
}

func init() {
	trigger.RegisterFactory(TRIGGER_REF, &TimerFactory{})
}

// Metadata implements trigger.Trigger.Metadata
func (t *TimerTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

func (t *TimerFactory) New(id string) trigger.Trigger2 {
	return &TimerTrigger{metadata: md, myId: id}
}

// Init implements ext.Trigger.Init
func (t *TimerTrigger) Init(config types.TriggerConfig, runner action.Runner) {
	log.Infof("In init, id: '%s', Metadata: '%+v', Config: '%+v'", t.myId, t.metadata, config)
	t.config = config
	t.runner = runner
}

// Start implements ext.Trigger.Start
func (t *TimerTrigger) Start() error {

	log.Debug("Start")
	t.timers = make(map[string]*scheduler.Job)
	endpoints := t.config.Handlers

	log.Debug("Processing endpoints")
	for _, endpoint := range endpoints {

		repeating := endpoint.Settings["repeating"]
		log.Debug("Repeating: ", repeating)
		if repeating == "false" {
			t.scheduleOnce(endpoint)
		} else if repeating == "true" {
			t.scheduleRepeating(endpoint)
		} else {
			log.Error("No match for repeating: ", repeating)
		}
		log.Debug("Settings repeating: ", endpoint.Settings["repeating"])
		log.Debugf("Processing Handler: %s", endpoint.ActionId)
	}

	return nil
}

// Stop implements ext.Trigger.Stop
func (t *TimerTrigger) Stop() error {

	log.Debug("Stopping endpoints")
	for k, v := range t.timers {
		if t.timers[k].IsRunning() {
			log.Debug("Stopping timer for : ", k)
			v.Quit <- true
		} else {
			log.Debugf("Timer: %s is not running", k)
		}
	}

	return nil
}

func (t *TimerTrigger) scheduleOnce(endpoint *types.TriggerHandler) {
	log.Info("Scheduling a run one time job")

	seconds := getInitialStartInSeconds(endpoint)
	log.Debug("Seconds till trigger fires: ", seconds)
	timerJob := scheduler.Every(int(seconds))

	if timerJob == nil {
		log.Error("timerJob is nil")
	}

	fn := func() {
		log.Debug("-- Starting \"Once\" timer process")

		action := action.Get2(endpoint.ActionId)
		log.Infof("Found action: '%+x'", action)
		log.Infof("ActionID: '%s'", endpoint.ActionId)
		_, _, err := t.runner.Run(context.Background(), action, endpoint.ActionId, nil)

		if err != nil {
			log.Error("Error starting action: ", err.Error())
		}
		timerJob.Quit <- true
	}

	timerJob, err := timerJob.Seconds().NotImmediately().Run(fn)
	if err != nil {
		log.Error("Error scheduleOnce flo err: ", err.Error())
	}

	t.timers[endpoint.ActionId] = timerJob
}

func (t *TimerTrigger) scheduleRepeating(endpoint *types.TriggerHandler) {
	log.Info("Scheduling a repeating job")

	seconds := getInitialStartInSeconds(endpoint)

	fn2 := func() {
		log.Debug("-- Starting \"Repeating\" (repeat) timer action")

		action := action.Get2(endpoint.ActionId)
		log.Infof("Found action: '%+x'", action)
		log.Infof("ActionID: '%s'", endpoint.ActionId)
		_, _, err := t.runner.Run(context.Background(), action, endpoint.ActionId, nil)

		if err != nil {
			log.Error("Error starting flow: ", err.Error())
		}
	}

	if endpoint.Settings["notImmediate"] == "false" {
		t.scheduleJobEverySecond(endpoint, fn2)
	} else {

		log.Debug("Seconds till trigger fires: ", seconds)
		timerJob := scheduler.Every(seconds)
		if timerJob == nil {
			log.Error("timerJob is nil")
		}

		t.scheduleJobEverySecond(endpoint, fn2)

		timerJob, err := timerJob.Seconds().NotImmediately().Run(fn2)
		if err != nil {
			log.Error("Error scheduleRepeating (first) flo err: ", err.Error())
		}
		if timerJob == nil {
			log.Error("timerJob is nil")
		}

		t.timers[endpoint.ActionId] = timerJob
	}
}

func getInitialStartInSeconds(endpoint *types.TriggerHandler) int {

	if startDate := endpoint.Settings["startDate"]; startDate == nil {
		return 0
	}

	layout := time.RFC3339
	startDate := endpoint.Settings["startDate"].(string)
	idx := strings.LastIndex(startDate, "Z")
	timeZone := startDate[idx+1 : len(startDate)]
	log.Debug("Time Zone: ", timeZone)
	startDate = strings.TrimSuffix(startDate, timeZone)
	log.Debug("startDate: ", startDate)

	// is timezone negative
	var isNegative bool
	isNegative = strings.HasPrefix(timeZone, "-")
	// remove sign
	timeZone = strings.TrimPrefix(timeZone, "-")

	triggerDate, err := time.Parse(layout, startDate)
	if err != nil {
		log.Error("Error parsing time err: ", err.Error())
	}
	log.Debug("Time parsed from settings: ", triggerDate)

	var hour int
	var minutes int

	sliceArray := strings.Split(timeZone, ":")
	if len(sliceArray) != 2 {
		log.Error("Time zone has wrong format: ", timeZone)
	} else {
		hour, _ = strconv.Atoi(sliceArray[0])
		minutes, _ = strconv.Atoi(sliceArray[1])

		log.Debug("Duration hour: ", time.Duration(hour)*time.Hour)
		log.Debug("Duration minutes: ", time.Duration(minutes)*time.Minute)
	}

	hours, _ := strconv.Atoi(timeZone)
	log.Debug("hours: ", hours)
	if isNegative {
		log.Debug("Adding to triggerDate")
		triggerDate = triggerDate.Add(time.Duration(hour) * time.Hour)
		triggerDate = triggerDate.Add(time.Duration(minutes) * time.Minute)
	} else {
		log.Debug("Subtracting to triggerDate")
		triggerDate = triggerDate.Add(time.Duration(hour * -1))
		triggerDate = triggerDate.Add(time.Duration(minutes))
	}

	currentTime := time.Now().UTC()
	log.Debug("Current time: ", currentTime)
	log.Debug("Setting start time: ", triggerDate)
	duration := time.Since(triggerDate)

	return int(math.Abs(duration.Seconds()))
}

type PrintJob struct {
	Msg string
}

func (j *PrintJob) Run() error {
	log.Debug(j.Msg)
	return nil
}

func (t *TimerTrigger) scheduleJobEverySecond(endpoint *types.TriggerHandler, fn func()) {

	var interval int = 0
	if seconds := endpoint.Settings["seconds"]; seconds != nil {
		seconds, _ := strconv.Atoi(seconds.(string))
		interval = interval + seconds
	}
	if minutes := endpoint.Settings["minutes"]; minutes != nil {
		minutes, _ := strconv.Atoi(minutes.(string))
		interval = interval + minutes*60
	}
	if hours := endpoint.Settings["hours"]; hours != nil {
		hours, _ := strconv.Atoi(hours.(string))
		interval = interval + hours*3600
	}

	log.Debug("Repeating seconds: ", interval)
	// schedule repeating
	timerJob, err := scheduler.Every(interval).Seconds().Run(fn)
	if err != nil {
		log.Error("Error scheduleRepeating (repeat seconds) flo err: ", err.Error())
	}
	if timerJob == nil {
		log.Error("timerJob is nil")
	}

	t.timers["r:"+endpoint.ActionId] = timerJob
}
