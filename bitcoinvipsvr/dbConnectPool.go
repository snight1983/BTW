package bitcoinvipsvr

import (
	"database/sql"
	"fmt"
)

var gdbconn *sql.DB

// 连接数据库
func onInitDbConnectPool() (bool, error) {
	var err error
	gdbconn, err = sql.Open("mysql", "vipmpool:vipmpool_123A#$@tcp(10.9.0.122:3306)/vipmpool?charset=utf8")
	if err != nil {
		fmt.Print(err)
		return false, err
	}
	//checkErr()
	gdbconn.SetMaxOpenConns(2000)
	gdbconn.SetMaxIdleConns(1000)
	err = gdbconn.Ping()
	if err != nil {
		fmt.Print(err)
		return false, err
	}
	return true, err
}
