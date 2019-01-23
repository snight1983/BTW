package bitcoinvipsvr

import (
	"fmt"
	"net"
	"time"
)

// 收到网络信息，更新注册信息
func netMinerRegedit(conn net.UDPAddr, data []byte, posbeg int, datalen int) bool {
	UserID, endpos, res := readString(data, posbeg, datalen)
	if res {
		UserName, endpos, res := readString(data, endpos, datalen)
		if res {
			RecvAddress, _, res := readString(data, endpos, datalen)
			if res {
				miner := gMinerRetMap.Get(UserID)
				tm := time.Now().Unix()
				if miner == nil {
					m := &Miner{-1, UserID, UserName, RecvAddress, tm, true, 0}
					gMinerRetMap.Set(m.UserID, m)
					gMinerRetMap.isInsert = true
				} else {
					if miner.(*Miner).UserName != UserName ||
						miner.(*Miner).RecvAddress != RecvAddress {
						miner.(*Miner).UserName = UserName
						miner.(*Miner).RecvAddress = RecvAddress
						miner.(*Miner).Ischange = true
						gMinerRetMap.isUpdate = true
					}
				}
				return true
			}
		}
	}
	return false
}

// 获取用户注册信息
func getMinerRegedit() error {
	if nil != gDbconn {
		rows, err := gDbconn.Query("select id, userid, recvaddress, username, createtm from t_miner_regedit")
		if err != nil {
			fmt.Print(err)
			return err
		}
		defer rows.Close()
		for rows.Next() {
			m := &Miner{}
			rows.Scan(&m.Dbid, &m.UserID, &m.RecvAddress, &m.UserName, &m.CreateTime)
			m.Ischange = false
			gMinerRetMap.Set(m.UserID, m)
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

// 更新挖矿注册用户
func updateMinerRegedit() error {
	tx, err := gDbconn.Begin()
	if nil != err {
		fmt.Print(err)
		return err
	}
	gMinerRetMap.lock.Lock()
	defer gMinerRetMap.lock.Unlock()

	for _, v := range gMinerRetMap.bm {
		if v.(*Miner).Ischange {
			if -1 != v.(*Miner).Dbid {
				stmt, err := tx.Prepare("update t_miner_regedit set userid=?,recvaddress=?,username=?,createtm=? where id=?")
				if err != nil {
					fmt.Print(err)
					continue
				}
				defer stmt.Close()
				{
					_, err := stmt.Exec(v.(*Miner).UserID, v.(*Miner).RecvAddress, v.(*Miner).UserName, v.(*Miner).CreateTime, v.(*Miner).Dbid)
					if err != nil {
						fmt.Print(err)
						continue
					}
					v.(*Miner).Ischange = false
				}
			}
		}
	}
	tx.Commit()
	return err
}

// 插入挖矿用户注册
func insertMinerRegedit() bool {
	gMinerRetMap.lock.Lock()
	defer gMinerRetMap.lock.Unlock()
	for _, v := range gMinerRetMap.bm {
		if v.(*Miner).Ischange {
			if -1 == v.(*Miner).Dbid {
				stmt, err := gDbconn.Prepare("insert t_miner_regedit set userid=?,recvaddress=?,username=?,createtm=?;")
				if err != nil {
					fmt.Print(err)
					return false
				}
				defer stmt.Close()
				{
					result, err := stmt.Exec(v.(*Miner).UserID, v.(*Miner).RecvAddress, v.(*Miner).UserName, v.(*Miner).CreateTime)
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
						v.(*Miner).Dbid = id
					}
				}
			}
		}
	}
	return true
}
