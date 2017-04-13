package instance

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

const defJSON = `
{
    "type": 1,
    "name": "test",
    "model": "test",
    "rootTask": {
      "id": 1,
      "type": 1,
      "activityType": "",
      "name": "root",
      "tasks": [
        {
          "id": 2,
          "type": 1,
          "name": "a"
        },
        {
          "id": 3,
          "type": 1,
          "name": "b"
        }
      ],
      "links": [
        { "id": 1, "type": 1,  "name": "","from": 2, "to": 3 }
      ]
    }
  }
`

//func TestFullSerialization(t *testing.T) {
//
//	defRep := &flowdef.DefinitionRep{}
//	json.Unmarshal([]byte(defJSON), defRep)
//
//	log.Infof("Def Rep: %v", defRep)
//
//	def, _ := flowdef.NewDefinition(defRep)
//
//	idGen, _ := util.NewGenerator()
//	id := idGen.NextAsString()
//
//	instance := NewFlowInstance(id, "uri1", def)
//
//	instance.Start(nil)
//
//	hasWork := true
//
//	for hasWork && instance.Status() < StatusCompleted {
//		hasWork = instance.DoStep()
//
//		json, _ := json.Marshal(instance)
//		log.Debugf("Snapshot: %s\n", string(json))
//	}
//}

//func TestIncrementalSerialization(t *testing.T) {
//
//	defRep := &flowdef.DefinitionRep{}
//	json.Unmarshal([]byte(defJSON), defRep)
//
//	idGen, _ := util.NewGenerator()
//	id := idGen.NextAsString()
//
//	def, _ := flowdef.NewDefinition(defRep)
//
//	instance := NewFlowInstance(id, "uri2", def)
//
//	instance.Start(nil)
//
//	hasWork := true
//
//	for hasWork && instance.Status() < StatusCompleted {
//		hasWork = instance.DoStep()
//
//		json, _ := json.Marshal(instance.GetChanges())
//		log.Debugf("Changes: %s\n", string(json))
//	}
//}
