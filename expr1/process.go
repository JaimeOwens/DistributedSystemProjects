package main

import (
	"log"
	"net"
	"strings"
	"time"
)

type Process struct {
	localID    string
	localAddr  string
	remotePros []Channel
	amount     int
	record     bool
}

type Channel struct {
	RemoteID   string
	RemoteAddr string
}

func NewProcess(config Config) Process {
	return Process{
		localID:    config.LocalID,
		localAddr:  config.LocalAddr,
		remotePros: config.RemotePros,
		amount:     101,
		record:     false,
	}
}

func (process *Process) getAmount() {
	log.Print("save snapshot in ", process.localAddr, " amount=", process.amount)
}

func (process *Process) SendMarker() {
	process.record = true
	process.getAmount()
	for _, remotePro := range process.remotePros {
		conn, err := net.Dial("tcp", remotePro.RemoteAddr)
		defer conn.Close()
		if err != nil {
			log.Print(err)
			return
		}
		_, err = conn.Write([]byte(process.localID + "#marker"))
		if err != nil {
			log.Print(err)
			return
		}
	}
}

func (process *Process) SendMessage() {
	for _, remotePro := range process.remotePros {
		conn, err := net.Dial("tcp", remotePro.RemoteAddr)
		defer conn.Close()
		if err != nil {
			log.Print(err)
			return
		}
		_, err = conn.Write([]byte(process.localID + "#message"))
		if err != nil {
			log.Print(err)
			return
		}
	}
}

func (process *Process) receiveAll(listener net.Listener) (bool, error) {
	conn, err := listener.Accept()
	if err != nil {
		log.Print(err)
		return false, err
	}

	buffer := make([]byte, 1024)
	length, err := conn.Read(buffer)
	if err != nil {
		return false, err
	}
	buffer = buffer[:length]

	split := strings.Split(string(buffer), "#")
	if len(split) < 2 {
		return false, err
	}

	msg := split[1]
	if msg == "message" {
		process.amount++
		return true, nil
	} else if msg == "marker" && process.record == false {
		process.SendMarker()
		return false, nil
	}

	return false, nil
}

func (process *Process) Loop() {
	listener, err := net.Listen("tcp", process.localAddr)
	if err != nil {
		log.Print(err.Error())
		return
	}
	for {
		flag, err := process.receiveAll(listener)
		if err != nil {
			continue
		}
		if flag == true {
			time.Sleep(1 * time.Second)
			process.SendMessage()
		}
	}
}
