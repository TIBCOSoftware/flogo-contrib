package fggos

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/definition"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

const defJSON = `
{
  "type": 1,
  "name": "Demo Flow",
  "model": "simple",
  "attributes": [
    { "name": "petMax", "type": "integer", "value": 5 }
  ],
  "rootTask": {
    "id": 1,
    "type": 1,
    "activityType": "",
    "name": "root",
    "tasks": [
      {
        "id": 2,
        "type": 1,
        "name": "A",
        "activityType": ""
      },
      {
        "id": 3,
        "type": 1,
        "name": "B",
        "activityType": ""
      },
      {
        "id": 4,
        "type": 1,
        "name": "C",
        "activityType": ""
      }
    ],
    "links": [
      { "id": 1, "type": 1, "name": "",  "from": 2, "to": 3, "value":"$sensorData.temp > 50" },
      { "id": 2, "type": 1, "name": "", "from": 2, "to": 4, "value":"$petId <= $petMax" },
      { "id": 3, "type": 1, "name": "", "from": 2, "to": 4, "value":"true" },
      { "id": 4, "type": 1, "name": "", "from": 2, "to": 4, "value":"isDefined($sensorData.ff)" },
      { "id": 5, "type": 1, "name": "", "from": 2, "to": 4, "value":"$sensorData.temp == 55" },
      { "id": 6, "type": 1, "name": "", "from": 2, "to": 4, "value":"${A3.result}.code == 1" }
    ]
  }
}
`

func TestLuaLinkExprManager_TestTransExpr(t *testing.T) {

	expr := "${A3.result}.code == 1"
	_, tExpr := transExpr(expr)
	fmt.Println(tExpr)

	expr = "${A3.result}.code==1"
	_, tExpr = transExpr(expr)
	fmt.Println(tExpr)

	expr = "$sensorData.temp == 1"
	_, tExpr = transExpr(expr)
	fmt.Println(tExpr)

	expr = "$sensorData.temp==1"
	_, tExpr = transExpr(expr)
	fmt.Println(tExpr)

	expr = "$petId<=$petMax"
	_, tExpr = transExpr(expr)
	fmt.Println(tExpr)

	expr = "$petId[0]<=$petMax"
	_, tExpr = transExpr(expr)
	fmt.Println(tExpr)

	expr = "isDefined($petId)"
	_, tExpr = transExpr(expr)
	fmt.Println(tExpr)

	expr = "isDefined($petId) || $x > 5"
	_, tExpr = transExpr(expr)
	fmt.Println(tExpr)

	expr = "isDefined($petId) || isDefined($x)"
	_, tExpr = transExpr(expr)
	fmt.Println(tExpr)
}

func TestLuaLinkExprManager_EvalLinkExpr(t *testing.T) {

	defRep := &definition.DefinitionRep{}
	json.Unmarshal([]byte(defJSON), defRep)

	def, _ := definition.NewDefinition(defRep)
	f := GosLinkExprManagerFactory{}
	mgr := f.NewLinkExprManager(def)

	link1 := def.GetLink(1)
	link2 := def.GetLink(2)
	link3 := def.GetLink(3)
	link4 := def.GetLink(4)
	link5 := def.GetLink(5)
	link6 := def.GetLink(6)

	sensorData := make(map[string]interface{})
	sensorData["temp"] = 55

	a3result := make(map[string]interface{})
	a3result["code"] = 1

	attrs := []*data.Attribute{
		data.NewAttribute("petMax", data.INTEGER, 4),
		data.NewAttribute("petId", data.INTEGER, 3),
		data.NewAttribute("sensorData", data.OBJECT, sensorData),
		data.NewAttribute("{A3.result}", data.OBJECT, a3result),
	}

	scope := data.NewSimpleScope(attrs, nil)

	result, err := mgr.EvalLinkExpr(link1, scope)
	if err != nil{
		t.Fatalf("Error evaluating expressions '%s'", err.Error())
	}
	fmt.Printf("Link 1 Result: %v\n", result)

	result, err = mgr.EvalLinkExpr(link2, scope)
	if err != nil{
		t.Fatalf("Error evaluating expressions '%s'", err.Error())
	}
	fmt.Printf("Link 2 Result: %v\n", result)

	result, err = mgr.EvalLinkExpr(link3, scope)
	if err != nil{
		t.Fatalf("Error evaluating expressions '%s'", err.Error())
	}
	fmt.Printf("Link 3 Result: %v\n", result)

	result, err = mgr.EvalLinkExpr(link4, scope)
	if err != nil{
		t.Fatalf("Error evaluating expressions '%s'", err.Error())
	}
	fmt.Printf("Link 4 Result: %v\n", result)

	result, err = mgr.EvalLinkExpr(link5, scope)
	if err != nil{
		t.Fatalf("Error evaluating expressions '%s'", err.Error())
	}
	fmt.Printf("Link 5 Result: %v\n", result)

	result, err = mgr.EvalLinkExpr(link6, scope)
	if err != nil{
		t.Fatalf("Error evaluating expressions '%s'", err.Error())
	}
	fmt.Printf("Link 6 Result: %v\n", result)

	scope.SetAttrValue("petId", 6)
	result, err = mgr.EvalLinkExpr(link2, scope)
	if err != nil{
		t.Fatalf("Error evaluating expressions '%s'", err.Error())
	}

	fmt.Printf("Link2 Result: %v\n", result)
}
