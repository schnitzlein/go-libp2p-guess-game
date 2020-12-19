/*
 *
 * The MIT License (MIT)
 *
 * Copyright (c) 2020 Christoph Schwalbe
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
 * This program is part of the demonstration a simple chat application using p2p communication.
 *
 */
package main

import (
  "strings"
  "strconv"
  "fmt"
  "os"
  "time"
  mrand "math/rand"
  "bufio"
  //"log"
)

//
// simple game beginn
//
// in production code use context

var my_secret_number int = 0

func gameRules() {
  fmt.Println("Guess the correct Number Game")
  fmt.Println("---------------------")
  fmt.Println("Game Rules: ")
  fmt.Println("  Guess the correct number by input a number.")
  fmt.Println("  The number must between 0 and 100.")
  fmt.Println("  The machine answer only with two indicators.")
  fmt.Println("  Indicators ::== [ above | below ] ")
  fmt.Println("  If the number is lower than your current guess,")
  fmt.Println("  than the 'below' indicator is shown.")
  fmt.Println("  Is the number above your current guess,")
  fmt.Println("  than the 'above' indicator is shown.")
  fmt.Println("")
  fmt.Println("Control Commands: ")
  fmt.Println("  exit")
  fmt.Println("  help")
}

func createMagicNumber() int{
  var limit int = 100

  // use always a new Seed
  mrand.Seed( time.Now().UnixNano() )
  var my_secret_number = mrand.Intn( limit )
  return my_secret_number
}

func isInt(s string) bool {
    _, err := strconv.ParseInt(s, 10, 64)
    return err == nil // err == nil means there is no error, so True means is valid Integer
}


func parseInput(s string) string {
  retVal := strings.Replace(s, "\n", "", -1)
  return retVal
}

func showMenu(reader *bufio.Reader) string {
  fmt.Println("guess a number between 0 and 100")
  fmt.Print("-> ")
  text, _ := reader.ReadString('\n')
  return text
}

// Game sequence diagram in text
// mainNode
// showMenu()
// parseInput()
// --> send to peerNode(s)
// <-- mainNode waits for answer
// (generell) if new node join --> send my_secret_number

// peer node
// (onetime at init) listen --> my_secret_number
// checkInput()
// --> send answer to mainNode

// todo: seperate this into pdu, and send sdu to higher layer
//       seperate the app code from p2p code
func checkInput(text string) string {
  var answer string = ""
  foo := isInt(text)
  //fmt.Println("DEBUG_text: ",text)
  //fmt.Println("DEBUG: ",foo)
  //fmt.Println("DEBUG_sec: ",my_secret_number)
  // This is a int
  if foo == true {
      guess, _ := strconv.Atoi(text)
      //fmt.Println("DEBUG_guess: ",guess)
      if guess < my_secret_number {
          answer = "above!"
      } else if guess > my_secret_number {
          answer = "below!"
      } else if guess == my_secret_number {
          answer = "=== success! ==="
          //os.Exit(0)
      }
    } else {
      // This is a text/string
      switch text  {
        case "foobar":
          answer = "foobar yes"
        case "barfoo":
          answer = "yes barfoo!"
        //case "debug":
        //  answer = string(my_secret_number)
        case "exit":
          os.Exit(0)
        case "help":
          gameRules()
        default:
          answer = "something completly different..."
      }
    }

  return answer
}

// onetime sender
func writeDataText(rw *bufio.ReadWriter, text string) {
     rw.WriteString(fmt.Sprintf("%s\n", text))
     rw.Flush()
}

// onetime receiver
func readDataText(rw *bufio.ReadWriter) {
  for {
     inbound, err := rw.ReadString('\n')
     if err != nil {
        panic(err)
     }
     if inbound == "" {
       //ignore
       return
     }
     if inbound != "\n" {
       //log.Fatalln("secret: ",inbound)
       foo := parseInput(inbound)
       my_secret_number, err = strconv.Atoi(foo)
       if err != nil {
          panic(err)
       }
       //fmt.Println("my_secret: ",my_secret_number)
       return
     }
  }

}
