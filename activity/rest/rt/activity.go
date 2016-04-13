package rest

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/ext/activity"
	"github.com/op/go-logging"
)

// log is the default package logger
var log = logging.MustGetLogger("activity-tibco-rest")

// RESTActivity is an Activity that is used to invoke a REST Operation
// inputs : {method,uri,params}
// outputs: {result}
type RESTActivity struct {
	metadata *activity.Metadata
}

// init create & register activity
func init() {
	md := activity.NewMetadata(jsonMetadata)
	activity.Register(&RESTActivity{metadata: md})
}

// Metadata returns the activity's metadata
func (a *RESTActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Invokes a REST Operation
func (a *RESTActivity) Eval(context activity.Context) bool {

	method, _ := context.GetAttrValue("method")
	uri, _ := context.GetAttrValue("uri")

	//for now assume parameter is in ActivityContext
	//params _:= context.GetAttrValue("parameters")

	//only doing get for now, so ignore payload
	//payload := nil

	path := BuildURIFromScope(uri, context)

	log.Debugf("REST Call: %s %s\n", method, path)

	req, err := http.NewRequest(method, path, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	log.Debug("response Status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)

	result := string(body)

	if log.IsEnabledFor(logging.DEBUG) {
		log.Debug("response Body:", result)
	}

	context.SetAttrValue("result", result)

	return true
}

// BuildURI is a temporary crude URI builder
func BuildURI(uri string, values map[string]string) string {

	var buffer bytes.Buffer
	buffer.Grow(len(uri))

	var i int
	for i < len(uri) {
		if uri[i] == '{' {
			j := i + 1
			for j < len(uri) && uri[j] != '}' {
				j++
			}

			param := uri[i+1 : j]
			log.Debugf("Param: %s\n", param)

			value := values[param]
			buffer.WriteString(value)

			i = j + 1
		} else {
			buffer.WriteByte(uri[i])
			i++
		}
	}

	return buffer.String()
}

// BuildURIFromScope is a temporary crude URI builder using Scope
func BuildURIFromScope(uri string, values data.Scope) string {

	var buffer bytes.Buffer
	buffer.Grow(len(uri))

	var i int
	for i < len(uri) {
		if uri[i] == '{' {
			j := i + 1
			for j < len(uri) && uri[j] != '}' {
				j++
			}

			param := uri[i+1 : j]
			log.Debugf("Param: %s\n", param)

			value, _ := values.GetAttrValue(param)
			buffer.WriteString(value)

			i = j + 1
		} else {
			buffer.WriteByte(uri[i])
			i++
		}
	}

	return buffer.String()
}
