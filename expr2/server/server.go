package main

import (
	"../lockrpc"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sync"

	"git.apache.org/thrift.git/lib/go/thrift"
)

const (
	HOST = "localhost"
	PORT = "9000"
)

var rwmutex *sync.RWMutex

var locktable map[int64]bool = make(map[int64]bool, 0)

func Addlock(id int64) {
	_, exits := locktable[id]
	if exits {
		return
	} else {
		rwmutex.Lock()
		locktable[id] = false
		rwmutex.Unlock()
		log.Print("add new client:", id)
	}
}

func Islock() bool {
	flag := false
	rwmutex.RLock()
	for _, v := range locktable {
		if v == true {
			flag = true
			break
		}
	}
	rwmutex.RUnlock()
	return flag
}

func Randomlock() int64 {
	var changeID int64 = -1
	if Islock() {
		return changeID
	}
	i := rand.Intn(len(locktable))
	rwmutex.RLock()
	for k := range locktable {
		if i == 0 {
			changeID = k
			break
		}
		i--
	}
	rwmutex.RUnlock()
	if changeID != -1 {
		rwmutex.Lock()
		locktable[changeID] = true
		rwmutex.Unlock()
	}
	return changeID
}

func Getlock() int64 {
	var res int64 = -1
	rwmutex.RLock()
	for key, val := range locktable {
		if val == true {
			res = key
		}
	}
	rwmutex.RUnlock()
	return res
}

func UnlockClient(id int64) {
	locktable[id] = false
}

type LockRPCImpl struct{}

func (lockrpcimpl *LockRPCImpl) DoLock(req *lockrpc.Req) (*lockrpc.Rsp, error) {
	Addlock(req.CliID)
	Randomlock()
	lockID := Getlock()
	var rsp lockrpc.Rsp
	rsp.CliID = lockID
	rsp.Operator = "acq"
	log.Print("DoLock ClientID=", rsp.CliID)
	return &rsp, nil
}

func (lockrpcimpl *LockRPCImpl) UnLock(req *lockrpc.Req) (*lockrpc.Rsp, error) {
	UnlockClient(req.CliID)
	var rsp lockrpc.Rsp
	rsp.CliID = req.CliID
	rsp.Operator = "rls"
	log.Print("UnLock ClientID=", rsp.CliID)
	return &rsp, nil
}

func (lockrpcimpl *LockRPCImpl) ClientStates() (*lockrpc.Rsp, error) {
	var rsp lockrpc.Rsp
	jsonString, err := json.Marshal(locktable)
	if err == nil {
		rsp.Buffer = string(jsonString)
	}
	return &rsp, nil
}

func StartServer() {
	handler := &LockRPCImpl{}
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
	rwmutex = new(sync.RWMutex)
	StartServer()
}
