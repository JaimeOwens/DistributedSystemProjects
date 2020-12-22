package main

import (
	"../lockrpc"
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
	PORT = "9000"
)

var RPCClient *lockrpc.GetLockClient

func ClientDoLock(num int64){
	for i:=0;i<int(num);i++{				
		req := lockrpc.Req{CliID:int64(i),Operator:"acq"}
		rsp, err := RPCClient.DoLock(&req); 
		if err != nil {
			log.Print(err, rsp)
			return
		}
		log.Print("DoLock ClientID:",rsp.CliID)
	}
}

func ClientUnLock(ID int64){
	req := lockrpc.Req{CliID:ID,Operator:"rls"}
	rsp, err := RPCClient.UnLock(&req);
	if err != nil{
		log.Fatal(err)
		return
	}
	log.Print("UnLock ClientID:",rsp.CliID)
}

func main() {
	fmt.Print("acq,5(acquire times) //Acquire client's lock\n")
	fmt.Print("rls,1(clientID)      //Release client's lock\n")
	fmt.Print("sts                  //Check global locks' status\n")

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
	fmt.Println("Connecting to:", HOST+":"+PORT)

	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		text = text[:len(text)-1]
		textList := strings.Split(text[:len(text)-1], ",")
		operator := textList[0]
		
		if len(textList) == 1 {
			if operator == "sts" {
				rsp, err := RPCClient.ClientStates()
				if err == nil {
					log.Print(rsp.Buffer)
				}
			}
			continue
		} else if len(textList) == 2 {
			ID := textList[1]
			num, err := strconv.ParseInt(ID, 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			if operator == "acq" {
				ClientDoLock(num)
				continue
			}
			if operator == "rls" {
				ClientUnLock(num)
				continue
			}
		}
		fmt.Print("incorrect cmd\n")
	}
	return
}
