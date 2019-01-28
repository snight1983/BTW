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

		now := time.Now()
		fmt.Printf("%d-%d-%d %d:%d:%d|", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
		fmt.Printf("Pool: New Header v:%d h:%d b:%d\n", gWorkHeader.nVersion, gWorkHeader.nHeight, gWorkHeader.nBitsBlock)
		return
	}
	now := time.Now()
	fmt.Printf("%d-%d-%d %d:%d:%d|", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	fmt.Printf("Pool: Old Header v:%d h:%d b:%d\n", gWorkHeader.nVersion, gWorkHeader.nHeight, gWorkHeader.nBitsBlock)
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
