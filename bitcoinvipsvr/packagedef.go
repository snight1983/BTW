package bitcoinvipsvr

import (
	"net"
	"sync"
)

const gMsgUserRegedit int32 = 1000
const gMsgMiningWorkRQ int32 = 1001
const gMsgMiningWorkRS int32 = 1002
const gMsgCheckWorkRQ int32 = 1003
const gMsgCheckWorkRS int32 = 1004
const gMsgMrkID int32 = 1005
const gHeaderMrkID int32 = 1006

type sMiner struct {
	n64Dbid       int64
	sUserID       string
	sUserName     string
	sRecvAddress  string
	n64CreateTime int64
	bIschange     bool
}

var gMinerRetMap *sSyncMap

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

var gAddrPayInfoMap *sSyncMap

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
	nBitsShare      uint32
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
	sRecvAddress     string
	nConfCnt         int
}

var gHeader []byte
var gWorkHeader sBlockHeader
var gMkrQueue *syncQueue
var gShareQueue *syncQueue
var gShareMap *sSyncMap
