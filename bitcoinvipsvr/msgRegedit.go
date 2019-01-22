package bitcoinvipsvr

import (
	"fmt"
	"net"
	"time"
)

// 收到网络信息，更新注册信息
func netMinerRegedit(conn net.UDPAddr, data []byte, posbeg int, datalen int) bool {

	res0, value0 := readInt32(data, posbeg, datalen)
	posbeg += 4
	posend := posbeg + int(value0)

	if res0 && (posend <= datalen) && (posend > posbeg) {
		UserID := string(data[posbeg:posend])
		posbeg = posend
		posend = posbeg + 4

		if posend <= datalen {
			res1, value1 := readInt32(data, posbeg, datalen)
			posbeg = posend
			posend = posbeg + int(value1)

			if res1 && (posend <= datalen) && (posend > posbeg) {
				UserName := string(data[posbeg:posend])
				posbeg = posend
				posend = posbeg + 4

				if posend <= datalen {
					res2, value2 := readInt32(data, posbeg, datalen)
					posbeg = posend
					posend += int(value2)

					if res2 && (posend <= datalen) && (posend > posbeg) {
						UserAddress := string(data[posbeg:posend])
						posbeg = posend
						posend = posbeg + 4

						if posend <= datalen {
							res3, value3 := readInt32(data, posbeg, datalen)
							posbeg = posend
							posend += int(value3)

							if res3 && (posend <= datalen) && (posend > posbeg) {
								UserDisName := string(data[posbeg:posend])

								miner := gMinerRetMap.Get(UserID)
								tm := time.Now().Unix()

								if miner == nil {
									m := &Miner{}
									m.UserID = UserID
									m.UserName = UserName
									m.UserAddress = UserAddress
									m.UserDisName = UserDisName
									m.Ischange = true
									m.LastTime = tm
									m.CreateTime = tm
									m.Dbid = -1
									gMinerRetMap.Set(m.UserID, m)

								} else {
									if tm-miner.(*Miner).LastTime > 5 {
										if miner.(*Miner).UserName != UserName ||
											miner.(*Miner).UserAddress != UserAddress ||
											miner.(*Miner).UserDisName != UserDisName {
											miner.(*Miner).UserName = UserName
											miner.(*Miner).UserAddress = UserAddress
											miner.(*Miner).UserDisName = UserDisName
											miner.(*Miner).Ischange = true
											gMinerRetMap.isUpdate = true
										}
										miner.(*Miner).LastTime = tm
									}
								}
								return true
							}
						}
					}
				}
			}
		}
	}
	return false
}

// 获取用户注册信息
func getMinerRegedit() (bool, error) {
	rows, err := gdbconn.Query("SELECT id, usermid, useraddress, username, createtm FROM t_miner_regedit")
	if err != nil {
		fmt.Print(err)
		return false, err
	}
	defer rows.Close()
	for rows.Next() {
		m := &Miner{}
		rows.Scan(&m.Dbid, &m.UserID, &m.UserAddress, &m.UserDisName, &m.CreateTime)
		m.Ischange = false
		m.LastTime = 0
		gMinerRetMap.Set(m.UserID, m)
	}
	if err = rows.Err(); err != nil {
		return false, err
	}
	gMinerRetMap.isUpdate = false
	gMinerRetMap.isInsert = false
	return true, err
}

// 更新挖矿注册用户
func updateMinerRegedit() bool {
	tx, _ := gdbconn.Begin()
	gMinerRetMap.lock.Lock()
	defer gMinerRetMap.lock.Unlock()

	for _, v := range gMinerRetMap.bm {
		if v.(*Miner).Ischange {
			if -1 != v.(*Miner).Dbid {
				stmt, err := tx.Prepare("update t_miner_regedit set usermid=?,useraddress=?,username=?,createtm=? where id=?")
				if err == nil {
					_, err := stmt.Exec(v.(*Miner).UserID, v.(*Miner).UserAddress, v.(*Miner).UserDisName, v.(*Miner).CreateTime, v.(*Miner).Dbid)
					if err == nil {
						v.(*Miner).Ischange = false
						//fmt.Printf("info:update new item %d uid:%s \n", Miner.Dbid, Miner.UserID)
					}
				}
			}
		}
	}
	tx.Commit()

	return true
}

// 插入挖矿用户注册
func insertMinerRegedit() bool {
	gMinerRetMap.lock.Lock()
	defer gMinerRetMap.lock.Unlock()
	for _, v := range gMinerRetMap.bm {
		if v.(*Miner).Ischange {
			if -1 == v.(*Miner).Dbid {
				stmt, err := gdbconn.Prepare("insert t_miner_regedit set usermid=?,useraddress=?,username=?,createtm=?;")
				if err != nil {
					fmt.Print(err)
					return false
				}
				result, err := stmt.Exec(v.(*Miner).UserID, v.(*Miner).UserAddress, v.(*Miner).UserDisName, v.(*Miner).CreateTime)
				if err != nil {
					fmt.Print(err)
					return false
				}

				id, err := result.LastInsertId()
				if err != nil {
					fmt.Print(err)
					return false
				}
				v.(*Miner).Dbid = id
				//fmt.Printf("info:Insert new item %d uid:%s \n", Miner.Dbid, Miner.UserID)
			}
		}
	}
	return true
}
