package main

import (
	"./networkModule"
	"fmt"
	//"net"
	"os/exec"
	"strconv"
	"time"
)

func main() {
	send_ch := make(chan networkModule.Udp_message)
	recieve_ch := make(chan networkModule.Udp_message)

	masterPort := 20004
	slavePort := 20003

	master := false

	err := networkModule.Udp_init(slavePort, masterPort, 1024, send_ch, recieve_ch)
	if err != nil {
		fmt.Print("main done. err = %s \n", err)
	}

	fmt.Println("slaveStart")

	lastRecieve := time.Now()

	for !master {
		select {
		case msg := <-recieve_ch:
			fmt.Println("recieving")
			fmt.Println(msg.Data)
			lastRecieve = time.Now()
		default:
			if time.Since(lastRecieve).Seconds() > 3.0 {
				fmt.Println("Taking ovah as master")
				master = true
			}
		}
	}
	//send_ch <- networkModule.Udp_message{Raddr: "terminate", Data: "", Length: 0}

	fmt.Println("masterStart")
	err = networkModule.Udp_init(masterPort, slavePort, 1024, send_ch, recieve_ch)
	if err != nil {
		fmt.Print("main done. err = %s \n", err)
	}

	cmd := exec.Command("gnome-terminal", "-x", "go", "run", "processPairs.go")
	cmd.Run()

	num := 0
	msg := networkModule.Udp_message{Raddr: "broadcast", Data: "", Length: 0}
	for {
		msg.Data = strconv.Itoa(num)
		msg.Length = len(msg.Data)
		send_ch <- msg
		fmt.Println("sennding")
		num++
		time.Sleep(time.Second * 1)
	}
}
