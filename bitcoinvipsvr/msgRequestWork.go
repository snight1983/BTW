package bitcoinvipsvr

import (
	"bytes"
	"encoding/binary"
	"net"
)

func msgMiningWorkRQ(pConn *net.UDPConn, paddr *net.UDPAddr, data []byte, posbeg int, datalen int) {
	var bNeedMkr bool
	ret, nHeight := readInt32(data, posbeg, datalen)
	if ret {
		if nHeight != gWorkHeader.nHeight {
			bNeedMkr = true
		} else {
			posbeg += 4
			ret, n8NeedMkr := readInt8(data, posbeg, datalen)
			if ret && n8NeedMkr == 1 {
				bNeedMkr = true
			}
		}
		msgSendBuffer := make([]byte, 0)
		msgID := bytes.NewBuffer([]byte{})
		binary.Write(msgID, binary.LittleEndian, gMsgMiningWorkRS)
		msgSendBuffer = append(msgSendBuffer, msgID.Bytes()...)
		msgSendBuffer = append(msgSendBuffer, gHeader...)
		if bNeedMkr {
			mkrBuf := gMkrQueue.Pop()
			msgSendBuffer = append(msgSendBuffer, mkrBuf.([]byte)...)
		}
		lens := len(msgSendBuffer)
		lens++
		pConn.WriteToUDP(msgSendBuffer, paddr)
	}
}

func msgCheckWorkRQ(pConn *net.UDPConn, paddr *net.UDPAddr, data []byte, posbeg int, datalen int) {
	share := gShareQueue.Pop()
	if nil != share {
		msgSendBuffer := make([]byte, 116)
		msgID := bytes.NewBuffer([]byte{})
		binary.Write(msgID, binary.LittleEndian, gMsgCheckWorkRS)
		copy(msgSendBuffer[0:4], msgID.Bytes()[0:4])
		nVersion := bytes.NewBuffer([]byte{})
		binary.Write(nVersion, binary.LittleEndian, share.(sShareData).nVersion)
		copy(msgSendBuffer[4:8], nVersion.Bytes()[0:4])
		nHeight := bytes.NewBuffer([]byte{})
		binary.Write(nHeight, binary.LittleEndian, share.(sShareData).nHeight)
		copy(msgSendBuffer[8:12], nHeight.Bytes()[0:4])
		nNonceLock := bytes.NewBuffer([]byte{})
		binary.Write(nNonceLock, binary.LittleEndian, share.(sShareData).nNonceLock)
		copy(msgSendBuffer[12:16], nNonceLock.Bytes()[0:4])
		nNonceBlock := bytes.NewBuffer([]byte{})
		binary.Write(nNonceBlock, binary.LittleEndian, share.(sShareData).nNonceBlock)
		copy(msgSendBuffer[16:20], nNonceBlock.Bytes()[0:4])
		copy(msgSendBuffer[20:52], share.(sShareData).byhashPervBlock[0:32])
		copy(msgSendBuffer[52:84], share.(sShareData).byhashMerkleRoot[0:32])
		copy(msgSendBuffer[84:116], share.(sShareData).byhashSeed[0:32])
		pConn.WriteToUDP(msgSendBuffer, paddr)
	}
}
