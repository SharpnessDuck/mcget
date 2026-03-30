package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/Tnze/go-mc/bot"
)

type response struct {
	Players struct {
		Max    int `json:"max"`
		Online int `json:"online"`
	} `json:"players"`
	Version struct {
		Name string `json:"name"`
	} `json:"version"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("write ip address")
		fmt.Println("usage: go run app.go <ip>")
		return
	}

	input := os.Args[1]
	var fullAddress string

	_, addrs, err := net.LookupSRV("minecraft", "tcp", input)
	if err == nil && len(addrs) > 0 {
		fullAddress = net.JoinHostPort(addrs[0].Target, fmt.Sprint(addrs[0].Port))
	} else {
		ip, port, err := net.SplitHostPort(input)
		if err != nil {
			ip = input
			port = "25565"
		}
		fullAddress = net.JoinHostPort(ip, port)
	}

	fmt.Printf("checking: %s...\n", fullAddress)

	resp, delay, err := bot.PingAndList(fullAddress)
	if err != nil {
		fmt.Printf("error: could not connect to %s\n", fullAddress)
		return
	}

	var status response
	if err := json.Unmarshal(resp, &status); err != nil {
		fmt.Println("error: server responded but sent incorrect information")
		return
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("ip:       %s\n", fullAddress)
	fmt.Printf("ping:     %v\n", delay.Truncate(time.Millisecond))
	fmt.Printf("online:   %d / %d\n", status.Players.Online, status.Players.Max)
	fmt.Printf("version:  %s\n", status.Version.Name)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}
