package ReadFile

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"bufio"
	"fmt"
    "os"
   // "io"
)

// log is the default package logger
var log = logger.GetLogger("Activity Akash-File Reader")

const (
	filename = "filename"
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

	filenameInput := context.GetInput(filename)

	ivfilename, ok := filenameInput.(string)
	if !ok {
		context.SetOutput("result", "FILENAME_NOT_SET")
		return true, fmt.Errorf("Filename not set")
	}

    fileHandle, _ := os.Open(ivfilename)
	defer fileHandle.Close()
	fileScanner := bufio.NewScanner(fileHandle)


//Current implementation is only for files containing one line only
//	for fileScanner.Scan() {
//		fmt.Println(fileScanner.Text())
//    }	
    
    context.SetOutput("result", fileScanner.Text())
    
    return true, nil
}
