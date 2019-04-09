package wits0

import (
	"bytes"
	"context"
	"strings"
	"time"

	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/tarm/serial"
)

type serialPort struct {
	settings *wits0Settings
	stream   *serial.Port
	trigger  *wits0Trigger
	endpoint *trigger.Handler
}

func (serialPort *serialPort) Init(trigger *wits0Trigger, endpoint *trigger.Handler) {
	serialPort.trigger = trigger
	serialPort.endpoint = endpoint
	serialPort.settings = &wits0Settings{}
	serialPort.settings.Init(trigger, endpoint)
}

func (serialPort *serialPort) createSerialConnection() {
	log.Info("Connecting to serial port: " + serialPort.settings.serialConfig.Name)
	stream, err := serial.OpenPort(serialPort.settings.serialConfig)
	if err != nil {
		log.Error(err)
		return
	}
	serialPort.stream = stream
	log.Info("Connected to serial port")

	serialPort.heartBeat()
	serialPort.readSerialData()
}

func (serialPort *serialPort) heartBeat() {
	if serialPort.settings.heartBeatInterval > 0 {
		duration := time.Second * time.Duration(serialPort.settings.heartBeatInterval)

		go func() {
			start := time.Now()
		writeLoop:
			for {
				time.Sleep(time.Millisecond * 100)
				select {
				case <-serialPort.trigger.stopCheck:
					break writeLoop
				default:
				}

				elapsed := time.Now().Sub(start)
				if elapsed > duration {
					log.Debug("Sending heartbeat")
					serialPort.stream.Write([]byte(serialPort.settings.heartBeatValue))
					start = time.Now()
				}
			}
		}()
	}
}

func (serialPort *serialPort) readSerialData() {
	buf := make([]byte, 1024)
	buffer := bytes.NewBufferString("")
readLoop:
	for {
		n, err := serialPort.stream.Read(buf)
		if err != nil {
			log.Error(err)
			break
		}
		if n > 0 {
			buffer.Write(buf[:n])
			buffer = serialPort.parseBuffer(buffer)
		}

		select {
		case <-serialPort.trigger.stopCheck:
			break readLoop
		default:
		}
	}
}

func (serialPort *serialPort) parseBuffer(buffer *bytes.Buffer) *bytes.Buffer {
	for buffer.Len() > 0 {
		check := buffer.String()
		indexStart := strings.Index(check, serialPort.settings.packetHeader)
		indexEnd := strings.Index(check, serialPort.settings.packetFooterWithLineSeparator)
		if indexEnd >= 0 && indexStart >= 0 && indexEnd < indexStart {
			log.Info("Dropping initial bad packet")
			return bytes.NewBufferString(check[indexStart:len(check)])
		}
		if indexStart >= 0 {
			startRemoved := string(check[indexStart:len(check)])
			indexStartSecond := indexStart + strings.Index(startRemoved, serialPort.settings.packetHeader)
			if indexStartSecond > 0 && indexStartSecond > indexStart && indexStartSecond < indexEnd {
				log.Info("Dropping bad packet")
				return bytes.NewBufferString(check[indexStartSecond+len(serialPort.settings.packetHeader) : len(check)])
			}
		}
		if indexStart >= 0 && indexEnd > indexStart {
			indexEndIncludeStopPacket := indexEnd + len(serialPort.settings.packetFooterWithLineSeparator)
			packet := check[indexStart:indexEndIncludeStopPacket]
			outputData := packet
			if !serialPort.settings.outputRaw {
				outputData = createJSON(packet, serialPort.settings)
			}
			if len(outputData) > 0 {
				trgData := make(map[string]interface{})
				trgData["data"] = outputData
				_, err := serialPort.endpoint.Handle(context.Background(), trgData)
				if err != nil {
					log.Error("Error starting action: ", err.Error())
				}
			}
			buffer = bytes.NewBufferString(check[indexEndIncludeStopPacket:len(check)])
		} else {
			return buffer
		}
	}
	return buffer
}
