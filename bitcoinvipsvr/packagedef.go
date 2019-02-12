package bitcoinvipsvr

import (
	"net"
	"sync"
)

const gMsgUserRegedit int32 = 1000
const gMsgMiningWorkRQ int32 = 1001
const gMsgMiningWorkRS int32 = 1002
const gMsgShareCheckRQ int32 = 1003
const gMsgShareCheckRS int32 = 1004
const gMsgMrkID int32 = 1005
const gHeaderMrkID int32 = 1006
const gMsgShareRQ int32 = 1007
const gMsgShareRS int32 = 1008
const gMsgShareCheckID int32 = 1009

type sMiner struct {
	n64Dbid       int64
	sUserID       string
	sUserName     string
	sRecvAddress  string
	n64CreateTime int64
	bIschange     bool
}

type sAddressPay struct {
	n64Dbid       int64
	sUserAddress  string
	n64UnPaid     int64
	n64Paid       int64
	n64CreateTime int64
	bIschange     bool
}

type sAddressPayMap struct {
	plock     *sync.RWMutex
	mBm       map[string]*sAddressPay
	bIsUpdate bool
	bIsInsert bool
}

type sJobUDPData struct {
	nLen   int
	byData []byte
	pAddr  *net.UDPAddr
	pConn  *net.UDPConn
}

type sBlockHeader struct {
	nVersion        int32
	nHeight         int32
	byhashPervBlock [32]byte
	nBitsBlock      uint32
}

type sShareData struct {
	sHexShare        string
	nVersion         int32
	nHeight          int32
	byhashPervBlock  []byte
	byhashMerkleRoot []byte
	byhashSeed       []byte
	nNonceLock       uint32
	nNonceBlock      uint32
	nNonceMrk        uint32
	nTime            uint32
	nBit             uint32
	sRecvAddress     string
	nConfCnt         int
	nType            uint32
}

var gnShareBit uint32
var gHeader []byte
var gWorkHeader sBlockHeader
var gMkrQueue *syncQueue
var gShareQueue *syncQueue
var gShareMap *sSyncMap
var gAddrPayInfoMap *sSyncMap
var gMinerRetMap *sSyncMap
var gPoolBlockIP string
var gJobUDPChannel = make(chan sJobUDPData, 2048)
