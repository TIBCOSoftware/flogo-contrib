package ReadFile

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// log is the default package logger
var log = logger.GetLogger("Activity Akash-File Reader")

const (
	filename = "filename"
	ivreadALine = "readALine"
	ivlineNumber = "lineNumber"

	ovresult = "result"
	
)
// MyActivity is a stub for your Activity implementation
type MyActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new AppActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &MyActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *MyActivity) Metadata() *activity.Metadata {
	return a.metadata
}


// Eval implements activity.Activity.Eval
func (a *MyActivity) Eval(context activity.Context) (done bool, err error) {

	// do eval

	//filenameInput, _ := context.GetInput(filename).(string)
	readALine, _ := toBool(context.GetInput(ivreadALine))
	lineNumber, _ := context.GetInput(ivreadALine).(int)

	//ivfilename, ok := filenameInput.(string)
	ivfilename, ok := context.GetInput(filename).(string)
	if !ok {
		context.SetOutput("result", "FILENAME_NOT_SET")
		return true, fmt.Errorf("Filename not set")
	}

  fileHandle, _ := os.Open(ivfilename)
	defer fileHandle.Close()
	fileScanner := bufio.NewScanner(fileHandle)


	if (readALine) {
		
		for fileScanner.Scan() {
		fmt.Println(fileScanner.Text())
    }	
    

	}

	lineNumber++;
	
  context.SetOutput("result", fileScanner.Text())
    
  return true, nil
}


func toBool(val interface{}) (bool, error) {

	b, ok := val.(bool)
	if !ok {
		s, ok := val.(string)

		if !ok {
			return false, fmt.Errorf("unable to convert to boolean")
		}

		var err error
		b, err = strconv.ParseBool(s)

		if err != nil {
			return false, err
		}
	}

	return b, nil
}