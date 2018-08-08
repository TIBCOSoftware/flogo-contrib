package mysqlfetch

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	// Imports the MySQL package so it can be used as a driver
	_ "github.com/go-sql-driver/mysql"

	"bytes"
	"encoding/json"
)

// log is the default logger for Project Flogo
var log = logger.GetLogger("activity-mysqlfetch")

const (
	ivHost     = "host"
	ivUsername = "username"
	ivPassword = "password"
	ivDatabase = "database"
	ivQuery    = "query"

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

	host := context.GetInput(ivHost).(string)
	if len(host) == 0 {
		log.Errorf("The host is not set")
		return true, fmt.Errorf("The host is not set")
	}
	log.Debugf("Host: [%s]", host)

	username := context.GetInput(ivUsername).(string)
	if len(username) == 0 {
		log.Errorf("The username is not set")
		return true, fmt.Errorf("The username is not set")
	}
	log.Debugf("Username: [%s]", username)

	password := context.GetInput(ivPassword).(string)
	if len(password) == 0 {
		log.Errorf("The password is not set")
		return true, fmt.Errorf("The password is not set")
	}
	log.Debugf("Password is set...")

	database := context.GetInput(ivDatabase).(string)
	if len(database) == 0 {
		log.Errorf("The database is not set")
		return true, fmt.Errorf("The database is not set")
	}
	log.Debugf("Database: [%s]", database)

	query := context.GetInput(ivQuery).(string)
	if len(query) == 0 {
		log.Errorf("The query is not set")
		return true, fmt.Errorf("The query is not set")
	}
	log.Debugf("Query: [%s]", query)

	connection := username + ":" + password + "@tcp(" + host + ")/" + database

	db, err := sql.Open("mysql", connection)
	if err != nil {
		return true, err
	}

	defer db.Close()

	f := make(map[string]interface{})
	// f := make(map[int]interface{})
	g := make(map[string]interface{})
	// g := make(map[int]interface{})

	sNo := 0

	rows, queryerr := db.Query(query)

	if queryerr != nil {
		return true, queryerr
	}

	cols, _ := rows.Columns()
	for rows.Next() {
		sNo++
		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			return true, err
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

			rowNo := "Row"+ strconv.Itoa(123)
			//f = map[int]interface{}{sNo: resultinterface}
			f = map[string]interface{}{rowNo: resultinterface}

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
	context.SetOutput(ovResult, h)

	return true, nil
}
