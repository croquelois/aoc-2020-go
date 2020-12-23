package main

import (
    "fmt"
    "io/ioutil"
    "strings"
    "strconv"
    "sort"
)

func parse(data string) []int {
  var arrStr = strings.Split(data, "\n")
  var arrInt = []int{}
  
  for _, s := range arrStr {
      i, err := strconv.Atoi(s)
      if err != nil {
          panic(err)
      }
      arrInt = append(arrInt, i)
  }
  return arrInt
}

func parseFile(filename string) []int {
  data, err := ioutil.ReadFile(filename)
  if err != nil {
    panic(err)
  }
  return parse(string(data))
}

func algo1(data []int) int {
  sort.Ints(data)
  n := 0
  count1 := 0
  count3 := 0
  for _, n1 := range data {
    diff := n1 - n
    if(diff == 1){ count1 += 1 }
    if(diff == 3){ count3 += 1 }
    n = n1
  }
  count3 += 1
  return count1 * count3 
}

type Chargers struct {
  present map[int]bool
  value map[int]int
  builtIn int
}

func Max(v []int) int {
  m := v[0]
  for i := 1; i < len(v); i++ {
      if v[i] > m {
          m = v[i]
      }
  }
  return m
}

func newChargers(data []int) Chargers {
  var c Chargers
  c.present = make(map[int]bool)
  c.value = make(map[int]int)
  c.present[0] = true
  for _, n := range data {
    c.present[n] = true
  }
  c.builtIn = Max(data)+3
  return c
}

func (c *Chargers) nbPossibility(n int) int {
  if n == c.builtIn { return 1 }
  _, ok := c.present[n]
  if !ok { return 0 }
  v, ok2 := c.value[n]
  if ok2 { return v }
  nv := c.nbPossibility(n+1)+c.nbPossibility(n+2)+c.nbPossibility(n+3)
  c.value[n] = nv
  return nv
}

func algo2(data []int) int {
  c := newChargers(data)
  return c.nbPossibility(0)
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
  assert_eq(algo1(parseFile("test1.txt")), 7*5, "1");
}

func test1_2() {
  assert_eq(algo1(parseFile("test2.txt")), 22*10, "2");
}

func question1() int {
  return algo1(parseFile("input.txt"));
}

func test2_1() {
  assert_eq(algo2(parseFile("test1.txt")), 8, "2-1");
}

func test2_2() {
  assert_eq(algo2(parseFile("test2.txt")), 19208, "2-2");
}

func question2() int {
  return algo2(parseFile("input.txt"));
}

func main() {
  test1_1()
  test1_2()
  fmt.Printf("Question1: %d\n", question1())
  test2_1()
  test2_2()
  fmt.Printf("Question2: %d\n", question2())
}
