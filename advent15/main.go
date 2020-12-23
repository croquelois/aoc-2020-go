package main

import (
    "fmt"
    "strings"
    "bufio"
    "os"
    "strconv"
)

func splitTrim(data string, sep string) []string {
  tmp := strings.Split(data, sep)
  for i := range tmp {
    tmp[i] = strings.TrimSpace(tmp[i]) 
  }
  return tmp;
}

func splitTrimInt(data string, sep string) []int {
  var arr = []int{}
  tmp := splitTrim(data, sep)
  for i := range tmp {
    value, err := strconv.Atoi(tmp[i])
    if err != nil {
        panic(err)
    }
    arr = append(arr, value)
  }
  return arr;
}

func parse(data string) []int {
  return splitTrimInt(data, ",")
}

func pause(prompt string){
  fmt.Println(prompt)
  input := bufio.NewScanner(os.Stdin)
  input.Scan()
}

func algo(nb []int, nbSpoken int) int {
  var ov, v int
  mem := make(map[int]int)
  for i:=0;i<nbSpoken;i++ {
    if i < len(nb) { v = nb[i] }
    age, ok := mem[v]
    ov = v
    if !ok {
      v = 0 
    } else {
      v = i - age
    }
    mem[ov] = i
  }
  return ov
}

func algo1(nb []int) int {
  return algo(nb, 2020)
}

func algo2(nb []int) int {
  return algo(nb, 30000000)
}

func assert(v bool, msg string){
  if !v {
    fmt.Printf("the test fail: %s\n", msg)
    panic("test failed !")
  }
}

func assert_eq(v int, e int, msg string){
  if v != e {
    fmt.Printf("the test %s fail, it give %d instead of %d\n", msg, v, e)
    panic("test failed !")
  }
}

func test1_1() {
  assert_eq(algo1(parse("0,3,6")), 436, "1");
}

func question1() int {
  return algo1(parse("16,11,15,0,1,7"));
}

func test2_1() {
  assert_eq(algo2(parse("0,3,6")), 175594, "2-1");
  assert_eq(algo2(parse("1,3,2")), 2578, "2-2");
  assert_eq(algo2(parse("2,1,3")), 3544142, "2-3");
  assert_eq(algo2(parse("1,2,3")), 261214, "2-4");
  assert_eq(algo2(parse("2,3,1")), 6895259, "2-5");
  assert_eq(algo2(parse("3,2,1")), 18, "2-6");
  assert_eq(algo2(parse("3,1,2")), 362, "2-7");
}

func question2() int {
  return algo2(parse("16,11,15,0,1,7"));
}

func main() {
  test1_1()
  fmt.Printf("Question1: %d\n", question1())
  test2_1()
  fmt.Printf("Question1: %d\n", question2())
}