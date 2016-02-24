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
	object := networkModule.Udp_message{Raddr: "broadcast", Data: 0, Length: 0}

	masterPort := 20011
	slavePort := 20012

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
			object.Data = msg.Data //strconv.Itoa
			fmt.Println("recieving")
			fmt.Println(strconv.Itoa(msg.Data)
			lastRecieve = time.Now()
		default:
			if time.Since(lastRecieve).Seconds() > 1.0 {
				fmt.Println("Taking ovah as master")
				master = true
			}
		}
	}
	//send_ch <- networkModule.Udp_message{Raddr: "", Data: "terminate", Length: 9}
	object.Raddr = "terminate"
	send_ch <- object
	//recieve_ch <- mg

	
	fmt.Println("masterStart")
	err = networkModule.Udp_init(masterPort, slavePort, 1024, send_ch, recieve_ch)
	if err != nil {
		fmt.Print("main done. err = %s \n", err)
	}

	cmd := exec.Command("gnome-terminal", "-x", "go", "run", "processPairs.go")
	cmd.Run()
	
	num := object.Data
	msg := networkModule.Udp_message{Raddr: "broadcast", Data: object.Data, Length: 0}
	for {
		fmt.Println("object data" + object.Data)
		msg.Data = strconv.Itoa(num)
		msg.Length = len(msg.Data)
		send_ch <- msg
		fmt.Println("sennding" + msg.Data)
		num++
		time.Sleep(time.Second * 1)
	}
}
