package bitcoinvipsvr

import "fmt"

func getAddresPayInfo() error {
	rows, err := gDbconn.Query("select id, recvaddress, unpay, paied from t_address_payinfo")
	if err != nil {
		fmt.Print(err)
		return err
	}
	defer rows.Close()
	for rows.Next() {
		m := &sAddressPay{}
		m.bIschange = false
		rows.Scan(&m.n64Dbid, &m.sUserAddress, &m.n64UnPaid, &m.n64Paid, &m.n64CreateTime)
		gAddrPayInfoMap.Set(m.sUserAddress, m)
	}
	if err = rows.Err(); err != nil {
		return err
	}
	gAddrPayInfoMap.isUpdate = false
	gAddrPayInfoMap.isInsert = false
	return err
}

func updateAddressPayInfo() error {
	tx, err := gDbconn.Begin()
	if err != nil {
		fmt.Print(err)
		return err
	}
	gAddrPayInfoMap.lock.Lock()
	defer gAddrPayInfoMap.lock.Unlock()
	for _, v := range gAddrPayInfoMap.bm {
		if v.(*sAddressPay).bIschange {
			if -1 != v.(*sAddressPay).n64Dbid {
				stmt, err := tx.Prepare("update t_address_payinfo set unpay=?,paied=? where recvaddress=?")
				if err != nil {
					fmt.Print(err)
					continue
				}
				defer stmt.Close()
				{
					_, err := stmt.Exec(v.(*sAddressPay).n64UnPaid, v.(*sAddressPay).n64Paid, v.(*sAddressPay).sUserAddress)
					if err != nil {
						fmt.Print(err)
						continue
					}
					v.(*sAddressPay).bIschange = false
				}
			}
		}
	}
	tx.Commit()
	return nil
}

func insertAddressPayInfo() bool {
	gMinerRetMap.lock.Lock()
	defer gMinerRetMap.lock.Unlock()
	for _, v := range gMinerRetMap.bm {
		if v.(*sMiner).bIschange {
			if -1 == v.(*sMiner).n64Dbid {
				stmt, err := gDbconn.Prepare("insert t_address_payinfo set recvaddress=?,unpay=?,paied=?,createtm=?;")
				if err != nil {
					fmt.Print(err)
					continue
				}
				defer stmt.Close()
				result, err := stmt.Exec(v.(*sAddressPay).sUserAddress, v.(*sAddressPay).n64UnPaid, v.(*sAddressPay).n64Paid, v.(*sAddressPay).n64CreateTime)
				if err != nil {
					fmt.Print(err)
					continue
				}
				id, err := result.LastInsertId()
				if err != nil {
					fmt.Print(err)
					continue
				}
				v.(*sMiner).n64Dbid = id
			}
		}
	}
	return true
}
