package wits0

import (
	"encoding/json"
	"strings"
)

// Records the WITS0 packet structure
type Records struct {
	Records []Record
}

// Record the WITS0 record structure
type Record struct {
	Record string
	Item   string
	Data   string
}

func createJSON(packet string, settings *wits0Settings) string {
	lines := strings.Split(packet, settings.lineEnding)
	records := make([]Record, len(lines)-3)
	parsingPackets := false
	index := 0
	for _, line := range lines {
		line = strings.Replace(line, settings.lineEnding, "", -1)
		if parsingPackets {
			if line == settings.packetFooter {
				parsingPackets = false
			} else {
				records[index].Record = line[0:2]
				records[index].Item = line[2:4]
				records[index].Data = line[4:len(line)]
				index = index + 1
			}
		} else if line == settings.packetHeader {
			parsingPackets = true
		}
	}
	jsonRecord, err := json.Marshal(Records{records})
	if err != nil {
		log.Error("Error converting packet to JSON: ", err)
		return ""
	}
	return string(jsonRecord)

}
