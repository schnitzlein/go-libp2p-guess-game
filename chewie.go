package main

import (
   "fmt"
   "math/rand"
   "time"
)

func helloworld(say string) {
  fmt.Println(say)
}

func getRandom(limit int) int {
  s1 := rand.NewSource(time.Now().UnixNano())
  r1 := rand.New(s1)
  retVal := r1.Intn(limit)
  return retVal
}
