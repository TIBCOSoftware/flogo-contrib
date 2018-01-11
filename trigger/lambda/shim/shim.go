package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"flag"

	fl "github.com/TIBCOSoftware/flogo-contrib/trigger/lambda"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

// Handle implements the Flogo Function handler
func Handle(ctx context.Context, evt json.RawMessage) (interface{}, error) {
	err := setupArgs(evt, &ctx)
	if err != nil {
		return nil, err
	}
	result, err := fl.Invoke()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func setupArgs(evt json.RawMessage, ctx *context.Context) error {
	// Setup environment argument
	evtJSON, err := json.Marshal(&evt)
	if err != nil {
		return err
	}

	evtFlag := flag.Lookup("evt")
	if evtFlag == nil {
		flag.String("evt", string(evtJSON), "Lambda Environment Arguments")
	} else {
		flag.Set("evt", string(evtJSON))
	}

	// Setup context argument
	ctxObj, _ := lambdacontext.FromContext(*ctx)
	var ctxBuff bytes.Buffer
	err = binary.Write(&ctxBuff, binary.BigEndian, ctxObj)
	if err != nil {
		return err
	}

	ctxFlag := flag.Lookup("ctx")
	if ctxFlag == nil {
		flag.String("ctx", ctxBuff.String(), "Lambda Context Arguments")
	} else {
		flag.Set("ctx", ctxBuff.String())
	}

	return nil
}

func main() {
	lambda.Start(Handle)
}
