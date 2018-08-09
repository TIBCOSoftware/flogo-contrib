package readfile

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"bufio"
	"fmt"
	"os"
)

// log is the default package logger
var log = logger.GetLogger("Activity Akash-File Reader")

const (
	filename = "filename"
	lineNumber = "lineNumber"

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
	ivlineNumber, ok := context.GetInput(lineNumber).(int)
	if !ok {
		context.SetOutput("result", "LINE NUMBER NOT SET")
		return true, fmt.Errorf("line number not set")
	}

	//ivfilename, ok := filenameInput.(string)
	ivfilename, ok := context.GetInput(filename).(string)
	if !ok {
		context.SetOutput("result", "FILENAME_NOT_SET")
		return true, fmt.Errorf("Filename not set")
	}

  fileHandle, _ := os.Open(ivfilename)
	defer fileHandle.Close()
	fileScanner := bufio.NewScanner(fileHandle)

	lastLine := 0
	line := ""
	
	for fileScanner.Scan() {
		lastLine++

		if lastLine == ivlineNumber {
			line = fileScanner.Text()
			break
		}
	}	
		
  context.SetOutput("result", line)
    
  return true, nil
}