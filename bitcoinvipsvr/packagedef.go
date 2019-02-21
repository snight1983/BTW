package bitcoinvipsvr

import (
	"net"
	"sync"
)

//const gMsgUserRegedit int32 = 1000
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

type sConfigJSON struct {
	RPCIP           string `json:"RpcIp"`
	RPCPORT         int    `json:"RpcPort"`
	RPCUSER         string `json:"RpcUser"`
	RPCPW           string `json:"RpcPw"`
	POOLWALLETSYNC  string `json:"PoolWalletSync"`
	POOLSHAREBIT    uint32 `json:"PoolShareBit"`
	POOLSHAREHASH   int32  `json:"PoolShareHash"`
	POOLMINPAY      int32  `json:"PoolMinPay"`
	POOLPAYINTERVAL int64  `json:"PoolPayInterval"`
	POOLPAYCALCULAT int64  `json:"PoolPayCalculat"`
	POOLINCOME      int32  `json:"PoolIncome"`
	POOLUDPLSN      string `json:"PoolUdpLsn"`
	POOLRECVADDR    string `json:"PoolRecvAddress"`
	POOLMINEADDR    string `json:"PoolMineAddress"`
	DBNAME          string `json:"DBName"`
	DBCONNECT       string `json:"DBConnect"`
}

type sDbBlocks struct {
	n64Blockid  int64
	sTxID       string
	nu64Amount  uint64
	n32IsHandle int
}

type sAddIncome struct {
	nID           int64
	sRecvaddress  string
	nu64Unpay     uint64
	nu64Paied     uint64
	f32Speed      float32
	n32Sharecur   int32
	n32Sharetotal int32
	nu64Createtm  int64
	bIsInDB       bool
}

type sMinerShare struct {
	nID       int64
	sAddr     string
	nIsHandle int
}

var gConfig *sConfigJSON
var gHeader []byte
var gWorkHeader sBlockHeader
var gMkrQueue *syncQueue
var gShareQueue *syncQueue
var gShareMap *sSyncMap
var gAddrPayInfoMap *sSyncMap
var gMinerRetMap *sSyncMap
var gJobUDPChannel = make(chan sJobUDPData, 2048)
