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

	userInput := context.GetInput(username)

	ivuser, ok := userInput.(string)
	if !ok {
		context.SetOutput("result", "user_NOT_SET")
		return true, fmt.Errorf("user not set")
	}

	ivpasswd, ok := context.GetInput(password).(int)

	if !ok {
		context.SetOutput("result", "passwd_NOT_SET")
		return true, fmt.Errorf("passwd not set")
	}

	dbInput := context.GetInput(database)

	ivdb, ok := dbInput.(string)
	if !ok {
		context.SetOutput("result", "DATABASE_NOT_SET")
		return true, fmt.Errorf("DataBase not set")
	}

	log.Debugf("All variables set")

	log.Debugf("Go MYSQL Connection")

	db, err := sql.Open("mysql", "flogo:password@tcp(localhost:3306)/testdb")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	log.Debugf("Successfully Connected to MySQL Database")

	return true, nil
}
