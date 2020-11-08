package main

import (
	"encoding/json"
	"expr2/thrift/lock"
	"fmt"
	"log"
	"math/rand"
	"sync"

	"git.apache.org/thrift.git/lib/go/thrift"
)

const (
	HOST = "127.0.0.1"
	PORT = "9001"
)

var MapRWMutex *sync.RWMutex
var CliState map[int64]bool = make(map[int64]bool, 0)

type FormatDataImpl struct{}

func MapAddNewClient(id int64) {
	_, exist := CliState[id]
	if exist {
		return
	} else {
		MapRWMutex.Lock()
		CliState[id] = false
		MapRWMutex.Unlock()
		log.Print("Add new client:", id)
	}
}

func IsLock() bool {
	flag := false
	MapRWMutex.RLock()
	for _, v := range CliState {
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
	if IsLock() {
		return changeID
	}
	i := rand.Intn(len(CliState))
	MapRWMutex.RLock()
	for k := range CliState {
		if i == 0 {
			changeID = k
			break
		}
		i--
	}
	MapRWMutex.RUnlock()
	if changeID != -1 {
		MapRWMutex.Lock()
		CliState[changeID] = true
		MapRWMutex.Unlock()
	}
	return changeID
}

func GetCurrentLock() int64 {
	var res int64 = -1
	MapRWMutex.RLock()
	for key, val := range CliState {
		if val == true {
			res = key
		}
	}
	MapRWMutex.RUnlock()
	return res
}

func UnlockClient(id int64) {
	_, exists := CliState[id]
	if exists {
		CliState[id] = false
	}
}

func LockClear() {
	var changeID int64 = -1
	MapRWMutex.RLock()
	for key, val := range CliState {
		if val == true {
			changeID = key
			break
		}
	}
	MapRWMutex.RUnlock()
	if changeID != -1 {
		MapRWMutex.Lock()
		CliState[changeID] = false
		MapRWMutex.Unlock()
	}
}

func (fdi *FormatDataImpl) DoLock(req *lock.Req) (*lock.Response, error) {
	MapAddNewClient(req.CliID)
	RandomLock()
	lockID := GetCurrentLock()
	var rsp lock.Response
	rsp.CliID = lockID
	rsp.Operator = "acquire"
	return &rsp, nil
}

func (fdi *FormatDataImpl) UnLock(req *lock.Req) (*lock.Response, error) {
	UnlockClient(req.CliID)
	var rsp lock.Response
	rsp.CliID = req.CliID
	rsp.Operator = "release"
	log.Print("ClientID:", rsp.CliID, "operator=", rsp.Operator)
	return &rsp, nil
}

func (fdi *FormatDataImpl) ClientStates() (*lock.Response, error) {
	var rsp lock.Response
	rsp.CliID = -1
	rsp.Operator = "-1"
	jsonString, err := json.Marshal(CliState)
	if err != nil {
		log.Print(err)
		return
	}
	rsp.Buffer = string(jsonString)
	return &rsp, nil
}

func RunServer() {
	handler := &FormatDataImpl
	processor := lock.NewGetLockProcessor(handler)
	serverTransport, err := thrift.NewTServerSocket(HOST + ":" + PORT)
	if err != nil {
		log.Fatal(err)
	}
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)
	fmt.Println("Start running server: ", HOST+":"+PORT)
	server.Serve()
}

func main() {
	MapRWMutex = new(sync.RWMutex)
	Runserver()
}
