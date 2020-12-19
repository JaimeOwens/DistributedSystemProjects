package main

import (
	"DistributedSystemProjects/expr2/lockrpc"
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"git.apache.org/thrift.git/lib/go/thrift"
)

const (
	HOST = "localhost"
	PORT = "6790"
)

var RPCClient *lockrpc.GetLockClient

package main

import (
	"../thrift/lockrpc"
	"log"
)

func MultiDoLock(num int64){
	for i:=0;i<int(num);i++{				
		req := lockrpc.Req{CliID:int64(i),Operator:"acquire"}
		rsp, err := RPCClient.DoLock(&req); if err != nil {
			log.Print(err, rsp)
			return
		}
		log.Print( "ClockID:",rsp.CliID," State:",rsp.Operator)
	}
}
func MultiUnLock(num int64){
	rspList := make([]*lockrpc.Rsp,num)
	for i:=0;i<int(num);i++{
		req := lockrpc.Req{CliID:num,Operator:"release"}
		rsp, err := RPCClient.UnLock(&req); if err != nil{
			rspList[i] = rsp
		}
	}
	log.Print("State:",rspList[0].Operator," clientID:",rspList[0].CliID)
}


func UnLockOne(ID int64){
	req := lockrpc.Req{CliID:ID,Operator:"release"}
	rsp, err := RPCClient.UnLock(&req);if err != nil{
		log.Fatal(err)
		return
	}
	log.Print("State:",rsp.Operator," clientID:",rsp.CliID)
}

func main() {
	tSocket, err := thrift.NewTSocket(net.JoinHostPort(HOST, PORT))
	if err != nil {
		log.Fatalln("tSocket error:", err)
	}
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	transport := transportFactory.GetTransport(tSocket)
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	RPCClient = lockrpc.NewGetLockClientFactory(transport, protocolFactory)
	if err := transport.Open(); err != nil {
		log.Fatalln("Error opening:", HOST+":"+PORT)
	}
	defer transport.Close()

	fmt.Print("e.g. acquire,5(不同client请求次数)\n")
	fmt.Print("e.g. release,1(relase clientID)\n")
	fmt.Print("e.g. status\n")
	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		text = text[:len(text)-1]
		textList := strings.Split(text, ",")
		if len(textList) == 1 {
			operator := strings.Split(text, ",")[0]
			operator = strings.Replace(operator, " ", "", -1)
			if operator == "status" {
				rsp, err := RPCClient.ClientStates()
				if err == nil {
					log.Print(rsp.Buffer)
				}
			}
			continue
		} else if len(textList) == 2 {
			operator := strings.Split(text, ",")[0]
			operator = strings.Replace(operator, " ", "", -1)
			ID := strings.Split(text, ",")[1]
			ID = strings.Replace(ID, " ", "", -1)
			num, err := strconv.ParseInt(ID, 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			if operator == "acquire" {
				MultiDoLock(num)
				continue
			}
			if operator == "release" {
				UnLockOne(num)
				continue
			}
		}

		fmt.Print("input correct cmd\n")

	}
	return
}
