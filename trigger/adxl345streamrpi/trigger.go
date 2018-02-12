package adxl345streamrpi

import (
	"context"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/rpi" // Mandatory for embd
	"strconv"
	"time"
)

// log is the default package logger
var log = logger.GetLogger("trigger-adxl345-rpi")

var interval = 500

// ADXL345Factory My Trigger factory
type ADXL345Factory struct {
	metadata *trigger.Metadata
}

//NewFactory create a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &ADXL345Factory{metadata: md}
}

//New Creates a new trigger instance for a given id
func (t *ADXL345Factory) New(config *trigger.Config) trigger.Trigger {
	return &ADXL345Trigger{metadata: t.metadata, config: config}
}

// ADXL345Trigger is a stub for your Trigger implementation
type ADXL345Trigger struct {
	metadata *trigger.Metadata
	runner   action.Runner
	config   *trigger.Config
}

func doEvery(d time.Duration, f func()) {
	for _ = range time.Tick(d) {
		f()
	}
}

// Init implements trigger.Trigger.Init
func (t *ADXL345Trigger) Init(runner action.Runner) {
	t.runner = runner

	if t.config.Settings == nil {
		log.Infof("No configuration set for the trigger... Using default configuration...")
	} else {
		if t.config.Settings["delay_ms"] != nil && t.config.Settings["delay_ms"] != "" {
			interval, _ = strconv.Atoi(t.config.GetSetting("delay_ms"))
		} else {
			log.Infof("No delay has been set. Using default value (", interval, "ms)")
		}
	}

	log.Infof("In init, id: '%s', Metadata: '%+v', Config: '%+v'", t.config.Id, t.metadata, t.config)
}

// Metadata implements trigger.Trigger.Metadata
func (t *ADXL345Trigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Start implements trigger.Trigger.Start
func (t *ADXL345Trigger) Start() error {
	// start the trigger
	log.Debug("Start Trigger ADXL345Stream for Raspberry PI")
	handlers := t.config.Handlers
	//t.timers = make(map[string]*scheduler.Job)

	log.Debug("Processing handlers")
	for _, handler := range handlers {
		t.scheduleRepeating(handler)
		log.Debugf("Processing Handler: %s", handler.ActionId)
	}
	return nil
}

// Stop implements trigger.Trigger.Start
func (t *ADXL345Trigger) Stop() error {
	// stop the trigger
	return nil
}

func (t *ADXL345Trigger) scheduleRepeating(endpoint *trigger.HandlerConfig) {

	log.Debug("Repeating every ", interval, "ms")

	fn2 := func() {
		act := action.Get(endpoint.ActionId)

		x, y, z,  err := t.getDataFromSensor(endpoint)
		if err != nil {
			log.Error("Error while reading sensor data. Err: ", err.Error())
		}

		data := make(map[string]interface{})
		data["x"] = x
		data["y"] = y
		data["z"] = z

		log.Debug("Acceleration: [x=", data["x"], " , y=", data["y"], ", z=", data["z"], "]")
		startAttrs, err := t.metadata.OutputsToAttrs(data, true)

		if err != nil {
			log.Errorf("After run error' %s'\n", err)
		}

		ctx := trigger.NewContext(context.Background(), startAttrs)
		results, err := t.runner.RunAction(ctx, act, nil)

		if err != nil {
			log.Errorf("An error occured while starting the flow. Err:", err)
		}
		log.Info("Exec: ", results)
	}

	// schedule repeating
	doEvery(time.Duration(interval)*time.Millisecond, fn2)
}

func (t *ADXL345Trigger) getDataFromSensor(endpoint *trigger.HandlerConfig) (x float64, y float64, z float64, err error) {

	opt := NewOpt()
	adxl345, errNew := New(embd.NewI2CBus(1), opt)

	if errNew != nil {
		log.Errorf("Error while creating the sensor reader !' %s'\n", errNew)
		err = errNew
	} else {
		log.Info("Sensor reader successfully created.")
	}

	accel, errRead := adxl345.Read()
	if errRead != nil {
		log.Errorf("Error while reading the data !!' %s'\n", errRead)
		err = errRead
	} else {
		log.Info("Data successfully read.")
	}
	x = accel.data[0]
	y = accel.data[1]
	z = accel.data[2]
	return x, y, z, err
}
