package main

import (
	"DistributedSystemProjects/expr2/lockrpc"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sync"

	"git.apache.org/thrift.git/lib/go/thrift"
)

const (
	HOST = "localhost"
	PORT = "6790"
)

var MapRWMutex *sync.RWMutex

var ClientState map[int64]bool = make(map[int64]bool, 0)

func MapAddNewClient(id int64) {
	_, exits := ClientState[id]
	if exits {
		return
	} else {
		MapRWMutex.Lock()
		ClientState[id] = false
		MapRWMutex.Unlock()
		log.Print("add new client:", id)
	}
}

func IsLockAlready() bool {
	flag := false
	MapRWMutex.RLock()
	for _, v := range ClientState {
		if v == true {
			flag = true
			break
		}
	}
	MapRWMutex.RUnlock()
	return flag
}

func RandomLock() int64 {
	var changeID int64 = -1
	if IsLockAlready() {
		return changeID
	}
	i := rand.Intn(len(ClientState))
	MapRWMutex.RLock()
	for k := range ClientState {
		if i == 0 {
			changeID = k
			break
		}
		i--
	}
	MapRWMutex.RUnlock()
	if changeID != -1 {
		MapRWMutex.Lock()
		ClientState[changeID] = true
		MapRWMutex.Unlock()
	}
	return changeID
}

func GetCurrentLock() int64 {
	var res int64 = -1
	MapRWMutex.RLock()
	for key, val := range ClientState {
		if val == true {
			res = key
		}
	}
	MapRWMutex.RUnlock()
	return res
}

func UnlockClient(id int64) {
	_, exits := ClientState[id]
	if exits {
		ClientState[id] = false
	}
}

func LockClear() {
	var changeID int64 = -1
	MapRWMutex.RLock()
	for key, val := range ClientState {
		if val == true {
			changeID = key
			break
		}
	}
	MapRWMutex.RUnlock()
	if changeID != -1 {
		MapRWMutex.Lock()
		ClientState[changeID] = false
		MapRWMutex.Unlock()
	}
}

type FormatDataImpl struct{}

func (fdi *FormatDataImpl) DoLock(req *lockrpc.Req) (*lockrpc.Rsp, error) {
	MapAddNewClient(req.CliID)
	RandomLock()
	lockID := GetCurrentLock()
	var rsp lockrpc.Rsp
	rsp.CliID = lockID
	rsp.Operator = "acquire"
	return &rsp, nil
}

func (fdi *FormatDataImpl) UnLock(req *lockrpc.Req) (*lockrpc.Rsp, error) {
	UnlockClient(req.CliID)
	var rsp lockrpc.Rsp
	rsp.CliID = req.CliID
	rsp.Operator = "release"
	log.Print("clockID=", rsp.CliID, " operator=", rsp.Operator)
	return &rsp, nil
}

func (fdi *FormatDataImpl) ClientStates() (*lockrpc.Rsp, error) {
	var rsp lockrpc.Rsp
	rsp.CliID = -1
	rsp.Operator = "-1"
	jsonString, err := json.Marshal(ClientState)
	if err == nil {
		rsp.Buffer = string(jsonString)
	}
	return &rsp, nil
}

func StartSever() {
	handler := &FormatDataImpl{}
	processor := lockrpc.NewGetLockProcessor(handler)

	serverTransport, err := thrift.NewTServerSocket(HOST + ":" + PORT)
	if err != nil {
		log.Fatalln("Error:", err)
	}
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)
	fmt.Println("Running at:", HOST+":"+PORT)
	server.Serve()
}

func main() {
	MapRWMutex = new(sync.RWMutex)
	StartSever()
}
