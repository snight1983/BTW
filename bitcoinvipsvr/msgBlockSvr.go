package bitcoinvipsvr

import (
	"fmt"
	"net"
	"time"
)

func msgHeaderID(pConn *net.UDPConn, paddr *net.UDPAddr, data []byte, posbeg int, datalen int) {

	if datalen != 48 {
		return
	}
	res, nVersion := readInt32(data, posbeg, datalen)
	if res == false {
		return
	}
	posbeg += 4
	res, nHeight := readInt32(data, posbeg, datalen)
	if res == false {
		return
	}
	posbeg += 4
	if nHeight != gWorkHeader.nHeight {
		gHeader = make([]byte, 44)
		copy(gHeader[0:44], data[4:48])
		gMkrQueue.Clear()
		gWorkHeader.nVersion = nVersion
		gWorkHeader.nHeight = nHeight
		copy(gWorkHeader.byhashPervBlock[0:32], data[posbeg:posbeg+32])
		posbeg += 32
		res, gWorkHeader.nBitsBlock = readUInt32(data, posbeg, datalen)
		posbeg += 4

		fmt.Printf("%s | Pool: New Header v:%d h:%d b:%d\n", time.Now().Format("2006-01-02 15:04:05"), gWorkHeader.nVersion, gWorkHeader.nHeight, gWorkHeader.nBitsBlock)
		//gWorkHeader.nBitsShare = 520159231
		/*
			"0000:": 520159231,
			"00000:": 504365055,
			"000000:": 503382015,
			"0000000:": 487587839,
			"00000000:": 486604799,
			"000000000:": 470810623
		*/
		return
	}
	//fmt.Printf("%s | Pool: Old Header v:%d h:%d b:%d\n", time.Now().Format("2006-01-02 15:04:05"), gWorkHeader.nVersion, gWorkHeader.nHeight, gWorkHeader.nBitsBlock)
}

func msgMrkID(pConn *net.UDPConn, paddr *net.UDPAddr, data []byte, posbeg int, datalen int) {
	if datalen != 1234 || gMkrQueue.size() >= 4096 {
		return
	}
	res, nHeight := readInt32(data, posbeg, datalen)
	if (res == false) || (nHeight != gWorkHeader.nHeight) {
		return
	}
	mrk := (make([]byte, 1296))
	copy(mrk[0:1296], data[8:1234])
	gMkrQueue.Push(mrk)
}
