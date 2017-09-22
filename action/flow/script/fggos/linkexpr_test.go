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
        "id": "C",
        "type": 1,
        "name": "C",
        "activityType": ""
      },
      {
        "id": "D",
        "type": 1,
        "name": "D",
        "activityType": ""
      },
      {
        "id": "4",
        "type": 1,
        "name": "E",
        "activityType": ""
      }
    ],
    "links": [
      { "id": 1, "type": 1, "name": "",  "from": 2, "to": 3, "value":"${petMax} > 2" },
      { "id": 2, "type": 1, "name": "",  "from": 2, "to": 3, "value":"${A2.result}.code == 1" },
      { "id": 3, "type": 1, "name": "",  "from": 2, "to": "C", "value":"isDefined(\"${activity.2.result}\")" },
      { "id": 4, "type": 1, "name": "", "from": "C", "to": "D", "value":"${activity.C.result} <= ${petMax}" },
      { "id": 5, "type": 1, "name": "", "from": "C", "to": 4, "value":"${petId}<=${petMax}" },
      { "id": 6, "type": 1, "name": "", "from": "C", "to": 4, "value":"${T.result}.code<=${petMax}" },
      { "id": 7, "type": 1, "name": "", "from": "C", "to": 4, "value":"isDefined(\"${petId}\")" }
    ]
  }
}
`

func TestGosLinkExprManager_TestTransExpr(t *testing.T) {

	expr := "${A3.result}.code == 1"
	_, tExpr := transExpr(expr)
	fmt.Println(expr)
	fmt.Println(tExpr)

	expr = "${A3.result}.code==1"
	_, tExpr = transExpr(expr)
	fmt.Println(expr)
	fmt.Println(tExpr)

	expr = "${activity.2.result}.code ==5"
	_, tExpr = transExpr(expr)
	fmt.Println(expr)
	fmt.Println(tExpr)

	expr = "isDefined(\"${activity.2.result}\")"
	_, tExpr = transExpr(expr)
	fmt.Println(expr)
	fmt.Println(tExpr)

	expr = "isDefined(\"${activity.2.result}\") && ${activity.2.result}.code > 5"
	_, tExpr = transExpr(expr)
	fmt.Println(expr)
	fmt.Println(tExpr)

	expr = "isDefined(\"${activity.2.result}\") || isDefined(\"${activity.1.result}\")"
	_, tExpr = transExpr(expr)
	fmt.Println(expr)
	fmt.Println(tExpr)

	expr = "${petId}<=${petMax}"
	_, tExpr = transExpr(expr)
	fmt.Println(expr)
	fmt.Println(tExpr)


	expr = "${activity.C.result} <= ${petMax}"
	_, tExpr = transExpr(expr)
	fmt.Println(expr)
	fmt.Println(tExpr)
}

func TestGosLinkExprManager_EvalLinkExpr(t *testing.T) {

	defRep := &definition.DefinitionRep{}
	err := json.Unmarshal([]byte(defJSON), defRep)

	if err != nil {
		t.Fatalf("Error evaluating expressions '%s'", err.Error())
	}

	def, _ := definition.NewDefinition(defRep)
	f := GosLinkExprManagerFactory{}
	mgr := f.NewLinkExprManager(def)

	link1 := def.GetLink(1)
	link2 := def.GetLink(2)
	link3 := def.GetLink(3)
	link4 := def.GetLink(4)
	link5 := def.GetLink(5)
	link6 := def.GetLink(6)
	link7 := def.GetLink(7)

	a2result := make(map[string]interface{})
	a2result["code"] = 1

	aCresult := 2

	attrs := []*data.Attribute{
		data.NewAttribute("petMax", data.INTEGER, 4),
		data.NewAttribute("petId", data.INTEGER, 3),
		data.NewAttribute("_A.C.result", data.OBJECT, aCresult),
		data.NewAttribute("_A.2.result", data.OBJECT, a2result),
		data.NewAttribute("_T.result", data.OBJECT, a2result),
	}

	scope := data.NewSimpleScope(attrs, nil)

	result, err := mgr.EvalLinkExpr(link1, scope)
	if err != nil {
		t.Fatalf("Error evaluating expressions '%s'", err.Error())
	}
	fmt.Printf("Link 1 Result: %v\n", result)

	result, err = mgr.EvalLinkExpr(link2, scope)
	if err != nil {
		t.Fatalf("Error evaluating expressions '%s'", err.Error())
	}
	fmt.Printf("Link 2 Result: %v\n", result)

	result, err = mgr.EvalLinkExpr(link3, scope)
	if err != nil {
		t.Fatalf("Error evaluating expressions '%s'", err.Error())
	}
	fmt.Printf("Link 3 Result: %v\n", result)

	result, err = mgr.EvalLinkExpr(link4, scope)
	if err != nil {
		t.Fatalf("Error evaluating expressions '%s'", err.Error())
	}
	fmt.Printf("Link 4 Result: %v\n", result)

	result, err = mgr.EvalLinkExpr(link5, scope)
	if err != nil {
		t.Fatalf("Error evaluating expressions '%s'", err.Error())
	}
	fmt.Printf("Link 5 Result: %v\n", result)

	scope.SetAttrValue("petId", 6)
	result, err = mgr.EvalLinkExpr(link2, scope)
	if err != nil {
		t.Fatalf("Error evaluating expressions '%s'", err.Error())
	}

	fmt.Printf("Link2 Result: %v\n", result)

	result, err = mgr.EvalLinkExpr(link6, scope)
	if err != nil {
		t.Fatalf("Error evaluating expressions '%s'", err.Error())
	}

	fmt.Printf("Link6 Result: %v\n", result)

	//todo fix so isDefined("petId") works
	result, err = mgr.EvalLinkExpr(link7, scope)
	if err != nil {
		t.Fatalf("Error evaluating expressions '%s'", err.Error())
	}

	fmt.Printf("Link7 Result: %v\n", result)
}
