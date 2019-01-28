package bitcoinvipsvr

import (
	"encoding/hex"
	"fmt"
	"net"
)

func msgShareReportID(pConn *net.UDPConn, paddr *net.UDPAddr, data []byte, posbeg int, datalen int) {
	buf := make([]byte, 32)
	copy(buf[0:32], data[posbeg:posbeg+32])
	posbeg += 32

	if posbeg <= datalen {
		return
	}

	hexHash := hex.EncodeToString(buf)
	share := gShareMap.Get(hexHash)
	if nil == share {
		share := &sShareData{}
		share.byhashPervBlock = make([]byte, 32)
		share.byhashMerkleRoot = make([]byte, 32)
		share.byhashSeed = make([]byte, 32)
		share.sHexShare = hexHash

		res, nVersion := readInt32(data, posbeg, datalen)
		if res == false {
			return
		}
		share.nVersion = nVersion
		posbeg += 4

		res, nHeight := readInt32(data, posbeg, datalen)
		if res == false || nHeight != gWorkHeader.nHeight {
			return
		}
		share.nHeight = nHeight
		posbeg += 4
		posend := posbeg + 32

		copy(share.byhashPervBlock[0:32], data[posbeg:posend])
		posbeg = posend
		posend += 32
		if posend > datalen {
			return
		}
		copy(share.byhashSeed[0:32], data[posbeg:posend])
		posbeg = posend
		res, share.nNonceLock = readUInt32(data, posbeg, datalen)
		if res == false {
			return
		}
		posbeg += 4
		res, share.nNonceBlock = readUInt32(data, posbeg, datalen)
		if res == false {
			return
		}
		posbeg += 4
		res, share.nNonceMrk = readUInt32(data, posbeg, datalen)
		if res == false {
			return
		}
		posbeg += 4

		share.sRecvAddress, posbeg, res = readString(data, posbeg, datalen)
		if res == false {
			return
		}
		res, nBit := readUInt32(data, posbeg, datalen)
		if res == false {
			return
		}
		if nBit >= gWorkHeader.nBitsBlock {
			///////////////////////////////////////////
			// br
		}
		share.sHexShare = hexHash
		share.nConfCnt = 0
		gShareMap.Set(hexHash, share)
		gShareQueue.Push(share)
	}
}

func msgShareCheckID(pConn *net.UDPConn, paddr *net.UDPAddr, data []byte, posbeg int, datalen int) {
	buf := make([]byte, 32)
	copy(buf[0:32], data[posbeg:posbeg+32])
	posbeg += 32
	if posbeg <= datalen {
		return
	}
	hexHash := hex.EncodeToString(buf)
	share := gShareMap.Get(hexHash)
	if nil != share {
		res, nConf := readInt8(data, posbeg, datalen)
		if res {
			if 1 == nConf {
				share.(*sShareData).nConfCnt++
				if share.(*sShareData).nConfCnt < 3 {
					gShareQueue.Push(share)
				}
			} else {
				gShareMap.Delete(hexHash)
			}
		}
	}
}

func insertShare(share *sShareData) bool {
	stmt, err := gDbconn.Prepare("insert t_miner_share set sharehash=?,address=?,height=?,ishandle=?;")
	if err != nil {
		fmt.Print(err)
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(share.sHexShare, share.sRecvAddress, share.nHeight, 0)
	if err != nil {
		fmt.Print(err)
		return false
	}
	return true
}
