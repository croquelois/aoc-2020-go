package main

import (
  "fmt"
  "bufio"
  "os"
  "strconv"
  "runtime/pprof"
  "flag"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var verbose = false

func itoa(i int) string {
  return strconv.Itoa(i)
}

func atoi(str string) int {
  value, err := strconv.Atoi(str)
  if err != nil { panic("Can't parse the number") }
  return value
}

func pause(prompt string){
  fmt.Println(prompt)
  input := bufio.NewScanner(os.Stdin)
  input.Scan()
}

func encKey(subject int, loopSize int) int {
  cur := 1
  for i:=0;i<loopSize;i++ {
    cur = (cur * subject) % 20201227
  }
  return cur
}

func algo1(subject int, card int, door int) int {
  cur := 1
  curLoop := 0
  cardLoop := -1
  doorLoop := -1
  for {
    cur = (cur * subject) % 20201227
    curLoop++
    if cardLoop <= 0 && card == cur { cardLoop = curLoop }
    if doorLoop <= 0 && door == cur { doorLoop = curLoop }
    if cardLoop >= 0 && doorLoop >= 0 { return encKey(card, doorLoop) }
  }
}

func assert(v bool, msg string){
  if !v {
    fmt.Printf("the test fail: %s\n", msg)
    panic("test failed !")
  }
}

func assert_eq(v int, e int, msg string){
  if v != e {
    fmt.Printf("the test '%s' fail, it give %d instead of %d\n", msg, v, e)
    panic("test failed !")
  }
}

func test1_1() {
  assert_eq(algo1(7,5764801,17807724), 14897079, "1");
}

func question1() int {
  return algo1(7,5099500,7648211);
}

func main() {
  flag.Parse()
  if *cpuprofile != "" {
      f, err := os.Create(*cpuprofile)
      if err != nil { panic(err) }
      pprof.StartCPUProfile(f)
      defer pprof.StopCPUProfile()
  }
  test1_1()
  fmt.Printf("Question1: %d\n", question1())
}