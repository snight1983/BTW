package bitcoinvipsvr

import "fmt"

// 获取用户总帐
func getAddresPayInfo() error {
	rows, err := gDbconn.Query("SELECT id, recvaddress, unpay, paied FROM t_address_payinfo")
	if err != nil {
		fmt.Print(err)
		return err
	}
	defer rows.Close()
	for rows.Next() {
		m := &sAddressPay{}
		m.Ischange = false
		rows.Scan(&m.Dbid, &m.UserAddress, &m.UnPaid, &m.Paid, &m.CreateTime)
		gAddrPayInfoMap.Set(m.UserAddress, m)
	}
	if err = rows.Err(); err != nil {
		return err
	}
	gAddrPayInfoMap.isUpdate = false
	gAddrPayInfoMap.isInsert = false
	return err
}

// 更新总帐信息
func updateAddressPayInfo() error {
	tx, err := gDbconn.Begin()
	if err != nil {
		fmt.Print(err)
		return err
	}
	gAddrPayInfoMap.lock.Lock()
	defer gAddrPayInfoMap.lock.Unlock()
	for _, v := range gAddrPayInfoMap.bm {
		if v.(*sAddressPay).Ischange {
			if -1 != v.(*sAddressPay).Dbid {
				stmt, err := tx.Prepare("update t_address_payinfo set unpay=?,paied=? where recvaddress=?")
				if err != nil {
					fmt.Print(err)
					continue
				}
				defer stmt.Close()
				{
					_, err := stmt.Exec(v.(*sAddressPay).UnPaid, v.(*sAddressPay).Paid, v.(*sAddressPay).UserAddress)
					if err != nil {
						fmt.Print(err)
						continue
					}
					v.(*sAddressPay).Ischange = false
				}
			}
		}
	}
	tx.Commit()
	return nil
}

// 插入新的总帐条目
func insertAddressPayInfo() bool {
	gMinerRetMap.lock.Lock()
	defer gMinerRetMap.lock.Unlock()
	for _, v := range gMinerRetMap.bm {
		if v.(*Miner).Ischange {
			if -1 == v.(*Miner).Dbid {
				stmt, err := gDbconn.Prepare("insert t_address_payinfo set recvaddress=?,unpay=?,paied=?,createtm=?;")
				if err != nil {
					fmt.Print(err)
					continue
				}
				defer stmt.Close()
				result, err := stmt.Exec(v.(*sAddressPay).UserAddress, v.(*sAddressPay).UnPaid, v.(*sAddressPay).Paid, v.(*sAddressPay).CreateTime)
				if err != nil {
					fmt.Print(err)
					continue
				}
				id, err := result.LastInsertId()
				if err != nil {
					fmt.Print(err)
					continue
				}
				v.(*Miner).Dbid = id
			}
		}
	}
	return true
}
