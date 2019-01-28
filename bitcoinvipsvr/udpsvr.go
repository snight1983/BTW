package bitcoinvipsvr

import (
	"fmt"
	"net"
	"time"
)

var gMsgcnt uint64
var gJobUDPChannel = make(chan sJobUDPData, 2048)

func jobUDPHandle(id int) {
	for {
		ljob := <-gJobUDPChannel
		gMsgcnt++
		//if gMsgcnt%1024 == 0 {
		//	now := time.Now()
		//	fmt.Printf("%d-%d-%d %d:%d:%d|", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
		//	fmt.Printf("Pool:PackageCnt:%d \n", gMsgcnt)
		//}
		res, lMsgid := readInt32(ljob.byData, 0, ljob.nLen)
		if true == res {
			switch lMsgid {
			case gMsgUserRegedit:
				netMinerRegedit(ljob.pAddr, ljob.byData, 4, ljob.nLen)
				break
			case gMsgMiningWorkRQ:
				msgMiningWorkRQ(ljob.pConn, ljob.pAddr, ljob.byData, 4, ljob.nLen)
				break
			case gMsgCheckWorkRQ:
				msgCheckWorkRQ(ljob.pConn, ljob.pAddr, ljob.byData, 4, ljob.nLen)
				break
			case gMsgMrkID:
				msgMrkID(ljob.pConn, ljob.pAddr, ljob.byData, 4, ljob.nLen)
				break
			case gHeaderMrkID:
				msgHeaderID(ljob.pConn, ljob.pAddr, ljob.byData, 4, ljob.nLen)
				break
			default:
				break
			}
		}
	} //end for
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
		job.byData = make([]byte, 1450)
		n, ctlAddr, err := listener.ReadFromUDP(job.byData)
		if err != nil {
			fmt.Print(err)
			continue
		}
		if n > 0 {
			job.nLen = n
			job.pAddr = ctlAddr
			job.pConn = listener
			gJobUDPChannel <- job
		}
	}
}

// StartSvr start udp lsn
func StartSvr() bool {
	/*
		{
			rpcClient, err := newClient("127.0.0.1", 8337, "Bitcoinvip", "Bitcoinvippw", false)
			if nil == err {
				reqJSON := "{\"method\":\"listunspent\",\"params\":[]}"
				retJSON, err2 := rpcClient.send(reqJSON)
				if err2 != nil || nil == retJSON {
					log.Fatalln(err2)
				}
			}
		}
	*/
	gMinerRetMap = newSyncMap()
	gShareMap = newSyncMap()
	gAddrPayInfoMap = newSyncMap()
	gMkrQueue = newQueue()
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
	tmCheckShare := tmCheckRetUpdate
	tmIncomeCheck := tmCheckRetUpdate

	for {
		tmCur := time.Now().Unix()
		if (tmCur-tmCheckRetInsert) > 120 && gMinerRetMap.isInsert {
			fmt.Println("insertMinerRegedit beg")
			tmCheckRetInsert = tmCur
			gMinerRetMap.isInsert = false
			start := time.Now()
			insertMinerRegedit()
			end := time.Now()
			fmt.Println("insertMinerRegedit total time:", end.Sub(start).Seconds())
		}

		if (tmCur-tmCheckRetUpdate) > 120 && gMinerRetMap.isUpdate {
			tmCheckRetUpdate = tmCur
			fmt.Println("updateMinerRegedit beg")
			gMinerRetMap.isUpdate = false
			start := time.Now()
			updateMinerRegedit()
			end := time.Now()
			fmt.Println("updateMinerRegedit total time:", end.Sub(start).Seconds())
		}

		if (tmCur - tmCheckShare) > 120 {
			tmCheckShare = tmCur
			fmt.Println("insertShare beg")
			start := time.Now()
			gShareMap.lock.RLock()
			defer gShareMap.lock.RUnlock()
			for key, value := range gShareMap.bm {
				if value.(*sShareData).nConfCnt >= 3 {
					insertShare(value.(*sShareData))
					delete(gShareMap.bm, key)
				} else if value.(*sShareData).nHeight < gWorkHeader.nHeight-6 {
					delete(gShareMap.bm, key)
				} else {
					gShareQueue.Push(value)
				}
			}
			end := time.Now()
			fmt.Println("insertShare total time:", end.Sub(start).Seconds())
		}

		if (tmCur - tmIncomeCheck) > 120 {
			tmIncomeCheck = tmCur

			fmt.Println("Income Check beg")
			start := time.Now()
			/*
				{
					rpcClient, err := newClient("127.0.0.1", 8337, "Bitcoinvip", "Bitcoinvippw", false)
					if nil == err {
						reqJSON := "{\"method\":\"listunspent\",\"params\":[]}"
						retJSON, err2 := rpcClient.send(reqJSON)
						if err2 != nil {
							log.Fatalln(err2)
						}
						log.Println("returnJson:", retJSON)
					}

				}
			*/
			end := time.Now()
			fmt.Println("insertShare total time:", end.Sub(start).Seconds())
		}

		time.Sleep(10 * time.Second)
	}
}
