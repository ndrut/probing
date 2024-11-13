package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

type Host struct {
	Name      string   `json:"name"`
	Addresses []string `json:"addresses"`
}

func pingAddress(name, address string) {
	pinger, err := probing.NewPinger(address)
	if err != nil {
		fmt.Printf("Failed to create pinger for %s (%s): %v\n", name, address, err)
		return
	}

	// Set privileged to true to allow sending ICMP packets.
	pinger.SetPrivileged(true)
	pinger.Interval = time.Second * 5 // 1-second interval between pings
	// Remove or comment out the following line:
	// pinger.Timeout = 0            // No overall timeout; run indefinitely
	pinger.Count = -1 // -1 for unlimited pings

	pinger.OnRecv = func(pkt *probing.Packet) {
		timestamp := time.Now().Format(time.RFC3339)
		fmt.Printf("[%s] Reply from %s (%s): RTT=%v\n", timestamp, name, pkt.IPAddr.String(), pkt.Rtt)
	}

	fmt.Printf("Starting to ping %s (%s)...\n", name, address)

	// Run the pinger in a separate goroutine since pinger.Run() is blocking
	go func() {
		err = pinger.Run()
		if err != nil {
			fmt.Printf("Error while pinging %s (%s): %v\n", name, address, err)
		}
	}()
}

func main() {
	// Read hosts from a JSON file named "hosts.json"
	data, err := ioutil.ReadFile("hosts.json")
	if err != nil {
		log.Fatal("Error reading hosts.json:", err)
	}

	var hosts []Host
	err = json.Unmarshal(data, &hosts)
	if err != nil {
		log.Fatal("Error parsing JSON:", err)
	}

	for _, host := range hosts {
		for _, address := range host.Addresses {
			// Start pinging each address
			pingAddress(host.Name, address)
		}
	}

	// Keep the main function running
	select {}
}
