package bitcoinvipsvr

import (
	"net"
	"sync"
)

// 用户注册msgID
const msgUserRegedit int32 = 1000

// Miner account
type Miner struct {
	Dbid        int64  // 数据库ID
	UserID      string // 机器码
	UserName    string // 用户名
	RecvAddress string // 用户收钱地址
	CreateTime  int64  // 创建时间
	Ischange    bool   // 是否发生变化
	ShareCnt    int    // 工作量数量(>0 才允许注册)
}

// 用户注册表
type sMinerRegMap struct {
	lock     *sync.RWMutex
	bm       map[string]*Miner
	IsUpdate bool
	IsInsert bool
}

// 用户注册表实例
var gMinerRetMap *sSyncMap

// AddressPay 地址支付总帐
type sAddressPay struct {
	Dbid        int64  // 数据库ID
	UserAddress string // 收钱地址
	UnPaid      int64  // 未支付
	Paid        int64  // 已支付
	CreateTime  int64  // 创建时间
	Ischange    bool
}

// 用户注册表
type sAddressPayMap struct {
	lock     *sync.RWMutex
	bm       map[string]*sAddressPay
	IsUpdate bool
	IsInsert bool
}

// 支付总信息
var gAddrPayInfoMap *sSyncMap

// JobUDPData data buf
type sJobUDPData struct {
	len  int         // 数据长度
	data []byte      // 数据
	addr net.UDPAddr // 地址
}
