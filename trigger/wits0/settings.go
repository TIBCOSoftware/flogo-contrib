package wits0

import (
	"strconv"
	"time"

	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/tarm/serial"
)

// wits0Record the WITS0 record structure
type wits0Settings struct {
	serialConfig                  *serial.Config
	packetHeader                  string
	packetFooter                  string
	lineEnding                    string
	heartBeatValue                string
	heartBeatInterval             int
	packetFooterWithLineSeparator string
	outputRaw                     bool
}

func (settings *wits0Settings) Init(t *wits0Trigger, endpoint *trigger.Handler) {
	settings.serialConfig = &serial.Config{
		Name:        GetSettingSafe(endpoint, "SerialPort", ""),
		Baud:        GetSafeNumber(endpoint, "BaudRate", 9600),
		ReadTimeout: time.Duration(GetSafeNumber(endpoint, "ReadTimeoutSeconds", 1)),
		Size:        byte(GetSafeNumber(endpoint, "DataBits", 8)),
		Parity:      serial.Parity(GetSafeNumber(endpoint, "Parity", 0)),
		StopBits:    serial.StopBits(GetSafeNumber(endpoint, "StopBits", 1)),
	}

	settings.packetHeader = GetSettingSafe(endpoint, "PacketHeader", "&&")
	settings.packetFooter = GetSettingSafe(endpoint, "PacketFooter", "!!")
	settings.lineEnding = GetSettingSafe(endpoint, "LineSeparator", "\r\n")
	settings.heartBeatValue = GetSettingSafe(endpoint, "HeartBeatValue", "&&\r\n0111-9999\r\n!!\r\n")
	settings.heartBeatInterval = GetSafeNumber(endpoint, "HeartBeatInterval", 30)
	settings.packetFooterWithLineSeparator = settings.packetFooter + settings.lineEnding
	settings.outputRaw = GetSafeBool(endpoint, "OutputRaw", false)

	log.Debug("Serial Config: ", settings.serialConfig)
	log.Debug("packetHeader: ", settings.packetHeader)
	log.Debug("packetFooter: ", settings.packetFooter)
	log.Debug("lineEnding: ", settings.lineEnding)
	log.Debug("heartBeatValue: ", settings.heartBeatValue)
	log.Debug("heartBeatInterval: ", settings.heartBeatInterval)
	log.Debug("outputRaw: ", settings.outputRaw)
}

// GetSettingSafe get a setting and returns default if not found
func GetSettingSafe(endpoint *trigger.Handler, setting string, defaultValue string) string {
	var retString string
	defer func() {
		if r := recover(); r != nil {
			retString = defaultValue
		}
	}()
	retString = endpoint.GetStringSetting(setting)
	return retString
}

// GetSafeNumber gets the number from the config checking for empty and using default
func GetSafeNumber(endpoint *trigger.Handler, setting string, defaultValue int) int {
	if settingString := GetSettingSafe(endpoint, setting, ""); settingString != "" {
		value, _ := strconv.Atoi(settingString)
		return value
	}
	return defaultValue
}

// GetSafeBool gets the bool from the config checking for empty and using default
func GetSafeBool(endpoint *trigger.Handler, setting string, defaultValue bool) bool {
	if settingString := GetSettingSafe(endpoint, setting, ""); settingString != "" {
		value, _ := strconv.ParseBool(settingString)
		return value
	}
	return defaultValue
}
