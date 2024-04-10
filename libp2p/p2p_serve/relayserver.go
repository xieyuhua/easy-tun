// relayserver.go
package main

import (
	"log"
	"fmt"
	"strings"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/relay"
)

func main() {
	run()
}

func run() {
// 	relay1, err := libp2p.New()
	relay1, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/8006") )
	if err != nil {
		log.Printf("Failed to create relay1: %v", err)
		return
	}

	_, err = relay.New(relay1)
	if err != nil {
		log.Printf("Failed to instantiate the relay: %v", err)
		return
	}
	
    replaced := strings.Replace( fmt.Sprintf("%v",relay1.Addrs()) , " ", "\n", -1)
    
	log.Printf("relay1Info ID: %v \n Addrs: \n %s",relay1.ID(), replaced )

	
	select {}
}