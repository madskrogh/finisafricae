package mysql

import (
	"database/sql"

	"github.com/madskrogh/finisafricae/util"
)

//InitDB creates and/or remakes mysql necesarry tables for the given database db
func InitDB(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS user(id varchar(64), uname varchar(32), email varchar(32), password varchar(64));")
	util.HandleError(err)
	_, err = db.Exec("DROP TABLE session;")
	util.HandleError(err)
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS session(id varchar(64), userid varchar(64), time varchar(64));")
	util.HandleError(err)
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS book(id varchar(64), userid varchar(64), title varchar(32), author varchar(32), year varchar(32), genre varchar(32), notes varchar(32));")
	util.HandleError(err)
}
