package bitcoinvipsvr

import (
	"fmt"
	"net"
	"time"
)

func netMinerRegedit(paddr *net.UDPAddr, data []byte, posbeg int, datalen int) bool {
	UserID, endpos, res := readString(data, posbeg, datalen)
	if res {
		UserName, endpos, res := readString(data, endpos, datalen)
		if res {
			RecvAddress, _, res := readString(data, endpos, datalen)
			if res {
				miner := gMinerRetMap.Get(UserID)
				tm := time.Now().Unix()
				if miner == nil {
					m := &Miner{-1, UserID, UserName, RecvAddress, tm, true}
					gMinerRetMap.Set(m.sUserID, m)
					gMinerRetMap.isInsert = true
				} else {
					if miner.(*Miner).sUserName != UserName ||
						miner.(*Miner).sRecvAddress != RecvAddress {
						miner.(*Miner).sUserName = UserName
						miner.(*Miner).sRecvAddress = RecvAddress
						miner.(*Miner).bIschange = true
						gMinerRetMap.isUpdate = true
					}
				}
				return true
			}
		}
	}
	return false
}

func getMinerRegedit() error {
	if nil != gDbconn {
		rows, err := gDbconn.Query("select id, userid, recvaddress, username from t_miner_regedit")
		if err != nil {
			fmt.Print(err)
			return err
		}
		defer rows.Close()
		for rows.Next() {
			m := &Miner{}
			rows.Scan(&m.n64Dbid, &m.sUserID, &m.sRecvAddress, &m.sUserName)
			m.bIschange = false
			gMinerRetMap.Set(m.sUserID, m)
		}
		if err = rows.Err(); err != nil {
			return err
		}
		gMinerRetMap.isUpdate = false
		gMinerRetMap.isInsert = false
		return err
	}
	return nil
}

func updateMinerRegedit() error {
	tx, err := gDbconn.Begin()
	if nil != err {
		fmt.Print(err)
		return err
	}
	gMinerRetMap.lock.Lock()
	defer gMinerRetMap.lock.Unlock()

	for _, v := range gMinerRetMap.bm {
		if v.(*Miner).bIschange {
			if -1 != v.(*Miner).n64Dbid {
				stmt, err := tx.Prepare("update t_miner_regedit set userid=?,recvaddress=?,username=? where id=?")
				if err != nil {
					fmt.Print(err)
					continue
				}
				defer stmt.Close()
				{
					_, err := stmt.Exec(v.(*Miner).sUserID, v.(*Miner).sRecvAddress, v.(*Miner).sUserName, v.(*Miner).n64Dbid)
					if err != nil {
						fmt.Print(err)
						continue
					}
					v.(*Miner).bIschange = false
				}
			}
		}
	}
	tx.Commit()
	return err
}

func insertMinerRegedit() bool {
	gMinerRetMap.lock.Lock()
	defer gMinerRetMap.lock.Unlock()
	for _, v := range gMinerRetMap.bm {
		if v.(*Miner).bIschange {
			if -1 == v.(*Miner).n64Dbid {
				stmt, err := gDbconn.Prepare("insert t_miner_regedit set userid=?,recvaddress=?,username=?,createtm=?;")
				if err != nil {
					fmt.Print(err)
					return false
				}
				defer stmt.Close()
				{
					result, err := stmt.Exec(v.(*Miner).sUserID, v.(*Miner).sRecvAddress, v.(*Miner).sUserName, v.(*Miner).n64CreateTime)
					if err != nil {
						fmt.Print(err)
						return false
					}
					{
						id, err := result.LastInsertId()
						if err != nil {
							fmt.Print(err)
							return false
						}
						v.(*Miner).n64Dbid = id
					}
				}
			}
		}
	}
	return true
}
