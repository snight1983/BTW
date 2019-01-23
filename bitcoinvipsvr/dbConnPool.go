package bitcoinvipsvr

import (
	"database/sql"
	"fmt"
)

var gDbconn *sql.DB

func onInitDbConnectPool() error {
	//temp , from conf later
	var err error
	gDbconn, err = sql.Open("mysql", "vipmpool:vipmpool_123A#$@tcp(10.9.0.122:3306)/vipmpool?charset=utf8")
	if err != nil {
		fmt.Print(err)
		return err
	}
	gDbconn.SetMaxOpenConns(2000)
	gDbconn.SetMaxIdleConns(1000)
	err = gDbconn.Ping()
	if err != nil {
		fmt.Print(err)
	}
	return err
}
