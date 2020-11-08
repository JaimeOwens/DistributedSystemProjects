package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	process1, err := LoadProcFromFile("Config/Process1.json")
	if err != nil {
		log.Fatal(err)
		return
	}
	process2, err := LoadProcFromFile("Config/Process2.json")
	if err != nil {
		log.Fatal(err)
		return
	}

	go process1.Loop()
	go process2.Loop()

	time.Sleep(1 * time.Second)
	process1.SendMessage()
	fmt.Print("s:start snapshot\nc:clear record\n")

	for {
		input := bufio.NewReader(os.Stdin)
		cmd, _ := input.ReadString('\n')
		if cmd[:len(cmd)-1] == "s" {
			process1.SendMarker()
			process2.getAmount()
		}
		if cmd[:len(cmd)-1] == "c" {
			process1.record = false
			process1.record = false
		}
	}
	return
}
