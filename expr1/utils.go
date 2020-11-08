package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type Config struct {
	LocalID    string
	LocalAddr  string
	RemotePros []Channel
}

func LoadProcFromFile(filename string) (*Process, error) {
	configFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	var config Config
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	process := NewProcess(config)
	fmt.Println(process.localID, " ", process.localAddr, " ", process.remotePros, " ", process.amount)
	return &process, nil
}
