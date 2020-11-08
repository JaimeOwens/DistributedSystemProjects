package main

import (
	"strconv"
	"bufio"
	"expr2/thrift/lock"
	"log"
	"net"
	"os"
	"strings"

	"github.com/apache/thrift/lib/go/thrift"
)

const (
	HOST = "127.0.0.1"
	POST = "9000"
)

var Client *lock.GetLockClient

func main() {
	tSocket, err := thrift.NewTSocket(net.JoinHostPort(HOST, PORT))
	if err != nil {
		log.Fatal(err)
	}
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	transport := transportFactory.GetTransport(tSocket)
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	Client = lock.NewGetLockClientFactory(transport, protocolFactory)
	err := transport.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer transport.Close()

	for {
		input := bufio.NewReader(os.Stdin)
		text, _ := input.ReadString('\n')
		text = text[:len(text)-1]
		textList := strings.Split(text, " ")

		if len(textList) == 1 {
			operator := textList[0]
			rsp, err := Client.ClientStates()
			if err != nil {
				log.Print(rsp.Buffer)
				log.Print(err)
			}
			continue
		} 
		
		if len(textList) == 2 {
			operator := textList[0]
			ID : =textList[1]
			num, err := strconv.ParseInt(ID, 10, 64)
			if err != nil {
				log.Print(err)
				continue
			}
			if operator == "acquire" {
				DoLock(num)
				continue
			}
			if operator == "release" {
				UnLock(num);
				continue
			}
		}
		
		fmt.Print("incorrect cmd\n")
	}
	return
}


func DoLock(ID int64) {
	req := lock.Request{CliID:int64(ID), Operator:"acquire"}
	rsp, err := Client.Dolock(&req)
	if err != nil {
		log.Print(err)
		log.Print(rsq)
		return
	}
	log.Print("ClientID:", rsp.CliID, "State:", rsp.Operator)
}

func UnLock(num int64) {
	req := lock.Request{CliID:int64(ID), Operator:"release"}
	rsp, err := Client.Unlock(&req)
	if err != nil {
		log.Print(err)
		log.Print(rsq)
		return
	}	
	log.Print("ClientID:", rsp.CliID, "State", rsq.Operator)
}