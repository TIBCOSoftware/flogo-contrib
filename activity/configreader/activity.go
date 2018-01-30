package configreader

import (
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/creamdog/gonfig"
	"os"
	"strconv"
	"sync"
)

// log is the default package logger
var log = logger.GetLogger("activity-tibco-configreader")

const (
	configFile       = "configFile"
	readEachTime     = "readEachTime"
	configName       = "configName"
	configValue      = "configValue"
	configType       = "configType"
	configDefaultVal = "defaultValue"
)

// ConfigReader structure
type ConfigReader struct {
	sync.Mutex
	metadata   *activity.Metadata
	gonfigConf gonfig.Gonfig
}

// NewActivity creates a new activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &ConfigReader{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *ConfigReader) Metadata() *activity.Metadata {
	return a.metadata
}

func (a *ConfigReader) getConfig(configFile string, reachEachTime bool, confName string, configType string, defaultValue interface{}) interface{} {
	a.Lock()
	defer a.Unlock()

	log.Debug("Variable readEachTime = ", reachEachTime)
	if reachEachTime || a.gonfigConf == nil {
		log.Debug("Need to read the configuration file...")
		f, err := os.Open(configFile)
		if err != nil {
			log.Error("Error while opening file ! ", err)
		}
		defer f.Close()
		a.gonfigConf, err = gonfig.FromJson(f)
		if err != nil {
			log.Error("Error while reading configuration file ! ", err)
		}
	}

	var confValue interface{}

	var err error

	switch configType {
	case "string":
		log.Debug("Reading STRING value")
		confValue, err = a.gonfigConf.GetString(confName, defaultValue)
	case "int":
		log.Debug("Reading INT value")
		confValue, err = a.gonfigConf.GetInt(confName, defaultValue)
	case "float":
		log.Debug("Reading FLOAT value")
		confValue, err = a.gonfigConf.GetFloat(confName, defaultValue)
	case "bool":
		log.Debug("Reading BOOL value")
		confValue, err = a.gonfigConf.GetBool(confName, defaultValue)
	default:
		log.Debug("Reading STRING value")
		confValue, err = a.gonfigConf.GetString(confName, defaultValue)
	}

	if err != nil {
		log.Error("Error while getting configuration value ! ", err)
	}

	log.Debug("Final value: ", confValue)

	return confValue
}

func toBool(val interface{}) (bool, error) {

	b, ok := val.(bool)
	if !ok {
		s, ok := val.(string)

		if !ok {
			return false, fmt.Errorf("Unable to convert to boolean")
		}

		var err error
		b, err = strconv.ParseBool(s)

		if err != nil {
			return false, err
		}
	}

	return b, nil
}

func (a *ConfigReader) setDefaultValue(defaultVal interface{}, configType string) interface{} {
	var configurationDefaultValueTmp interface{}
	var confValue interface{}

	var err error

	if defaultVal != nil {
		configurationDefaultValueTmp = defaultVal
	}

	switch configType {
	case "string":
		if configurationDefaultValueTmp != nil {
			confValue = configurationDefaultValueTmp.(string)
		} else {
			confValue = ""
		}
	case "int":
		if configurationDefaultValueTmp != nil {
			confValue, err = strconv.ParseInt(configurationDefaultValueTmp.(string), 10, 64)
		} else {
			confValue = 0
		}
	case "float":
		if configurationDefaultValueTmp != nil {
			confValue, err = strconv.ParseFloat(configurationDefaultValueTmp.(string), 64)
		} else {
			confValue = 0
		}
	case "bool":
		if configurationDefaultValueTmp != nil {
			confValue, err = strconv.ParseBool(configurationDefaultValueTmp.(string))
		} else {
			confValue = true
		}
	default:
		if configurationDefaultValueTmp != nil {
			confValue = configurationDefaultValueTmp.(string)
		} else {
			confValue = ""
		}
	}

	if err != nil {
		log.Error("Error while setting default value !")
	}

	log.Debugf("Input default value is [%s]", confValue)
	return confValue
}

// Eval implements activity.Activity.Eval
func (a *ConfigReader) Eval(context activity.Context) (done bool, err error) {

	configFile := context.GetInput(configFile).(string)
	log.Debugf("Config file [%s]", configFile)

	var readEachTimeB bool

	if context.GetInput(readEachTime) != nil {
		log.Debug("Variable readEachTime is not null.")
		readEachTimeB, _ = toBool(context.GetInput(readEachTime))
	}
	if context.GetInput(configName) != nil {
		var configurationName string
		var configurationType string

		configurationName = context.GetInput(configName).(string)

		if context.GetInput(configType) != nil {
			configurationType = context.GetInput(configType).(string)
		} else {
			log.Debug("Using default value (string) for configuration type.")
			configurationType = "string"
		}
		log.Debugf("Configuration name [%s], Configuration type [%s]", configurationName, configurationType)

		configurationDefaultValue := a.setDefaultValue(context.GetInput(configDefaultVal), configurationType)

		log.Debug("Getting config value...")
		confValue := a.getConfig(configFile, readEachTimeB, configurationName, configurationType, configurationDefaultValue)
		log.Debugf("Final value returned [%s]", confValue)

		context.SetOutput(configValue, confValue)
	} else {
		return false, fmt.Errorf("no configuration name has been set")
	}
	return true, nil
}
