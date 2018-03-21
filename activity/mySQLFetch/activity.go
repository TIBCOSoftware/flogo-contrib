package mySQLFetch

import (
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// log is the default package logger
var log = logger.GetLogger("activity-akash-mySql")

const (
	host     = "host"
	username = "username"
	password = "password"
	database = "database"
	query	 = "query"
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

	log.Debugf("Go MYSQL Connection")

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

	rows, _ := db.Query(ivquery)
	fmt.Println("Check " + ivquery + "\n" + conn_str)

	return true, nil
}
