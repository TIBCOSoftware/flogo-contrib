package timer

import (
	"context"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/carlescere/scheduler"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var log = logger.GetLogger("timer-trigger")

type TimerTrigger struct {
	metadata *trigger.Metadata
	runner   action.Runner
	config   *trigger.Config
	timers   map[string]*scheduler.Job
}

//NewFactory create a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &TimerFactory{metadata: md}
}

// TimerFactory Timer Trigger factory
type TimerFactory struct {
	metadata *trigger.Metadata
}

//New Creates a new trigger instance for a given id
func (t *TimerFactory) New(config *trigger.Config) trigger.Trigger {
	return &TimerTrigger{metadata: t.metadata, config: config}
}

// Metadata implements trigger.Trigger.Metadata
func (t *TimerTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Init implements ext.Trigger.Init
func (t *TimerTrigger) Init(runner action.Runner) {
	t.runner = runner
	//log.Infof("In init, id: '%s', Metadata: '%+v', Config: '%+v'", t.config.Id, t.metadata, t.config)
}

// Start implements ext.Trigger.Start
func (t *TimerTrigger) Start() error {

	log.Debug("Start")
	t.timers = make(map[string]*scheduler.Job)
	handlers := t.config.Handlers

	log.Debug("Processing handlers")
	for _, handler := range handlers {

		repeating := handler.Settings["repeating"]
		log.Debug("Repeating: ", repeating)
		if repeating == "false" {
			t.scheduleOnce(handler)
		} else if repeating == "true" {
			t.scheduleRepeating(handler)
		} else {
			log.Error("No match for repeating: ", repeating)
		}
		log.Debug("Settings repeating: ", handler.Settings["repeating"])
		log.Debugf("Processing Handler: %s", handler.ActionId)
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

func (t *TimerTrigger) scheduleOnce(endpoint *trigger.HandlerConfig) {
	log.Info("Scheduling a run one time job")

	seconds := getInitialStartInSeconds(endpoint)
	log.Debug("Seconds until trigger fires: ", seconds)

	var timerJob *scheduler.Job

	if seconds != 0 {
		timerJob = scheduler.Every(int(seconds))
	} else {
		log.Debug("Start Date not specified, executing action immediately")
	}

	fn := func() {
		log.Debug("Executing \"Once\" timer trigger")

		act := action.Get(endpoint.ActionId)

		if act != nil {
			log.Debugf("Running action: %s", endpoint.ActionId)

			_, err := t.runner.RunAction(context.Background(), act, nil )

			if err != nil {
				log.Error("Error running action: ", err.Error())
			}
		} else {
			log.Errorf("Action '%s' not found", endpoint.ActionId)
		}

		if timerJob != nil {
			timerJob.Quit <- true
		}
	}

	if seconds == 0 {
		fn()
	} else {
		timerJob, err := timerJob.Seconds().NotImmediately().Run(fn)
		if err != nil {
			log.Error("Error performing scheduleOnce: ", err.Error())
		}

		t.timers[endpoint.ActionId] = timerJob
	}
}

func (t *TimerTrigger) scheduleRepeating(endpoint *trigger.HandlerConfig) {
	log.Info("Scheduling a repeating job")

	seconds := getInitialStartInSeconds(endpoint)

	fn2 := func() {
		log.Debug("-- Starting \"Repeating\" (repeat) timer action")

		act := action.Get(endpoint.ActionId)
		log.Debugf("Found action: '%+x'", act)
		log.Debugf("ActionID: '%s'", endpoint.ActionId)
		_, _, err := t.runner.Run(context.Background(), act, endpoint.ActionId, nil)

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

func getInitialStartInSeconds(handlerCfg *trigger.HandlerConfig) int {

	if _, ok := handlerCfg.Settings["startDate"]; !ok {
		return 0
	}

	layout := time.RFC3339
	startDate := handlerCfg.GetSetting("startDate")

	if startDate == "" {
		return 0
	}

	idx := strings.LastIndex(startDate, "Z")
	timeZone := startDate[idx+1: ]
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

func (t *TimerTrigger) scheduleJobEverySecond(handlerCfg *trigger.HandlerConfig, fn func()) {

	var interval int = 0
	if seconds := handlerCfg.GetSetting("seconds"); seconds != "" {
		seconds, _ := strconv.Atoi(seconds)
		interval = interval + seconds
	}
	if minutes := handlerCfg.GetSetting("minutes"); minutes != "" {
		minutes, _ := strconv.Atoi(minutes)
		interval = interval + minutes*60
	}
	if hours := handlerCfg.GetSetting("hours"); hours != "" {
		hours, _ := strconv.Atoi(hours)
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

	t.timers["r:"+handlerCfg.ActionId] = timerJob
}
