package bitcoinvipsvr

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/buger/jsonparser"
)

func jobUDPHandle(id int) {
	for {
		ljob := <-gJobUDPChannel
		res, lMsgid := readInt32(ljob.byData, 0, ljob.nLen)
		if true == res {
			switch lMsgid {
			case gMsgMiningWorkRQ:
				msgMiningWorkRQ(ljob.pConn, ljob.pAddr, ljob.byData, 4, ljob.nLen)
				break
			case gMsgShareCheckRQ:
				msgShareCheckRQ(ljob.pConn, ljob.pAddr, ljob.byData, 4, ljob.nLen)
				break
			case gMsgMrkID:
				msgMrkID(ljob.pConn, ljob.pAddr, ljob.byData, 4, ljob.nLen)
				break
			case gHeaderMrkID:
				msgHeaderID(ljob.pConn, ljob.pAddr, ljob.byData, 4, ljob.nLen)
				break
			case gMsgShareRQ:
				msgShareReportRQ(ljob.pConn, ljob.pAddr, ljob.byData, 4, ljob.nLen)
				break
			case gMsgShareCheckID:
				msgShareCheckID(ljob.pConn, ljob.pAddr, ljob.byData, 4, ljob.nLen)
				break
			default:
				break
			}
		}
	} //end for
}

func svrLsn() bool {
	addr, err := net.ResolveUDPAddr("udp", gConfig.POOLUDPLSN)
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

func load(file string) (*sConfigJSON, error) {
	gConfig = &sConfigJSON{}
	buf, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	jsonParser := json.NewDecoder(buf)
	err = jsonParser.Decode(gConfig)
	return gConfig, err
}

func insertBlockTask() {
	start := time.Now()
	rpcClient, err := newClient(gConfig.RPCIP, gConfig.RPCPORT, gConfig.RPCUSER, gConfig.RPCPW, false)
	if nil == err {
		reqJSON := "{\"method\":\"listunspent\",\"params\":[6,288,[\"vLhwX9u9UUhBtxYdZWSPKxukv6fZC9f7xW\"]]}"
		retJSON, err2 := rpcClient.send(reqJSON)
		if err2 != nil || nil == retJSON {
			log.Fatalln(err2)
		}
		jsonparser.ArrayEach(retJSON, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			sTxid, err := jsonparser.GetString(value, "txid")
			if nil == err {
				sAmoult, err := jsonparser.GetFloat(value, "amount")
				sAmoult *= 1000000
				string1 := strconv.FormatFloat(sAmoult, 'f', -1, 64)
				n64Amount, err := strconv.ParseInt(string1, 10, 64)
				if nil == err && n64Amount > 0 {
					bSpendable, err := jsonparser.GetBoolean(value, "spendable")
					if nil == err && bSpendable {
						bSolvable, err := jsonparser.GetBoolean(value, "solvable")
						if nil == err && bSolvable {
							bSafe, err := jsonparser.GetBoolean(value, "safe")
							if nil == err && bSafe {
								res := insertBlock(sTxid, n64Amount)
								if res {
									fmt.Printf("%s | txid:%s amount:%d\n", time.Now().Format("2006-01-02 15:04:05"), sTxid, n64Amount)
								}
							}
						}
					}
				}
			}
		}, "result")
	}
	end := time.Now()
	fmt.Printf("%s | insertBlockTask total time:%f\n", time.Now().Format("2006-01-02 15:04:05"), end.Sub(start).Seconds())
}

func insertShareTask() {
	start := time.Now()
	gShareMap.lock.RLock()
	for key, value := range gShareMap.bm {
		if nil != value {
			if value.(*sShareData).nConfCnt >= 3 {
				insertShare(value.(*sShareData))
				delete(gShareMap.bm, key)
			} else if value.(*sShareData).nHeight < gWorkHeader.nHeight-3 {
				delete(gShareMap.bm, key)
			} else {
				gShareQueue.Push(value)
				fmt.Printf("%s | Push Share Check:%s Height:%d Conf:%d\n", time.Now().Format("2006-01-02 15:04:05"), key, value.(*sShareData).nHeight, value.(*sShareData).nConfCnt)
			}
		}
	}
	gShareMap.lock.RUnlock()
	end := time.Now()
	fmt.Printf("%s | insertShareTask total time:%f\n", time.Now().Format("2006-01-02 15:04:05"), end.Sub(start).Seconds())

}

func incomeCalculate(tm int64) {
	start := time.Now()
	oBlocks, err := getBlockIncome()
	if nil == err {
		var n64Total uint64
		if nil == err && nil != oBlocks {
			for _, block := range oBlocks {
				if block.n32IsHandle == 0 {
					block.n32IsHandle = 1
					err = updateBlockIncome(block)
					if nil == err {
						n64Total += block.nu64Amount
					}
				}
			}
		}
		n64RealTotal := n64Total
		n64PoolRecv := n64Total
		if n64Total > 0 {
			n64PoolRecv /= 1000
			n64PoolRecv *= uint64(gConfig.POOLINCOME)
			n64Total -= n64PoolRecv
			shareTotal := 0
			oIncomeMap, err := getAddressIncome()
			if nil == err {
				oSHares, err := getMinerShare()
				for _, value := range oIncomeMap {
					value.n32Sharecur = 0
				}
				var poolAddr *sAddIncome
				poolAddr = nil
				if _, ok := oIncomeMap[gConfig.POOLRECVADDR]; !ok {
					poolAddr = &sAddIncome{}
					poolAddr.sRecvaddress = gConfig.POOLRECVADDR
					poolAddr.nu64Unpay = 0
					poolAddr.nu64Paied = 0
					poolAddr.f32Speed = 0
					poolAddr.n32Sharecur = 0
					poolAddr.n32Sharetotal = 0
					poolAddr.nu64Createtm = start.Unix()
					poolAddr.bIsInDB = false
					oIncomeMap[gConfig.POOLRECVADDR] = poolAddr
				} else {
					poolAddr = oIncomeMap[gConfig.POOLRECVADDR]
				}

				if nil == err && nil != oSHares {
					for _, share := range oSHares {
						share.nIsHandle = 1
						if _, ok := oIncomeMap[share.sAddr]; !ok {
							obj := &sAddIncome{}
							obj.sRecvaddress = share.sAddr
							obj.nu64Unpay = 0
							obj.nu64Paied = 0
							obj.f32Speed = 0
							obj.n32Sharecur = 1
							obj.n32Sharetotal = 1
							obj.nu64Createtm = start.Unix()
							obj.bIsInDB = false
							oIncomeMap[share.sAddr] = obj
							shareTotal++
						} else {
							oIncomeMap[share.sAddr].n32Sharecur++
							oIncomeMap[share.sAddr].n32Sharetotal++
							shareTotal++
						}
					}
					updateMinerShare(oSHares)
				}
				fPart := float64(n64Total) / float64(shareTotal)
				var fSpeedTotal float32
				var u64AllPay uint64
				for _, value := range oIncomeMap {
					value.nu64Unpay = uint64(fPart * float64(value.n32Sharecur))
					fPeerHash := float32(value.n32Sharecur * gConfig.POOLSHAREHASH)
					fPeerTime := float32(start.Unix() - tm)
					value.f32Speed = fPeerHash / fPeerTime
					fSpeedTotal += value.f32Speed
					u64AllPay += value.nu64Unpay
				}
				poolAddr.nu64Unpay += (n64RealTotal - u64AllPay)
				insertPoolStatu(len(oIncomeMap), fSpeedTotal, time.Now().Unix(), n64RealTotal, poolAddr.nu64Unpay)
				updateAddrIncome(oIncomeMap)
			}
		}
	}

	end := time.Now()
	fmt.Printf("%s | IncomeCalculate() total time:%f\n", time.Now().Format("2006-01-02 15:04:05"), end.Sub(start).Seconds())
}

// StartSvr start udp lsn
func StartSvr() bool {
	path, err := GetCurrentPath()
	path += "/config.json"
	_, err = load(path)
	if nil != err {
		return false
	}

	fmt.Printf("%s | Config Path:%s\n", time.Now().Format("2006-01-02 15:04:05"), path)
	fmt.Printf("%s | Config RPCIP:%s\n", time.Now().Format("2006-01-02 15:04:05"), gConfig.RPCIP)
	fmt.Printf("%s | Config RPCPORT:%d\n", time.Now().Format("2006-01-02 15:04:05"), gConfig.RPCPORT)
	fmt.Printf("%s | Config RPCUSER:%s\n", time.Now().Format("2006-01-02 15:04:05"), gConfig.RPCUSER)
	fmt.Printf("%s | Config RPCPW:%s\n", time.Now().Format("2006-01-02 15:04:05"), gConfig.RPCPW)
	fmt.Printf("%s | Config POOLWALLETSYNC:%s\n", time.Now().Format("2006-01-02 15:04:05"), gConfig.POOLWALLETSYNC)
	fmt.Printf("%s | Config POOLSHAREBIT:%d\n", time.Now().Format("2006-01-02 15:04:05"), gConfig.POOLSHAREBIT)
	fmt.Printf("%s | Config POOLSHAREHASH:%d\n", time.Now().Format("2006-01-02 15:04:05"), gConfig.POOLSHAREHASH)
	fmt.Printf("%s | Config POOLMINPAY:%d\n", time.Now().Format("2006-01-02 15:04:05"), gConfig.POOLMINPAY)
	fmt.Printf("%s | Config POOLPAYINTERVAL:%d\n", time.Now().Format("2006-01-02 15:04:05"), gConfig.POOLPAYINTERVAL)
	fmt.Printf("%s | Config POOLPAYCALCULAT:%d\n", time.Now().Format("2006-01-02 15:04:05"), gConfig.POOLPAYCALCULAT)
	fmt.Printf("%s | Config POOLINCOME:%d\n", time.Now().Format("2006-01-02 15:04:05"), gConfig.POOLINCOME)
	fmt.Printf("%s | Config POOLUDPLSN:%s\n", time.Now().Format("2006-01-02 15:04:05"), gConfig.POOLUDPLSN)
	fmt.Printf("%s | Config POOLRECVADDR:%s\n", time.Now().Format("2006-01-02 15:04:05"), gConfig.POOLRECVADDR)
	fmt.Printf("%s | Config POOLMINEADDR:%s\n", time.Now().Format("2006-01-02 15:04:05"), gConfig.POOLMINEADDR)
	fmt.Printf("%s | Config DBNAME:%s\n", time.Now().Format("2006-01-02 15:04:05"), gConfig.DBNAME)
	fmt.Printf("%s | Config DBCONNECT:%s\n", time.Now().Format("2006-01-02 15:04:05"), gConfig.DBCONNECT)

	gMkrQueue = newQueue()
	gShareQueue = newQueue()
	gShareMap = newSyncMap()
	gMinerRetMap = newSyncMap()
	gAddrPayInfoMap = newSyncMap()
	err = onInitDbConnectPool()

	if err != nil {
		fmt.Print(err)
		return false
	}
	go svrLsn()
	tmCheckShare := time.Now().Unix()
	tmCheckBlock := tmCheckShare
	tmIncomeCalculate := tmCheckShare
	tmPayForMiner := tmCheckShare

	for {
		tmCur := time.Now().Unix()

		if (tmCur - tmCheckShare) > 120 {
			tmCheckShare = tmCur
			insertShareTask()
		}

		if (tmCur - tmCheckBlock) > 1800 {
			tmCheckBlock = tmCur
			insertBlockTask()
		}

		if (tmCur - tmIncomeCalculate) > gConfig.POOLPAYCALCULAT {
			incomeCalculate(tmIncomeCalculate)
			tmIncomeCalculate = tmCur
		}

		if (tmCur - tmPayForMiner) > gConfig.POOLPAYINTERVAL {
			tmPayForMiner = tmCur

		}
		time.Sleep(10 * time.Second)
	}
}
