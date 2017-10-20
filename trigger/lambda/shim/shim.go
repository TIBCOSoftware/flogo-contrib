package main

import (
	"encoding/json"
	"flag"

	"github.com/TIBCOSoftware/flogo-contrib/trigger/lambda"
	"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
)


func Handle(evt json.RawMessage, ctx *runtime.Context) (interface{}, error) {
	err := setupArgs(evt, ctx)
	if err != nil {
		return nil, err
	}
	result, err := lambda.Invoke()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func setupArgs(evt json.RawMessage, ctx *runtime.Context) error {
	// Setup environment argument
	evtJson, err := json.Marshal(&evt)
	if err != nil {
		return err
	}
	
	evtFlag := flag.Lookup("evt")
	if evtFlag == nil {
		flag.String("evt", string(evtJson), "Lambda Environment Arguments")
	} else {
		flag.Set("evt", string(evtJson))
	}
	
	// Setup context argument
	ctxJson, err := json.Marshal(ctx)
	if err != nil {
		return err
	}
	
	ctxFlag := flag.Lookup("ctx")
	if ctxFlag == nil {
		flag.String("ctx", string(ctxJson), "Lambda Context Arguments")
	} else {
		flag.Set("ctx", string(ctxJson))
	}
	
	return nil
}

func main() {
	// No Op
}
