package bitcoinvipsvr

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
	"time"
)

func msgShareReportRQ(pConn *net.UDPConn, paddr *net.UDPAddr, data []byte, posbeg int, datalen int) {
	buf := make([]byte, 32)
	posend := posbeg + 32
	if posend >= datalen {
		return
	}
	copy(buf[0:32], data[posbeg:posend])
	posbeg = posend
	hexHash := hex.EncodeToString(buf)
	share := gShareMap.Get(hexHash)
	if nil == share {
		share := &sShareData{}
		share.byhashPervBlock = make([]byte, 32)
		share.byhashMerkleRoot = make([]byte, 32)
		share.byhashSeed = make([]byte, 32)
		share.sHexShare = hexHash
		res, nVersion := readInt32(data, posbeg, datalen)
		posbeg += 4
		if false == res || posbeg > datalen {
			return
		}
		share.nVersion = nVersion

		res, nHeight := readInt32(data, posbeg, datalen)
		posbeg += 4
		if (posbeg > datalen) || (false == res) || (nHeight != gWorkHeader.nHeight) {
			return
		}
		share.nHeight = nHeight
		posend = posbeg + 32
		if posend > datalen {
			return
		}
		copy(share.byhashPervBlock[0:32], data[posbeg:posend])
		posbeg = posend
		posend += 32

		if posend > datalen {
			return
		}
		copy(share.byhashSeed[0:32], data[posbeg:posend])
		posbeg = posend
		res, share.nNonceMrk = readUInt32(data, posbeg, datalen)
		posbeg += 4

		if false == res || posbeg > datalen {
			return
		}
		posend = posbeg + 32

		if posend > datalen {
			return
		}
		copy(share.byhashMerkleRoot[0:32], data[posbeg:posend])
		posbeg = posend

		res, share.nTime = readUInt32(data, posbeg, datalen)
		posbeg += 4
		if false == res || posbeg > datalen {
			return
		}

		res, share.nBit = readUInt32(data, posbeg, datalen)
		posbeg += 4
		if res == false || posbeg > datalen {
			return
		}

		res, share.nNonceLock = readUInt32(data, posbeg, datalen)
		posbeg += 4
		if res == false || posbeg > datalen {
			return
		}

		res, share.nNonceBlock = readUInt32(data, posbeg, datalen)
		posbeg += 4
		if false == res || posbeg > datalen {
			return
		}

		res, share.nType = readUInt32(data, posbeg, datalen)
		posbeg += 4
		if false == res || posbeg > datalen {
			return
		}

		share.sRecvAddress, posbeg, res = readString(data, posbeg, datalen)
		if res == false {
			return
		}

		if share.nType == 1 {
			msgSendBuffer := make([]byte, 0)

			byMsgID := bytes.NewBuffer([]byte{})
			binary.Write(byMsgID, binary.LittleEndian, gMsgShareRQ)
			msgSendBuffer = append(msgSendBuffer, byMsgID.Bytes()...)

			byMrkNnoce := bytes.NewBuffer([]byte{})
			binary.Write(byMrkNnoce, binary.LittleEndian, share.nNonceMrk)
			msgSendBuffer = append(msgSendBuffer, byMrkNnoce.Bytes()...)

			msgSendBuffer = append(msgSendBuffer, share.byhashMerkleRoot...)
			msgSendBuffer = append(msgSendBuffer, share.byhashSeed...)

			byTime := bytes.NewBuffer([]byte{})
			binary.Write(byTime, binary.LittleEndian, share.nTime)
			msgSendBuffer = append(msgSendBuffer, byTime.Bytes()...)

			byNonceLock := bytes.NewBuffer([]byte{})
			binary.Write(byNonceLock, binary.LittleEndian, share.nNonceLock)
			msgSendBuffer = append(msgSendBuffer, byNonceLock.Bytes()...)

			byNonceBlock := bytes.NewBuffer([]byte{})
			binary.Write(byNonceBlock, binary.LittleEndian, share.nNonceBlock)
			msgSendBuffer = append(msgSendBuffer, byNonceBlock.Bytes()...)

			conn, err := net.Dial("udp", gConfig.POOLWALLETSYNC)
			defer conn.Close()
			if err == nil {
				conn.Write(msgSendBuffer)
			}
		}
		share.sHexShare = hexHash
		share.nConfCnt = 0
		gShareMap.Set(hexHash, share)
		fmt.Printf("%s | Pool: Push Share New:%s\n", time.Now().Format("2006-01-02 15:04:05"), hexHash)
		gShareQueue.Push(share)
	}

	byMsgSendBuffer := make([]byte, 0)
	byMsgID := bytes.NewBuffer([]byte{})
	binary.Write(byMsgID, binary.LittleEndian, gMsgShareCheckRQ)
	byMsgSendBuffer = append(byMsgSendBuffer, byMsgID.Bytes()...)
	pConn.WriteToUDP(byMsgSendBuffer, paddr)
}

func msgShareCheckRQ(pConn *net.UDPConn, paddr *net.UDPAddr, data []byte, posbeg int, datalen int) {
	share := gShareQueue.Pop()
	if nil != share {
		fmt.Printf("%s | Pool: msgShareCheckRQ QueueSize:%d\n", time.Now().Format("2006-01-02 15:04:05"), gShareQueue.size())
		msgSendBuffer := make([]byte, 0)
		byMsgID := bytes.NewBuffer([]byte{})

		binary.Write(byMsgID, binary.LittleEndian, gMsgShareCheckRS)
		msgSendBuffer = append(msgSendBuffer, byMsgID.Bytes()...)

		byVersion := bytes.NewBuffer([]byte{})
		binary.Write(byVersion, binary.LittleEndian, share.(*sShareData).nVersion)
		msgSendBuffer = append(msgSendBuffer, byVersion.Bytes()...)

		bHeight := bytes.NewBuffer([]byte{})
		binary.Write(bHeight, binary.LittleEndian, share.(*sShareData).nHeight)
		msgSendBuffer = append(msgSendBuffer, bHeight.Bytes()...)

		bynNonceLock := bytes.NewBuffer([]byte{})
		binary.Write(bynNonceLock, binary.LittleEndian, share.(*sShareData).nNonceLock)
		msgSendBuffer = append(msgSendBuffer, bynNonceLock.Bytes()...)

		bynNonceBlock := bytes.NewBuffer([]byte{})
		binary.Write(bynNonceBlock, binary.LittleEndian, share.(*sShareData).nNonceBlock)
		msgSendBuffer = append(msgSendBuffer, bynNonceBlock.Bytes()...)

		bynTime := bytes.NewBuffer([]byte{})
		binary.Write(bynTime, binary.LittleEndian, share.(*sShareData).nTime)
		msgSendBuffer = append(msgSendBuffer, bynTime.Bytes()...)

		bynBit := bytes.NewBuffer([]byte{})
		binary.Write(bynBit, binary.LittleEndian, share.(*sShareData).nBit)
		msgSendBuffer = append(msgSendBuffer, bynBit.Bytes()...)

		msgSendBuffer = append(msgSendBuffer, share.(*sShareData).byhashPervBlock...)
		msgSendBuffer = append(msgSendBuffer, share.(*sShareData).byhashMerkleRoot...)
		msgSendBuffer = append(msgSendBuffer, share.(*sShareData).byhashSeed...)

		pConn.WriteToUDP(msgSendBuffer, paddr)
	}
}

func msgShareCheckID(pConn *net.UDPConn, paddr *net.UDPAddr, data []byte, posbeg int, datalen int) {
	buf := make([]byte, 32)
	copy(buf[0:32], data[posbeg:posbeg+32])
	posbeg += 32
	if posbeg >= datalen {
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

func insertBlock(txid string, amount int64) bool {
	if nil != gDbconn {
		stmt, err := gDbconn.Prepare("insert t_miner_block set txid=?,amount=?,finish=?;")
		if err != nil {
			return false
		}
		defer stmt.Close()
		_, err = stmt.Exec(txid, amount, 0)
		if err != nil {
			return false
		}
		return true
	}
	return false
}
