package bitcoinvipsvr

import (
	"fmt"
	"net"
	"time"
)

var gMsgcnt uint64
var gJobUDPChannel = make(chan sJobUDPData, 1024)

func jobUDPHandle(id int) {
	for {
		job := <-gJobUDPChannel
		gMsgcnt++
		res, msgid := readInt32(job.data, 0, job.len)
		if true == res {
			switch msgid {
			case msgUserRegedit:
				netMinerRegedit(job.addr, job.data, 4, job.len)
				break
			default:
				break
			}
		}
	}
}

func svrLsn() bool {

	gMsgcnt = 0

	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:10008")
	if err != nil {
		fmt.Print(err)
		return false
	}
	listener, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Print(err)
		return false
	}
	defer listener.Close()

	for i := 0; i < 4; i++ {
		go jobUDPHandle(i)
	}

	for {
		var job sJobUDPData
		job.data = make([]byte, 1450)
		n, ctlAddr, err := listener.ReadFromUDP(job.data)
		if err != nil {
			fmt.Print(err)
			continue
		}
		if n > 0 {
			job.len = n
			job.addr = *ctlAddr
			gJobUDPChannel <- job
		}
	}
}

// StartSvr start udp lsn
func StartSvr() bool {
	gMinerRetMap = newSyncMap()
	gAddrPayInfoMap = newSyncMap()
	err := onInitDbConnectPool()
	if err != nil {
		fmt.Print(err)
		return false
	}
	err = getMinerRegedit()
	if err != nil {
		fmt.Print(err)
		return false
	}
	go svrLsn()

	tmCheckRetUpdate := time.Now().Unix()
	tmCheckRetInsert := tmCheckRetUpdate

	for {
		tmCur := time.Now().Unix()
		fmt.Println("StartSvr wait in")
		// 检查是否需要插入用户注册信息
		if (tmCur-tmCheckRetInsert) > 120 && gMinerRetMap.isInsert {
			fmt.Println("insertMinerRegedit beg")
			tmCheckRetInsert = tmCur
			gMinerRetMap.isInsert = false
			start := time.Now()
			insertMinerRegedit()
			end := time.Now()
			fmt.Println("insertMinerRegedit total time:", end.Sub(start).Seconds())
		}
		// 检查是否需要更新用户注册信息
		if (tmCur-tmCheckRetUpdate) > 120 && gMinerRetMap.isUpdate {
			fmt.Println("updateMinerRegedit beg")
			tmCheckRetUpdate = tmCur
			gMinerRetMap.isUpdate = false
			start := time.Now()
			updateMinerRegedit()
			end := time.Now()
			fmt.Println("updateMinerRegedit total time:", end.Sub(start).Seconds())
		}
		time.Sleep(10 * time.Second)
		fmt.Println("StartSvr wait out")
	}
}
