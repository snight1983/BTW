package bitcoinvipsvr

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// 钱包连接参数
type rpcClient struct {
	serverAddr string
	user       string
	passwd     string
	httpClient *http.Client
}

/*
// 请求信息
type rpcRequest struct {
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      int64       `json:"id"`
	JSONRpc string      `json:"jsonrpc"`
}

type rpcResponse struct {
	ID     int64           `json:"id"`
	Result json.RawMessage `json:"result"`
	Err    interface{}     `json:"error"`
}
*/
//连接配置
func newClient(host string, port int, user, passwd string, useSSL bool) (c *rpcClient, err error) {
	if len(host) == 0 {
		err = errors.New("Bad call missing argument host")
		return
	}
	var serverAddr string
	var httpClient *http.Client
	if useSSL {
		serverAddr = "https://"
		t := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		httpClient = &http.Client{Transport: t}
	} else {
		serverAddr = "http://"
		httpClient = &http.Client{}
	}
	c = &rpcClient{serverAddr: fmt.Sprintf("%s%s:%d", serverAddr, host, port), user: user, passwd: passwd, httpClient: httpClient}
	return
}

// 超时处理
func (c *rpcClient) doTimeoutRequest(timer *time.Timer, req *http.Request) (*http.Response, error) {
	type result struct {
		resp *http.Response
		err  error
	}
	done := make(chan result, 1)
	go func() {
		resp, err := c.httpClient.Do(req)
		done <- result{resp, err}
	}()
	// Wait for the read or the timeout
	select {
	case r := <-done:
		return r.resp, r.err
	case <-timer.C:
		return nil, errors.New("Timeout reading data from server")
	}
}

//通信
func (c *rpcClient) send(reqJSON string) ([]byte, error) {
	connectTimer := time.NewTimer(30 * time.Second)
	reqJSONByte := []byte(reqJSON)
	payloadBuffer := bytes.NewReader(reqJSONByte)
	req, err := http.NewRequest("POST", c.serverAddr, payloadBuffer)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	req.Header.Add("Accept", "application/json")
	if len(c.user) > 0 || len(c.passwd) > 0 {
		req.SetBasicAuth(c.user, c.passwd)
	}
	resp, err := c.doTimeoutRequest(connectTimer, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		err = errors.New("HTTP error: " + resp.Status)
		return nil, err
	}
	return data, nil
}
