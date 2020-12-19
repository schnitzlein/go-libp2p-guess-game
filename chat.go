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
 *
 */
package main

import (
        "encoding/json"
        "io/ioutil"
        //"guessgame" broken import shit with go everything relativ to online? wtf

        "time"
        "strings"
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

//
// simple game beginn
//
func gameRules() String{
  var textblob string = ""
  textblob += "Guess the correct Number Game"
  textblob += "---------------------"
  textblob += "Game Rules: "
  textblob += "  Guess the correct number by input a number."
  textblob += "  The number must between 0 and 100."
  textblob += "  The machine answer only with two indicators."
  textblob += "  Indicators ::== [ above | below ] "
  textblob += "  If the number is lower than your current guess,"
  textblob += "  than the 'below' indicator is shown."
  textblob += "  Is the number above your current guess,"
  textblob += "  than the 'above' indicator is shown."
  textblob += ""
  textblob += "Control Commands: "
  textblob += "  exit"
  textblob += "  help"

}

func isInt(s string) bool {
    _, err := strconv.ParseInt(s, 10, 64)
    return err == nil // err == nil means there is no error, so True means is valid Integer
}

func runGame() {
  var limit int = 100

  // use always a new Seed
  rand.Seed( time.Now().UnixNano() )
  var my_secret_number int = rand.Intn( limit )

  // todo: share this secret and post it to all new peers
  // todo: awaiting loop in main
  // automatic loop in peers for receiving a number

  reader := bufio.NewReader(os.Stdin)

  gameRules()

  // endless loop
  for {
    fmt.Println("guess a number between 0 and ", limit)
    fmt.Print("-> ")
    text, _ := reader.ReadString('\n')
    // convert CRLF to LF
    text = strings.Replace(text, "\n", "", -1)

    foo := isInt(text)
    if foo == true {
      guess, _ := strconv.Atoi(text)

      if guess < my_secret_number {
          fmt.Println("above!")
      } else if guess > my_secret_number {
          fmt.Println("below!")
      } else if guess == my_secret_number {
          fmt.Println("================")
          fmt.Println("=== success! ===")
          fmt.Println("================")
          os.Exit(0)
      }
    } else {
      // This is a text/string
      switch text  {
        case "foobar":
          fmt.Println("foobar yes")
        case "barfoo":
          fmt.Println("yes barfoo!")
        //case "debug":
        //  fmt.Println("secret", my_secret_number)
        case "exit":
          os.Exit(0)
        case "help":
          gameRules()
        default:
         fmt.Println("something completly different...")
      }
    }

  } // for
}
//
// simple game ends
//

type MainNode struct {
     NodeName, NodeID string
}

type MagicNumber struct {
     secret int
}

func handleStream(s network.Stream) {
	log.Println("Got a new stream!")
        log.Println(s.Stat())

	// Create a buffer stream for non blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

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

func writeDataPeer(rw *bufio.ReadWriter, NodeName string) {
         stdReader := bufio.NewReader(os.Stdin)
         for {

           // add here the guess function
           game.runGame()
           fmt.Print(NodeName, " > ")
           var sendData = "after " + string(2) + "secs, "
           foo, err := stdReader.ReadString('\n')
           // above fills from guesser

           if err != nil {
               panic(err)
           }
           sendData += foo

           // send after Timer fired
           timer1 := time.NewTimer(2 * time.Second)
           <-timer1.C


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

        // my stuff
        data := MainNode {
          NodeName: "MyNodeName",
          NodeID: host.ID().Pretty(),
        }
        //
        fmt.Println("This is nodeid: ",host.ID().Pretty())
        fmt.Println(data) // useless data

        file, _ := json.MarshalIndent(data, "", " ")

	_ = ioutil.WriteFile("test.json", file, 0644)

        //

	if *dest == "" {
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

		fmt.Printf("Run './chat -d /ip4/127.0.0.1/tcp/%v/p2p/%s' on another console.\n", port, host.ID().Pretty())
		fmt.Println("You can replace 127.0.0.1 with public IP as well.")
		fmt.Printf("\nWaiting for incoming connection\n\n")

		// Hang forever
		<-make(chan struct{})
	} else {
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

		// Create a thread to read and write data.

                fmt.Println("You are: ",host.ID().Pretty())
		go writeDataPeer(rw, host.ID().Pretty())
		go readData(rw, "other")

		// Hang forever.
		select {}
	}
}
