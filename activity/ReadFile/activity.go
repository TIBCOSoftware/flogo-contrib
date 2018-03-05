package ReadFile

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"bufio"
	"fmt"
    "os"
    "io"
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

func ReadString(filename string) {
    f, err := os.Open(filename)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer f.Close()
    r := bufio.NewReader(f)
    line, err := r.ReadString('\n')
    for err == nil {
        fmt.Print(line)
        line, err = r.ReadString('\n')
    }
    if err != io.EOF {
        fmt.Println(err)
        return
    }
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

	ReadString(ivfilename)
	return true, nil
}
