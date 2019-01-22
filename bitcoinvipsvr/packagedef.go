package bitcoinvipsvr

import (
	"bytes"
	"encoding/binary"
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
	UserAddress string // 用户收钱地址
	UserDisName string // 用户名
	CreateTime  int64  // 创建时间
	Ischange    bool   // 是否发生变化
	LastTime    int64  // 时间戳
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
//var gminerMap sMinerRegMap
var gMinerRetMap *sSyncMap

// AddressPay 地址支付总帐
type sAddressPay struct {
	Dbid        int64  // 数据库ID
	UserAddress string // 收钱地址
	UnPaid      int64  // 创建时间
	Paid        int64  // 创建时间
	CreateTime  int64  // 创建时间
}

// 用户注册表
type sAddressPayMap struct {
	lock     *sync.RWMutex
	bm       map[string]*sAddressPay
	IsUpdate bool
	IsInsert bool
}

// JobUDPData data buf
type sJobUDPData struct {
	len  int         // 数据长度
	data [1450]byte  // 数据
	addr net.UDPAddr // 地址
}

func readInt32(buf []byte, begpos int, total int) (bool, int32) {
	endpos := begpos + 4
	if endpos <= total {
		var lvalue32 int32
		readbuf := bytes.NewReader(buf[begpos : begpos+4])
		binary.Read(readbuf, binary.LittleEndian, &lvalue32)
		return true, lvalue32
	}
	return false, 0
}

func readInt64(buf []byte, begpos int, total int) (bool, int64) {
	endpos := begpos + 8
	if endpos <= total {
		var lvalue64 int64
		readbuf := bytes.NewReader(buf[begpos : begpos+8])
		binary.Read(readbuf, binary.BigEndian, &lvalue64)
		return true, lvalue64
	}
	return false, 0
}

func readuint64(buf []byte, begpos int, total int) (bool, uint64) {
	endpos := begpos + 8
	if endpos <= total {
		var lvalue64 uint64
		readbuf := bytes.NewReader(buf[begpos : begpos+8])
		binary.Read(readbuf, binary.LittleEndian, &lvalue64)
		return true, lvalue64
	}
	return false, 0
}
