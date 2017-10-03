package lambda

import (
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

var log = logger.GetLogger("activity-tibco-lambda")

const (
	ivArn       = "arn"
	ivRegion    = "region"
	ivAccessKey = "accessKey"
	ivSecretKey = "secretKey"
	ivPayload   = "payload"

	ovValue = "value"
)

// LambdaActivity is a App Activity implementation
type LambdaActivity struct {
	sync.Mutex
	metadata *activity.Metadata
}

// NewActivity creates a new LambdaActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &LambdaActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *LambdaActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *LambdaActivity) Eval(context activity.Context) (done bool, err error) {
	arn := context.GetInput(ivArn).(string)
	var accessKey, secretKey = "", ""
	if context.GetInput(ivAccessKey) != nil {
		accessKey = context.GetInput(ivAccessKey).(string)
	}
	if context.GetInput(ivSecretKey) != nil {
		secretKey = context.GetInput(ivSecretKey).(string)
	}
	payload := context.GetInput(ivPayload).(string)

	var config *aws.Config
	if accessKey != "" && secretKey != "" {
		config = aws.NewConfig().WithRegion(context.GetInput(ivRegion).(string)).WithCredentials(credentials.NewStaticCredentials(accessKey, secretKey, ""))
	} else {
		config = aws.NewConfig().WithRegion(context.GetInput(ivRegion).(string))
	}
	aws := lambda.New(session.New(config))

	out, awsErr := aws.Invoke(&lambda.InvokeInput{
		FunctionName: &arn,
		Payload:      []byte(payload)})

	if awsErr != nil {
		log.Error(err)

		return true, awsErr
	}

	log.Info(out)

	return true, nil
}
