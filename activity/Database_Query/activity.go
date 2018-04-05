package Database_Query

import (
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mxk/go-sqlite/sqlite3"

	"bytes"
	"encoding/json"
	"strings"
)

// log is the default package logger
var log = logger.GetLogger("activity-akash-Database_Query")

const (
	driverName     = "driverName"
	datasourceName = "datasourceName"
	query          = "query"

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

	queryInput := context.GetInput(query)

	ivquery, ok := queryInput.(string)
	if !ok {
		context.SetOutput("result", "QUERY_NOT_SET")
		return true, fmt.Errorf("Query not set")
	}

	// Check if it is a select query or not
	query_check := strings.Fields(ivquery)
	if strings.ToLower(query_check[0]) != "select" {
		context.SetOutput("result", "NOT_A_SELECT_QUERY")
		return true, fmt.Errorf("Query not a select query")
	}
	//////////////////////////////////////////////////

	log.Debugf("query" + ivquery)

	log.Debugf("All Parameters set")

	log.Debugf("Go SQL Connection Initiated...")

	db, err := sql.Open(ivdriverName, ivdatasourceName)
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	fmt.Println("Successfully Connected to Database")

	//////////////////////////////////////////////////////////

	f := make(map[int]interface{})
	var g = make(map[int]interface{})

	sNo := 0

	rows, queryerr := db.Query(ivquery)

	if queryerr != nil {
		panic(queryerr.Error())
	}

	cols, _ := rows.Columns()
	for rows.Next() {
		sNo += 1
		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			//return err
			panic(err.Error())
		}
		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
			m[colName] = fmt.Sprintf("%s", m[colName])
			jsonString, _ := json.Marshal(m)
			var resultinterface interface{}

			d := json.NewDecoder(bytes.NewReader(jsonString))
			d.UseNumber()
			err = d.Decode(&resultinterface)
			f = map[int]interface{}{sNo: resultinterface}

		}
		for k, v := range f {
			g[k] = v
		}

	}

	//Preparing the output result

	jsonString, _ := json.Marshal(g)
	var resultinterface interface{}
	d := json.NewDecoder(bytes.NewReader(jsonString))
	d.UseNumber()
	err = d.Decode(&resultinterface)
	h := map[string]interface{}{"results": resultinterface}

	//context.SetOutput(ovResult, h)
	jsonString1, _ := json.Marshal(h)
	//js := fmt.Sprintf("%v", jsonString1)

	js := string(jsonString1)

	context.SetOutput(ovResult, js)

	return true, nil
}
