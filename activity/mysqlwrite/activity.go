package mysqlwrite

import (
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// log is the default package logger
var log = logger.GetLogger("activity-akash-Database_Query")

const (
	driverName     = "driverName"
	datasourceName = "datasourceName"
	preparequery   = "preparequery"
	queryvalue    = "queryvalue"

	ovResult = "result"
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

	////////  Set DriverName of the driver //////////

	driverNameInput := context.GetInput(driverName)

	ivdriverName, ok := driverNameInput.(string)
	if !ok {
		context.SetOutput("result", "driverNameSET")
		return true, fmt.Errorf("driverName not set")
	}
	log.Debugf("driverNamename" + ivdriverName)

	////////  END - Set DriverName of the driver //////////

	////////  Set connection String of the driver //////////

	datasourceNameInput := context.GetInput(datasourceName)

	ivdatasourceName, ok := datasourceNameInput.(string)
	if !ok {
		context.SetOutput("result", "datasourceNameSET")
		return true, fmt.Errorf("datasourceName not set")
	}
	log.Debugf("datasourceNamename" + ivdatasourceName)

	////////  END - Set connection String of the driver //////////

	preparequeryInput := context.GetInput(preparequery)

	ivpreparequery, ok := preparequeryInput.(string)
	if !ok {
		context.SetOutput("result", "QUERY_NOT_SET")
		return true, fmt.Errorf("Query not set")
	}

	queryvalueInput := context.GetInput(queryvalue)

	ivqueryvalue, ok := queryvalueInput.(string)
	if !ok {
		context.SetOutput("result", "QUERY_NOT_SET")
		return true, fmt.Errorf("Query not set")
	}

	//////////////////////////////////////////////////

	log.Debugf("query" + ivpreparequery)

	log.Debugf("All Parameters set")

	log.Debugf("Go SQL Connection Initiated...")

	db, err := sql.Open(ivdriverName, ivdatasourceName)
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	fmt.Println("Successfully Connected to Database")

	//////////////////////////////////////////////////////////

	// insert
	stmt, err := db.Prepare(ivpreparequery)
	if err != nil {
		panic(err.Error())
	}

	res, err := stmt.Exec(ivqueryvalue)
	if err != nil {
		panic(err.Error())
	}

	id, err := res.LastInsertId()
	if err != nil {
		panic(err.Error())
	}

	//_, queryerr := db.Query(ivquery)

	// if queryerr != nil {
	// 	panic(queryerr.Error())
	// }

	context.SetOutput(ovResult, id)

	return true, nil
}
