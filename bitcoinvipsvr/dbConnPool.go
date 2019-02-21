package bitcoinvipsvr

import (
	"database/sql"
	"fmt"
)

var gDbconn *sql.DB

func onInitDbConnectPool() error {
	var err error
	gDbconn, err = sql.Open(gConfig.DBNAME, gConfig.DBCONNECT)
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

func getBlockIncome() ([]*sDbBlocks, error) {
	if nil != gDbconn {
		var oBlocks []*sDbBlocks
		rows, err := gDbconn.Query("select id, txid, amount, finish from t_miner_block where finish=0")
		if err != nil {
			fmt.Print(err)
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			block := &sDbBlocks{}
			rows.Scan(&block.n64Blockid, &block.sTxID, &block.nu64Amount, &block.n32IsHandle)
			oBlocks = append(oBlocks, block)
		}
		if err = rows.Err(); err != nil {
			return nil, err
		}
		return oBlocks, err
	}
	return nil, nil
}

func updateBlockIncome(block *sDbBlocks) error {
	if nil != gDbconn {
		stmt, err := gDbconn.Prepare("update t_miner_block set txid=?,amount=?,finish=? where id=?")
		if err != nil {
			fmt.Print(err)
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(block.sTxID, block.nu64Amount, block.n32IsHandle, block.n64Blockid)
		if err != nil {
			fmt.Print(err)
			return err
		}
	}
	return nil
}

func getAddressIncome() (map[string]*sAddIncome, error) {
	if nil != gDbconn {
		addrIcomeMap := make(map[string]*sAddIncome)
		rows, err := gDbconn.Query("select id, recvaddress, unpay, paied, speed, sharecur, sharetotal,createtm from t_miner_income")
		if err != nil {
			fmt.Print(err)
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			oIncome := &sAddIncome{}
			err := rows.Scan(&oIncome.nID, &oIncome.sRecvaddress, &oIncome.nu64Unpay, &oIncome.nu64Paied, &oIncome.f32Speed, &oIncome.n32Sharecur, &oIncome.n32Sharetotal, &oIncome.nu64Createtm)
			fmt.Print(err)
			if nil == err {
				oIncome.bIsInDB = true
				addrIcomeMap[oIncome.sRecvaddress] = oIncome
			}
		}
		if err = rows.Err(); err != nil {
			return nil, err
		}
		return addrIcomeMap, err
	}
	return nil, nil
}

func updateAddrIncome(oIncomeMap map[string]*sAddIncome) error {
	if nil != gDbconn {
		tx, err := gDbconn.Begin()
		if nil != err {
			fmt.Print(err)
			return err
		}
		for _, income := range oIncomeMap {
			if income.bIsInDB {
				stmt, err := tx.Prepare("update t_miner_income set unpay=?,speed=?,sharecur=? where id=?")
				if err != nil {
					fmt.Print(err)
					continue
				}
				defer stmt.Close()
				_, err = stmt.Exec(income.nu64Unpay, income.f32Speed, income.n32Sharecur, income.nID)
				if err != nil {
					fmt.Print(err)
					continue
				}
			} else {
				stmt, err := tx.Prepare("insert t_miner_income set recvaddress=?,unpay=?,paied=?,speed=?,sharecur=?,sharetotal=?,createtm=?")
				if err != nil {
					fmt.Print(err)
					continue
				}
				defer stmt.Close()
				_, err = stmt.Exec(income.sRecvaddress, income.nu64Unpay, income.nu64Paied, income.f32Speed, income.n32Sharecur, income.n32Sharetotal, income.nu64Createtm)
				if err != nil {
					fmt.Print(err)
					continue
				}
			}
		}
		return tx.Commit()
	}
	return nil
}

func getMinerShare() ([]*sMinerShare, error) {
	if nil != gDbconn {
		var oShares []*sMinerShare
		rows, err := gDbconn.Query("select id, address from t_miner_share where ishandle=0")
		if err != nil {
			fmt.Print(err)
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			share := &sMinerShare{}
			rows.Scan(&share.nID, &share.sAddr)
			oShares = append(oShares, share)
		}
		if err = rows.Err(); err != nil {
			return nil, err
		}
		return oShares, err
	}
	return nil, nil
}

func updateMinerShare(allShare []*sMinerShare) error {
	if nil != gDbconn {
		tx, err := gDbconn.Begin()
		if nil != err {
			fmt.Print(err)
			return err
		}
		for _, share := range allShare {
			stmt, err := tx.Prepare("update t_miner_share set ishandle=1 where id=?")
			if err != nil {
				fmt.Print(err)
				continue
			}
			defer stmt.Close()
			_, err = stmt.Exec(share.nID)
			if err != nil {
				fmt.Print(err)
				continue
			}
		}
		return tx.Commit()
	}
	return nil
}

func insertPoolStatu(addrcnt int, speed float32, tm int64, total uint64, pool uint64) error {
	if nil != gDbconn {
		stmt, err := gDbconn.Prepare("insert t_miner_pool_statu set addrcnt=?,speed=?,createtm=?,totalincome=?,poolincome=?;")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(addrcnt, speed, tm, total, pool)
		if err != nil {
			return err
		}
	}
	return nil
}
