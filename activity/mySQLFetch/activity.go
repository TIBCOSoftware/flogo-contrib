package mySQLFetch

import (
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"bytes"
	"encoding/json"
)

// log is the default package logger
var log = logger.GetLogger("activity-akash-mySql")

const (
	host     = "host"
	username = "username"
	password = "password"
	database = "database"
	query	 = "query"

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

	hostInput := context.GetInput(host)

	ivhost, ok := hostInput.(string)
	if !ok {
		context.SetOutput("result", "HOSTSET")
		return true, fmt.Errorf("Host not set")
	}
	log.Debugf("Hostname" + ivhost)

	userInput := context.GetInput(username)

	ivuser, ok := userInput.(string)
	if !ok {
		context.SetOutput("result", "user_NOT_SET")
		return true, fmt.Errorf("user not set")
	}
	log.Debugf("username" + ivuser)

	ivpasswd, ok := context.GetInput(password).(string)

	if !ok {
		context.SetOutput("result", "passwd_NOT_SET")
		return true, fmt.Errorf("passwd not set")
	}
	log.Debugf("password" + ivpasswd)
	

	dbInput := context.GetInput(database)

	ivdb, ok := dbInput.(string)
	if !ok {
		context.SetOutput("result", "DATABASE_NOT_SET")
		return true, fmt.Errorf("DataBase not set")
	}
	log.Debugf("database" + ivdb)

	queryInput := context.GetInput(query)

	ivquery, ok := queryInput.(string)
	if !ok {
		context.SetOutput("result", "QUERY_NOT_SET")
		return true, fmt.Errorf("Query not set")
	}
	log.Debugf("query" + ivquery)
	

	log.Debugf("All variables set")

	log.Debugf("Go MYSQL Connection Initiated")

	conn_str := ivuser + ":" +ivpasswd + "@tcp(" + ivhost + ")/" + ivdb
	log.Debugf(conn_str)

//db, err := sql.Open("mysql", "flogo:password@tcp(localhost:3306)/testdb")
	db, err := sql.Open("mysql", conn_str)
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	log.Debugf("Successfully Connected to MySQL Database")

	//////////////////////////////////////////////////////////
	f := make(map[int]interface{})
	var g = make(map[int]interface{})

	sNo := 0

	rows, _ := db.Query(ivquery)

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
		//m[colName] = fmt.Sprintf("%s",m[colName])
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
	
	context.SetOutput(ovResult, h)

	return true, nil
}
