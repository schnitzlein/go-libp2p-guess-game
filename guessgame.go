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
)

//
// simple game beginn
//
// in production code use context

var my_secret_number int = 0

func gameRules() string{
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
  return textblob
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
  // This is a int
  if foo == true {
      guess, _ := strconv.Atoi(text)
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
          answer = gameRules()
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

