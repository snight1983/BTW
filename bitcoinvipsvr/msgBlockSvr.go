package bitcoinvipsvr

import "net"

func msgHeaderID(pConn *net.UDPConn, paddr *net.UDPAddr, data []byte, posbeg int, datalen int) {
	if datalen != 52 {
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
		gHeader = make([]byte, 48)
		copy(gHeader[0:48], data[4:52])
	}
	gWorkHeader.nVersion = nVersion
	gWorkHeader.nHeight = nHeight

	copy(gWorkHeader.byhashPervBlock[0:32], data[posbeg:posbeg+32])
	posbeg += 32
	res, gWorkHeader.nBitsBlock = readUInt32(data, posbeg, datalen)
	// back?
}

func msgMrkID(pConn *net.UDPConn, paddr *net.UDPAddr, data []byte, posbeg int, datalen int) {
	if datalen != 1234 {
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
