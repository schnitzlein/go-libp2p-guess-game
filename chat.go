/*
 *
 * The MIT License (MIT)
 *
 * Copyright (c) 2014 Juan Batiz-Benet
 *               2020 Christoph Schwalbe
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 *
 * This program demonstrate a simple chat application using p2p communication.
 * With backoff Timer demostration and uses Juan Batiz-Benet Code as Code base from 2014.
 * Thanks to you Juan Batiz-Benet.
 * Intention here was to use libp2p and transfer the secret or config parameters for a new node.
 * The game could also or simplier done with a simple socket communication ...
 *
 */
package main

import (
        "encoding/json"
        "io/ioutil"

        "time"
        //"strings"
        "strconv"

	"bufio"
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"log"
	mrand "math/rand"
	"os"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"

	"github.com/multiformats/go-multiaddr"
)


// initial starting game config
func initGame() {
  my_secret_number = createMagicNumber()
}



type MainNode struct {
     NodeName, NodeID string
}

type MagicNumber struct {
     secret int
}

func handleStream(s network.Stream) {
	log.Println("Got a new stream!")
        log.Println(s.Stat())

        // send one time data
        //rw_init := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
        //go writeDataText(rw_init, string(my_secret_number))

	// Create a buffer stream for non blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

        //var mysendstr string = ""+ strconv.Itoa(77)
        go writeDataText(rw, strconv.Itoa(my_secret_number) )

	go readData(rw, "other")
	go writeData(rw, "main")

	// stream 's' will stay open until you close it (or the other side closes it).
}
func readData(rw *bufio.ReadWriter, NodeName string) {
	for {
		str, _ := rw.ReadString('\n')

		if str == "" {
			return
		}
		if str != "\n" {
			// Green console colour: 	\x1b[32m
			// Reset console colour: 	\x1b[0m
			fmt.Printf("%s: \x1b[32m%s\x1b[0m> ", NodeName, str)
		}

	}
}

func readDataPeer(rw *bufio.ReadWriter, NodeName string, wc chan string) {
        for {
                str, _ := rw.ReadString('\n')

                if str == "" {
                        return
                }
                if str != "\n" {
                        // Green console colour:        \x1b[32m
                        // Reset console colour:        \x1b[0m
                        fmt.Printf("%s: \x1b[32m%s\x1b[0m> ", NodeName, str)
                        val := parseInput( str )
                        retVal := checkInput( val )
                        wc <- retVal // drop data into channel
                }

        }
}


func writeData(rw *bufio.ReadWriter, NodeName string) {
	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(NodeName," > ")
		sendData, err := stdReader.ReadString('\n')

		if err != nil {
			panic(err)
		}

		rw.WriteString(fmt.Sprintf("%s\n", sendData))
		rw.Flush()
	}

}

func writeDataPeer(rw *bufio.ReadWriter, NodeName string, c chan string) {
         for {

           // add here the guess function
           // check if readerstream got new bytes?
           fmt.Print(NodeName, " > ")

           sendData := <- c // receive from channel
           //fmt.Println(sendData)
           rw.WriteString(fmt.Sprintf("%s\n", sendData))
           rw.Flush()
         }
}


func main() {
	sourcePort := flag.Int("sp", 0, "Source port number")
	dest := flag.String("d", "", "Destination multiaddr string")
	help := flag.Bool("help", false, "Display help")
	debug := flag.Bool("debug", false, "Debug generates the same node ID on every execution")

	flag.Parse()

	if *help {
		fmt.Printf("This program demonstrates a simple p2p chat application using libp2p\n\n")
		fmt.Println("Usage: Run './chat -sp <SOURCE_PORT>' where <SOURCE_PORT> can be any port number.")
		fmt.Println("Now run './chat -d <MULTIADDR>' where <MULTIADDR> is multiaddress of previous listener host.")

		os.Exit(0)
	}

	// If debug is enabled, use a constant random source to generate the peer ID. Only useful for debugging,
	// off by default. Otherwise, it uses rand.Reader.
	var r io.Reader
	if *debug {
		// Use the port number as the randomness source.
		// This will always generate the same host ID on multiple executions, if the same port number is used.
		// Never do this in production code.
		r = mrand.New(mrand.NewSource(int64(*sourcePort)))
	} else {
		r = rand.Reader
	}

	// Creates a new RSA key pair for this host.
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		panic(err)
	}

	// 0.0.0.0 will listen on any interface device.
	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", *sourcePort))

	// libp2p.New constructs a new libp2p Host.
	// Other options can be added here.
	host, err := libp2p.New(
		context.Background(),
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(prvKey),
	)

	if err != nil {
		panic(err)
	}

        // my stuff with Tests
        data := MainNode {
          NodeName: "MyNodeName",
          NodeID: host.ID().Pretty(),
        }
        // testing something ...
        fmt.Println("This is nodeid: ",host.ID().Pretty())
        fmt.Println(data) // useless data

        file, _ := json.MarshalIndent(data, "", " ")

	_ = ioutil.WriteFile("test.json", file, 0644)
        // testing ends

        // channels for go routine communication
        channel_peer_reader := make(chan string)
        //channel_peer_reader_notify := make(chan bool)


	if *dest == "" {
                // Main Node
		// Set a function as stream handler.
		// This function is called when a peer connects, and starts a stream with this protocol.
		// Only applies on the receiving side.
		host.SetStreamHandler("/chat/1.0.0", handleStream)

		// Let's get the actual TCP port from our listen multiaddr, in case we're using 0 (default; random available port).
		var port string
		for _, la := range host.Network().ListenAddresses() {
			if p, err := la.ValueForProtocol(multiaddr.P_TCP); err == nil {
				port = p
				break
			}
		}

		if port == "" {
			panic("was not able to find actual local port")
		}

                // init Config + Game Rules
                initGame()
                gameRules()
                //fmt.Println("TEST: ",my_secret_number)

		fmt.Printf("Run './chat -d /ip4/127.0.0.1/tcp/%v/p2p/%s' on another console.\n", port, host.ID().Pretty())
		fmt.Println("You can replace 127.0.0.1 with public IP as well.")
		fmt.Printf("\nWaiting for incoming connection\n\n")

		// Hang forever
		<-make(chan struct{})
	} else {
                // Peer Node
		fmt.Println("This node's multiaddresses:")
		for _, la := range host.Addrs() {
			fmt.Printf(" - %v\n", la)
		}
		fmt.Println()

		// Turn the destination into a multiaddr.
		maddr, err := multiaddr.NewMultiaddr(*dest)
		if err != nil {
			log.Fatalln(err)
		}

		// Extract the peer ID from the multiaddr.
		info, err := peer.AddrInfoFromP2pAddr(maddr)
		if err != nil {
			log.Fatalln(err)
		}

		// Add the destination's peer multiaddress in the peerstore.
		// This will be used during connection and stream creation by libp2p.
		host.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)

		// Start a stream with the destination.
		// Multiaddress of the destination peer is fetched from the peerstore using 'peerId'.
		s, err := host.NewStream(context.Background(), info.ID, "/chat/1.0.0")
		if err != nil {
			panic(err)
		}

		// Create a buffered stream so that read and writes are non blocking.
		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

                // initial state receiving config/secret
                go readDataText(rw)
                time.Sleep(1 * time.Second)  // dirty solution, use sync with mutex and lock or check if secret != 0
                // better here check if my_secret_number != 0 OR
                // use mutex with lock and block the time ...
                // solve the issue that another peer could steal the session
                // also there is a lot of possiblities to retrieve the secret
		// Create a thread to read and write data.

                fmt.Println("You are: ",host.ID().Pretty())
		go writeDataPeer(rw, host.ID().Pretty(), channel_peer_reader)
		go readDataPeer(rw, "other", channel_peer_reader)

		// Hang forever.
		select {}
	}
}
