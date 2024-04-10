//main.go

package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/client"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
)

func main() {
	dist := flag.String("d", "", "your fanal nodeID")
	sid  := flag.String("sid", "12D3KooWAxYf4KA2pEhwzoz4cs65biWREhqp6nTrdFTMNnY9vk1b", "relay1Info ID")
	ip  := flag.String("ip", "43.143.2.144/tcp/34301", "serve ip ")
	flag.Parse()
	run(dist, sid , ip)
}

func run(dist,relayServerID,ip *string) {
	unreachableNode, err := libp2p.New(
		libp2p.NoListenAddrs,
		// Usually EnableRelay() is not required as it is enabled by default
		// but NoListenAddrs overrides this, so we're adding it in explictly again.
		libp2p.EnableRelay(),
	)
	if err != nil {
		log.Printf("Failed to create unreachable1: %v", err)
		return
	}
	
	log.Printf("/ip4/%v/p2p/%v", *ip , *relayServerID)
	
	ms1, _ := ma.NewMultiaddr(fmt.Sprintf("/ip4/%v/p2p/%v", *ip , *relayServerID))
	relay1info, err := peer.AddrInfoFromP2pAddr(ms1)

	if err := unreachableNode.Connect(context.Background(), *relay1info); err != nil {
		log.Printf("Failed to connect unreachable1 and relay1: %v", err)
		return
	}
	_, err = client.Reserve(context.Background(), unreachableNode, *relay1info)
	if err != nil {
		log.Printf("unreachable2 failed to receive a relay reservation from relay1. %v", err)
		return
	}

	unreachableNode.SetStreamHandler("/chatStream", func(s network.Stream) {
		log.Println("Awesome! We're now communicating via the relay!")
		log.Println("Got a new stream! p2p ")

		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

		go readData(rw)
		go writeData(rw)
	})

	if *dist == "" {
		//说明本节点是等待被另一个节点call的，只需要与中专节点连接一下注册到即可
		//写一个简单的阻塞
		fmt.Printf("WatingNode is onLine!")
		fmt.Printf("use \" go run .\\main.go -d %v \" in other CLI", unreachableNode.ID())
		select {}
	}

	relayaddr, err := ma.NewMultiaddr("/p2p/" + relay1info.ID.String() + "/p2p-circuit/p2p/" + *dist)
	if err != nil {
		log.Println(err)
		return
	}

	useRelayToDistInfo, err := peer.AddrInfoFromP2pAddr(relayaddr)
	if err := unreachableNode.Connect(context.Background(), *useRelayToDistInfo); err != nil {
		log.Printf("Unexpected error here. Failed to connect unreachable1 and unreachable2: %v", err)
		return
	}

	log.Println("Yep, that worked!")

	s, err := unreachableNode.NewStream(network.WithUseTransient(context.Background(), "chatStream"), useRelayToDistInfo.ID, "/chatStream")
	if err != nil {
		log.Println("Whoops, this should have worked...: ", err)
		return
	}

	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	go writeData(rw)
	go readData(rw)

	select {}
}

func readData(rw *bufio.ReadWriter) {
	for {
		str, _ := rw.ReadString('\n')

		if str == "" {
			return
		}
		if str != "\n" {
			// Green console colour: 	\x1b[32m
			// Reset console colour: 	\x1b[0m
			fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
		}

	}
}

func writeData(rw *bufio.ReadWriter) {
	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			log.Println(err)
			return
		}

		rw.WriteString(fmt.Sprintf("%s\n", sendData))
		rw.Flush()
	}
}