package bitcoinvipsvr

import (
	"fmt"
	"net"
	"time"
)

var gmsgcnt uint64
var gJobUDPChannel = make(chan sJobUDPData, 1024)

func jobUDPHandle(id int) {
	for {
		job := <-gJobUDPChannel
		res, msgid := readInt32(job.data[0:4], 0, job.len)
		if true == res {
			switch msgid {
			case msgUserRegedit:
				netMinerRegedit(job.addr, job.data[0:], 4, job.len)
				break
			default:
				break
			}
		}
	}
}

func svrLsn() bool {

	gmsgcnt = 0

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

	for i := 0; i < 1; i++ {
		go jobUDPHandle(i)
	}

	for {
		var job sJobUDPData
		n, ctlAddr, err := listener.ReadFromUDP(job.data[0:1449])
		if err != nil {
			fmt.Print(err)
			continue
		}
		if n > 0 {
			job.len = n
			job.addr = *ctlAddr
			gJobUDPChannel <- job
			//gmsgcnt++
		}
	}
}

// StartSvr start udp lsn
func StartSvr() bool {
	gMinerRetMap = newSyncMap()

	//gminerMap.lock = new(sync.RWMutex)
	//gminerMap.bm = make(map[string]*Miner)

	_, err := onInitDbConnectPool()
	if err != nil {
		fmt.Print(err)
		return false
	}
	_, err = getMinerRegedit()
	if err != nil {
		fmt.Print(err)
		return false
	}

	// 开始网络监听
	go svrLsn()

	tmCheckRetUpdate := time.Now().Unix()
	tmCheckRetInsert := tmCheckRetUpdate

	for {
		tmCur := time.Now().Unix()
		// 检查是否需要插入用户登录信息
		if (tmCur-tmCheckRetInsert) > 120 && gMinerRetMap.isInsert {
			fmt.Println("insertMinerRegedit beg")
			tmCheckRetInsert = tmCur
			gMinerRetMap.isInsert = false
			start := time.Now()
			insertMinerRegedit()
			end := time.Now()
			fmt.Println("insertMinerRegedit total time:", end.Sub(start).Seconds())
		}
		// 检查是否需要插入用户登录信息
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
	}
	//return true
}
