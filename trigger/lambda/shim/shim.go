package main

import (
	"encoding/json"
	"flag"

	"github.com/TIBCOSoftware/flogo-contrib/trigger/lambda"
	"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
)


func Handle(evt json.RawMessage, ctx *runtime.Context) (string, error) {
	err := setupArgs(evt, ctx)
	if err != nil {
		return "", err
	}
	result, err := lambda.Invoke()
	if err != nil {
		return "", err
	}
	return result, nil
}

func setupArgs(evt json.RawMessage, ctx *runtime.Context) error {
	evtFlag := flag.Lookup("evt")
	if evtFlag == nil {
		// Setup environment argument
		evtJson, err := json.Marshal(&evt)
		if err != nil {
			return err
		}
		flag.String("evt", string(evtJson), "Lambda Environment Arguments")
	}
	ctxFlag := flag.Lookup("ctx")
	if ctxFlag == nil {
		// Setup context argument
		ctxJson, err := json.Marshal(ctx)
		if err != nil {
			return err
		}
		flag.String("ctx", string(ctxJson), "Lambda Context Arguments")
	}
	return nil
}

func main() {
	// No Op
}