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

		shareBit := bytes.NewBuffer([]byte{})
		binary.Write(shareBit, binary.LittleEndian, gConfig.POOLSHAREBIT)
		msgSendBuffer = append(msgSendBuffer, shareBit.Bytes()...)

		msgSendBuffer = append(msgSendBuffer, gHeader...)
		if bNeedMkr {
			mkrBuf := gMkrQueue.Pop()
			if nil != mkrBuf {
				msgSendBuffer = append(msgSendBuffer, mkrBuf.([]byte)...)
			}
		}
		pConn.WriteToUDP(msgSendBuffer, paddr)
	}
}
