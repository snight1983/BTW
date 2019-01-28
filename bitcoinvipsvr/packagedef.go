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

// Miner account
type Miner struct {
	n64Dbid       int64  // 数据库ID
	sUserID       string // 机器码
	sUserName     string // 用户名
	sRecvAddress  string // 用户收钱地址
	n64CreateTime int64  // 创建时间
	bIschange     bool   // 是否发生变化
}

// 用户注册实例
var gMinerRetMap *sSyncMap

// AddressPay 地址支付总帐
type sAddressPay struct {
	n64Dbid       int64  // 数据库ID
	sUserAddress  string // 收钱地址
	n64UnPaid     int64  // 未支付
	n64Paid       int64  // 已支付
	n64CreateTime int64  // 创建时间
	bIschange     bool
}

// 用户总帐表
type sAddressPayMap struct {
	plock     *sync.RWMutex
	mBm       map[string]*sAddressPay
	bIsUpdate bool
	bIsInsert bool
}

// 总帐信息Map
var gAddrPayInfoMap *sSyncMap

// JobUDPData data buf
type sJobUDPData struct {
	nLen   int          // 数据长度
	byData []byte       // 数据
	pAddr  *net.UDPAddr // 地址
	pConn  *net.UDPConn // 连接
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
